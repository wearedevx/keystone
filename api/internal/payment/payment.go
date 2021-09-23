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
	SubscriptionStatusUnpaid    SubscriptionStatus = "unpaid"
	SubscriptionStatusPaid                         = "paid"
	SubscriptionStatusCancelled                    = "cancelled"
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
	StartCheckout(organization *models.Organization, seats int) (string, error)
	//
	HandleEvent(r *http.Request) (Event, error)
	//
	GetSubscription(subscriptionID SubscriptionID) (Subscription, error)
	// Updates the subscription (changes the number of seats
	UpdateSubscription(subscriptionID SubscriptionID, seats int) error
	// Cancels the subscription
	CancelSubscription(subscriptionID SubscriptionID) error
}

type EventType string

const (
	EventNothing               EventType = "nothing"
	EventCheckoutComplete                = "checkout.complete"
	EventSubscriptionPaid                = "subscription.paid"
	EventSubscriptionUnpaid              = "subscription.unpaid"
	EventSubscriptionCancelled           = "subscription.cancelled"
)

type Event struct {
	Type           EventType
	SubscriptionID SubscriptionID `json:"subscriptionId"`
}
