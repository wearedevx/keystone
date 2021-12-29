package controllers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/julienschmidt/httprouter"
	"github.com/wearedevx/keystone/api/internal/router"
	"gorm.io/gorm"

	// . "github.com/wearedevx/keystone/api/pkg/apierrors"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func TestGetUser(t *testing.T) {
	user, org := seedSingleUser()
	defer teardownUserAndOrganization(user, org)

	type args struct {
		in0  router.Params
		in1  io.ReadCloser
		in2  repo.IRepo
		user models.User
	}
	tests := []struct {
		name    string
		args    args
		want    router.Serde
		want1   int
		wantErr bool
	}{
		{
			name: "it works",
			args: args{
				in0:  router.Params{},
				in1:  nil,
				in2:  nil,
				user: user,
			},
			want:    &user,
			want1:   http.StatusOK,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetUser(
				tt.args.in0,
				tt.args.in1,
				tt.args.in2,
				tt.args.user,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetUser() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPostUserToken(t *testing.T) {
	device := seedOnlyDevice()
	defer teardownOnlyDevice(device)

	type args struct {
		w   http.ResponseWriter
		r   *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name              string
		args              args
		wantStatus        int
		wantAuthorization bool
	}{
		{
			name: "it works",
			args: args{
				w: newMockResponse(),
				r: &http.Request{Body: io.NopCloser(bytes.NewBufferString(`
							{
								"AccountType": "github",
								"Token": {
									"access_token": "YSB0b2tlbg=="
								},
								"PublicKey": "YSB2ZXJ5IHB1YmxpYyBrZXk=",
								"Device": "a-device",
								"DeviceUID": "a-device-uid"
							}`))},
				in2: []httprouter.Param{},
			},
			wantAuthorization: true,
		},
		{
			name: "bad device name",
			args: args{
				w: newMockResponse(),
				r: &http.Request{Body: io.NopCloser(bytes.NewBufferString(`
							{
								"AccountType": "github",
								"Token": {
									"access_token": "YSB0b2tlbg=="
								},
								"PublicKey": "YSB2ZXJ5IHB1YmxpYyBrZXk=",
								"Device": "is that such a bad device name ?",
								"DeviceUID": "a-device-uid"
							}`))},
				in2: []httprouter.Param{},
			},
			wantStatus:        http.StatusConflict,
			wantAuthorization: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PostUserToken(tt.args.w, tt.args.r, tt.args.in2)
			got := tt.args.w.(*mockResponseWriter)

			if got.status != tt.wantStatus {
				t.Errorf(
					"PostUserToken() got status %v, want %v",
					got.status,
					tt.wantStatus,
				)
				return
			}

			if tt.wantAuthorization {
				if _, ok := got.headers["Authorization"]; !ok {
					t.Error("PostUserToken() missing authorization header")
					return
				}
			}
		})
	}
}

func TestGetAuthRedirect(t *testing.T) {
	validRequest := seedLoginRequest()
	defer teardownLoginRequest(validRequest)

	goodState, _ := models.AuthState{
		TemporaryCode: validRequest.TemporaryCode,
		Version:       "test-version",
	}.Encode()

	tooShortState, _ := models.AuthState{
		TemporaryCode: "1234",
		Version:       "test-version",
	}.Encode()

	notFoundState, _ := models.AuthState{
		TemporaryCode: "aisretnriasentiasretniasrent",
		Version:       "test-version",
	}.Encode()

	okUrl, _ := url.Parse(
		fmt.Sprintf(
			"http://tests.com/redirect?state=%s&code=123456",
			goodState,
		),
	)
	tooShortCodeUrl, _ := url.Parse(
		fmt.Sprintf(
			"http://tests.com/redirect?state=%s&code=123456",
			tooShortState,
		),
	)
	notFoundUrl, _ := url.Parse(
		fmt.Sprintf(
			"http://tests.com/redirect?state=%s&code=123456",
			notFoundState,
		),
	)

	type args struct {
		w   http.ResponseWriter
		r   *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name         string
		args         args
		wantStatus   int
		wantResponse string
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
			wantResponse: `You have been successfully authenticated.
You may now return to your terminal and start using Keystone.

Thank you!`,
		},
		{
			name: "temporary code is too short",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					URL: tooShortCodeUrl,
				},
				in2: []httprouter.Param{},
			},
			wantStatus:   http.StatusBadRequest,
			wantResponse: "Bad Request\n",
		},
		{
			name: "no such request",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					URL: notFoundUrl,
				},
				in2: []httprouter.Param{},
			},
			wantStatus:   http.StatusNotFound,
			wantResponse: "not found\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetAuthRedirect(tt.args.w, tt.args.r, tt.args.in2)

			got := tt.args.w.(*mockResponseWriter)

			if got.status != tt.wantStatus {
				t.Errorf(
					"GetAuthRedirect() got status %v, want %v",
					got.status,
					tt.wantStatus,
				)
				return
			}

			gotResponse := got.body.String()
			if gotResponse != tt.wantResponse {
				t.Errorf(
					"GetAuthRedirec() got response %v, want %v",
					gotResponse,
					tt.wantResponse,
				)
				return
			}
		})
	}
}

func TestPostLoginRequest(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		in1 *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			name: "it work",
			args: args{
				w:   newMockResponse(),
				in1: &http.Request{},
				in2: []httprouter.Param{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PostLoginRequest(tt.args.w, tt.args.in1, tt.args.in2)

			got := tt.args.w.(*mockResponseWriter)

			if got.status != tt.wantStatus {
				t.Errorf(
					"PostLoginRequest() got status %v, want %v",
					got.status,
					tt.wantStatus,
				)
				return
			}

			var lr models.LoginRequest
			if err := lr.Deserialize(got.body); err != nil {
				t.Errorf(
					"PostLoginRequest() invalid respons body: %v",
					got.body.String(),
				)
			}
		})
	}
}

func TestGetLoginRequest(t *testing.T) {
	lr := seedLoginRequest()
	defer teardownLoginRequest(lr)

	okUrl, _ := url.Parse(
		fmt.Sprintf("http://tests.com/login-request?code=%s", lr.TemporaryCode),
	)
	notFoundUrl, _ := url.Parse("http://tests.com/login-request?code=notacodebutitcouldhavebeenjusthavetobelongenough")

	type args struct {
		w   http.ResponseWriter
		r   *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name        string
		args        args
		wantStatus  int
		want        *models.LoginRequest
		wantMessage string
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
			wantStatus:  http.StatusOK,
			want:        &lr,
			wantMessage: "",
		},
		{
			name: "not found",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					URL: notFoundUrl,
				},
				in2: []httprouter.Param{},
			},
			wantStatus:  http.StatusNotFound,
			want:        nil,
			wantMessage: "not found\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetLoginRequest(tt.args.w, tt.args.r, tt.args.in2)

			got := tt.args.w.(*mockResponseWriter)

			if got.status != tt.wantStatus {
				t.Errorf(
					"GetLoginRequest() got status %v, want %v",
					got.status,
					tt.wantStatus,
				)
				return
			}

			switch {
			case tt.want != nil && got.body.Len() != 0:
				var lr models.LoginRequest
				lr.Deserialize(got.body)

				if lr.ID != tt.want.ID ||
					lr.TemporaryCode != tt.want.TemporaryCode ||
					lr.AuthCode != tt.want.AuthCode ||
					lr.Answered != tt.want.Answered {
					t.Errorf("GetLoginRequest() got %v, want %v", lr, tt.want)
				}
			case tt.want == nil &&
				tt.wantMessage != "":
				if got.body.String() != tt.wantMessage {
					t.Errorf("GetLoginRequest() got %v, want %v", got.body.String(), tt.wantMessage)
				}
			default:
				// t.Errorf("GetLoginRequest() got %v, want %v", got.body.String(), tt.want)
			}
		})
	}
}

func TestGetUserKeys(t *testing.T) {
	user, org := seedSingleUser()
	defer teardownUserAndOrganization(user, org)

	otherUser, otherOrg := seedSingleUser()
	defer teardownUserAndOrganization(otherUser, otherOrg)

	type args struct {
		params router.Params
		in1    io.ReadCloser
		Repo   repo.IRepo
		in3    models.User
	}
	tests := []struct {
		name       string
		args       args
		want       *models.UserDevices
		wantStatus int
		wantErr    string
	}{
		{
			name: "it works",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"userID": otherUser.UserID,
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				in3:  models.User{},
			},
			want: &models.UserDevices{
				UserID:  otherUser.ID,
				UserUID: "",
				Devices: otherUser.Devices,
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "bad request if no user id",
			args: args{
				params: router.Params{},
				in1:    nil,
				Repo:   repo.NewRepo(),
				in3:    user,
			},
			want:       &models.UserDevices{},
			wantStatus: http.StatusBadRequest,
			wantErr:    "bad request",
		},
		{
			name: "not found if the target user does not exist",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"userID": "not a member id",
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				in3:  user,
			},
			want:       &models.UserDevices{},
			wantStatus: http.StatusNotFound,
			wantErr:    "not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotStatus, err := GetUserKeys(
				tt.args.params,
				tt.args.in1,
				tt.args.Repo,
				tt.args.in3,
			)
			if err.Error() != tt.wantErr {
				t.Errorf(
					"GetUserKeys() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if gotStatus != tt.wantStatus {
				t.Errorf(
					"GetUserKeys() gotStatus = %v, want %v",
					gotStatus,
					tt.wantStatus,
				)
			}

			gotResponse := got.(*models.UserDevices)
			gotUserID := gotResponse.UserID
			gotDevices := gotResponse.Devices

			if gotUserID != tt.want.UserID {
				t.Errorf(
					"GetUserKeys() got UserID = %v, want %v",
					gotUserID,
					tt.want.UserID,
				)
				return
			}

			if len(gotDevices) != len(tt.want.Devices) {
				t.Errorf(
					"GetUserKeys() got Devices %v, want %v",
					gotDevices,
					tt.want.UserID,
				)
				return
			}

			for _, wantDevice := range tt.want.Devices {
				found := false

				for _, gotDevice := range gotDevices {
					if gotDevice.ID == wantDevice.ID {
						found = true
						break
					}
				}

				if !found {
					t.Errorf(
						"GetUserKeys() got Devices %v, want %v",
						gotDevices,
						tt.want.UserID,
					)
					return
				}
			}
		})
	}
}

func seedOnlyDevice() (device models.Device) {
	err := repo.NewRepo().GetDb().Transaction(func(db *gorm.DB) error {
		faker.FakeData(&device)
		db.Create(&device)

		return db.Error
	})
	if err != nil {
		panic(err)
	}

	return device
}

func teardownOnlyDevice(device models.Device) {
	err := repo.NewRepo().GetDb().Transaction(func(db *gorm.DB) error {
		db.Delete(&device)
		return db.Error
	})
	if err != nil {
		panic(err)
	}
}

func seedLoginRequest() (lr models.LoginRequest) {
	err := repo.NewRepo().GetDb().Transaction(func(db *gorm.DB) error {
		lr.TemporaryCode = faker.UUIDDigit()

		db.Create(&lr)

		return db.Error
	})
	if err != nil {
		panic(err)
	}

	return lr
}

func teardownLoginRequest(lr models.LoginRequest) {
	err := repo.NewRepo().GetDb().Transaction(func(db *gorm.DB) error {
		db.Delete(&lr)

		return db.Error
	})
	if err != nil {
		panic(err)
	}
}
