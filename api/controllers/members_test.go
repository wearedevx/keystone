// +build test

package controllers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
	"gorm.io/gorm"
)

func TestDoUsersExist(t *testing.T) {
	project, members := seedMembers(true /*paid*/)
	defer teardownMembers(project, members)

	type args struct {
		in0  router.Params
		body io.ReadCloser
		Repo repo.IRepo
		user models.User
	}

	tests := []struct {
		name       string
		args       args
		want       *models.CheckMembersResponse
		wantStatus int
		wantErr    string
	}{
		{
			name: "finds-a-user",
			args: args{
				in0: router.Params{},
				body: io.NopCloser(
					bytes.NewBufferString(
						fmt.Sprintf(
							`{"MemberIDs":["%s"]}`,
							members["developer"].UserID,
						),
					),
				),
				Repo: new(repo.Repo),
				user: members["admin"],
			},
			want: &models.CheckMembersResponse{
				Success: true,
				Error:   "",
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "does not find a user",
			args: args{
				in0: router.Params{},
				body: io.NopCloser(
					bytes.NewBufferString(
						`{"MemberIDs":["idontexist"]}`,
					),
				),
				Repo: new(repo.Repo),
				user: members["admin"],
			},
			want: &models.CheckMembersResponse{
				Success: false,
				Error:   "idontexist do not exists",
			},
			wantStatus: http.StatusNotFound,
			wantErr:    "",
		},
		{
			name: "does not find multiple users",
			args: args{
				in0: router.Params{},
				body: io.NopCloser(
					bytes.NewBufferString(
						`
                        {
                            "MemberIDs": [
                                "idontexist",
                                "idontexisteither"
                            ]
                        }`,
					),
				),
				Repo: new(repo.Repo),
				user: members["admin"],
			},
			want: &models.CheckMembersResponse{
				Success: false,
				Error:   "idontexist, idontexisteither do not exists",
			},
			wantStatus: http.StatusNotFound,
			wantErr:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotStatus, err := DoUsersExist(
				tt.args.in0,
				tt.args.body,
				tt.args.Repo,
				tt.args.user,
			)
			if err.Error() != tt.wantErr {
				t.Errorf(
					"DoUsersExist() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoUsersExist() got = %v, want %v", got, tt.want)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf(
					"DoUsersExist() gotStatus = %v, want %v",
					gotStatus,
					tt.wantStatus,
				)
			}
		})
	}
}

func TestPutMembersSetRole(t *testing.T) {
	project, members := seedMembers(true /* paid */)
	unpaidProject, unpaidMembers := seedMembers(false /* unpaid */)
	defer teardownMembers(unpaidProject, unpaidMembers)
	defer teardownMembers(project, members)

	type args struct {
		params router.Params
		body   io.ReadCloser
		Repo   repo.IRepo
		user   models.User
	}
	tests := []struct {
		name         string
		args         args
		wantResponse router.Serde
		wantStatus   int
		wantErr      string
	}{
		{
			name: "admin-sets-role-on-dev",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
				}),
				body: io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`
                {
                    "MemberID": "%s",
                    "RoleName": "lead-dev"
                }
                `, members["developer"].UserID))),
				Repo: new(repo.Repo),
				user: members["admin"],
			},
			wantResponse: nil,
			wantStatus:   http.StatusOK,
			wantErr:      "",
		},
		{
			name: "dev-cannot-set-role-of-admin",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
				}),
				body: io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`
                {
                    "MemberID": "%s",
                    "RoleName": "lead-dev"
                }
                `, members["admin"].UserID))),
				Repo: new(repo.Repo),
				user: members["developer"],
			},
			wantResponse: nil,
			wantStatus:   http.StatusForbidden,
			wantErr:      "permission denied",
		},
		{
			name: "returns not found for project",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": "not-a-project",
				}),
				body: io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`
                {
                    "MemberID": "%s",
                    "RoleName": "lead-dev"
                }
                `, members["admin"].UserID))),
				Repo: new(repo.Repo),
				user: members["admin"],
			},
			wantResponse: nil,
			wantStatus:   http.StatusNotFound,
			wantErr:      "not found",
		},
		{
			name: "returns not found for member",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
				}),
				body: io.NopCloser(bytes.NewBufferString(`
                {
                    "MemberID": "not-a-member",
                    "RoleName": "lead-dev"
                }
                `)),
				Repo: new(repo.Repo),
				user: members["developer"],
			},
			wantResponse: nil,
			wantStatus:   http.StatusNotFound,
			wantErr:      "not found",
		},
		{
			name: "needs a paid organization",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": unpaidProject.UUID,
				}),
				body: io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`
                {
                    "MemberID": "%s",
                    "RoleName": "lead-dev"
                }
                `, unpaidMembers["developer"].UserID))),
				Repo: new(repo.Repo),
				user: members["admin"],
			},
			wantResponse: nil,
			wantStatus:   http.StatusForbidden,
			wantErr:      "needs upgrade",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, gotStatus, err := PutMembersSetRole(
				tt.args.params,
				tt.args.body,
				tt.args.Repo,
				tt.args.user,
			)
			if err.Error() != tt.wantErr {
				t.Errorf(
					"PutMembersSetRole() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf(
					"PutMembersSetRole() gotResponse = %v, want %v",
					gotResponse,
					tt.wantResponse,
				)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf(
					"PutMembersSetRole() gotStatus = %v, want %v",
					gotStatus,
					tt.wantStatus,
				)
			}
		})
	}
}

func seedMembers(
	paid bool,
) (project models.Project, users map[string]models.User) {
	new(repo.Repo).GetDb().Transaction(func(db *gorm.DB) error {
		roles := make([]models.Role, 5)

		db.Model(&models.Role{}).Find(&roles)

		roleWithName := func(n string) *models.Role {
			for _, role := range roles {
				if n == role.Name {
					return &role
				}
			}
			return nil
		}

		users = map[string]models.User{
			"admin":     {},
			"devops":    {},
			"lead-dev":  {},
			"developer": {},
		}

		organization := models.Organization{}
		faker.FakeData(&organization)
		organization.Paid = paid
		db.Save(&organization)

		faker.FakeData(&project)
		project.OrganizationID = organization.ID
		project.Organization = organization
		db.Save(&project)

		for roleName, user := range users {
			role := roleWithName(roleName)
			faker.FakeData(&user)
			db.Save(&user)

			db.Save(&models.ProjectMember{
				ProjectID: project.ID,
				UserID:    user.ID,
				RoleID:    role.ID,
			})

			users[roleName] = user
		}

		return db.Error
	})

	return project, users
}

func teardownMembers(project models.Project, users map[string]models.User) {
	new(repo.Repo).GetDb().Transaction(func(db *gorm.DB) error {
		db.Begin().
			Exec("delete from project_members where project_id = ?", project.ID).
			Exec(`
            delete from projects where id = ?;
            delete from organizations where id = ?;
            `, project.ID, project.OrganizationID)

		for _, user := range users {
			db.Exec("delete from users where id = ?", user.ID)
		}

		return db.Error
	})
}
