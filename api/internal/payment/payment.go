package payment

import "time"

type CustomerID string
type SubscriptionID string

type Subscription struct {
	ID                 SubscriptionID
	CustomerID         CustomerID
	Seats              int
	CurrentPeriodStart time.Time
	CurrentPeriodEnd   time.Time
}

type Payment interface {
	// New Subscription should start a checkout process
	// the quantity params should be the number of
	// unique users associated with the organization
	// that subscription is for
	NewSubscription(seats int) (Subscription, error)
	GetSubscription(subscriptionID SubscriptionID) (Subscription, error)
	// Updates the subscription (changes the number of seats
	UpdateSubscription(subscriptionID SubscriptionID, seats int) error
	// Cancels the subscription
	CancelSubscription(subscriptionID SubscriptionID) error
}
