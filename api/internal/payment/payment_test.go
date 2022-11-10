package payment

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	stripe "github.com/stripe/stripe-go/v72"
	"github.com/wearedevx/keystone/api/pkg/models"
)

func TestNewStripePayment(t *testing.T) {
	tests := []struct {
		name string
		want Payment
	}{
		{
			name: "instanciates-stripe",
			want: &stripePayment{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStripePayment(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStripePayment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stripePayment_StartCheckout(t *testing.T) {
	type args struct {
		organization *models.Organization
		seats        int64
	}
	tests := []struct {
		name           string
		sp             *stripePayment
		args           args
		wantSessionID  string
		wantCustomerID string
		wantURL        string
		wantErr        bool
	}{
		{
			name: "starts-checkout-session-creates-customer",
			sp:   &stripePayment{},
			args: args{
				organization: &models.Organization{
					ID:             12,
					Name:           "orga-12",
					Paid:           false,
					Private:        false,
					CustomerID:     "",
					SubscriptionID: "",
					UserID:         0,
					User: models.User{
						Email: "memberx@mail.com",
					},
				},
				seats: 0,
			},
			wantSessionID:  "session-id",
			wantCustomerID: "customer-id",
			wantURL:        "http://session.url",
			wantErr:        false,
		},
		{
			name: "starts-checkout-session-has-customer",
			sp:   &stripePayment{},
			args: args{
				organization: &models.Organization{
					ID:             12,
					Name:           "orga-12",
					Paid:           false,
					Private:        false,
					CustomerID:     "customer-12",
					SubscriptionID: "",
					UserID:         0,
					User: models.User{
						Email: "memberx@mail.com",
					},
				},
				seats: 0,
			},
			wantSessionID:  "session-id",
			wantCustomerID: "customer-12",
			wantURL:        "http://session.url",
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := &stripePayment{}
			gotSessionID, gotCustomerID, gotURL, err := sp.StartCheckout(tt.args.organization, tt.args.seats)
			if (err != nil) != tt.wantErr {
				t.Errorf("stripePayment.StartCheckout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSessionID != tt.wantSessionID {
				t.Errorf("stripePayment.StartCheckout() gotSessionID = %v, want %v", gotSessionID, tt.wantSessionID)
			}
			if gotCustomerID != tt.wantCustomerID {
				t.Errorf("stripePayment.StartCheckout() gotCustomerID = %v, want %v", gotCustomerID, tt.wantCustomerID)
			}
			if gotURL != tt.wantURL {
				t.Errorf("stripePayment.StartCheckout() gotURL = %v, want %v", gotURL, tt.wantURL)
			}
		})
	}
}

func Test_stripePayment_GetManagementLink(t *testing.T) {
	type args struct {
		organization *models.Organization
	}
	tests := []struct {
		name    string
		sp      *stripePayment
		args    args
		wantURL string
		wantErr bool
	}{
		{
			name: "",
			sp:   &stripePayment{},
			args: args{
				organization: &models.Organization{
					CustomerID: "customer-id",
				},
			},
			wantURL: "http://portal-session.url",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := &stripePayment{}
			gotURL, err := sp.GetManagementLink(tt.args.organization)
			if (err != nil) != tt.wantErr {
				t.Errorf("stripePayment.GetManagementLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotURL != tt.wantURL {
				t.Errorf("stripePayment.GetManagementLink() = %v, want %v", gotURL, tt.wantURL)
			}
		})
	}
}

func Test_stripePayment_HandleEvent(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name             string
		sp               *stripePayment
		args             args
		wantPaymentEvent Event
		wantErr          bool
	}{
		{
			name: "ignored-event",
			sp:   &stripePayment{},
			args: args{
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBuffer([]byte{})),
				},
			},
			wantPaymentEvent: Event{
				Type:           "nothing",
				OrganizationID: 0,
				SessionID:      "",
				CustomerID:     "",
				SubscriptionID: "",
			},
			wantErr: false,
		},
		{
			name: "error-if-bad-event",
			sp:   &stripePayment{},
			args: args{
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBufferString("bad-event")),
				},
			},
			wantPaymentEvent: Event{},
			wantErr:          true,
		},
		{
			name: "ignores-a-specific-parsing-error",
			sp:   &stripePayment{},
			args: args{
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBufferString("bad-event-to-ignore")),
				},
			},
			wantPaymentEvent: Event{
				Type:           "nothing",
				OrganizationID: 0,
				SessionID:      "",
				CustomerID:     "",
				SubscriptionID: "",
			},
			wantErr: true,
		},
		{
			name: "handles-checkout.session.completed",
			sp:   &stripePayment{},
			args: args{
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBufferString(`
{
    "type": "checkout.session.completed",
    "client_reference_id": "12",
    "id": "session-id",
    "customer": "cus_srietnsirent",
    "subscription": "sub_auesrntausrent"
}
`)),
				},
			},
			wantPaymentEvent: Event{
				Type:           EventCheckoutComplete,
				OrganizationID: 12,
				SessionID:      "session-id",
				CustomerID:     "cus_srietnsirent",
				SubscriptionID: "sub_auesrntausrent",
			},
			wantErr: false,
		},
		{
			name: "handles-invoice.paid",
			sp:   &stripePayment{},
			args: args{
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBufferString(`
{
    "type": "invoice.paid",
    "customer": "cus_srietnsirent",
    "subscription": "sub_auesrntausrent"
}
`)),
				},
			},
			wantPaymentEvent: Event{
				Type:           EventSubscriptionPaid,
				OrganizationID: 0,
				SessionID:      "",
				CustomerID:     "cus_srietnsirent",
				SubscriptionID: "sub_auesrntausrent",
			},
			wantErr: false,
		},
		{
			name: "handles-invoice.payment_failed",
			sp:   &stripePayment{},
			args: args{
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBufferString(`
{
    "type": "invoice.payment_failed",
    "customer": "cus_srietnsirent",
    "subscription": "sub_auesrntausrent"
}
`)),
				},
			},
			wantPaymentEvent: Event{
				Type:           EventSubscriptionUnpaid,
				OrganizationID: 0,
				SessionID:      "",
				CustomerID:     "cus_srietnsirent",
				SubscriptionID: "sub_auesrntausrent",
			},
			wantErr: false,
		},
		{
			name: "handles-customer.subscription.updated:active",
			sp:   &stripePayment{},
			args: args{
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBufferString(`
{
    "type": "customer.subscription.updated",
    "customer": "cus_srietnsirent",
    "id": "sub_auesrntausrent",
    "status": "active"
}
`)),
				},
			},
			wantPaymentEvent: Event{
				Type:           EventSubscriptionPaid,
				OrganizationID: 0,
				SessionID:      "",
				CustomerID:     "cus_srietnsirent",
				SubscriptionID: "sub_auesrntausrent",
			},
			wantErr: false,
		},
		{
			name: "handles-customer.subscription.updated:trialing",
			sp:   &stripePayment{},
			args: args{
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBufferString(`
{
    "type": "customer.subscription.updated",
    "customer": "cus_srietnsirent",
    "id": "sub_auesrntausrent",
    "status": "trialing"
}
`)),
				},
			},
			wantPaymentEvent: Event{
				Type:           EventSubscriptionPaid,
				OrganizationID: 0,
				SessionID:      "",
				CustomerID:     "cus_srietnsirent",
				SubscriptionID: "sub_auesrntausrent",
			},
			wantErr: false,
		},
		{
			name: "handles-customer.subscription.updated:incomplete",
			sp:   &stripePayment{},
			args: args{
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBufferString(`
{
    "type": "customer.subscription.updated",
    "customer": "cus_srietnsirent",
    "id": "sub_auesrntausrent",
    "status": "incomplete"
}
        `)),
				},
			},
			wantPaymentEvent: Event{
				Type:           EventSubscriptionPaid,
				OrganizationID: 0,
				SessionID:      "",
				CustomerID:     "cus_srietnsirent",
				SubscriptionID: "sub_auesrntausrent",
			},
			wantErr: false,
		},
		{
			name: "handles-customer.subscription.updated:past_due",
			sp:   &stripePayment{},
			args: args{
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBufferString(`
{
    "type": "customer.subscription.updated",
    "customer": "cus_srietnsirent",
    "id": "sub_auesrntausrent",
    "status": "past_due"
}
                `)),
				},
			},
			wantPaymentEvent: Event{
				Type:           EventSubscriptionUnpaid,
				OrganizationID: 0,
				SessionID:      "",
				CustomerID:     "cus_srietnsirent",
				SubscriptionID: "sub_auesrntausrent",
			},
			wantErr: false,
		},
		{
			name: "handles-customer.subscription.updated:incomplete_expired",
			sp:   &stripePayment{},
			args: args{
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBufferString(`
{
    "type": "customer.subscription.updated",
    "customer": "cus_srietnsirent",
    "id": "sub_auesrntausrent",
    "status": "incomplete_expired"
}
                        `)),
				},
			},
			wantPaymentEvent: Event{
				Type:           EventSubscriptionUnpaid,
				OrganizationID: 0,
				SessionID:      "",
				CustomerID:     "cus_srietnsirent",
				SubscriptionID: "sub_auesrntausrent",
			},
			wantErr: false,
		},
		{
			name: "handles-customer.subscription.updated:canceled",
			sp:   &stripePayment{},
			args: args{
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBufferString(`
{
    "type": "customer.subscription.updated",
    "customer": "cus_srietnsirent",
    "id": "sub_auesrntausrent",
    "status": "canceled"
}
                                `)),
				},
			},
			wantPaymentEvent: Event{
				Type:           EventSubscriptionCanceled,
				OrganizationID: 0,
				SessionID:      "",
				CustomerID:     "cus_srietnsirent",
				SubscriptionID: "sub_auesrntausrent",
			},
			wantErr: false,
		},
		{
			name: "handles-customer.subscription.deleted",
			sp:   &stripePayment{},
			args: args{
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBufferString(`
{
    "type": "customer.subscription.deleted",
    "customer": "cus_srietnsirent",
    "id": "sub_auesrntausrent"
}
                                `)),
				},
			},
			wantPaymentEvent: Event{
				Type:           EventSubscriptionCanceled,
				OrganizationID: 0,
				SessionID:      "",
				CustomerID:     "cus_srietnsirent",
				SubscriptionID: "sub_auesrntausrent",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := &stripePayment{}
			gotPaymentEvent, err := sp.HandleEvent(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("stripePayment.HandleEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPaymentEvent, tt.wantPaymentEvent) {
				t.Errorf("stripePayment.HandleEvent() = %v, want %v", gotPaymentEvent, tt.wantPaymentEvent)
			}
		})
	}
}

func Test_stripePayment_GetSubscription(t *testing.T) {
	type args struct {
		subscriptionID SubscriptionID
	}
	ti, _ := time.Parse("2006-01-02 01:00:00", "1970-01-01 01:00:00")
	ti = ti.In(time.FixedZone("CET", 3600))
	tests := []struct {
		name             string
		sp               *stripePayment
		args             args
		wantSubscription Subscription
		wantErr          bool
	}{
		{
			name: "",
			sp:   &stripePayment{},
			args: args{
				subscriptionID: "subscription-id",
			},
			wantSubscription: Subscription{
				ID:                 "subscription-id",
				CustomerID:         "customer-id",
				Seats:              1,
				Status:             "paid",
				CurrentPeriodStart: ti,
				CurrentPeriodEnd:   ti,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := &stripePayment{}
			got, err := sp.GetSubscription(tt.args.subscriptionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("stripePayment.GetSubscription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			want := tt.wantSubscription

			if got.ID != want.ID ||
				got.CustomerID != want.CustomerID ||
				got.Seats != want.Seats ||
				got.Status != want.Status {
				t.Errorf("stripePayment.GetSubscription() = %v, want %v", got, want)
			}
		})
	}
}

func Test_stripePayment_UpdateSubscription(t *testing.T) {
	type args struct {
		subscriptionID SubscriptionID
		seats          int64
	}
	tests := []struct {
		name    string
		sp      *stripePayment
		args    args
		wantErr bool
	}{
		{
			name: "",
			sp:   &stripePayment{},
			args: args{
				subscriptionID: "subscription-id",
				seats:          2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := &stripePayment{}
			if err := sp.UpdateSubscription(tt.args.subscriptionID, tt.args.seats); (err != nil) != tt.wantErr {
				t.Errorf("stripePayment.UpdateSubscription() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_stripePayment_CancelSubscription(t *testing.T) {
	type args struct {
		subscriptionID SubscriptionID
	}
	tests := []struct {
		name    string
		sp      *stripePayment
		args    args
		wantErr bool
	}{
		{
			name: "",
			sp:   &stripePayment{},
			args: args{
				subscriptionID: "subscription-id",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := &stripePayment{}
			if err := sp.CancelSubscription(tt.args.subscriptionID); (err != nil) != tt.wantErr {
				t.Errorf("stripePayment.CancelSubscription() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_stripeSubscriptionStatus(t *testing.T) {
	type args struct {
		in stripe.SubscriptionStatus
	}
	tests := []struct {
		name    string
		args    args
		wantOut SubscriptionStatus
	}{
		{
			name: "active",
			args: args{
				in: "active",
			},
			wantOut: SubscriptionStatusPaid,
		},
		{
			name: "trialing",
			args: args{
				in: "trialing",
			},
			wantOut: SubscriptionStatusPaid,
		},
		{
			name: "incomplete",
			args: args{
				in: "incomplete",
			},
			wantOut: SubscriptionStatusPaid,
		},
		{
			name: "past_due",
			args: args{
				in: "past_due",
			},
			wantOut: SubscriptionStatusUnpaid,
		},
		{
			name: "unpaid",
			args: args{
				in: "unpaid",
			},
			wantOut: SubscriptionStatusUnpaid,
		},
		{
			name: "canceled",
			args: args{
				in: "canceled",
			},
			wantOut: SubscriptionStatusCanceled,
		},
		{
			name: "anything else",
			args: args{
				in: "anything else",
			},
			wantOut: SubscriptionStatusUnpaid,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := stripeSubscriptionStatus(tt.args.in); gotOut != tt.wantOut {
				t.Errorf("stripeSubscriptionStatus() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
