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

func TestPostProject(t *testing.T) {
	user, organization := seedSingleUser()
	defer teardownUserAndOrganization(user, organization)

	projectName := faker.Sentence()

	type args struct {
		in0  router.Params
		body io.ReadCloser
		Repo repo.IRepo
		user models.User
	}
	tests := []struct {
		name       string
		args       args
		want       *models.Project
		wantStatus int
		wantErr    string
	}{
		{
			name: "creates a project",
			args: args{
				in0: router.Params{},
				body: io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`
                {
                    "name": "%s",
                    "organization_id": %d
                }
                `, projectName, organization.ID))),
				Repo: repo.NewRepo(),
				user: user,
			},
			want: &models.Project{
				TTL:                 7,
				DaysBeforeTTLExpiry: 2,
				Name:                projectName,
				Members:             []models.ProjectMember{},
				UserID:              user.ID,
				User:                user,
				Environments:        []models.Environment{},
				OrganizationID:      organization.ID,
				Organization:        organization,
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "does not create a project without an organization",
			args: args{
				in0: router.Params{},
				body: io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`
                {
                    "name": "%s"
                }
                `, projectName))),
				Repo: repo.NewRepo(),
				user: user,
			},
			want:       nil,
			wantStatus: http.StatusNotFound,
			wantErr:    "not found",
		},
		// TODO: test that organization owners are added as admins to projects created by other organization members
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotStatus, err := PostProject(
				tt.args.in0,
				tt.args.body,
				tt.args.Repo,
				tt.args.user,
			)
			if err.Error() != tt.wantErr {
				t.Errorf(
					"PostProject() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if gotStatus != tt.wantStatus {
				t.Errorf(
					"PostProject() gotStatus = %v, want %v",
					gotStatus,
					tt.wantStatus,
				)
				return
			}

			gotProject := got.(*models.Project)

			if tt.want == nil {
				if gotProject != nil {
					t.Errorf(
						"PostProject() got %v, want %v",
						gotProject,
						tt.want,
					)
				}
				return
			}

			if len(gotProject.UUID) == 0 {
				t.Errorf(
					"PostProject() got UUID %v, want valid uuid",
					gotProject.UUID,
				)
			}

			if gotProject.Name != tt.want.Name {
				t.Errorf(
					"PostProject() got Name %v, want %v",
					gotProject.UUID,
					tt.want.Name,
				)
			}

			if gotProject.UserID != tt.want.UserID {
				t.Errorf(
					"PostProject() got UserID %v, want %v",
					gotProject.UserID,
					tt.want.UserID,
				)
			}

			if gotProject.OrganizationID != tt.want.OrganizationID {
				t.Errorf(
					"PostProject() got OrganizationID %v, want %v",
					gotProject.OrganizationID,
					tt.want.OrganizationID,
				)
			}

			if len(gotProject.Environments) != 3 {
				t.Errorf(
					"PostProject() got Environments %v, want 3 of them",
					gotProject.Environments,
				)
			}
		})
	}
}

func TestGetProjects(t *testing.T) {
	user, organization, projects := seedManyProjectsForOneUser()
	defer teardownUserAndOrganization(user, organization)
	defer teardownManyProjects(projects)

	type args struct {
		in0  router.Params
		in1  io.ReadCloser
		Repo repo.IRepo
		user models.User
	}
	tests := []struct {
		name       string
		args       args
		want       *models.GetProjectsResponse
		wantStatus int
		wantErr    string
	}{
		{
			name: "gets projects for the user",
			args: args{
				in0:  router.Params{},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotStatus, err := GetProjects(
				tt.args.in0,
				tt.args.in1,
				tt.args.Repo,
				tt.args.user,
			)
			if err.Error() != tt.wantErr {
				t.Errorf(
					"GetProjects() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if gotStatus != tt.wantStatus {
				t.Errorf(
					"GetProjects() gotStatus = %v, want %v",
					gotStatus,
					tt.wantStatus,
				)
				return
			}

			gotResponse := got.(*models.GetProjectsResponse)

			if tt.want == nil {
				if gotResponse != nil {
					t.Errorf("GetProjects() got = %v, want nil", gotResponse)
				}
				return
			}

			gotIds := make([]uint, len(gotResponse.Projects))
			wantIds := make([]uint, len(tt.want.Projects))

			for i, p := range gotResponse.Projects {
				gotIds[i] = p.UserID
			}

			for i, p := range tt.want.Projects {
				wantIds[i] = p.UserID
			}

			gotProjects := gotResponse.Projects
			if len(gotProjects) != len(tt.want.Projects) {
				t.Errorf(
					"GetProjects() got len = %v, want %v",
					len(gotProjects),
					len(tt.want.Projects),
				)
				return
			}

			for index, gotProject := range gotProjects {
				wantProject := tt.want.Projects[index]
				if gotProject.UUID != wantProject.UUID {
					t.Errorf("GetProjects() got = %v, want %v",
						gotProject,
						wantProject,
					)
				}
			}
		})
	}
}

func TestGetProjectsMembers(t *testing.T) {
	Repo := repo.NewRepo()
	users := make([]models.User, 5)

	for i := 1; i < 5; i++ {
		user, organization := seedSingleUser()
		users[i] = user
		defer teardownUserAndOrganization(user, organization)
	}

	user, org, project := seedOneProjectForOneUser()
	defer teardownProject(project)
	defer teardownUserAndOrganization(user, org)
	users[0] = user

	roles := testsGetRoles()

	projectMembers := seedProjectMembers(project, users, roles["developer"])
	defer teardownProjectMembers(projectMembers)

	externalUser, externalOrganization := seedSingleUser()
	defer teardownUserAndOrganization(externalUser, externalOrganization)

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
			name: "gets project members",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
				}),
				in1:  nil,
				Repo: Repo,
				user: user,
			},
			want: &models.GetMembersResponse{
				Members: projectMembers,
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "not found if project does not exist",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": "not a valid project id",
				}),
				in1:  nil,
				Repo: Repo,
				user: user,
			},
			want:       &models.GetMembersResponse{},
			wantStatus: http.StatusNotFound,
			wantErr:    "not found",
		},
		{
			name: "not found if user is not a member of the project",
			args: args{
				params: router.ParamsFrom(map[string]string{"projectID": project.UUID}),
				in1:    nil,
				Repo:   Repo,
				user:   externalUser,
			},
			want:       &models.GetMembersResponse{},
			wantStatus: http.StatusNotFound,
			wantErr:    "not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotStatus, err := GetProjectsMembers(
				tt.args.params,
				tt.args.in1,
				tt.args.Repo,
				tt.args.user,
			)
			if err.Error() != tt.wantErr {
				t.Errorf(
					"GetProjectsMembers() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if gotStatus != tt.wantStatus {
				t.Errorf(
					"GetProjectsMembers() gotStatus = %v, want %v",
					gotStatus,
					tt.wantStatus,
				)
				return
			}

			gotResponse := got.(*models.GetMembersResponse)

			if len(gotResponse.Members) != len(tt.want.Members) {
				t.Errorf(
					"GetProjectMembers() got Members = %v, want %v",
					gotResponse.Members,
					tt.want.Members,
				)
				return
			}

			for _, wantMember := range tt.want.Members {
				found := false

				for _, gotMember := range gotResponse.Members {
					if gotMember.ID == wantMember.ID {
						found = true
						break
					}
				}

				if !found {
					t.Errorf(
						"GetProjectMembers() member not found = %v in %v",
						wantMember,
						gotResponse.Members,
					)
					return
				}
			}
		})
	}
}

func TestPostProjectsMembers(t *testing.T) {
	userToAdd, organizationToAdd := seedSingleUser()
	defer teardownUserAndOrganization(userToAdd, organizationToAdd)

	adminToAdd, arganizationToAdd := seedSingleUser()
	defer teardownUserAndOrganization(adminToAdd, arganizationToAdd)

	addedUser, addedOrganization := seedSingleUser()
	defer teardownUserAndOrganization(addedUser, addedOrganization)

	paidUser, paidOrganization, paidProject := seedOneProjectForOneUser()
	testsSetOrganisationPaid(&paidOrganization)
	defer teardownUserAndOrganization(paidUser, paidOrganization)
	defer teardownProject(paidProject)

	freeUser, freeOrganization, freeProject := seedOneProjectForOneUser()
	defer teardownUserAndOrganization(freeUser, freeOrganization)
	defer teardownProject(freeProject)

	roles := testsGetRoles()

	paidProjectMember := seedProjectMember(paidProject, paidUser, roles["admin"])
	defer teardownProjectMember(paidProjectMember)

	existingMember := seedProjectMember(paidProject, addedUser, roles["developer"])
	defer teardownProjectMember(existingMember)

	freeProjectMember := seedProjectMember(freeProject, freeUser, roles["admin"])
	defer teardownProjectMember(freeProjectMember)

	type args struct {
		params router.Params
		body   io.ReadCloser
		Repo   repo.IRepo
		user   models.User
	}
	tests := []struct {
		name       string
		args       args
		want       *models.AddMembersResponse
		wantStatus int
		wantErr    string
	}{
		{
			name: "adds a project member",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": paidProject.UUID,
				}),
				body: io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{
  "Members": [
		{
			"MemberID": "%s",
			"RoleID": %d
		}
  ]
}`, userToAdd.UserID, roles["developer"].ID))),
				Repo: repo.NewRepo(),
				user: paidUser,
			},
			want: &models.AddMembersResponse{
				Success: true,
				Error:   "",
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "bad request if no project id",
			args: args{
				params: router.ParamsFrom(map[string]string{}),
				body: io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{
  "Members": [
		{
			"MemberID": "%s",
			"RoleID": %d
		}
  ]
}`, userToAdd.UserID, roles["developer"].ID))),
				Repo: repo.NewRepo(),
				user: paidUser,
			},
			want:       &models.AddMembersResponse{},
			wantStatus: http.StatusBadRequest,
			wantErr:    "bad request",
		},
		{
			name: "not found if bad project_id",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": "not a real project id",
				}),
				body: io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{
  "Members": [
		{
			"MemberID": "%s",
			"RoleID": %d
		}
  ]
}`, userToAdd.UserID, roles["developer"].ID))),
				Repo: repo.NewRepo(),
				user: paidUser,
			},
			want:       &models.AddMembersResponse{},
			wantStatus: http.StatusNotFound,
			wantErr:    "not found",
		},
		{
			name: "ask for upgrade if organization is free and pm is not admin",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": freeProject.UUID,
				}),
				body: io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{
  "Members": [
		{
			"MemberID": "%s",
			"RoleID": %d
		}
  ]
}`, userToAdd.UserID, roles["developer"].ID))),
				Repo: repo.NewRepo(),
				user: freeUser,
			},
			want:       &models.AddMembersResponse{},
			wantStatus: http.StatusForbidden,
			wantErr:    "needs upgrade",
		},
		{
			name: "conflict if user already in project",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": paidProject.UUID,
				}),
				body: io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{
  "Members": [
		{
			"MemberID": "%s",
			"RoleID": %d
		}
  ]
}`, addedUser.UserID, roles["developer"].ID))),
				Repo: repo.NewRepo(),
				user: paidUser,
			},
			want: &models.AddMembersResponse{
				Success: false,
				Error:   "user already in project",
			},
			wantStatus: http.StatusConflict,
			wantErr:    "member already in project",
		},
		{
			name: "developer user cannot add admin user",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": paidProject.UUID,
				}),
				body: io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{
		  "Members": [
				{
					"MemberID": "%s",
					"RoleID": %d
				}
			]
		}`, adminToAdd.UserID, roles["admin"].ID))),
				Repo: repo.NewRepo(),
				user: addedUser,
			},
			want: &models.AddMembersResponse{
				Success: false,
				Error:   "permission denied",
			},
			wantStatus: http.StatusForbidden,
			wantErr:    "permission denied",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotStatus, err := PostProjectMembers(
				tt.args.params,
				tt.args.body,
				tt.args.Repo,
				tt.args.user,
			)
			if err.Error() != tt.wantErr {
				t.Errorf(
					"PostProjectsMembers() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"PostProjectsMembers() got = %v, want %v",
					got,
					tt.want,
				)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf(
					"PostProjectsMembers() gotStatus = %v, want %v",
					gotStatus,
					tt.wantStatus,
				)
			}
		})
	}
}

func TestDeleteProjectsMembers(t *testing.T) {
	adminUser, adminOrganization, project := seedOneProjectForOneUser()
	defer teardownUserAndOrganization(adminUser, adminOrganization)
	defer teardownProject(project)

	userToRemove, orgToRemove := seedSingleUser()
	defer teardownUserAndOrganization(userToRemove, orgToRemove)

	userOtherToRemove, orgOtherToRemove := seedSingleUser()
	defer teardownUserAndOrganization(userOtherToRemove, orgOtherToRemove)

	userNotAdmin, orgNotAdmin := seedSingleUser()
	defer teardownUserAndOrganization(userNotAdmin, orgNotAdmin)

	roles := testsGetRoles()

	memberAdmin := seedProjectMember(project, adminUser, roles["admin"])
	defer teardownProjectMember(memberAdmin)

	memberToRemove := seedProjectMember(project, userToRemove, roles["developer"])
	defer teardownProjectMember(memberToRemove)

	memberOtherToRemove := seedProjectMember(project, userOtherToRemove, roles["developer"])
	defer teardownProjectMember(memberOtherToRemove)

	memberNotAdmin := seedProjectMember(project, userNotAdmin, roles["developer"])
	defer teardownProjectMember(memberNotAdmin)

	type args struct {
		params router.Params
		body   io.ReadCloser
		Repo   repo.IRepo
		user   models.User
	}
	tests := []struct {
		name       string
		args       args
		want       *models.RemoveMembersResponse
		wantStatus int
		wantErr    string
	}{
		{
			name: "removes a member",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
				}),
				body: io.NopCloser(
					bytes.NewBufferString(fmt.Sprintf(`{
  "Members": ["%s"]
}`,
						userToRemove.UserID,
					))),
				Repo: repo.NewRepo(),
				user: adminUser,
			},
			want: &models.RemoveMembersResponse{
				Success: true,
				Error:   "",
			},
			wantStatus: http.StatusOK,
			wantErr:    "",
		},
		{
			name: "fails because no such project",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": "this is not a project",
				}),
				body: io.NopCloser(
					bytes.NewBufferString(fmt.Sprintf(`{
  "Members": ["%s"]
}`,
						userToRemove.UserID,
					))),
				Repo: repo.NewRepo(),
				user: adminUser,
			},
			want: &models.RemoveMembersResponse{
				Success: false,
				Error:   "No such project",
			},
			wantStatus: http.StatusNotFound,
			wantErr:    "not found",
		},
		{
			name: "fails because not priviledged enough",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
				}),
				body: io.NopCloser(
					bytes.NewBufferString(fmt.Sprintf(`{
  "Members": ["%s"]
}`,
						userOtherToRemove.UserID,
					))),
				Repo: repo.NewRepo(),
				user: userNotAdmin,
			},
			want: &models.RemoveMembersResponse{
				Success: false,
				Error:   "permission denied",
			},
			wantStatus: http.StatusForbidden,
			wantErr:    "permission denied",
		},
		{
			name: "fails because members not in project",
			args: args{
				params: router.ParamsFrom(map[string]string{
					"projectID": project.UUID,
				}),
				body: io.NopCloser(
					bytes.NewBufferString(`{
  "Members": ["not@member", "me@neither"]
}`)),
				Repo: repo.NewRepo(),
				user: adminUser,
			},
			want: &models.RemoveMembersResponse{
				Success: false,
				Error:   "not a member",
			},
			wantStatus: http.StatusConflict,
			wantErr:    "not a member",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotStatus, err := DeleteProjectsMembers(
				tt.args.params,
				tt.args.body,
				tt.args.Repo,
				tt.args.user,
			)
			if err.Error() != tt.wantErr {
				t.Errorf(
					"DeleteProjectsMembers() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"DeleteProjectsMembers() got = %v, want %v",
					got,
					tt.want,
				)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf(
					"DeleteProjectsMembers() gotStatus = %v, want %v",
					gotStatus,
					tt.wantStatus,
				)
			}
		})
	}
}

func TestGetAccessibleEnvironments(t *testing.T) {
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
			got, gotStatus, err := GetAccessibleEnvironments(
				tt.args.params,
				tt.args.in1,
				tt.args.Repo,
				tt.args.user,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"GetAccessibleEnvironments() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"GetAccessibleEnvironments() got = %v, want %v",
					got,
					tt.want,
				)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf(
					"GetAccessibleEnvironments() gotStatus = %v, want %v",
					gotStatus,
					tt.wantStatus,
				)
			}
		})
	}
}

func TestDeleteProject(t *testing.T) {
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
			got, gotStatus, err := DeleteProject(
				tt.args.params,
				tt.args.in1,
				tt.args.Repo,
				tt.args.user,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"DeleteProject() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteProject() got = %v, want %v", got, tt.want)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf(
					"DeleteProject() gotStatus = %v, want %v",
					gotStatus,
					tt.wantStatus,
				)
			}
		})
	}
}

func TestGetProjectsOrganization(t *testing.T) {
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
			got, gotStatus, err := GetProjectsOrganization(
				tt.args.params,
				tt.args.in1,
				tt.args.Repo,
				tt.args.user,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"GetProjectsOrganization() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"GetProjectsOrganization() got = %v, want %v",
					got,
					tt.want,
				)
			}
			if gotStatus != tt.wantStatus {
				t.Errorf(
					"GetProjectsOrganization() gotStatus = %v, want %v",
					gotStatus,
					tt.wantStatus,
				)
			}
		})
	}
}

func seedSingleUser() (user models.User, organization models.Organization) {
	err := repo.NewRepo().GetDb().Transaction(func(db *gorm.DB) error {
		faker.FakeData(&user)
		db.Create(&user)

		faker.FakeData(&organization)
		organization.UserID = user.ID
		organization.User = user

		db.Create(&organization)

		return db.Error
	})

	if err != nil {
		panic(err)
	}

	return user, organization
}

func teardownUserAndOrganization(
	user models.User,
	organization models.Organization,
) {
	repo.NewRepo().GetDb().Transaction(func(db *gorm.DB) error {
		db.Delete(&user)
		db.Delete(&organization)
		return db.Error
	})
}

func seedOneProjectForOneUser() (user models.User, organization models.Organization, project models.Project) {
	user, organization = seedSingleUser()
	faker.FakeData(&project)
	project.UserID = user.ID
	project.User = user
	project.OrganizationID = organization.ID
	project.Organization = organization

	err := repo.NewRepo().GetDb().Transaction(func(db *gorm.DB) error {
		db.Create(&project)

		return db.Error
	})

	if err != nil {
		panic(err)
	}

	return user, organization, project
}

func testsSetOrganisationPaid(organization *models.Organization) {
	err := repo.NewRepo().GetDb().Transaction(func(db *gorm.DB) error {
		organization.Paid = true
		return db.Save(organization).Error
	})

	if err != nil {
		panic(err)
	}
}

func testsGetRoles() map[string]models.Role {
	r := make(map[string]models.Role)

	err := repo.NewRepo().GetDb().Transaction(func(db *gorm.DB) error {
		roles := []models.Role{}
		db.Find(&roles)

		for _, n := range roles {
			r[n.Name] = n
		}

		return db.Error
	})

	if err != nil {
		panic(err)
	}

	return r
}

func seedProjectMember(project models.Project, user models.User, role models.Role) (projectMember models.ProjectMember) {

	err := repo.NewRepo().GetDb().Transaction(func(db *gorm.DB) error {
		projectMember := models.ProjectMember{
			ProjectID: project.ID,
			Project:   project,
			UserID:    user.ID,
			User:      user,
			RoleID:    role.ID,
			Role:      role,
		}

		return db.
			Model(&models.ProjectMember{}).
			Create(&projectMember).
			Error
	})

	if err != nil {
		panic(err)
	}

	return projectMember
}

func seedProjectMembers(project models.Project, users []models.User, role models.Role) []models.ProjectMember {
	pms := make([]models.ProjectMember, len(users))

	err := repo.NewRepo().GetDb().Transaction(func(db *gorm.DB) error {
		for idx, user := range users {
			projectMember := models.ProjectMember{
				ProjectID: project.ID,
				Project:   project,
				UserID:    user.ID,
				User:      user,
				RoleID:    role.ID,
				Role:      role,
			}
			pms[idx] = projectMember
		}

		return db.
			Model(&models.ProjectMember{}).
			Create(&pms).
			Error
	})

	if err != nil {
		panic(err)
	}

	return pms
}

func teardownProjectMember(projectMember models.ProjectMember) {
	err := repo.NewRepo().GetDb().Transaction(func(db *gorm.DB) error {
		return db.
			Model(&models.ProjectMember{}).
			Delete("id = ?", projectMember.ID).
			Error
	})

	if err != nil {
		panic(err)
	}
}
func teardownProjectMembers(projectMembers []models.ProjectMember) {
	for _, m := range projectMembers {
		teardownProjectMember(m)
	}
}

func seedManyProjectsForOneUser() (user models.User, organization models.Organization, projects []models.Project) {
	err := repo.NewRepo().GetDb().Transaction(func(db *gorm.DB) error {
		user, organization = seedSingleUser()
		projects = make([]models.Project, 10)

		for i, project := range projects {
			faker.FakeData(&project)
			project.UserID = user.ID
			project.User = user
			project.OrganizationID = organization.ID
			project.Organization = organization

			if err := db.Create(&project).Error; err != nil {
				return err
			}
			if err := db.Create(&models.ProjectMember{
				UserID:    user.ID,
				ProjectID: project.ID,
				RoleID:    4,
			}).Error; err != nil {
				return err
			}

			projects[i] = project
		}

		db.Joins(
			"inner join project_members pm on pm.user_id = ? and pm.project_id = projects.ID",
			user.ID,
		).
			Find(&projects)

		return db.Error
	})

	if err != nil {
		panic(err)
	}

	return user, organization, projects
}

func teardownProject(project models.Project) {
	teardownManyProjects([]models.Project{project})
}

func teardownManyProjects(projects []models.Project) {
	repo.NewRepo().GetDb().Transaction(func(db *gorm.DB) error {
		for _, project := range projects {
			db.Delete(
				&models.ProjectMember{},
				"project_id = ?",
				project.ID,
			)
		}
		db.Delete(projects)

		return db.Error
	})
}
