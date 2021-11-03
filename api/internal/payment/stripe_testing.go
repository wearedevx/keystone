// +build test

package payment

import (
	"errors"
	"net/http"

	stripe "github.com/stripe/stripe-go/v72"
	"github.com/wearedevx/keystone/api/pkg/models"
)

var (
	stripeKey               string
	stripeWebhookSecret     string
	stripeSubscriptionPrice string
)

var ErrorNotAStripeEvent = errors.New("not a stripe event")

func init() {
	stripe.Key = stripeKey
	stripeSubscriptionPrice = "all_is_fake"
}

type stripePayment struct{}

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
	return "session_id", "customer_id", "http://fake_payment", nil
}

// GetManagementLink returns a URL to a management page, for the user
// to cancel their subscription.
func (sp *stripePayment) GetManagementLink(
	organization *models.Organization,
) (url string, err error) {
	return "http://fake_portal", nil
}

// HandleEvent returns a usable payment.Event from a webhook call
func (sp *stripePayment) HandleEvent(
	r *http.Request,
) (paymentEvent Event, err error) {
	return paymentEvent, err
}

// Returns subcription information
func (sp *stripePayment) GetSubscription(
	subscriptionID SubscriptionID,
) (subscription Subscription, err error) {
	return subscription, err
}

// Updates the subscription (changes the number of seats)
func (sp *stripePayment) UpdateSubscription(
	subscriptionID SubscriptionID,
	seats int64,
) (err error) {
	return nil
}

// Cancels the subscription
func (sp *stripePayment) CancelSubscription(
	subscriptionID SubscriptionID,
) (err error) {
	return nil
}
