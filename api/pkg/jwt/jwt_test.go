package jwt

import (
	"testing"
	"time"

	"github.com/wearedevx/keystone/api/pkg/models"
)

func TestMakeToken(t *testing.T) {
	salt = "testing-salt"
	type args struct {
		user      models.User
		deviceUID string
		when      time.Time
	}
	when, _ := time.Parse("2006-01-02T15:04:05Z", "2021-11-11T16:01:01")
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "",
			args: args{
				user: models.User{
					UserID: "memberx@github",
				},
				deviceUID: "unique-device-id",
				when:      when,
			},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkZXZpY2VfdWlkIjoidW5pcXVlLWRldmljZS1pZCIsImV4cCI6LTYyMTMzMDA0ODAwLCJpYXQiOi02MjEzNTU5NjgwMCwiaXNzIjoia2V5c3RvbmUiLCJzdWIiOiJtZW1iZXJ4QGdpdGh1YiJ9.QlPGFzAUn_yhlNoWwEJ3hao0YZQERX5v21H5sfZAnxs",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MakeToken(tt.args.user, tt.args.deviceUID, tt.args.when)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MakeToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifyToken(t *testing.T) {
	type args struct {
		token string
	}
	okToken, _ := MakeToken(models.User{UserID: "memberx@github"}, "unique-device-id", time.Now())
	tests := []struct {
		name         string
		args         args
		wantUserID   string
		wantDeviceID string
		wantErr      bool
	}{
		{
			name: "verify-token-ok",
			args: args{
				token: okToken,
			},
			wantUserID:   "memberx@github",
			wantDeviceID: "unique-device-id",
			wantErr:      false,
		},

		{
			name: "verify-token-Bearer-ok",
			args: args{
				token: "Bearer " + okToken,
			},
			wantUserID:   "memberx@github",
			wantDeviceID: "unique-device-id",
			wantErr:      false,
		},
		{
			name: "verify-token-invalid",
			args: args{
				token: "Bearer eyJhbGcwEJ3hao0YZQERX5v21H5sfZAnxs",
			},
			wantUserID:   "",
			wantDeviceID: "",
			wantErr:      true,
		},
		{
			name: "verify-token-expired",
			args: args{
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkZXZpY2VfdWlkIjoidW5pcXVlLWRldmljZS1pZCIsImV4cCI6LTYyMTMzMDA0ODAwLCJpYXQiOi02MjEzNTU5NjgwMCwiaXNzIjoia2V5c3RvbmUiLCJzdWIiOiJtZW1iZXJ4QGdpdGh1YiJ9.QlPGFzAUn_yhlNoWwEJ3hao0YZQERX5v21H5sfZAnxs",
			},
			wantUserID:   "",
			wantDeviceID: "",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := VerifyToken(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantUserID {
				t.Errorf("VerifyToken() got = %v, want %v", got, tt.wantUserID)
			}
			if got1 != tt.wantDeviceID {
				t.Errorf("VerifyToken() got1 = %v, want %v", got1, tt.wantDeviceID)
			}
		})
	}
}
