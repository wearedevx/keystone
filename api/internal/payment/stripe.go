// +build !test

package payment

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	stripe "github.com/stripe/stripe-go/v72"
	portalsession "github.com/stripe/stripe-go/v72/billingportal/session"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/invoice"
	"github.com/stripe/stripe-go/v72/sub"
	"github.com/stripe/stripe-go/v72/usagerecord"
	"github.com/stripe/stripe-go/v72/webhook"
	"github.com/wearedevx/keystone/api/internal/constants"
	"github.com/wearedevx/keystone/api/pkg/models"
)

var stripeKey string
var stripeWebhookSecret string
var stripeSubscriptionPrice string

var (
	ErrorNotAStripeEvent = errors.New("not a stripe event")
)

func init() {
	stripe.Key = stripeKey
	stripeSubscriptionPrice = "price_1JcAYwHJIcXvsTNiofdvI5VQ"
}

type stripePayment struct {
}

func NewStripePayment() Payment {
	return new(stripePayment)
}

// StartCheckout should start a checkout process
// the quantity params should be the number of
// unique users associated with the organization
// that subscription is for
func (sp *stripePayment) StartCheckout(
	organization *models.Organization,
	seats int64,
) (sessionID, customerID, url string, err error) {
	organizationID := strconv.FormatUint(uint64(organization.ID), 10)
	name := organization.Name
	email := organization.User.Email
	var cus *stripe.Customer
	var ses *stripe.CheckoutSession

	successUrl := fmt.Sprintf("https://%s/checkout-success?session_id={CHECKOUT_SESSION_ID}", constants.Domain)
	cancelUrl := fmt.Sprintf("https://%s/checkout-cancel?session_id={CHECKOUT_SESSION_ID}", constants.Domain)

	sessionParams := stripe.CheckoutSessionParams{
		SuccessURL:         stripe.String(successUrl),
		CancelURL:          stripe.String(cancelUrl),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price: stripe.String(stripeSubscriptionPrice),
			},
		},
		Params: stripe.Params{
			Metadata: map[string]string{
				"seats":        strconv.FormatInt(seats, 10),
				"organization": organizationID,
			},
		},
	}

	// Create a customer if there isn’t one already
	if organization.CustomerID == "" {
		customerParams := stripe.CustomerParams{
			Name:  stripe.String(name),
			Email: stripe.String(email),
		}

		cus, err = customer.New(&customerParams)
		if err != nil {
			goto done
		}

		customerID = cus.ID
	} else {
		customerID = organization.CustomerID
	}

	sessionParams.Customer = stripe.String(customerID)
	sessionParams.ClientReferenceID = stripe.String(
		strconv.FormatUint(uint64(organization.ID), 10),
	)

	ses, err = session.New(&sessionParams)
	if err != nil {
		goto done
	}

	sessionID = ses.ID
	url = ses.URL

done:
	return sessionID, customerID, url, err
}

// GetManagementLink returns a URL to a management page, for the user
// to cancel their subscription.
func (sp *stripePayment) GetManagementLink(
	organization *models.Organization,
) (url string, err error) {
	// Authenticate your user.
	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(organization.CustomerID),
		ReturnURL: stripe.String(fmt.Sprintf("https://%s/", constants.Domain)),
	}

	ps, err := portalsession.New(params)
	if err != nil {
		return "", err
	}

	url = ps.URL

	return url, nil
}

// HandleEvent returns a usable payment.Event from a webhook call
func (sp *stripePayment) HandleEvent(
	r *http.Request,
) (paymentEvent Event, err error) {
	var b []byte
	var event stripe.Event

	b, err = ioutil.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf(
			"cannot read event body %w",
			ErrorNotAStripeEvent,
		)
		goto done
	}

	event, err = webhook.ConstructEvent(b, r.Header.Get("Stripe-Signature"), stripeWebhookSecret)

	if err != nil {
		// That error can be safely ignored.
		if !strings.Contains(err.Error(), "cannot unmarshal string into Go struct field Event.request of type stripe.EventRequest") {
			err = fmt.Errorf(
				"cannot construct event %w",
				ErrorNotAStripeEvent,
			)
			goto done
		}
	}

	switch event.Type {
	case "checkout.session.completed":
		ses := event.Data.Object
		oid, _ := strconv.ParseUint(ses["client_reference_id"].(string), 10, 16)

		paymentEvent.Type = EventCheckoutComplete
		paymentEvent.OrganizationID = uint(oid)
		paymentEvent.SessionID = ses["id"].(string)
		paymentEvent.CustomerID = CustomerID(ses["customer"].(string))
		paymentEvent.SubscriptionID = SubscriptionID(ses["subscription"].(string))

	case "invoice.paid":
		invoice := event.Data.Object
		paymentEvent.Type = EventSubscriptionPaid
		paymentEvent.CustomerID = CustomerID(invoice["customer"].(string))
		paymentEvent.SubscriptionID = SubscriptionID(invoice["subscription"].(string))

	case "invoice.payment_failed":
		invoice := event.Data.Object
		paymentEvent.Type = EventSubscriptionUnpaid
		paymentEvent.CustomerID = CustomerID(invoice["customer"].(string))
		paymentEvent.SubscriptionID = SubscriptionID(invoice["subscription"].(string))

	case "customer.subscription.updated":
		subscription := event.Data.Object
		paymentEvent.CustomerID = CustomerID(subscription["customer"].(string))
		paymentEvent.SubscriptionID = SubscriptionID(subscription["id"].(string))
		status := subscription["status"].(string)

		switch {
		case status == "active" || status == "trialing" || status == "incomplete":
			paymentEvent.Type = EventSubscriptionPaid

		case status == "past_due" || status == "unpaid" || status == "incomplete_expired":
			paymentEvent.Type = EventSubscriptionUnpaid

		case status == "canceled":
			paymentEvent.Type = EventSubscriptionCanceled
		}

	case "customer.subscription.deleted":
		subscription := event.Data.Object
		paymentEvent.Type = EventSubscriptionCanceled
		paymentEvent.CustomerID = CustomerID(subscription["customer"].(string))
		paymentEvent.SubscriptionID = SubscriptionID(subscription["id"].(string))

	default:
		paymentEvent.Type = EventNothing
	}

done:
	return paymentEvent, err
}

// Returns subcription information
func (sp *stripePayment) GetSubscription(
	subscriptionID SubscriptionID,
) (subscription Subscription, err error) {
	s, err := sub.Get(string(subscriptionID), &stripe.SubscriptionParams{
		Params: stripe.Params{
			Expand: []*string{
				stripe.String("customer"),
			},
		},
	})

	seats, err := stripeGetSeats(s.ID)
	if err != nil {
		goto done
	}

	subscription.ID = SubscriptionID(s.ID)
	subscription.CustomerID = CustomerID(s.Customer.ID)
	subscription.Seats = seats
	subscription.Status = stripeSubscriptionStatus(s.Status)
	subscription.CurrentPeriodStart = time.Unix(s.CurrentPeriodStart, 0)
	subscription.CurrentPeriodEnd = time.Unix(s.CurrentPeriodEnd, 0)

done:
	return subscription, err
}

//
func stripeGetSeats(subscriptionID string) (seats int, err error) {
	var inv *stripe.Invoice
	params := stripe.InvoiceParams{
		Subscription: stripe.String(subscriptionID),
	}

	params.AddExpand("subscription")

	inv, err = invoice.GetNext(&params)
	if err != nil {
		goto done
	}

	for _, v := range inv.Lines.Data {
		seats = seats + int(v.Quantity)
	}

done:
	return seats, err
}

//
func stripeSubscriptionStatus(in stripe.SubscriptionStatus) (out SubscriptionStatus) {
	switch {
	case in == "active" || in == "trialing" || in == "incomplete":
		out = SubscriptionStatusPaid

	case in == "past_due" || in == "unpaid":
		out = SubscriptionStatusUnpaid

	case in == "canceled":
		out = SubscriptionStatusCanceled

	default:
		out = SubscriptionStatusUnpaid
	}

	return out
}

// Updates the subscription (changes the number of seats)
func (sp *stripePayment) UpdateSubscription(subscriptionID SubscriptionID, seats int64) (err error) {
	params := stripe.SubscriptionParams{}
	var s *stripe.Subscription
	var urParams *stripe.UsageRecordParams

	s, err = sub.Get(string(subscriptionID), &params)
	if err != nil {
		goto done
	}

	for _, si := range s.Items.Data {
		urParams = &stripe.UsageRecordParams{
			SubscriptionItem: stripe.String(si.ID),
			Quantity:         stripe.Int64(int64(seats)),
			Timestamp:        stripe.Int64(time.Now().Unix()),
			// no increment, the seat value is the new value
			Action: stripe.String(string(stripe.UsageRecordActionSet)),
		}

		_, err = usagerecord.New(urParams)
	}

done:
	return err
}

// Cancels the subscription
func (sp *stripePayment) CancelSubscription(subscriptionID SubscriptionID) (err error) {
	_, err = sub.Cancel(string(subscriptionID), nil)
	if err != nil {
		return err
	}

	return nil
}
