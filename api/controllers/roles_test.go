package controllers

import (
	"io"
	"reflect"
	"testing"

	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func TestGetRoles(t *testing.T) {
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
			got, gotStatus, err := GetRoles(tt.args.params, tt.args.in1, tt.args.Repo, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRoles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRoles() got = %v, want %v", got, tt.want)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("GetRoles() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}
