package controllers

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func TestGetEnvironmentPublicKeys(t *testing.T) {
	Repo := new(repo.Repo)
	adminUser, devUser, environments := seedEnvironmentPublicKeys(Repo)
	// defer teardownEnvironment(adminUser, devUser, environments)

	type args struct {
		params router.Params
		in1    io.ReadCloser
		Repo   repo.IRepo
		user   models.User
	}
	tests := []struct {
		name       string
		args       args
		want       models.PublicKeys
		wantStatus int
		wantErr    string
	}{
		{
			name: "returns public keys",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"envID": environments["dev"].EnvironmentID,
				}),
				in1:  nil,
				Repo: newFakeRepo(noCrashers),
				user: adminUser,
			},
			want: models.PublicKeys{
				Keys: []models.UserDevices{
					{
						UserID:  adminUser.ID,
						UserUID: adminUser.UserID,
						Devices: adminUser.Devices,
					},
					{
						UserID:  devUser.ID,
						UserUID: devUser.UserID,
						Devices: devUser.Devices,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "returns not found",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"envID": "that environment is not one",
				}),
				in1:  nil,
				Repo: newFakeRepo(noCrashers),
				user: adminUser,
			},
			want: models.PublicKeys{
				Keys: []models.UserDevices{},
			},
			wantStatus: http.StatusNotFound,
			wantErr:    "not found",
		},
		{
			name: "fails getting the environment",
			args: args{
				params: router.Params{},
				in1:    nil,
				Repo: newFakeRepo(map[string]error{
					"GetEnvironment": errors.New("unexpected error"),
				}),
				user: adminUser,
			},
			want: models.PublicKeys{
				Keys: []models.UserDevices{},
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    "failed to get: unexpected error",
		},
		{
			name: "fails to get the environment public keys",
			args: args{
				params: router.Params{},
				in1:    nil,
				Repo: newFakeRepo(map[string]error{
					"GetEnvironmentPublicKeys": errors.New("unexpected error"),
				}),
				user: adminUser,
			},
			want: models.PublicKeys{
				Keys: []models.UserDevices{},
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    "failed to get: unexpected error",
		},
		{
			name: "fails to check the rights",
			args: args{
				params: router.Params{},
				in1:    nil,
				Repo: newFakeRepo(map[string]error{
					"GetProjectMember": errors.New("unexpected error"),
				}),
				user: adminUser,
			},
			want: models.PublicKeys{
				Keys: []models.UserDevices{},
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    "unexpected error",
		},
		{
			name: "dev user cannot see the prod keys",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"envID": environments["prod"].EnvironmentID,
				}),
				in1:  nil,
				Repo: newFakeRepo(noCrashers),
				user: devUser,
			},
			want: models.PublicKeys{
				Keys: []models.UserDevices{},
			},
			wantStatus: http.StatusForbidden,
			wantErr:    "permission denied",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotStatus, err := GetEnvironmentPublicKeys(
				tt.args.params,
				tt.args.in1,
				tt.args.Repo,
				tt.args.user,
			)
			if err.Error() != tt.wantErr {
				t.Errorf(
					"GetEnvironmentPublicKeys() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}

			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("GetEnvironmentPublicKeys() got = %v, want %v", got, tt.want)
			// }

			if gotStatus != tt.wantStatus {
				t.Errorf(
					"GetEnvironmentPublicKeys() gotStatus = %v, want %v",
					gotStatus,
					tt.wantStatus,
				)
				return
			}

			//
			gotResponse := got.(*models.PublicKeys)
			gotKeys := gotResponse.Keys
			wantKeys := tt.want.Keys

			gotLen := len(gotKeys)
			wantLen := len(wantKeys)
			if gotLen != wantLen {
				t.Errorf(
					"GetEnvironmentPublicKeys() got len = %d, want %d",
					gotLen,
					wantLen,
				)
				return
			}

			for index, gotKey := range gotKeys {
				wantKey := tt.want.Keys[index]

				if gotKey.UserID != wantKey.UserID ||
					gotKey.UserUID != wantKey.UserUID {
					t.Errorf(
						"GetEnvironmentPublicKeys() key = %v, want %v",
						gotKey,
						wantKey,
					)
				}

			}
		})

	}

}

func seedEnvironmentPublicKeys(
	Repo *repo.Repo,
) (models.User, models.User, map[string]models.Environment) {
	var project models.Project
	var organization models.Organization
	var adminUser models.User
	var devUser models.User

	faker.FakeData(&adminUser)
	faker.FakeData(&devUser)

	adminUser.Devices = []models.Device{{
		PublicKey: []byte{},
		Name:      faker.Word(),
		UID:       faker.UUIDHyphenated(),
	}}
	devUser.Devices = []models.Device{{
		PublicKey: []byte{},
		Name:      faker.Word(),
		UID:       faker.UUIDHyphenated(),
	}}
	Repo.GetOrCreateUser(&adminUser)
	Repo.GetOrCreateUser(&devUser)

	Repo.GetDb().
		Model(&adminUser).
		Association("Devices").
		Find(&adminUser.Devices)
	Repo.GetDb().Model(&devUser).Association("Devices").Find(&devUser.Devices)

	organization = models.Organization{
		UserID:  adminUser.ID,
		Name:    adminUser.UserID,
		Private: true,
	}

	Repo.GetOrganization(&organization)
	// Logs require that organization to be paid
	Repo.OrganizationSetPaid(&organization, false)

	project = models.Project{
		UUID:                faker.UUIDHyphenated(),
		TTL:                 30,
		DaysBeforeTTLExpiry: 7,
		Name:                "activity-logs-project",
		UserID:              adminUser.ID,
		Environments:        []models.Environment{},
		OrganizationID:      organization.ID,
		Organization:        organization,
	}
	Repo.GetOrCreateProject(&project)
	dev := project.GetEnvironment("dev")
	staging := project.GetEnvironment("staging")
	prod := project.GetEnvironment("prod")

	Repo.ProjectAddMembers(project, []models.MemberRole{
		{
			MemberID: devUser.UserID,
			RoleID:   4,
		},
	}, adminUser)

	if Repo.Err() != nil {
		panic(Repo.Err())
	}

	return adminUser, devUser, map[string]models.Environment{
		"dev":     *dev,
		"staging": *staging,
		"prod":    *prod,
	}
}

func teardownEnvironment(
	adminUser, devUser models.User,
	environments map[string]models.Environment,
) {
	db := new(repo.Repo).GetDb()
	projectID := environments["dev"].ProjectID

	db.Exec("delete from environments where id = ?", environments["dev"].ID)
	db.Exec("delete from environments where id = ?", environments["staging"].ID)
	db.Exec("delete from environments where id = ?", environments["prod"].ID)

	db.Exec("delete from project_members where project_id = ?", projectID)

	db.Exec("delete from projects where id = ?", projectID)
	db.Exec(
		"delete from users where id = ? or id = ?",
		adminUser.ID,
		devUser.ID,
	)
}
