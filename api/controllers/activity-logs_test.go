// +build test

package controllers

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func TestGetActivityLogs(t *testing.T) {
	project, organization, user := seedActivityLogs(true) // paid
	unpaidProject, uo, uu := seedActivityLogs(false)      // unpaid

	type args struct {
		params router.Params
		body   io.ReadCloser
		Repo   repo.IRepo
		in3    models.User
	}
	tests := []struct {
		name       string
		args       args
		want       *models.GetActivityLogResponse
		wantStatus int
		wantErr    bool
	}{
		{
			name: "returns-all-logs",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
				}),
				body: ioutil.NopCloser(strings.NewReader(`
{
    "actions": [],
    "environments": [],
    "users": [],
    "limit": 200
}
`)),
				Repo: &repo.Repo{},
				in3:  models.User{},
			},
			want: &models.GetActivityLogResponse{
				// NOTE: Entries must be in reverse order from their
				// insertion
				Logs: []models.ActivityLogLite{
					{
						UserID:          user.UserID,
						ProjectName:     project.Name,
						EnvironmentName: "staging",
						Action:          "GetAccessibleEnvironments",
						Success:         true,
						ErrorMessage:    "",
					},
					{
						UserID:          user.UserID,
						ProjectName:     project.Name,
						EnvironmentName: "dev",
						Action:          "GetEnvironmentPublicKeys",
						Success:         false,
						ErrorMessage:    "An error occured while getting the environment public keys",
					},
					{
						UserID:          user.UserID,
						ProjectName:     project.Name,
						EnvironmentName: "dev",
						Action:          "GetEnvironmentPublicKeys",
						Success:         true,
						ErrorMessage:    "",
					},
					{
						UserID:          user.UserID,
						ProjectName:     project.Name,
						EnvironmentName: "prod",
						Action:          "GetMessagesFromProjectByUser",
						Success:         true,
						ErrorMessage:    "",
					},
					{
						UserID:          user.UserID,
						ProjectName:     project.Name,
						EnvironmentName: "prod",
						Action:          "GetMessagesFromProjectByUser",
						Success:         true,
						ErrorMessage:    "",
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "filters-by-action",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
				}),
				body: ioutil.NopCloser(strings.NewReader(`
{
    "actions": ["GetEnvironmentPublicKeys"],
    "environments": [],
    "users": [],
    "limit": 200
}
`)),
				Repo: &repo.Repo{},
				in3:  models.User{},
			},
			want: &models.GetActivityLogResponse{
				// NOTE: Entries must be in reverse order from their
				// insertion
				Logs: []models.ActivityLogLite{
					{
						UserID:          user.UserID,
						ProjectName:     project.Name,
						EnvironmentName: "dev",
						Action:          "GetEnvironmentPublicKeys",
						Success:         false,
						ErrorMessage:    "An error occured while getting the environment public keys",
					},
					{
						UserID:          user.UserID,
						ProjectName:     project.Name,
						EnvironmentName: "dev",
						Action:          "GetEnvironmentPublicKeys",
						Success:         true,
						ErrorMessage:    "",
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "filters-by-environment",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
				}),
				body: ioutil.NopCloser(strings.NewReader(`
{
    "actions": [],
    "environments": ["dev"],
    "users": [],
    "limit": 200
}
`)),
				Repo: &repo.Repo{},
				in3:  models.User{},
			},
			want: &models.GetActivityLogResponse{
				// NOTE: Entries must be in reverse order from their
				// insertion
				Logs: []models.ActivityLogLite{
					{
						UserID:          user.UserID,
						ProjectName:     project.Name,
						EnvironmentName: "dev",
						Action:          "GetEnvironmentPublicKeys",
						Success:         false,
						ErrorMessage:    "An error occured while getting the environment public keys",
					},
					{
						UserID:          user.UserID,
						ProjectName:     project.Name,
						EnvironmentName: "dev",
						Action:          "GetEnvironmentPublicKeys",
						Success:         true,
						ErrorMessage:    "",
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "errors-on-bad-project",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": "not-a-project-uui",
				}),
				body: ioutil.NopCloser(strings.NewReader(`
{
    "actions": [],
    "environments": [],
    "users": [],
    "limit": 200
}
`)),
				Repo: &repo.Repo{},
				in3:  models.User{},
			},
			want: &models.GetActivityLogResponse{
				Logs: []models.ActivityLogLite{},
			},
			wantStatus: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name: "errors-on-organization-not-paid",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": unpaidProject.UUID,
				}),
				body: ioutil.NopCloser(strings.NewReader(`
{
    "actions": [],
    "environments": [],
    "users": [],
    "limit": 200
}
`)),
				Repo: &repo.Repo{},
				in3:  models.User{},
			},
			want: &models.GetActivityLogResponse{
				Logs: []models.ActivityLogLite{},
			},
			wantStatus: http.StatusForbidden,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotStatus, err := GetActivityLogs(tt.args.params, tt.args.body, tt.args.Repo, tt.args.in3)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetActivityLogs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("GetActivityLogs() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
			gotR := got.(*models.GetActivityLogResponse)
			if len(gotR.Logs) != len(tt.want.Logs) {
				t.Errorf("GetActivityLogs() len = %d, want %d", len(gotR.Logs), len(tt.want.Logs))
				return
			}

			for i, gotEntry := range gotR.Logs {
				wantEntry := tt.want.Logs[i]

				if gotEntry.UserID != wantEntry.UserID ||
					gotEntry.ProjectName != wantEntry.ProjectName ||
					gotEntry.EnvironmentName != wantEntry.EnvironmentName ||
					gotEntry.ErrorMessage != wantEntry.ErrorMessage ||
					gotEntry.Action != wantEntry.Action ||
					gotEntry.Success != wantEntry.Success {
					t.Errorf(
						"GetActivityLogs() entry %d = %v, want %v",
						i,
						gotEntry,
						wantEntry,
					)
				}
			}
		})
	}

	teardownActivityLogs(project, organization, user)
	teardownActivityLogs(unpaidProject, uo, uu)
}

func seedActivityLogs(paid bool) (models.Project, models.Organization, models.User) {
	var project models.Project
	var organization models.Organization
	var user models.User
	var device models.Device

	Repo := new(repo.Repo)
	db := Repo.GetDb()

	user = models.User{}
	device = models.Device{}
	faker.FakeData(&user)
	faker.FakeData(&device)

	user.Devices = []models.Device{device}
	db.Save(&user)

	organization = models.Organization{}
	faker.FakeData(&organization)
	organization.UserID = user.ID
	organization.Name = user.UserID
	organization.Paid = paid
	db.Save(&organization)

	project = models.Project{
		UUID:                faker.UUIDHyphenated(),
		TTL:                 30,
		DaysBeforeTTLExpiry: 7,
		Name:                "activity-logs-project",
		UserID:              user.ID,
		Environments:        []models.Environment{},
		OrganizationID:      organization.ID,
		Organization:        organization,
	}
	Repo.GetOrCreateProject(&project)
	dev := project.GetEnvironment("dev")
	staging := project.GetEnvironment("staging")
	prod := project.GetEnvironment("prod")

	if dev == nil {
		panic(fmt.Sprintf("dev environment not found for %d", project.ID))
	}
	if staging == nil {
		panic(fmt.Sprintf("staging environment not found for %d", project.ID))
	}
	if prod == nil {
		panic(fmt.Sprintf("prod environment not found for %d", project.ID))
	}

	err := Repo.GetDb().Exec(
		`insert into activity_logs (user_id, project_id, environment_id, action, success, message, created_at, updated_at)
values
(@userID, @projectID, @prodID,    "GetMessagesFromProjectByUser", true,  "",                                                           CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(@userID, @projectID, @prodID,    "GetMessagesFromProjectByUser", true,  "",                                                           CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(@userID, @projectID, @devID,     "GetEnvironmentPublicKeys",     true,  "",                                                           CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(@userID, @projectID, @devID,     "GetEnvironmentPublicKeys",     false, "An error occured while getting the environment public keys", CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(@userID, @projectID, @stagingID, "GetAccessibleEnvironments",    true,  "",                                                           CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
`,
		sql.Named("userID", user.ID),
		sql.Named("projectID", project.ID),
		sql.Named("prodID", prod.ID),
		sql.Named("stagingID", staging.ID),
		sql.Named("devID", dev.ID),
	).Error

	if err != nil {
		panic(err)
	}

	if err := Repo.Err(); err != nil {
		panic(err)
	}

	return project, organization, user
}

func teardownActivityLogs(p models.Project, o models.Organization, u models.User) {
	Repo := new(repo.Repo)
	db := Repo.GetDb()

	db.Exec(
		`delete from activity_logs where project_id = @project_id`,
		sql.Named("project_id", p.ID),
	)
	db.Exec(
		`delete from environments where project_id = @project_id`,
		sql.Named("project_id", p.ID),
	)
	db.Exec(
		`delete from projects where id = @project_id`,
		sql.Named("project_id", p.ID),
	)
	db.Exec(
		`delete from organizations where id = @project_id`,
		sql.Named("organization_id", p.ID),
	)
	db.Exec(
		`delete from users where id = @user_id`,
		sql.Named("user_id", p.ID),
	)

	if db.Error != nil {
		panic(db.Error)
	}
}
