package payment

var stripeKey string

type stripePayment struct {
}

func NewStripePayment() Payment {
	return new(stripePayment)
}

// New Subscription should start a checkout process
// the quantity params should be the number of
// unique users associated with the organization
// that subscription is for
func (sp *stripePayment) NewSubscription(seats int) (subscription Subscription, err error) {
	return subscription, err
}

// Returns subcription information
func (sp *stripePayment) GetSubscription(
	subscriptionID SubscriptionID,
) (subscription Subscription, err error) {
	return subscription, err
}

// Updates the subscription (changes the number of seats)
func (sp *stripePayment) UpdateSubscription(subscriptionID SubscriptionID, seats int) error {
	return nil
}

// Cancels the subscription
func (sp *stripePayment) CancelSubscription(subscriptionID SubscriptionID) error {
	return nil
}
