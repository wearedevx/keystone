package controllers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func TestGetOrganizations(t *testing.T) {
	mainUser, org, project := seedOneProjectForOneUser()
	usera, orga, projecta := seedOneProjectForOneUser()
	userb, orgb, projectb := seedOneProjectForOneUser()
	userc, orgc, projectc := seedOneProjectForOneUser()

	roles := testsGetRoles()

	mainMember := seedProjectMember(project, mainUser, roles["admin"])
	membera := seedProjectMember(projecta, mainUser, roles["devops"])
	memberb := seedProjectMember(projectb, mainUser, roles["developer"])
	memberc := seedProjectMember(projectc, mainUser, roles["developer"])

	defer func() {
		teardownProjectMember(memberc)
		teardownProjectMember(memberb)
		teardownProjectMember(membera)
		teardownProjectMember(mainMember)

		teardownProject(projectc)
		teardownProject(projectb)
		teardownProject(projecta)
		teardownProject(project)

		teardownUserAndOrganization(userc, orgc)
		teardownUserAndOrganization(userb, orgb)
		teardownUserAndOrganization(usera, orga)
		teardownUserAndOrganization(mainUser, org)
	}()

	type args struct {
		params router.Params
		in1    io.ReadCloser
		Repo   repo.IRepo
		user   models.User
	}
	tests := []struct {
		name       string
		args       args
		want       *models.GetOrganizationsResponse
		wantStatus int
		wantErr    string
	}{
		{
			name: "lists all organisations the user is member of",
			args: args{
				params: router.Params{},
				in1:    nil,
				Repo:   repo.NewRepo(),
				user:   mainUser,
			},
			want: &models.GetOrganizationsResponse{
				Organizations: []models.Organization{
					org,
					orga,
					orgb,
					orgc,
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "lists all organisations the user owns",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"owned": "1",
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				user: mainUser,
			},
			want: &models.GetOrganizationsResponse{
				Organizations: []models.Organization{
					org,
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "find an organization by name among those the user is a member of",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"name": orga.Name,
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				user: mainUser,
			},
			want: &models.GetOrganizationsResponse{
				Organizations: []models.Organization{
					orga,
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "find an organization by name among those the user owns",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"name":  org.Name,
					"owned": "1",
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				user: mainUser,
			},
			want: &models.GetOrganizationsResponse{
				Organizations: []models.Organization{
					org,
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "does not an organization by name among those the user owns",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"name":  orga.Name,
					"owned": "1",
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				user: mainUser,
			},
			want: &models.GetOrganizationsResponse{
				Organizations: []models.Organization{},
			},
			wantStatus: http.StatusNotFound,
			wantErr:    "not found",
		},
		{
			name: "does not an organization by name among those the user is member of",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"name": "wait... that's no organization name!",
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				user: mainUser,
			},
			want: &models.GetOrganizationsResponse{
				Organizations: []models.Organization{},
			},
			wantStatus: http.StatusNotFound,
			wantErr:    "not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotStatus, err := GetOrganizations(
				tt.args.params,
				tt.args.in1,
				tt.args.Repo,
				tt.args.user,
			)
			if err.Error() != tt.wantErr {
				t.Errorf(
					"GetOrganizations() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if gotStatus != tt.wantStatus {
				t.Errorf(
					"GetOrganizations() gotStatus = %v, want %v",
					gotStatus,
					tt.wantStatus,
				)
				return
			}

			if tt.want != nil && got != nil {
				gotResponse := got.(*models.GetOrganizationsResponse)

				if len(
					gotResponse.Organizations,
				) != len(
					tt.want.Organizations,
				) {
					t.Errorf(
						"GetOrganizations() got = %v, want %v",
						got,
						tt.want,
					)
					return
				}

				for _, wantOrganization := range tt.want.Organizations {
					found := false

					for _, gotOrganization := range gotResponse.Organizations {
						if wantOrganization.ID == gotOrganization.ID {
							found = true
							break
						}
					}

					if !found {
						t.Errorf(
							"GetOrganizations() got = %v, want %v",
							got,
							tt.want,
						)
						return
					}
				}
			}
		})
	}
}

func TestPostOrganization(t *testing.T) {
	user, org := seedSingleUser()
	defer teardownUserAndOrganization(user, org)

	type args struct {
		in0  router.Params
		body io.ReadCloser
		Repo repo.IRepo
		user models.User
	}
	tests := []struct {
		name       string
		args       args
		want       *models.Organization
		wantStatus int
		wantErr    string
	}{
		{
			name: "creates an organization",
			args: args{
				in0: router.Params{},
				body: io.NopCloser(bytes.NewBufferString(`{
  "name": "organization-name"
}`)),
				Repo: repo.NewRepo(),
				user: user,
			},
			want: &models.Organization{
				Name: "organization-name",
			},
			wantStatus: http.StatusCreated,
			wantErr:    "",
		},
		{
			name: "the paid property is always false on creation",
			args: args{
				in0: router.Params{},
				body: io.NopCloser(bytes.NewBufferString(`{
  "name": "hacky-organization-name",
	"paid": true
}`)),
				Repo: repo.NewRepo(),
				user: user,
			},
			want: &models.Organization{
				Name: "hacky-organization-name",
				Paid: false,
			},
			wantStatus: http.StatusCreated,
			wantErr:    "",
		},
		{
			name: "no spaces in organization name",
			args: args{
				in0: router.Params{},
				body: io.NopCloser(bytes.NewBufferString(`{
  "name": "organization name"
}`)),
				Repo: repo.NewRepo(),
				user: user,
			},
			want:       nil,
			wantStatus: http.StatusBadRequest,
			wantErr:    "bad organization name",
		},
		{
			name: "cant create orga with same name",
			args: args{
				in0: router.Params{},
				body: io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{
  "name": "%s"
}`, org.Name))),
				Repo: repo.NewRepo(),
				user: user,
			},
			want:       nil,
			wantStatus: http.StatusConflict,
			wantErr:    "organization name already taken",
		},
		{
			name: "input must be valid json",
			args: args{
				in0: router.Params{},
				body: io.NopCloser(bytes.NewBufferString(`{
  "name": "a-valid-name",
}`)),
				Repo: repo.NewRepo(),
				user: user,
			},
			want:       nil,
			wantStatus: http.StatusBadRequest,
			wantErr:    "bad request: invalid character '}' looking for beginning of object key string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotStatus, err := PostOrganization(
				tt.args.in0,
				tt.args.body,
				tt.args.Repo,
				tt.args.user,
			)
			if err.Error() != tt.wantErr {
				t.Errorf(
					"PostOrganization() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if gotStatus != tt.wantStatus {
				t.Errorf(
					"PostOrganization() gotStatus = %v, want %v",
					gotStatus,
					tt.wantStatus,
				)
				return
			}

			if tt.want != nil && got != nil {
				gotOrganization := got.(*models.Organization)

				if gotOrganization.Name != tt.want.Name {
					t.Errorf(
						"PostOrganization() got.Name = %v, want %v",
						gotOrganization.Name,
						tt.want.Name,
					)
					return
				}

				if gotOrganization.Paid != tt.want.Paid {
					t.Errorf(
						"PostOrganization() got.Paid = %v, want %v",
						gotOrganization.Paid,
						tt.want.Paid,
					)
					return
				}
			}

		})
	}
}

func TestUpdateOrganization(t *testing.T) {
	user, orga := seedSingleUser()
	otherUser, otherOrga := seedSingleUser()

	defer teardownUserAndOrganization(otherUser, otherOrga)
	defer teardownUserAndOrganization(user, orga)

	type args struct {
		in0  router.Params
		body io.ReadCloser
		Repo repo.IRepo
		user models.User
	}
	tests := []struct {
		name       string
		args       args
		want       *models.Organization
		wantStatus int
		wantErr    string
	}{
		{
			name: "updates an organization name",
			args: args{
				in0: router.Params{},
				body: io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{
	"id": %d,
  "name": "anew-name"
}`, orga.ID))),
				Repo: repo.NewRepo(),
				user: user,
			},
			want: &models.Organization{
				Name:    "anew-name",
				Private: false,
				Paid:    false,
				UserID:  user.ID,
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "updates an organization privacy",
			args: args{
				in0: router.Params{},
				body: io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{
	"id": %d,
	"private": true
}`, orga.ID))),
				Repo: repo.NewRepo(),
				user: user,
			},
			want: &models.Organization{
				Name:    "anew-name",
				Private: true,
				Paid:    false,
				UserID:  user.ID,
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "do not update the paid flag",
			args: args{
				in0: router.Params{},
				body: io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{
	"id": %d,
	"paid": true,
	"private": true
}`, orga.ID))),
				Repo: repo.NewRepo(),
				user: user,
			},
			want: &models.Organization{
				Name:    "anew-name",
				Paid:    false,
				Private: true,
				UserID:  user.ID,
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "do not update the user_id",
			args: args{
				in0: router.Params{},
				body: io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{
	"id": %d,
	"user_id": %d,
	"private": true
}`, orga.ID, otherUser.ID))),
				Repo: repo.NewRepo(),
				user: user,
			},
			want: &models.Organization{
				Name:    "anew-name",
				Paid:    false,
				Private: true,
				UserID:  user.ID,
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "not if the name is taken",
			args: args{
				in0: router.Params{},
				body: io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{
	"id": %d,
  "name": "%s"
}`, orga.ID, otherOrga.Name))),
				Repo: repo.NewRepo(),
				user: user,
			},
			want:       nil,
			wantStatus: http.StatusConflict,
			wantErr:    "organization name already taken",
		},
		{
			name: "not if the name is bad",
			args: args{
				in0: router.Params{},
				body: io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{
	"id": %d,
  "name": "a name with spaces is a bad name"
}`, orga.ID))),
				Repo: repo.NewRepo(),
				user: user,
			},
			want:       nil,
			wantStatus: http.StatusBadRequest,
			wantErr:    "bad organization name",
		},
		{
			name: "must own the orga",
			args: args{
				in0: router.Params{},
				body: io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{
					"id": %d,
					"name": "a-name-that-is-a-valid-one"
				}`, orga.ID))),
				Repo: repo.NewRepo(),
				user: otherUser,
			},
			want:       nil,
			wantStatus: http.StatusForbidden,
			wantErr:    "not organization owner",
		},
		{
			name: "body must be valid json",
			args: args{
				in0: router.Params{},
				body: io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{
					"id": %d,
					"name": "a-name-that-is-a-valid-one",
				}`, orga.ID))),
				Repo: repo.NewRepo(),
				user: user,
			},
			want:       nil,
			wantStatus: http.StatusBadRequest,
			wantErr:    "bad request: invalid character '}' looking for beginning of object key string",
		},
		{
			name: "organization must exists",
			args: args{
				in0: router.Params{},
				body: io.NopCloser(bytes.NewBufferString(`{
							"id": 123049,
							"name": "a-name-that-is-a-valid-one"
						}`)),
				Repo: repo.NewRepo(),
				user: user,
			},
			want:       nil,
			wantStatus: http.StatusNotFound,
			wantErr:    "not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotStatus, err := UpdateOrganization(
				tt.args.in0,
				tt.args.body,
				tt.args.Repo,
				tt.args.user,
			)
			if err.Error() != tt.wantErr {
				t.Errorf(
					"UpdateOrganization() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if gotStatus != tt.wantStatus {
				t.Errorf(
					"UpdateOrganization() gotStatus = %v, want %v",
					gotStatus,
					tt.wantStatus,
				)
				return
			}

			if tt.want != nil {
				gotOrganization := got.(*models.Organization)

				if gotOrganization.Name != tt.want.Name {
					t.Errorf(
						"UpdateOrganization() got.Name = %v, want %v",
						gotOrganization.Name,
						tt.want.Name,
					)
				}

				if gotOrganization.Private != tt.want.Private {
					t.Errorf(
						"UpdateOrganization() got.Private = %v, want %v",
						gotOrganization.Private,
						tt.want.Private,
					)
				}

				if gotOrganization.Paid != tt.want.Paid {
					t.Errorf(
						"UpdateOrganization() got.Paid = %v, want %v",
						gotOrganization.Paid,
						tt.want.Paid,
					)
				}

				if gotOrganization.UserID != tt.want.UserID {
					t.Errorf(
						"UpdateOrganization() got.UserID = %v, want %v",
						gotOrganization.UserID,
						tt.want.UserID,
					)
				}
			} else if tt.want != got {
				t.Errorf(
					"UpdateOrganization() got = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestGetOrganizationProjects(t *testing.T) {
	user, org, projects := seedManyProjectsForOneUser()
	defer teardownUserAndOrganization(user, org)
	defer teardownManyProjects(projects)

	type args struct {
		params router.Params
		in1    io.ReadCloser
		Repo   repo.IRepo
		user   models.User
	}
	tests := []struct {
		name       string
		args       args
		want       *models.GetProjectsResponse
		wantStatus int
		wantErr    string
	}{
		{
			name: "gets projects related to an organization",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"orgaID": strconv.FormatUint(uint64(org.ID), 10),
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				user: user,
			},
			want: &models.GetProjectsResponse{
				Projects: projects,
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "bad request if bad id",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"orgaID": "that's not a number",
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				user: user,
			},
			want:       &models.GetProjectsResponse{},
			wantStatus: http.StatusBadRequest,
			wantErr:    `bad request: strconv.ParseUint: parsing "that's not a number": invalid syntax`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotStatus, err := GetOrganizationProjects(
				tt.args.params,
				tt.args.in1,
				tt.args.Repo,
				tt.args.user,
			)
			if err.Error() != tt.wantErr {
				t.Errorf(
					"GetOrganizationProjects() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if gotStatus != tt.wantStatus {
				t.Errorf(
					"GetOrganizationProjects() gotStatus = %v, want %v",
					gotStatus,
					tt.wantStatus,
				)
				return
			}

			if tt.want != nil {
				gotResponse := got.(*models.GetProjectsResponse)

				if len(gotResponse.Projects) != len(tt.want.Projects) {
					t.Errorf(
						"GetOrganizationProjects() got = %v, want %v",
						got,
						tt.want,
					)
					return
				}

				for _, wantProject := range tt.want.Projects {
					found := false

					for _, gotProject := range gotResponse.Projects {
						if gotProject.ID == wantProject.ID {
							found = true
							break
						}
					}

					if !found {
						t.Errorf(
							"GetOrganizationProjects() got = %v, want %v",
							got,
							tt.want,
						)
						return
					}
				}
			}

		})
	}
}

func TestGetOrganizationMembers(t *testing.T) {
	user, org, project := seedOneProjectForOneUser()
	defer teardownUserAndOrganization(user, org)
	defer teardownProject(project)

	usera, orga := seedSingleUser()
	userb, orgb := seedSingleUser()
	userc, orgc := seedSingleUser()
	userd, orgd := seedSingleUser()

	roles := testsGetRoles()

	member := seedProjectMember(project, user, roles["admin"])
	membera := seedProjectMember(project, usera, roles["developer"])
	memberb := seedProjectMember(project, userb, roles["developer"])
	memberc := seedProjectMember(project, userc, roles["developer"])
	memberd := seedProjectMember(project, userd, roles["developer"])

	defer func() {
		teardownUserAndOrganization(usera, orga)
		teardownUserAndOrganization(userb, orgb)
		teardownUserAndOrganization(userc, orgc)
		teardownUserAndOrganization(userd, orgd)

		teardownProjectMember(member)
		teardownProjectMember(membera)
		teardownProjectMember(memberb)
		teardownProjectMember(memberc)
		teardownProjectMember(memberd)
	}()

	type args struct {
		params router.Params
		in1    io.ReadCloser
		Repo   repo.IRepo
		user   models.User
	}
	tests := []struct {
		name       string
		args       args
		want       *models.GetMembersResponse
		wantStatus int
		wantErr    string
	}{
		{
			name: "gets all members of an organization",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"orgaID": strconv.FormatUint(uint64(org.ID), 10),
				}),
				in1:  nil,
				Repo: repo.NewRepo(),
				user: user,
			},
			want: &models.GetMembersResponse{
				Members: []models.ProjectMember{
					member,
					membera,
					memberb,
					memberc,
					memberd,
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotStatus, err := GetOrganizationMembers(
				tt.args.params,
				tt.args.in1,
				tt.args.Repo,
				tt.args.user,
			)
			if err.Error() != tt.wantErr {
				t.Errorf(
					"GetOrganizationMembers() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if gotStatus != tt.wantStatus {
				t.Errorf(
					"GetOrganizationMembers() gotStatus = %v, want %v",
					gotStatus,
					tt.wantStatus,
				)
				return
			}

			if tt.want != nil {
				gotResponse := got.(*models.GetMembersResponse)

				if len(gotResponse.Members) != len(tt.want.Members) {
					t.Errorf(
						"GetOrganizationMembers() got = %v, want %v",
						got,
						tt.want,
					)
				}

			}
		})
	}
}
