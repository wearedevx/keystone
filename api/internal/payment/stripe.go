package payment

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	stripe "github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/invoice"
	"github.com/stripe/stripe-go/v72/sub"
	"github.com/stripe/stripe-go/v72/usagerecord"
	"github.com/stripe/stripe-go/v72/webhook"
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
	seats int,
) (url string, err error) {
	organizationID := strconv.FormatUint(uint64(organization.ID), 10)
	name := organization.Name
	email := organization.Owner.Email
	var cus *stripe.Customer
	var ses *stripe.CheckoutSession

	customerParams := stripe.CustomerParams{
		Name:  stripe.String(name),
		Email: stripe.String(email),
	}

	sessionParams := stripe.CheckoutSessionParams{
		SuccessURL:         stripe.String("localhost:9001/subscription/success?sesssion_id={CHECKOUT_SESSION_ID}"),
		CancelURL:          stripe.String("localhost:9001/subscription/canceled"),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price: stripe.String(stripeSubscriptionPrice),
			},
		},
		Params: stripe.Params{
			Metadata: map[string]string{
				"seats":        strconv.Itoa(seats),
				"organization": organizationID,
			},
		},
	}

	cus, err = customer.New(&customerParams)
	if err != nil {
		goto done
	}

	sessionParams.Customer = stripe.String(cus.ID)
	sessionParams.CustomerEmail = stripe.String(email)
	sessionParams.ClientReferenceID = stripe.String(organization.Name)

	ses, err = session.New(&sessionParams)
	if err != nil {
		goto done
	}

	url = ses.URL

done:
	return url, err
}

// HandleEvent returns a usable payment.Event from a webhook call
func (sp *stripePayment) HandleEvent(
	r *http.Request,
) (paymentEvent Event, err error) {
	var b []byte
	var event stripe.Event

	if r.Method != http.MethodPost {
		err = fmt.Errorf(
			"method %s not allowed %w",
			r.Method,
			ErrorNotAStripeEvent,
		)
		goto done
	}

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
		err = fmt.Errorf(
			"cannot construct event %w",
			ErrorNotAStripeEvent,
		)
		goto done
	}

	switch event.Type {
	case "checkout.session.completed":
		ses := event.Data.Object
		paymentEvent.Type = EventCheckoutComplete
		paymentEvent.SubscriptionID = SubscriptionID(ses["subscription"].(string))

	case "invoice.paid":
		invoice := event.Data.Object
		paymentEvent.Type = EventSubscriptionPaid
		paymentEvent.SubscriptionID = SubscriptionID(invoice["subscription"].(string))

	case "invoice.payment_failed":
		invoice := event.Data.Object
		paymentEvent.Type = EventSubscriptionUnpaid
		paymentEvent.SubscriptionID = SubscriptionID(invoice["subscription"].(string))

	case "customer.subscription.updated":
		subscription := event.Data.Object
		paymentEvent.SubscriptionID = SubscriptionID(subscription["id"].(string))
		status := subscription["status"].(string)

		switch {
		case status == "active" || status == "trialing" || status == "incomplete":
			paymentEvent.Type = EventSubscriptionPaid

		case status == "past_due" || status == "unpaid":
			paymentEvent.Type = EventSubscriptionUnpaid

		case status == "canceled":
			paymentEvent.Type = EventSubscriptionCancelled
		}

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
		out = SubscriptionStatusCancelled

	default:
		out = SubscriptionStatusUnpaid
	}

	return out
}

// Updates the subscription (changes the number of seats)
func (sp *stripePayment) UpdateSubscription(subscriptionID SubscriptionID, seats int) (err error) {
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
			Action:           stripe.String(string(stripe.UsageRecordActionSet)),
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
