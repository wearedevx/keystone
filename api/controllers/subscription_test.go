package controllers

import (
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func TestPostSubscription(t *testing.T) {
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
		wantErr      bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, gotStatus, err := PostSubscription(tt.args.params, tt.args.in1, tt.args.Repo, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostSubscription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("PostSubscription() gotResponse = %v, want %v", gotResponse, tt.wantResponse)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("PostSubscription() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}

func TestGetPollSubscriptionSuccess(t *testing.T) {
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
		wantErr      bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, gotStatus, err := GetPollSubscriptionSuccess(tt.args.params, tt.args.in1, tt.args.Repo, tt.args.user)
			if (err != nil) != tt.wantErr {
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
	type args struct {
		w   http.ResponseWriter
		r   *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetCheckoutSuccess(tt.args.w, tt.args.r, tt.args.in2)
		})
	}
}

func TestGetCheckoutCancel(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		r   *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetCheckoutCancel(tt.args.w, tt.args.r, tt.args.in2)
		})
	}
}

func TestPostStripeWebhook(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		r   *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PostStripeWebhook(tt.args.w, tt.args.r, tt.args.in2)
		})
	}
}

func TestManageSubscription(t *testing.T) {
	type args struct {
		params router.Params
		in1    io.ReadCloser
		Repo   repo.IRepo
		user   models.User
	}
	tests := []struct {
		name       string
		args       args
		want       router.Serde
		wantStatus int
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotStatus, err := ManageSubscription(tt.args.params, tt.args.in1, tt.args.Repo, tt.args.user)
			if (err != nil) != tt.wantErr {
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
