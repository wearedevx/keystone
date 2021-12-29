// +build !test

package payment

import (
	stripe "github.com/stripe/stripe-go/v72"
	portalsession "github.com/stripe/stripe-go/v72/billingportal/session"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/invoice"
	"github.com/stripe/stripe-go/v72/sub"
	"github.com/stripe/stripe-go/v72/usagerecord"
	"github.com/stripe/stripe-go/v72/webhook"
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

func getNextInvoice(params *stripe.InvoiceParams) (*stripe.Invoice, error) {
	return invoice.GetNext(params)
}

func getSubscription(id string, params *stripe.SubscriptionParams) (*stripe.Subscription, error) {
	return sub.Get(id, params)
}

func cancelSubscription(id string, params *stripe.SubscriptionCancelParams) (*stripe.Subscription, error) {
	return sub.Cancel(id, params)
}

func constructWebhookEvent(payload []byte, header, secret string) (stripe.Event, error) {
	return webhook.ConstructEvent(payload, header, secret)
}
