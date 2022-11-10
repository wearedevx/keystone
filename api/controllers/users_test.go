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
	"strings"
	"testing"

	"github.com/bxcodec/faker/v3"
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
		w    http.ResponseWriter
		r    *http.Request
		in2  router.Params
		Repo repo.IRepo
	}
	tests := []struct {
		name              string
		args              args
		wantStatus        int
		wantAuthorization bool
		wantErr           bool
		wantMsg           string
	}{
		{
			name: "it works",
			args: args{
				w: newMockResponse(),
				r: &http.Request{Body: ioutil.NopCloser(bytes.NewBufferString(`
							{
								"AccountType": "github",
								"Token": {
									"access_token": "YSB0b2tlbg=="
								},
								"PublicKey": "YSB2ZXJ5IHB1YmxpYyBrZXk=",
								"Device": "a-device",
								"DeviceUID": "a-device-uid"
							}`))},
				in2:  router.Params{},
				Repo: newFakeRepo(noCrashers),
			},
			wantStatus:        http.StatusOK,
			wantAuthorization: true,
			wantErr:           false,
		},
		{
			name: "bad device name",
			args: args{
				w: newMockResponse(),
				r: &http.Request{Body: ioutil.NopCloser(bytes.NewBufferString(`
							{
								"AccountType": "github",
								"Token": {
									"access_token": "YSB0b2tlbg=="
								},
								"PublicKey": "YSB2ZXJ5IHB1YmxpYyBrZXk=",
								"Device": "is that such a bad device name ?",
								"DeviceUID": "a-device-uid"
							}`))},
				in2:  router.Params{},
				Repo: newFakeRepo(noCrashers),
			},
			wantStatus:        http.StatusConflict,
			wantAuthorization: false,
			wantErr:           true,
			wantMsg:           "bad device name",
		},
		{
			name: "bad payload",
			args: args{
				w: newMockResponse(),
				r: &http.Request{Body: ioutil.NopCloser(bytes.NewBufferString(`
									{
										"AccountType": "github",
										"Token": 
										"PublicKey": "YSB2ZXJ5IHB1YmxpYyBrZXk=",
										"Device": "is that such a bad device name ?",
										"DeviceUID": "a-device-uid"
									}`))},
				in2:  router.Params{},
				Repo: newFakeRepo(noCrashers),
			},
			wantStatus:        http.StatusBadRequest,
			wantAuthorization: false,
			wantErr:           true,
			wantMsg:           "Bad Request",
		},
		{
			name: "missing public key",
			args: args{
				w: newMockResponse(),
				r: &http.Request{Body: ioutil.NopCloser(bytes.NewBufferString(`
				{
					"AccountType": "github",
					"Token": {
						"access_token": "YSB0b2tlbg=="
					},
					"PublicKey": "",
					"Device": "is that such a bad device name ?",
					"DeviceUID": "a-device-uid"
				}`))},
				in2:  router.Params{},
				Repo: newFakeRepo(noCrashers),
			},
			wantStatus:        http.StatusBadRequest,
			wantAuthorization: false,
			wantErr:           true,
			wantMsg:           "Bad Request",
		},
		{
			name: "fails to create a user",
			args: args{
				w: newMockResponse(),
				r: &http.Request{Body: ioutil.NopCloser(bytes.NewBufferString(`
						{
							"AccountType": "github",
							"Token": {
								"access_token": "YSB0b2tlbg=="
							},
							"PublicKey": "YSB2ZXJ5IHB1YmxpYyBrZXk=",
							"Device": "a-good-device-name",
							"DeviceUID": "a-device-uid"
						}`))},
				in2: router.Params{},
				Repo: newFakeRepo(map[string]error{
					"GetOrCreateUser": errors.New("unexpected error"),
				}),
			},
			wantStatus:        http.StatusInternalServerError,
			wantAuthorization: false,
			wantErr:           true,
			wantMsg:           "Internal Server Error",
		},
		{
			name: "fails to get newly created devices",
			args: args{
				w: newMockResponse(),
				r: &http.Request{Body: ioutil.NopCloser(bytes.NewBufferString(`
								{
									"AccountType": "github",
									"Token": {
										"access_token": "YSB0b2tlbg=="
									},
									"PublicKey": "YSB2ZXJ5IHB1YmxpYyBrZXk=",
									"Device": "another-unique-device-name",
									"DeviceUID": "a-device-uid"
								}`))},
				in2: router.Params{},
				Repo: newFakeRepo(map[string]error{
					"GetNewlyCreatedDevices": errors.New("unexpected error"),
				}),
			},
			wantStatus:        http.StatusInternalServerError,
			wantAuthorization: false,
			wantErr:           true,
			wantMsg:           "Internal Server Error",
		},
		{
			name: "fails to find admins of projects",
			args: args{
				w: newMockResponse(),
				r: &http.Request{Body: ioutil.NopCloser(bytes.NewBufferString(`
										{
											"AccountType": "github",
											"Token": {
												"access_token": "YSB0b2tlbg=="
											},
											"PublicKey": "YSB2ZXJ5IHB1YmxpYyBrZXk=",
											"Device": "yet-another-device-name",
											"DeviceUID": "a-device-uid"
										}`))},
				in2: router.Params{},
				Repo: newFakeRepo(map[string]error{
					"GetAdminsFromUserProjects": errors.New("unexpected error"),
				}),
			},
			wantStatus:        http.StatusInternalServerError,
			wantAuthorization: false,
			wantErr:           true,
			wantMsg:           "unexpected error",
		},
		{
			name: "fais setting newly created device flags",
			args: args{
				w: newMockResponse(),
				r: &http.Request{Body: ioutil.NopCloser(bytes.NewBufferString(`
												{
													"AccountType": "github",
													"Token": {
														"access_token": "YSB0b2tlbg=="
													},
													"PublicKey": "YSB2ZXJ5IHB1YmxpYyBrZXk=",
													"Device": "yet-another-unique-device-name",
													"DeviceUID": "a-device-uid"
												}`))},
				in2: router.Params{},
				Repo: newFakeRepo(map[string]error{
					"SetNewlyCreatedDevice": errors.New("unexpected error"),
				}),
			},
			wantStatus:        http.StatusInternalServerError,
			wantAuthorization: false,
			wantErr:           true,
			wantMsg:           "unexpected error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := PostUserToken(tt.args.w, tt.args.r, tt.args.in2, tt.args.Repo)
			got := tt.args.w.(*mockResponseWriter)

			if (err != nil) != tt.wantErr {
				t.Errorf(
					"PostUserToken() err %v, want %v",
					err,
					tt.wantErr,
				)
			}

			if status != tt.wantStatus {
				t.Errorf(
					"PostUserToken() got status %v, want %v",
					status,
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

			if tt.wantErr {
				gotMsg := got.body.String()
				if gotMsg != tt.wantMsg {
					t.Errorf("PostUserToken() msg = %v, want %v", gotMsg, tt.wantMsg)
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

	okURL, _ := url.Parse(
		fmt.Sprintf(
			"http://tests.com/redirect?state=%s&code=123456",
			goodState,
		),
	)
	tooShortCodeURL, _ := url.Parse(
		fmt.Sprintf(
			"http://tests.com/redirect?state=%s&code=123456",
			tooShortState,
		),
	)
	notFoundURL, _ := url.Parse(
		fmt.Sprintf(
			"http://tests.com/redirect?state=%s&code=123456",
			notFoundState,
		),
	)
	badURL, _ := url.Parse("http://tests.com/redirect?state=rubbish")

	type args struct {
		w    http.ResponseWriter
		r    *http.Request
		in2  router.Params
		Repo repo.IRepo
	}
	tests := []struct {
		name         string
		args         args
		wantStatus   int
		wantResponse []string
		wantErr      bool
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
			wantStatus: http.StatusOK,
			wantResponse: []string{
				"You have been successfully authenticated",
				"You may now return to your terminal and start using Keystone",
			},
		},
		{
			name: "bad url",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					URL: badURL,
				},
				in2:  router.Params{},
				Repo: newFakeRepo(noCrashers),
			},
			wantStatus: http.StatusBadRequest,
			wantResponse: []string{
				"Bad Request",
				"The link used is malformed",
			},
		},
		{
			name: "temporary code is too short",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					URL: tooShortCodeURL,
				},
				in2:  router.Params{},
				Repo: newFakeRepo(noCrashers),
			},
			wantStatus: http.StatusBadRequest,
			wantResponse: []string{
				"Bad Request",
				"The provided code is invalid",
			},
		},
		{
			name: "no such request",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					URL: notFoundURL,
				},
				in2:  router.Params{},
				Repo: newFakeRepo(noCrashers),
			},
			wantStatus: http.StatusBadRequest,
			wantResponse: []string{
				"Bad Request",
				"The link used is invalid or expired",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := GetAuthRedirect(tt.args.w, tt.args.r, tt.args.in2, tt.args.Repo)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAuthRedirect() err = %v, want %v", err, tt.wantErr)
			}

			got := tt.args.w.(*mockResponseWriter)

			if status != tt.wantStatus {
				t.Errorf(
					"GetAuthRedirect() got status %v, want %v",
					status,
					tt.wantStatus,
				)
				return
			}

			gotResponse := got.body.String()
			for _, want := range tt.wantResponse {
				if !strings.Contains(gotResponse, want) {
					t.Errorf(
						"GetAuthRedirec() want %v, not found in response %s",
						want,
						gotResponse,
					)
					return
				}
			}
		})
	}
}

func TestPostLoginRequest(t *testing.T) {
	type args struct {
		w    http.ResponseWriter
		in1  *http.Request
		in2  router.Params
		Repo repo.IRepo
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
		wantMsg    string
	}{
		{
			name: "it work",
			args: args{
				w:    newMockResponse(),
				in1:  &http.Request{},
				in2:  router.Params{},
				Repo: newFakeRepo(noCrashers),
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "it fails interting in db",
			args: args{
				w:   newMockResponse(),
				in1: &http.Request{},
				in2: router.Params{},
				Repo: newFakeRepo(map[string]error{
					"CreateLoginRequest": errors.New("unexpected error"),
				}),
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    true,
			wantMsg:    "Status Internal Server Error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := PostLoginRequest(tt.args.w, tt.args.in1, tt.args.in2, tt.args.Repo)

			if (err != nil) != tt.wantErr {
				t.Errorf(
					"PostLoginRequest() err = %v, want %v",
					err,
					tt.wantErr,
				)

				return
			}

			got := tt.args.w.(*mockResponseWriter)

			if status != tt.wantStatus {
				t.Errorf(
					"PostLoginRequest() got status %v, want %v",
					got.status,
					tt.wantStatus,
				)
				return
			}

			if !tt.wantErr {
				var lr models.LoginRequest
				if err := lr.Deserialize(got.body); err != nil {
					t.Errorf(
						"PostLoginRequest() invalid respons body: %v",
						got.body.String(),
					)
				}
			} else {
				var msg = got.body.String()
				if msg != tt.wantMsg {
					t.Errorf(
						"PostLoginRequest() msg = %v, want %v",
						msg,
						tt.wantMsg,
					)
				}
			}
		})
	}
}

func TestGetLoginRequest(t *testing.T) {
	lr := seedLoginRequest()
	defer teardownLoginRequest(lr)

	okURL, _ := url.Parse(
		fmt.Sprintf("http://tests.com/login-request?code=%s", lr.TemporaryCode),
	)
	notFoundURL, _ := url.Parse(
		"http://tests.com/login-request?code=notacodebutitcouldhavebeenjusthavetobelongenough",
	)
	tooShortURL, _ := url.Parse(
		"http://tests.com/login-request?code=notaco",
	)

	type args struct {
		w    http.ResponseWriter
		r    *http.Request
		in2  router.Params
		Repo repo.IRepo
	}
	tests := []struct {
		name        string
		args        args
		wantStatus  int
		want        *models.LoginRequest
		wantMessage string
		wantErr     bool
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
			wantStatus:  http.StatusOK,
			want:        &lr,
			wantMessage: "",
			wantErr:     false,
		},
		{
			name: "not found",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					URL: notFoundURL,
				},
				in2:  router.Params{},
				Repo: newFakeRepo(noCrashers),
			},
			wantStatus:  http.StatusNotFound,
			want:        nil,
			wantMessage: "not found",
			wantErr:     true,
		},
		{
			name: "temporary code is too short",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					URL: tooShortURL,
				},
				in2:  router.Params{},
				Repo: newFakeRepo(noCrashers),
			},
			wantStatus:  http.StatusBadRequest,
			want:        nil,
			wantMessage: "Bad Request",
			wantErr:     true,
		},
		{
			name: "it fails",
			args: args{
				w: newMockResponse(),
				r: &http.Request{
					URL: okURL,
				},
				in2: router.Params{},
				Repo: newFakeRepo(map[string]error{
					"GetLoginRequest": errors.New("unexpected error"),
				}),
			},
			wantStatus:  http.StatusInternalServerError,
			want:        nil,
			wantMessage: "unexpected error",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := GetLoginRequest(tt.args.w, tt.args.r, tt.args.in2, tt.args.Repo)

			if (err != nil) != tt.wantErr {
				t.Errorf(
					"GetLoginRequest() err = %v, want %v",
					err,
					tt.wantErr,
				)
				return
			}

			got := tt.args.w.(*mockResponseWriter)

			if status != tt.wantStatus {
				t.Errorf(
					"GetLoginRequest() got status %v, want %v",
					status,
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
					t.Errorf(
						"GetLoginRequest() got %v, want %v",
						got.body.String(),
						tt.wantMessage,
					)
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
				Repo: newFakeRepo(noCrashers),
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
				Repo:   newFakeRepo(noCrashers),
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
				Repo: newFakeRepo(noCrashers),
				in3:  user,
			},
			want:       &models.UserDevices{},
			wantStatus: http.StatusNotFound,
			wantErr:    "not found",
		},
		{
			name: "fails to get the user",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"userID": user.UserID,
				}),
				in1: nil,
				Repo: newFakeRepo(map[string]error{
					"GetUser": errors.New("unexpected error"),
				}),
				in3: models.User{},
			},
			want:       &models.UserDevices{},
			wantStatus: http.StatusInternalServerError,
			wantErr:    "failed to get: unexpected error",
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
	err := repo.NewRepo().GetDB().Transaction(func(db *gorm.DB) error {
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
	err := repo.NewRepo().GetDB().Transaction(func(db *gorm.DB) error {
		db.Delete(&device)
		return db.Error
	})
	if err != nil {
		panic(err)
	}
}

func seedLoginRequest() (lr models.LoginRequest) {
	err := repo.NewRepo().GetDB().Transaction(func(db *gorm.DB) error {
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
	err := repo.NewRepo().GetDB().Transaction(func(db *gorm.DB) error {
		db.Delete(&lr)

		return db.Error
	})
	if err != nil {
		panic(err)
	}
}
