// +build test

package controllers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func TestPostInvite(t *testing.T) {
	user, otherUser, project := seedPostInvite()
	defer teardownPostInvite(user, otherUser, project)

	type args struct {
		in0  router.Params
		body io.ReadCloser
		Repo repo.IRepo
		user models.User
	}
	tests := []struct {
		name       string
		args       args
		want       *models.GetInviteResponse
		wantStatus int
		wantErr    string
	}{
		{
			name: "it invites some user",
			args: args{
				in0: router.Params{},
				body: io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`
{
    "Email": "some-user@some-mail-service.com",
    "ProjectName": "%s"
}
`, project.Name))),
				Repo: new(repo.Repo),
				user: user,
			},
			want: &models.GetInviteResponse{
				UserUIDs: []string{},
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "user has an account",
			args: args{
				in0: router.Params{},
				body: io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`
{
    "Email": "%s",
    "ProjectName": "%s"
}
`, otherUser.Email, project.Name))),
				Repo: new(repo.Repo),
				user: user,
			},
			want: &models.GetInviteResponse{
				UserUIDs: []string{
					otherUser.UserID,
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "error if project does not exist",
			args: args{
				in0: router.Params{},
				body: io.NopCloser(bytes.NewBufferString(`
{
    "Email": "some.user@some.mail.com",
    "ProjectName": "that-is-not-a-project"
}
`)),
				Repo: new(repo.Repo),
				user: user,
			},
			want: &models.GetInviteResponse{
				UserUIDs: []string{},
			},
			wantStatus: http.StatusNotFound,
			wantErr:    "failed to get: not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotStatus, err := PostInvite(
				tt.args.in0,
				tt.args.body,
				tt.args.Repo,
				tt.args.user,
			)
			if err.Error() != tt.wantErr {
				t.Errorf("PostInvite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStatus != tt.wantStatus {
				t.Errorf(
					"PostInvite() gotStatus = %v, want %v",
					gotStatus,
					tt.wantStatus,
				)
				return
			}

			gotResponse := got.(*models.GetInviteResponse)
			gotUserUIDs := gotResponse.UserUIDs
			wantUserUIDs := tt.want.UserUIDs
			gotLen := len(gotUserUIDs)
			wantLen := len(wantUserUIDs)

			if gotLen != wantLen {
				t.Errorf(
					"PostInvite() got len = %v, want len %v",
					gotLen,
					wantLen,
				)
				return
			}

			for index, gotUID := range gotUserUIDs {
				wantUID := wantUserUIDs[index]

				if gotUID != wantUID {
					t.Errorf("PostInvite() got = %v, want %v", gotUID, wantUID)
				}
			}
		})
	}
}

func seedPostInvite() (user, otherUser models.User, project models.Project) {
	db := new(repo.Repo).GetDb()
	var projectMember models.ProjectMember

	faker.FakeData(&user)
	faker.FakeData(&otherUser)
	faker.FakeData(&project)

	db.Create(&project)
	db.Create(&user)
	db.Create(&otherUser)

	projectMember = models.ProjectMember{
		User:      user,
		UserID:    user.ID,
		Project:   project,
		ProjectID: project.ID,
		RoleID:    1, //admin
	}
	db.Create(&projectMember)

	return user, otherUser, project
}

func teardownPostInvite(user, otherUser models.User, project models.Project) {
	db := new(repo.Repo).GetDb()

	db.Exec(
		"delete from project_members where project_id = ?",
		project.ID,
	)
	db.Exec("delete from projects where id = ?", project.ID)
	db.Exec("delete from users where id = ?", user.ID)
	db.Exec("delete from users where id = ?", otherUser.ID)
}
