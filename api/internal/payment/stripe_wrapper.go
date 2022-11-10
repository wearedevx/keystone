//go:build !test
// +build !test

package payment

import (
	stripe "github.com/stripe/stripe-go/v73"
	portalsession "github.com/stripe/stripe-go/v73/billingportal/session"
	"github.com/stripe/stripe-go/v73/checkout/session"
	"github.com/stripe/stripe-go/v73/customer"
	"github.com/stripe/stripe-go/v73/invoice"
	"github.com/stripe/stripe-go/v73/subscription"
	"github.com/stripe/stripe-go/v73/usagerecord"
	"github.com/stripe/stripe-go/v73/webhook"
)

func newPortalSession(params *stripe.BillingPortalSessionParams) (*stripe.BillingPortalSession, error) {
	return portalsession.New(params)
}

func newSession(params *stripe.CheckoutSessionParams) (*stripe.CheckoutSession, error) {
	return session.New(params)
}

func newCustomer(params *stripe.CustomerParams) (*stripe.Customer, error) {
	return customer.New(params)
}

func newUsageRecord(params *stripe.UsageRecordParams) (*stripe.UsageRecord, error) {
	return usagerecord.New(params)
}

func getNextInvoice(params *stripe.InvoiceUpcomingParams) (*stripe.Invoice, error) {
	return invoice.Upcoming(params)
}

func getSubscription(id string, params *stripe.SubscriptionParams) (*stripe.Subscription, error) {
	return subscription.Get(id, params)
}

func cancelSubscription(id string, params *stripe.SubscriptionCancelParams) (*stripe.Subscription, error) {
	return subscription.Cancel(id, params)
}

func constructWebhookEvent(payload []byte, header, secret string) (stripe.Event, error) {
	return webhook.ConstructEvent(payload, header, secret)
}
