package controllers

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/julienschmidt/httprouter"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
	"gorm.io/gorm"
)

func TestPostSubscription(t *testing.T) {
	user, organization, project := seedOneProjectForOneUser()
	defer teardownUserAndOrganization(user, organization)
	defer teardownProject(project)

	paidUser, paidOrganization := seedSingleUser()
	defer teardownUserAndOrganization(paidUser, paidOrganization)

	testsSetOrganisationPaid(&paidOrganization)

	type args struct {
		params router.Params
		in1    io.ReadCloser
		Repo   repo.IRepo
		user   models.User
	}
	tests := []struct {
		name         string
		args         args
		wantResponse *models.StartSubscriptionResponse
		wantStatus   int
		wantErr      string
	}{
		{
			name: "creates a subscription",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"organizationName": organization.Name,
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				user: user,
			},
			wantResponse: &models.StartSubscriptionResponse{
				SessionID: "session-id",
				Url:       "http://session.url",
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "the organization must exists",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"organizationName": "not an organizationName",
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				user: user,
			},
			wantResponse: nil,
			wantStatus:   http.StatusNotFound,
			wantErr:      "not found",
		},
		{
			name: "the user must own the organization",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"organizationName": organization.Name,
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				user: paidUser,
			},
			wantResponse: nil,
			wantStatus:   http.StatusForbidden,
			wantErr:      "permission denied",
		},
		{
			name: "ther organiation must not be paid",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"organizationName": paidOrganization.Name,
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				user: paidUser,
			},
			wantResponse: nil,
			wantStatus:   http.StatusConflict,
			wantErr:      "already subscribed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, gotStatus, err := PostSubscription(tt.args.params, tt.args.in1, tt.args.Repo, tt.args.user)
			if err.Error() != tt.wantErr {
				t.Errorf("PostSubscription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantResponse != nil && gotResponse != nil {
				if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
					t.Errorf("PostSubscription() gotResponse = %v, want %v", gotResponse, tt.wantResponse)
				}
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("PostSubscription() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}

func TestGetPollSubscriptionSuccess(t *testing.T) {
	user, organization := seedSingleUser()
	defer teardownUserAndOrganization(user, organization)

	csession := seedCheckoutSession()
	defer teardownCheckoutSession(csession)

	type args struct {
		params router.Params
		in1    io.ReadCloser
		Repo   repo.IRepo
		user   models.User
	}
	tests := []struct {
		name         string
		args         args
		wantResponse router.Serde
		wantStatus   int
		wantErr      string
	}{
		{
			name: "it works",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"sessionID": csession.SessionID,
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				user: user,
			},
			wantResponse: nil,
			wantStatus:   http.StatusOK,
			wantErr:      "",
		},
		{
			name: "not found",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"sessionID": "this is not a session id",
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				user: user,
			},
			wantResponse: nil,
			wantStatus:   http.StatusNotFound,
			wantErr:      "not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, gotStatus, err := GetPollSubscriptionSuccess(tt.args.params, tt.args.in1, tt.args.Repo, tt.args.user)
			if err.Error() != tt.wantErr {
				t.Errorf("GetPollSubscriptionSuccess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("GetPollSubscriptionSuccess() gotResponse = %v, want %v", gotResponse, tt.wantResponse)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("GetPollSubscriptionSuccess() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}

func TestGetCheckoutSuccess(t *testing.T) {
	cs := seedCheckoutSession()
	defer teardownCheckoutSession(cs)

	okUrl, _ := url.Parse(fmt.Sprintf("http://tests.com?session_id=%s", cs.SessionID))
	koUrl, _ := url.Parse("http://tests.com?session_id=not-a-session-id")

	type args struct {
		w   http.ResponseWriter
		r   *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name       string
		args       args
		want       string
		wantStatus int
	}{
		{
			name: "it works",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					URL: okUrl,
				},
				in2: []httprouter.Param{},
			},
			want:       "Thank you for subscribing to Keystone!",
			wantStatus: http.StatusOK,
		},
		{
			name: "not found",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					URL: koUrl,
				},
				in2: []httprouter.Param{},
			},
			want:       "No such checkout session",
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetCheckoutSuccess(tt.args.w, tt.args.r, tt.args.in2)

			got := tt.args.w.(*mockResponseWriter)
			if got.body.String() != tt.want {
				t.Errorf("GetCheckoutSuccess() got.body %v, want %v", got.body.String(), tt.want)
				return
			}

			if got.status != tt.wantStatus {
				t.Errorf("GetCheckoutSuccess() got.status %v, want %v", got.status, tt.wantStatus)
			}
		})
	}
}

func TestGetCheckoutCancel(t *testing.T) {
	cs := seedCheckoutSession()
	defer teardownCheckoutSession(cs)

	okUrl, _ := url.Parse(fmt.Sprintf("http://tests.com?session_id=%s", cs.SessionID))
	koUrl, _ := url.Parse("http://tests.com?session_id=not-a-session-id")

	type args struct {
		w   http.ResponseWriter
		r   *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name       string
		args       args
		want       string
		wantStatus int
	}{
		{
			name: "it works",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					URL: okUrl,
				},
				in2: []httprouter.Param{},
			},
			want:       "You cancelled your subscription",
			wantStatus: http.StatusOK,
		},
		{
			name: "not found",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					URL: koUrl,
				},
				in2: []httprouter.Param{},
			},
			want:       "No such checkout session",
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetCheckoutCancel(tt.args.w, tt.args.r, tt.args.in2)

			got := tt.args.w.(*mockResponseWriter)
			if got.body.String() != tt.want {
				t.Errorf("GetCheckoutCancel() got.body %v, want %v", got.body.String(), tt.want)
				return
			}

			if got.status != tt.wantStatus {
				t.Errorf("GetCheckoutCancel() got.status %v, want %v", got.status, tt.wantStatus)
			}
		})
	}
}

func TestPostStripeWebhook(t *testing.T) {
	defaultUser, defaultOrg := seedSingleUser()
	defer teardownUserAndOrganization(defaultUser, defaultOrg)

	defaultSession := seedCheckoutSession()
	defer teardownCheckoutSession(defaultSession)

	// completed
	completedUser, completedOrg := seedSingleUser()
	defer teardownUserAndOrganization(completedUser, completedOrg)

	completedSession := seedCheckoutSession()
	defer teardownCheckoutSession(completedSession)

	// paid
	paidUser, paidOrg := seedSingleUser()
	defer teardownUserAndOrganization(paidUser, paidOrg)

	paidSession := seedCheckoutSession()
	defer teardownCheckoutSession(paidSession)

	// unpaid
	unpaidUser, unpaidOrg := seedSingleUser()
	defer teardownUserAndOrganization(unpaidUser, unpaidOrg)

	unpaidSession := seedCheckoutSession()
	defer teardownCheckoutSession(unpaidSession)

	// canceled
	canceledUser, canceledOrg := seedSingleUser()
	defer teardownUserAndOrganization(canceledUser, canceledOrg)

	canceledSession := seedCheckoutSession()
	defer teardownCheckoutSession(canceledSession)

	type args struct {
		w   http.ResponseWriter
		r   *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name             string
		args             args
		wantStatus       int
		wantOrganization *models.Organization
	}{
		{
			name: "ignores event",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBuffer([]byte{})),
				},
				in2: []httprouter.Param{},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "error-if-bad-event",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBufferString("bad-event")),
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "ignores-a-specific-parsing-error",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBufferString("bad-event-to-ignore")),
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "handles checkout.session.completed",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBufferString(fmt.Sprintf(
						`{
    "type": "checkout.session.completed",
    "client_reference_id": "%d",
    "id": "%s",
    "customer": "cus_srietnsirent",
    "subscription": "sub_auesrntausrent"
}`,
						completedOrg.ID,
						completedSession.SessionID,
					))),
				},
			},
			wantStatus: http.StatusOK,
			wantOrganization: &models.Organization{
				ID:             completedOrg.ID,
				CustomerID:     "cus_srietnsirent",
				SubscriptionID: "sub_auesrntausrent",
				Paid:           true,
			},
		},
		{
			name: "handles invoice.paid",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBufferString(fmt.Sprintf(
						`{
						"type": "invoice.paid",
						"customer": "%s",
						"subscription": "%s"
				}`,
						paidOrg.CustomerID,
						paidOrg.SubscriptionID,
					)))},
			},
			wantStatus: http.StatusOK,
			wantOrganization: &models.Organization{
				ID:             paidOrg.ID,
				CustomerID:     paidOrg.CustomerID,
				SubscriptionID: paidOrg.SubscriptionID,
				Paid:           true,
			},
		},
		{
			name: "handles invoice.payment_failed",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBufferString(fmt.Sprintf(
						`{
				"type": "invoice.payment_failed",
				"customer": "%s",
				"subscription": "%s"
		}`,
						unpaidOrg.CustomerID,
						unpaidOrg.SubscriptionID,
					))),
				},
			},
			wantOrganization: &models.Organization{
				ID:             unpaidOrg.ID,
				CustomerID:     unpaidOrg.CustomerID,
				SubscriptionID: unpaidOrg.SubscriptionID,
				Paid:           false,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "handles-customer.subscription.updated:active",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBufferString(fmt.Sprintf(
						`{
    "type": "customer.subscription.updated",
    "customer": "%s",
    "id": "%s",
    "status": "active"
}`,
						paidOrg.CustomerID,
						paidOrg.SubscriptionID,
					))),
				},
			},
			wantOrganization: &models.Organization{
				ID:             paidOrg.ID,
				CustomerID:     paidOrg.CustomerID,
				SubscriptionID: paidOrg.SubscriptionID,
				Paid:           true,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "handles-customer.subscription.updated:trialing",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBufferString(fmt.Sprintf(
						`{
	"type": "customer.subscription.updated",
	"customer": "%s",
	"id": "%s",
	"status": "trialing"
}`,
						paidOrg.CustomerID,
						paidOrg.SubscriptionID,
					))),
				},
			},
			wantOrganization: &models.Organization{
				ID:             paidOrg.ID,
				CustomerID:     paidOrg.CustomerID,
				SubscriptionID: paidOrg.SubscriptionID,
				Paid:           true,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "handles-customer.subscription.updated:incomplete",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBufferString(fmt.Sprintf(
						`{
	"type": "customer.subscription.updated",
	"customer": "%s",
	"id": "%s",
	"status": "incomplete"
}`,
						paidOrg.CustomerID,
						paidOrg.SubscriptionID,
					))),
				},
			},
			wantOrganization: &models.Organization{
				ID:             paidOrg.ID,
				CustomerID:     paidOrg.CustomerID,
				SubscriptionID: paidOrg.SubscriptionID,
				Paid:           true,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "handles-customer.subscription.updated:past_due",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBufferString(fmt.Sprintf(
						`{
				"type": "customer.subscription.updated",
				"customer": "%s",
				"id": "%s",
				"status": "past_due"
		}`,
						unpaidOrg.CustomerID,
						unpaidOrg.SubscriptionID,
					))),
				},
			},
			wantOrganization: &models.Organization{
				ID:             unpaidOrg.ID,
				CustomerID:     unpaidOrg.CustomerID,
				SubscriptionID: unpaidOrg.SubscriptionID,
				Paid:           false,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "handles-customer.subscription.updated:incomplete_expired",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBufferString(fmt.Sprintf(
						`{
				"type": "customer.subscription.updated",
				"customer": "%s",
				"id": "%s",
				"status": "incomplete_expired"
		}`,
						unpaidOrg.CustomerID,
						unpaidOrg.SubscriptionID,
					))),
				},
			},
			wantOrganization: &models.Organization{
				ID:             unpaidOrg.ID,
				CustomerID:     unpaidOrg.CustomerID,
				SubscriptionID: unpaidOrg.SubscriptionID,
				Paid:           false,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "handles-customer.subscription.updated:canceled",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBufferString(fmt.Sprintf(
						`{
				"type": "customer.subscription.updated",
				"customer": "%s",
				"id": "%s",
				"status": "canceled"
		}`,
						canceledOrg.CustomerID,
						canceledOrg.SubscriptionID,
					))),
				},
			},
			wantOrganization: &models.Organization{
				ID:             canceledOrg.ID,
				CustomerID:     canceledOrg.CustomerID,
				SubscriptionID: canceledOrg.SubscriptionID,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "handles-customer.subscription.deleted",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBufferString(fmt.Sprintf(
						`{
				"type": "customer.subscription.deleted",
				"customer": "%s",
				"id": "%s"
		}`,
						canceledOrg.CustomerID,
						canceledOrg.SubscriptionID,
					))),
				},
			},
			wantOrganization: &models.Organization{
				ID:             canceledOrg.ID,
				CustomerID:     canceledOrg.CustomerID,
				SubscriptionID: canceledOrg.SubscriptionID,
			},
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PostStripeWebhook(tt.args.w, tt.args.r, tt.args.in2)

			got := tt.args.w.(*mockResponseWriter)
			if got.status != tt.wantStatus {
				t.Errorf("PostStripeWebhook() got status %v, want %v", got.status, tt.wantStatus)
				return
			}

			if tt.wantOrganization != nil {
				gotOrganization := models.Organization{}
				repo.NewRepo().GetDb().First(&gotOrganization, tt.wantOrganization.ID)

				if gotOrganization.CustomerID != tt.wantOrganization.CustomerID ||
					gotOrganization.SubscriptionID != tt.wantOrganization.SubscriptionID ||
					gotOrganization.Paid != tt.wantOrganization.Paid {
					t.Errorf("PostStripeWebhook() got organization %v, want %v", gotOrganization, tt.wantOrganization)
					return
				}
			}
		})
	}
}

func TestManageSubscription(t *testing.T) {
	user, org := seedSingleUser()
	otherUser, otherOrg := seedSingleUser()
	defer teardownUserAndOrganization(user, org)
	defer teardownUserAndOrganization(otherUser, otherOrg)

	testsSetOrganisationPaid(&org)

	type args struct {
		params router.Params
		in1    io.ReadCloser
		Repo   repo.IRepo
		user   models.User
	}
	tests := []struct {
		name       string
		args       args
		want       *models.ManageSubscriptionResponse
		wantStatus int
		wantErr    string
	}{
		{
			name: "returns a management link",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"organizationName": org.Name,
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				user: user,
			},
			want: &models.ManageSubscriptionResponse{
				Url: "http://portal-session.url",
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "organization must exist",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"organizationName": "not an organization name",
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				user: user,
			},
			want:       nil,
			wantStatus: http.StatusNotFound,
			wantErr:    "not found",
		},
		{
			name: "the user must own the organization",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"organizationName": org.Name,
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				user: otherUser,
			},
			want:       nil,
			wantStatus: http.StatusForbidden,
			wantErr:    "permission denied",
		},
		{
			name: "the organization must be paid",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"organizationName": otherOrg.Name,
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				user: otherUser,
			},
			want:       nil,
			wantStatus: http.StatusForbidden,
			wantErr:    "needs upgrade",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotStatus, err := ManageSubscription(tt.args.params, tt.args.in1, tt.args.Repo, tt.args.user)
			if err.Error() != tt.wantErr {
				t.Errorf("ManageSubscription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ManageSubscription() got = %v, want %v", got, tt.want)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("ManageSubscription() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}

func seedCheckoutSession() models.CheckoutSession {
	cs := models.CheckoutSession{}

	err := repo.NewRepo().GetDb().Transaction(func(db *gorm.DB) error {
		faker.FakeData(&cs)

		db.Create(&cs)

		return db.Error
	})

	if err != nil {
		panic(err)
	}

	return cs
}

func teardownCheckoutSession(cs models.CheckoutSession) {
	err := repo.NewRepo().GetDb().Transaction(func(db *gorm.DB) error {
		db.Delete(&cs)
		return db.Error
	})

	if err != nil {
		panic(err)
	}
}
