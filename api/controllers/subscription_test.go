package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/bxcodec/faker/v3"
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
				Repo: newFakeRepo(noCrashers),
				user: user,
			},
			wantResponse: &models.StartSubscriptionResponse{
				SessionID: "session-id",
				URL:       "http://session.url",
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "must have an organziation name",
			args: args{
				params: router.Params{},
				in1:    nil,
				Repo:   newFakeRepo(noCrashers),
				user:   user,
			},
			wantResponse: nil,
			wantStatus:   http.StatusBadRequest,
			wantErr:      "bad request",
		},
		{
			name: "the organization must exists",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"organizationName": "not an organizationName",
				}),
				in1:  nil,
				Repo: newFakeRepo(noCrashers),
				user: user,
			},
			wantResponse: nil,
			wantStatus:   http.StatusNotFound,
			wantErr:      "not found",
		},
		{
			name: "fails to get the organization or count its memebers",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"organizationName": organization.Name,
				}),
				in1: nil,
				Repo: newFakeRepo(map[string]error{
					"OrganizationCountMembers": errors.New("unexpected error"),
				}),
				user: user,
			},
			wantResponse: nil,
			wantStatus:   http.StatusInternalServerError,
			wantErr:      "failed to get: unexpected error",
		},
		{
			name: "the user must own the organization",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"organizationName": organization.Name,
				}),
				in1:  nil,
				Repo: newFakeRepo(noCrashers),
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
				Repo: newFakeRepo(noCrashers),
				user: paidUser,
			},
			wantResponse: nil,
			wantStatus:   http.StatusConflict,
			wantErr:      "already subscribed",
		},
		{
			name: "create checkout session fail",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"organizationName": organization.Name,
				}),
				in1: nil,
				Repo: newFakeRepo(map[string]error{
					"OrganizationSetCustomer": errors.New("unexpected error"),
				}),
				user: user,
			},
			wantResponse: nil,
			wantStatus:   http.StatusInternalServerError,
			wantErr:      "failed to create resource: unexpected error",
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
				Repo: newFakeRepo(noCrashers),
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
				Repo: newFakeRepo(noCrashers),
				user: user,
			},
			wantResponse: nil,
			wantStatus:   http.StatusNotFound,
			wantErr:      "not found",
		},
		{
			name: "other errors",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"sessionID": csession.SessionID,
				}),
				in1: nil,
				Repo: newFakeRepo(map[string]error{
					"GetCheckoutSession": errors.New("unexpected error"),
				}),
				user: user,
			},
			wantResponse: nil,
			wantStatus:   http.StatusInternalServerError,
			wantErr:      "failed to get: unexpected error",
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

	csd := seedCheckoutSession()
	defer teardownCheckoutSession(csd)

	okURL, _ := url.Parse(fmt.Sprintf("http://tests.com?session_id=%s", cs.SessionID))
	koURL, _ := url.Parse("http://tests.com?session_id=not-a-session-id")

	dURL, _ := url.Parse(fmt.Sprintf("http://tests.com?sessions_id=%s", csd.SessionID))

	type args struct {
		w    http.ResponseWriter
		r    *http.Request
		in2  router.Params
		Repo repo.IRepo
	}
	tests := []struct {
		name       string
		args       args
		want       string
		wantStatus int
		wantErr    bool
	}{
		{
			name: "it works",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					URL: okURL,
				},
				in2:  router.Params{},
				Repo: newFakeRepo(noCrashers),
			},
			want:       "Thank you for subscribing to Keystone!",
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "not found",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					URL: koURL,
				},
				in2:  router.Params{},
				Repo: newFakeRepo(noCrashers),
			},
			want:       "No such checkout session",
			wantStatus: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name: "fails to delete the resource",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					URL: dURL,
				},
				in2: router.Params{},
				Repo: newFakeRepo(map[string]error{
					"DeleteCheckoutSession": errors.New("unexpected error"),
				}),
			},
			want:       "An error occurred: unexpected error",
			wantStatus: http.StatusInternalServerError,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := GetCheckoutSuccess(tt.args.w, tt.args.r, tt.args.in2, tt.args.Repo)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetCheckoutSuccess() err = %v, want %v", err, tt.wantErr)
				return
			}

			got := tt.args.w.(*mockResponseWriter)
			if got.body.String() != tt.want {
				t.Errorf("GetCheckoutSuccess() got.body %v, want %v", got.body.String(), tt.want)
				return
			}

			if status != tt.wantStatus {
				t.Errorf("GetCheckoutSuccess() status %v, want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestGetCheckoutCancel(t *testing.T) {
	cs := seedCheckoutSession()
	defer teardownCheckoutSession(cs)

	csd := seedCheckoutSession()
	defer teardownCheckoutSession(csd)

	okURL, _ := url.Parse(fmt.Sprintf("http://tests.com?session_id=%s", cs.SessionID))
	koURL, _ := url.Parse("http://tests.com?session_id=not-a-session-id")

	dURL, _ := url.Parse(fmt.Sprintf("http://tests.com?session_id=%s", csd.SessionID))

	type args struct {
		w    http.ResponseWriter
		r    *http.Request
		in2  router.Params
		Repo repo.IRepo
	}
	tests := []struct {
		name       string
		args       args
		want       string
		wantStatus int
		wantErr    bool
	}{
		{
			name: "it works",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					URL: okURL,
				},
				in2:  router.Params{},
				Repo: newFakeRepo(noCrashers),
			},
			want:       "You cancelled your subscription",
			wantStatus: http.StatusOK,
		},
		{
			name: "not found",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					URL: koURL,
				},
				in2:  router.Params{},
				Repo: newFakeRepo(noCrashers),
			},
			want:       "No such checkout session",
			wantStatus: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name: "if fails at deleting the resource",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					URL: dURL,
				},
				in2: router.Params{},
				Repo: newFakeRepo(map[string]error{
					"DeleteCheckoutSession": errors.New("unexpected error"),
				}),
			},
			want:       "An error occurred: unexpected error",
			wantStatus: http.StatusInternalServerError,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := GetCheckoutCancel(tt.args.w, tt.args.r, tt.args.in2, tt.args.Repo)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetCheckoutCancel() err = %v, want %v", err, tt.wantErr)
				return
			}

			got := tt.args.w.(*mockResponseWriter)
			if got.body.String() != tt.want {
				t.Errorf("GetCheckoutCancel() got.body %v, want %v", got.body.String(), tt.want)
				return
			}

			if status != tt.wantStatus {
				t.Errorf("GetCheckoutCancel() status %v, want %v", status, tt.wantStatus)
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
		w    http.ResponseWriter
		r    *http.Request
		in2  router.Params
		Repo repo.IRepo
	}
	tests := []struct {
		name             string
		args             args
		wantStatus       int
		wantOrganization *models.Organization
		wantErr          bool
	}{
		{
			name: "ignores event",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBuffer([]byte{})),
				},
				in2:  router.Params{},
				Repo: newFakeRepo(noCrashers),
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
				Repo: newFakeRepo(noCrashers),
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name: "ignores-a-specific-parsing-error",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					Body: ioutil.NopCloser(bytes.NewBufferString("bad-event-to-ignore")),
				},
				Repo: newFakeRepo(noCrashers),
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    true,
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
				Repo: newFakeRepo(noCrashers),
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
			name: "fails checkout.session.completed",
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
				Repo: newFakeRepo(map[string]error{
					"GetCheckoutSession": errors.New("unexpected error"),
				}),
			},
			wantStatus:       http.StatusInternalServerError,
			wantOrganization: nil,
			wantErr:          true,
		},
		{
			name: "fails checkout.session.completed II",
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
				Repo: newFakeRepo(map[string]error{
					"UpdateCheckoutSession": errors.New("unexpected error"),
				}),
			},
			wantStatus:       http.StatusInternalServerError,
			wantOrganization: nil,
			wantErr:          true,
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
					))),
				},
				Repo: newFakeRepo(noCrashers),
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
			name: "fails invoice.paid",
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
					))),
				},
				Repo: newFakeRepo(map[string]error{
					"GetOrganization": errors.New("unexpected error"),
				}),
			},
			wantStatus:       http.StatusInternalServerError,
			wantOrganization: nil,
			wantErr:          true,
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
				Repo: newFakeRepo(noCrashers),
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
			name: "fails invoice.payment_failed",
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
				Repo: newFakeRepo(map[string]error{
					"GetOrganization": errors.New("unexpected error"),
				}),
			},
			wantOrganization: nil,
			wantStatus:       http.StatusInternalServerError,
			wantErr:          true,
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
				Repo: newFakeRepo(noCrashers),
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
			name: "fails-customer.subscription.updated:active",
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
				Repo: newFakeRepo(map[string]error{
					"GetOrganization": errors.New("unexpected error"),
				}),
			},
			wantOrganization: nil,
			wantStatus:       http.StatusInternalServerError,
			wantErr:          true,
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
				Repo: newFakeRepo(noCrashers),
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
			name: "fails-customer.subscription.updated:trialing",
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
				Repo: newFakeRepo(map[string]error{
					"GetOrganization": errors.New("unexpected error"),
				}),
			},
			wantOrganization: nil,
			wantStatus:       http.StatusInternalServerError,
			wantErr:          true,
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
				Repo: newFakeRepo(noCrashers),
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
			name: "fails-customer.subscription.updated:incomplete",
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
				Repo: newFakeRepo(map[string]error{
					"GetOrganization": errors.New("unexpected error"),
				}),
			},
			wantOrganization: nil,
			wantStatus:       http.StatusInternalServerError,
			wantErr:          true,
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
				Repo: newFakeRepo(noCrashers),
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
			name: "fails-customer.subscription.updated:past_due",
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
				Repo: newFakeRepo(map[string]error{
					"GetOrganization": errors.New("unexpected error"),
				}),
			},
			wantOrganization: nil,
			wantStatus:       http.StatusInternalServerError,
			wantErr:          true,
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
				Repo: newFakeRepo(noCrashers),
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
			name: "fails-customer.subscription.updated:incomplete_expired",
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
				Repo: newFakeRepo(map[string]error{
					"GetOrganization": errors.New("unexpected error"),
				}),
			},
			wantOrganization: nil,
			wantStatus:       http.StatusInternalServerError,
			wantErr:          true,
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
				Repo: newFakeRepo(noCrashers),
			},
			wantOrganization: &models.Organization{
				ID:             canceledOrg.ID,
				CustomerID:     canceledOrg.CustomerID,
				SubscriptionID: canceledOrg.SubscriptionID,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "fails-customer.subscription.updated:canceled",
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
				Repo: newFakeRepo(map[string]error{
					"GetOrganization": errors.New("unexpected error"),
				}),
			},
			wantOrganization: nil,
			wantStatus:       http.StatusInternalServerError,
			wantErr:          true,
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
				Repo: newFakeRepo(noCrashers),
			},
			wantOrganization: &models.Organization{
				ID:             canceledOrg.ID,
				CustomerID:     canceledOrg.CustomerID,
				SubscriptionID: canceledOrg.SubscriptionID,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "fails-customer.subscription.deleted",
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
				Repo: newFakeRepo(map[string]error{
					"GetOrganization": errors.New("unexpected error"),
				}),
			},
			wantOrganization: nil,
			wantStatus:       http.StatusInternalServerError,
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := PostStripeWebhook(tt.args.w, tt.args.r, tt.args.in2, tt.args.Repo)

			if (err != nil) != tt.wantErr {
				t.Errorf("PostStripeWebhook() error: %v, want %v", err, tt.wantErr)
				return
			}

			if status != tt.wantStatus {
				t.Errorf("PostStripeWebhook() status %v, want %v", status, tt.wantStatus)
				return
			}

			if tt.wantOrganization != nil {
				gotOrganization := models.Organization{}
				repo.NewRepo().GetDB().First(&gotOrganization, tt.wantOrganization.ID)

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
				Repo: newFakeRepo(noCrashers),
				user: user,
			},
			want: &models.ManageSubscriptionResponse{
				URL: "http://portal-session.url",
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
				Repo: newFakeRepo(noCrashers),
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
				Repo: newFakeRepo(noCrashers),
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
				Repo: newFakeRepo(noCrashers),
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

	err := repo.NewRepo().GetDB().Transaction(func(db *gorm.DB) error {
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
	err := repo.NewRepo().GetDB().Transaction(func(db *gorm.DB) error {
		db.Delete(&cs)
		return db.Error
	})

	if err != nil {
		panic(err)
	}
}
