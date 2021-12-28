package controllers

import (
	"bytes"
	"io"
	"net/http"
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
			got, got1, err := GetUser(tt.args.in0, tt.args.in1, tt.args.in2, tt.args.user)
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
				in2: []httprouter.Param{}},
			wantStatus:        http.StatusOK,
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
				in2: []httprouter.Param{}},
			wantStatus:        http.StatusConflict,
			wantAuthorization: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PostUserToken(tt.args.w, tt.args.r, tt.args.in2)
			got := tt.args.w.(*mockResponseWriter)

			if got.status != tt.wantStatus {
				t.Errorf("PostUserToken() got status %v, want %v", got.status, tt.wantStatus)
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
	type args struct {
		w   http.ResponseWriter
		r   *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name string
		args args
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetAuthRedirect(tt.args.w, tt.args.r, tt.args.in2)
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
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PostLoginRequest(tt.args.w, tt.args.in1, tt.args.in2)
		})
	}
}

func TestGetLoginRequest(t *testing.T) {
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
			GetLoginRequest(tt.args.w, tt.args.r, tt.args.in2)
		})
	}
}

func TestGetUserKeys(t *testing.T) {
	type args struct {
		params router.Params
		in1    io.ReadCloser
		Repo   repo.IRepo
		in3    models.User
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
			got, gotStatus, err := GetUserKeys(tt.args.params, tt.args.in1, tt.args.Repo, tt.args.in3)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserKeys() got = %v, want %v", got, tt.want)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("GetUserKeys() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
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
