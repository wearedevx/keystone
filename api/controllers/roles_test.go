package controllers

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/wearedevx/keystone/api/db/seed"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
	"gorm.io/gorm"
)

func TestGetRoles(t *testing.T) {
	user, organization, project := seedOneProjectForOneUser()
	defer teardownUserAndOrganization(user, organization)
	defer teardownProject(project)

	_, paidOrganization, paidProject := seedOneProjectForOneUser()
	defer teardownUserAndOrganization(user, organization)
	defer teardownProject(project)

	testsSetOrganisationPaid(&paidOrganization)

	type args struct {
		params router.Params
		in1    io.ReadCloser
		Repo   repo.IRepo
		user   models.User
	}
	tests := []struct {
		name               string
		args               args
		want               *models.GetRolesResponse
		wantRoles          []string
		wantStatus         int
		wantErr            string
		wantTearDownBefore bool
	}{
		{
			name: "returns only admin role because not paid organization",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
				}),
				Repo: newFakeRepo(noCrashers),
				user: models.User{},
			},
			wantRoles:          []string{"admin"},
			wantStatus:         http.StatusOK,
			wantErr:            "",
			wantTearDownBefore: false,
		},
		{
			name: "returns all 4 roles because paid organization",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": paidProject.UUID,
				}),
				Repo: newFakeRepo(noCrashers),
				user: models.User{},
			},
			wantRoles:          []string{"admin", "developer", "devops", "lead-dev"},
			wantStatus:         http.StatusOK,
			wantErr:            "",
			wantTearDownBefore: false,
		},
		{
			name: "returns error because project does not exist",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": "doesnotexist",
				}),
				Repo: newFakeRepo(noCrashers),
				user: models.User{},
			},
			wantRoles:          []string{},
			wantStatus:         http.StatusInternalServerError,
			wantErr:            "not found",
			wantTearDownBefore: false,
		},
		{
			name: "returns all 4 roles because no project specified",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": "",
				}),
				Repo: newFakeRepo(noCrashers),
				user: models.User{},
			},
			wantRoles:          []string{"admin", "developer", "devops", "lead-dev"},
			wantStatus:         http.StatusOK,
			wantErr:            "",
			wantTearDownBefore: false,
		},
		{
			name: "returns error because no roles in db",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
				}),
				Repo: newFakeRepo(noCrashers),
				user: models.User{},
			},
			wantRoles:          []string{},
			wantStatus:         http.StatusNotFound,
			wantErr:            "not found",
			wantTearDownBefore: true,
		},
		{
			name: "fails to get roles on a free project",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
				}),
				Repo: newFakeRepo(map[string]error{
					"GetRole": errors.New("unexpected error"),
				}),
				user: models.User{},
			},
			wantRoles:          []string{},
			wantStatus:         http.StatusInternalServerError,
			wantErr:            "failed to get: unexpected error",
			wantTearDownBefore: false,
		},
		{
			name: "fails to get roles on a paid project",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": paidProject.UUID,
				}),
				Repo: newFakeRepo(map[string]error{
					"GetRoles": errors.New("unexpected error"),
				}),
				user: models.User{},
			},
			wantRoles:          []string{},
			wantStatus:         http.StatusInternalServerError,
			wantErr:            "failed to get: unexpected error",
			wantTearDownBefore: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantTearDownBefore == true {
				teardownRoles()
			}
			got, gotStatus, err := GetRoles(tt.args.params, tt.args.in1, tt.args.Repo, tt.args.user)

			gotResponse := got.(*models.GetRolesResponse)

			if err.Error() != tt.wantErr {
				t.Errorf("GetRoles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(gotResponse.Roles) != len(tt.wantRoles) {
				t.Errorf("GetRoles() len got = %v, want len %v", len(gotResponse.Roles), len(tt.wantRoles))
			}
			for _, wr := range tt.wantRoles {
				found := false
				for _, r := range gotResponse.Roles {
					if r.Name == wr {
						found = true
					}
				}
				if !found {
					t.Errorf("GetRoles() want %v, but not found", wr)
				}
			}
			if gotStatus != tt.wantStatus {
				t.Errorf("GetRoles() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
			if tt.wantTearDownBefore {
				seedRoles()
			}
		})
	}
}

func teardownRoles() {
	repo.NewRepo().GetDB().Transaction(func(db *gorm.DB) error {
		db.Delete(
			&models.Roles{},
			"1 = 1",
		)

		return db.Error
	})
}
func seedRoles() {
	repo.NewRepo().GetDB().Transaction(func(db *gorm.DB) error {
		seed.Seed(db)
		return db.Error
	})
}
