//go:build test
// +build test

package payment

import (
	"encoding/json"
	"errors"
	"time"

	stripe "github.com/stripe/stripe-go/v73"
)

func newPortalSession(params *stripe.BillingPortalSessionParams) (*stripe.BillingPortalSession, error) {
	return &stripe.BillingPortalSession{
		Customer:  *params.Customer,
		ID:        "portal-session",
		ReturnURL: "http://portal-session-return.url",
		URL:       "http://portal-session.url",
	}, nil
}

func newSession(params *stripe.CheckoutSessionParams) (*stripe.CheckoutSession, error) {
	return &stripe.CheckoutSession{
		ID:  "session-id",
		URL: "http://session.url",
	}, nil
}

func newCustomer(params *stripe.CustomerParams) (*stripe.Customer, error) {
	return &stripe.Customer{
		ID: "customer-id",
	}, nil
}

func newUsageRecord(params *stripe.UsageRecordParams) (*stripe.UsageRecord, error) {
	return &stripe.UsageRecord{}, nil
}

func getNextInvoice(params *stripe.InvoiceUpcomingParams) (*stripe.Invoice, error) {
	return &stripe.Invoice{
		Lines: &stripe.InvoiceLineList{
			Data: []*stripe.InvoiceLine{{Quantity: 1}},
		},
	}, nil
}

func getSubscription(id string, params *stripe.SubscriptionParams) (*stripe.Subscription, error) {
	return &stripe.Subscription{
		Customer: &stripe.Customer{
			ID: "customer-id",
		},
		ID:     id,
		Status: "active",
		Items: &stripe.SubscriptionItemList{
			APIResource: stripe.APIResource{},
			ListMeta:    stripe.ListMeta{},
			Data: []*stripe.SubscriptionItem{{
				APIResource:       stripe.APIResource{},
				BillingThresholds: stripe.SubscriptionItemBillingThresholds{},
				Created:           0,
				Deleted:           false,
				ID:                "item-id",
				Metadata:          map[string]string{},
				Object:            "",
				Plan:              &stripe.Plan{},
				Price:             &stripe.Price{},
				Quantity:          0,
				Subscription:      "",
				TaxRates:          []*stripe.TaxRate{},
			}},
		},
	}, nil
}

func cancelSubscription(id string, params *stripe.SubscriptionCancelParams) (*stripe.Subscription, error) {
	return &stripe.Subscription{
		CurrentPeriodEnd:   time.Now().Add(30 * 24 * time.Hour).Unix(),
		CurrentPeriodStart: time.Now().Add(-24 * time.Hour).Unix(),
		Customer: &stripe.Customer{
			ID: "customer-id",
		},
		ID:     id,
		Status: "canceled",
	}, nil
}

func constructWebhookEvent(
	payload []byte,
	header,
	secret string,
) (_ stripe.Event, err error) {
	if string(payload) == "bad-event" {
		return stripe.Event{}, errors.New("webhook error")
	}

	if string(payload) == "bad-event-to-ignore" {
		err = errors.
			New("error: cannot unmarshal string into Go struct field Event.request of type stripe.EventRequest")
	}

	data := make(map[string]interface{})
	json.Unmarshal(payload, &data)

	event_type, ok := data["type"]
	if !ok {
		goto nothing
	}

	return stripe.Event{
		Data: &stripe.EventData{
			Object: data,
		},
		ID:   "event-id",
		Type: event_type.(string),
	}, err

nothing:
	return stripe.Event{
		Data: &stripe.EventData{},
		ID:   "event-id",
		Type: "nothing",
	}, err

}
