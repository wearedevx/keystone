//go:build test
// +build test

package controllers

import (
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func TestGetDevices(t *testing.T) {
	Repo := new(repo.Repo)
	user, device := seedDevice(Repo)
	defer teardownDevice(user, device)

	type args struct {
		in0  router.Params
		in1  io.ReadCloser
		Repo repo.IRepo
		user models.User
	}
	tests := []struct {
		name           string
		args           args
		want           *models.GetDevicesResponse
		wantDevicesLen int
		wantStatus     int
		wantErr        string
	}{
		{
			name: "returns-user-device",
			args: args{
				in0:  router.Params{},
				in1:  nil,
				Repo: Repo,
				user: user,
			},
			want: &models.GetDevicesResponse{
				Devices: []models.Device{device},
			},
			wantDevicesLen: 1,
			wantStatus:     http.StatusOK,
			wantErr:        "",
		},
		{
			name: "returns-empty list if bad user",
			args: args{
				in0:  router.Params{},
				in1:  nil,
				Repo: Repo,
				user: models.User{
					ID: 110010,
				},
			},
			want: &models.GetDevicesResponse{
				Devices: []models.Device{device},
			},
			wantDevicesLen: 0,
			wantStatus:     http.StatusOK,
			wantErr:        "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotStatus, err := GetDevices(tt.args.in0, tt.args.in1, tt.args.Repo, tt.args.user)
			if err.Error() != tt.wantErr {
				t.Errorf("GetDevices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotStatus != tt.wantStatus {
				t.Errorf("GetDevices() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}

			gotResponse := got.(*models.GetDevicesResponse)
			gotDevicesLen := len(gotResponse.Devices)
			if gotDevicesLen != tt.wantDevicesLen {
				t.Errorf("GetDevices() len = %d, want %d", gotDevicesLen, tt.wantDevicesLen)
				return
			}

			if tt.wantDevicesLen == 0 {
				return
			}

			gotDevice := gotResponse.Devices[0]
			wantDevice := tt.want.Devices[0]

			if gotDevice.ID != wantDevice.ID ||
				!reflect.DeepEqual(gotDevice.PublicKey, wantDevice.PublicKey) ||
				gotDevice.Name != wantDevice.Name ||
				gotDevice.UID != wantDevice.UID {
				t.Errorf("GetDevices() got = %v, want %v", gotDevice, wantDevice)
			}

		})
	}
}

func TestDeleteDevice(t *testing.T) {
	Repo := new(repo.Repo)
	user, device := seedDevice(Repo)
	defer teardownDevice(user, device)

	type args struct {
		params router.Params
		in1    io.ReadCloser
		Repo   repo.IRepo
		user   models.User
	}
	tests := []struct {
		name       string
		args       args
		want       *models.RemoveDeviceResponse
		wantStatus int
		wantErr    string
	}{
		{
			name: "deletes a device",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"uid": device.UID,
				}),
				in1:  nil,
				Repo: Repo,
				user: user,
			},
			want:       nil,
			wantStatus: http.StatusNoContent,
			wantErr:    "",
		},
		{
			name: "returns not found",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"uid": "that is no uid",
				}),
				in1:  nil,
				Repo: Repo,
				user: user,
			},
			want:       &models.RemoveDeviceResponse{Success: false, Error: "not found"},
			wantStatus: http.StatusNotFound,
			wantErr:    "no device",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotStatus, err := DeleteDevice(tt.args.params, tt.args.in1, tt.args.Repo, tt.args.user)
			if err.Error() != tt.wantErr {
				t.Errorf("DeleteDevice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteDevice() got = %v, want %v", got, tt.want)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("DeleteDevice() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}

}

func seedDevice(Repo *repo.Repo) (user models.User, device models.Device) {
	db := Repo.GetDB()

	user = models.User{}
	device = models.Device{}

	faker.FakeData(&user)
	faker.FakeData(&device)

	db.Omit("Devices").Save(&user)
	db.Omit("Users").Save(&device)

	user.Devices = append(user.Devices, device)

	db.Save(&user)

	return user, device
}

func teardownDevice(user models.User, device models.Device) {
	db := new(repo.Repo).GetDB()

	db.Exec(
		"delete from user_devices where user_id = ? and device_id = ?",
		user.ID,
		device.ID,
	).
		Exec(
			"delete from devices where id = ?",
			device.ID,
		).
		Exec(
			"delete from users where id = ?",
			user.ID,
		)
}
