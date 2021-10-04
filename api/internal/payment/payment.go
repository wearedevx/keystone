package payment

import (
	"net/http"
	"time"

	"github.com/wearedevx/keystone/api/pkg/models"
)

type CustomerID string
type SubscriptionID string

type SubscriptionStatus string

const (
	SubscriptionStatusUnpaid   SubscriptionStatus = "unpaid"
	SubscriptionStatusPaid                        = "paid"
	SubscriptionStatusCanceled                    = "canceled"
)

type Subscription struct {
	ID                 SubscriptionID
	CustomerID         CustomerID
	Seats              int
	Status             SubscriptionStatus
	CurrentPeriodStart time.Time
	CurrentPeriodEnd   time.Time
}

type Payment interface {
	// StartCheckout should start a checkout process
	// the quantity params should be the number of
	// unique users associated with the organization
	// that subscription is for
	// returns a session id, a url for the user proceed, and an error
	StartCheckout(organization *models.Organization, seats int64) (string, string, string, error)
	GetManagementLink(organization *models.Organization) (string, error)
	// For webhooks, returns a unified event from a request
	HandleEvent(r *http.Request) (Event, error)
	// returns subscription information
	GetSubscription(subscriptionID SubscriptionID) (Subscription, error)
	// Updates the subscription (changes the number of seats
	UpdateSubscription(subscriptionID SubscriptionID, seats int64) error
	// Cancels the subscription
	CancelSubscription(subscriptionID SubscriptionID) error
}

type EventType string

const (
	EventNothing              EventType = "nothing"
	EventCheckoutComplete               = "checkout.complete"
	EventSubscriptionPaid               = "subscription.paid"
	EventSubscriptionUnpaid             = "subscription.unpaid"
	EventSubscriptionCanceled           = "subscription.canceled"
)

type Event struct {
	Type           EventType      `json:"type"`
	OrganizationID uint           `json:"organization_id"`
	SessionID      string         `json:"session_id"`
	CustomerID     CustomerID     `json:"customer_id"`
	SubscriptionID SubscriptionID `json:"subscription_id"`
}
