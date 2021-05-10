package rights

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wearedevx/keystone/api/pkg/models"
	. "github.com/wearedevx/keystone/api/pkg/models"
	. "github.com/wearedevx/keystone/api/pkg/repo"
	"gorm.io/gorm"
)

type FakeRepo struct{}

func (f *FakeRepo) CreateEnvironment(_ *models.Environment) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) CreateEnvironmentType(_ *models.EnvironmentType) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) CreateLoginRequest() models.LoginRequest {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) CreateProjectMember(_ *models.ProjectMember, _ *models.Role) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) CreateRole(_ *models.Role) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) CreateRoleEnvironmentType(_ *models.RolesEnvironmentType) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) CreateSecret(_ *models.Secret) {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) DeleteLoginRequest(_ string) bool {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) FindUsers(_ []string) (map[string]models.User, []string) {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetDb() *gorm.DB {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetEnvironment(_ *models.Environment) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetEnvironmentType(_ *models.EnvironmentType) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetLoginRequest(_ string) (models.LoginRequest, bool) {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetOrCreateEnvironment(_ *models.Environment) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetOrCreateEnvironmentType(_ *models.EnvironmentType) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetOrCreateProject(_ *models.Project) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetOrCreateProjectMember(_ *models.ProjectMember, _ string) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetOrCreateRole(_ *models.Role) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetOrCreateRoleEnvType(_ *models.RolesEnvironmentType) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetOrCreateUser(_ *models.User) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetProject(_ *models.Project) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetProjectByUUID(_ string, _ *models.Project) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetRoleByID(_ uint, _ *models.Role) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetRoleByName(_ string, _ *models.Role) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetRoles(_ *[]models.Role) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetSecretByName(_ string, _ *models.Secret) {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetUser(_ *models.User) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) ProjectAddMembers(_ models.Project, _ []models.MemberRole) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) ProjectGetMembers(_ *models.Project, _ *[]models.ProjectMember) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) ProjectLoadUsers(_ *models.Project) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) ProjectRemoveMembers(_ models.Project, _ []string) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) ProjectSetRoleForUser(_ models.Project, _ models.User, _ models.Role) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) SetLoginRequestCode(_ string, _ string) models.LoginRequest {
	panic("not implemented") // TODO: Implement
}

func getRoleByEnvironmentTypeAndRole(environmentType *EnvironmentType, role *Role) RolesEnvironmentType {
	switch {
	case environmentType.Name == "dev" && role.Name == "dev":
		return RolesEnvironmentType{
			Read:   true,
			Write:  true,
			Invite: false,
		}
	case environmentType.Name == "staging" && role.Name == "dev":
		return RolesEnvironmentType{
			Read:   false,
			Write:  false,
			Invite: false,
		}
	case environmentType.Name == "prod" && role.Name == "dev":
		return RolesEnvironmentType{
			Read:   false,
			Write:  false,
			Invite: false,
		}

	case environmentType.Name == "dev" && role.Name == "devops":
		return RolesEnvironmentType{
			Read:   true,
			Write:  true,
			Invite: true,
		}
	case environmentType.Name == "staging" && role.Name == "devops":
		return RolesEnvironmentType{
			Read:   true,
			Write:  true,
			Invite: false,
		}
	case environmentType.Name == "prod" && role.Name == "devops":
		return RolesEnvironmentType{
			Read:   false,
			Write:  false,
			Invite: false,
		}

	case environmentType.Name == "dev" && role.Name == "admin":
		return RolesEnvironmentType{
			Read:   true,
			Write:  true,
			Invite: true,
		}
	case environmentType.Name == "staging" && role.Name == "admin":
		return RolesEnvironmentType{
			Read:   true,
			Write:  true,
			Invite: true,
		}
	case environmentType.Name == "prod" && role.Name == "admin":
		return RolesEnvironmentType{
			Read:   true,
			Write:  true,
			Invite: true,
		}

	default:
		return RolesEnvironmentType{
			Read:   false,
			Write:  false,
			Invite: false,
		}
	}
}

func getRoleByUsername(userName string) Role {
	switch {
	case userName == "dev":
		return Role{
			Name: "dev",
		}
	case userName == "devops":
		return Role{
			Name: "devops",
		}
	case userName == "admin":
		return Role{
			Name: "admin",
		}
	default:
		return Role{
			Name: "nothing",
		}
	}
}

func (fakeRepo *FakeRepo) Err() error {
	return nil
}

func (fakeRepo *FakeRepo) GetRolesEnvironmentType(rolesEnvironmentType *RolesEnvironmentType) IRepo {
	*rolesEnvironmentType = getRoleByEnvironmentTypeAndRole(&rolesEnvironmentType.EnvironmentType, &rolesEnvironmentType.Role)

	return fakeRepo
}

func (fakeRepo *FakeRepo) GetProjectMember(projectMember *ProjectMember) IRepo {
	role := getRoleByUsername(projectMember.User.Username)
	*projectMember = ProjectMember{
		Role: role,
	}

	return fakeRepo
}
func (fakeRepo *FakeRepo) GetInvitableRoles(role Role, roles *[]Role) IRepo {
	// return fakeRepo
	// TODO
	return nil
}

func TestCanUserHasRightEnvironment(t *testing.T) {
	fakeRepo := &FakeRepo{}
	project := &Project{}

	userDev := &User{Username: "dev"}
	userDevops := &User{Username: "devops"}
	userAdmin := &User{Username: "admin"}

	environmentDev := &Environment{Name: "dev"}
	environmentStaging := &Environment{Name: "staging"}
	environmentProd := &Environment{Name: "prod"}

	// DEV env

	// Dev user
	canDevReadDev, _ := CanUserReadEnvironment(fakeRepo, userDev, project, environmentDev)
	assert.True(t, canDevReadDev, "Oops! "+userDev.Username+" should be able to read "+environmentDev.Name+" environment")
	canDevWriteDev, _ := CanUserWriteOnEnvironment(fakeRepo, userDev, project, environmentDev)
	assert.True(t, canDevWriteDev, "Oops! "+userDev.Username+" should be able to write "+environmentDev.Name+" environment")
	canDevInviteDev, _ := CanUserInviteOnEnvironment(fakeRepo, userDev, project, environmentDev)
	assert.False(t, canDevInviteDev, "Oops! "+userDev.Username+" can't invite on "+environmentDev.Name+" environment")

	// Devops user
	canDevopsReadDev, _ := CanUserReadEnvironment(fakeRepo, userDevops, project, environmentDev)
	assert.True(t, canDevopsReadDev, "Oops! "+userDevops.Username+" user should be able to read "+environmentDev.Name+" environment")
	canDevopsWriteDev, _ := CanUserWriteOnEnvironment(fakeRepo, userDevops, project, environmentDev)
	assert.True(t, canDevopsWriteDev, "Oops! "+userDevops.Username+" user should be able to write "+environmentDev.Name+" environment")
	canDevopsInviteDev, _ := CanUserInviteOnEnvironment(fakeRepo, userDevops, project, environmentDev)
	assert.True(t, canDevopsInviteDev, "Oops! "+userDevops.Username+" should be able to "+environmentDev.Name+" environment")

	// Admin user
	canAdminReadDev, _ := CanUserReadEnvironment(fakeRepo, userAdmin, project, environmentDev)
	assert.True(t, canAdminReadDev, "Oops! "+userAdmin.Username+" user should be able to read "+environmentDev.Name+" environment")
	canAdminWriteDev, _ := CanUserWriteOnEnvironment(fakeRepo, userAdmin, project, environmentDev)
	assert.True(t, canAdminWriteDev, "Oops! "+userAdmin.Username+" user should be able to write "+environmentDev.Name+" environment")
	canAdminInviteDev, _ := CanUserInviteOnEnvironment(fakeRepo, userAdmin, project, environmentDev)
	assert.True(t, canAdminInviteDev, "Oops! "+userAdmin.Username+" should be able to "+environmentDev.Name+" environment")

	// Staging env

	// Dev user
	canDevReadStaging, _ := CanUserReadEnvironment(fakeRepo, userDev, project, environmentStaging)
	assert.False(t, canDevReadStaging, "Oops! "+userDev.Username+" can't read "+environmentStaging.Name+" environment")
	canDevWriteStaging, _ := CanUserWriteOnEnvironment(fakeRepo, userDev, project, environmentStaging)
	assert.False(t, canDevWriteStaging, "Oops! "+userDev.Username+" user can't write "+environmentStaging.Name+" environment")
	canDevInviteStaging, _ := CanUserInviteOnEnvironment(fakeRepo, userDev, project, environmentStaging)
	assert.False(t, canDevInviteStaging, "Oops! "+userDev.Username+" user can't invite on "+environmentStaging.Name+" environment")

	// Devops user
	canDevopsReadCI, _ := CanUserReadEnvironment(fakeRepo, userDevops, project, environmentStaging)
	assert.True(t, canDevopsReadCI, "Oops! "+userDevops.Username+" user should be able to read "+environmentStaging.Name+" environment")
	canDevopsWriteCI, _ := CanUserWriteOnEnvironment(fakeRepo, userDevops, project, environmentStaging)
	assert.True(t, canDevopsWriteCI, "Oops! "+userDevops.Username+" user should be able to write "+environmentStaging.Name+" environment")
	canDevopsInviteCI, _ := CanUserInviteOnEnvironment(fakeRepo, userDevops, project, environmentStaging)
	assert.False(t, canDevopsInviteCI, "Oops! "+userDevops.Username+" can't invite to "+environmentStaging.Name+" environment")

	// Admin user
	canAdminReadCI, _ := CanUserReadEnvironment(fakeRepo, userAdmin, project, environmentStaging)
	assert.True(t, canAdminReadCI, "Oops! "+userAdmin.Username+" user should be able to read "+environmentStaging.Name+" environment")
	canAdminWriteCI, _ := CanUserWriteOnEnvironment(fakeRepo, userAdmin, project, environmentStaging)
	assert.True(t, canAdminWriteCI, "Oops! "+userAdmin.Username+" user should be able to write "+environmentStaging.Name+" environment")
	canAdminInviteCI, _ := CanUserInviteOnEnvironment(fakeRepo, userAdmin, project, environmentStaging)
	assert.True(t, canAdminInviteCI, "Oops! "+userAdmin.Username+" should be able to "+environmentStaging.Name+" environment")

	// PROD env

	// Dev user
	canDevReadProd, _ := CanUserReadEnvironment(fakeRepo, userDev, project, environmentProd)
	assert.False(t, canDevReadProd, "Oops! "+userDev.Username+" can't read "+environmentProd.Name+" environment")
	canDevWriteProd, _ := CanUserWriteOnEnvironment(fakeRepo, userDev, project, environmentProd)
	assert.False(t, canDevWriteProd, "Oops! "+userDev.Username+" user can't write "+environmentProd.Name+" environment")
	canDevInviteProd, _ := CanUserInviteOnEnvironment(fakeRepo, userDev, project, environmentProd)
	assert.False(t, canDevInviteProd, "Oops! "+userDev.Username+" user can't invite on "+environmentProd.Name+" environment")

	// Devops user
	canDevopsReadProd, _ := CanUserReadEnvironment(fakeRepo, userDevops, project, environmentProd)
	assert.False(t, canDevopsReadProd, "Oops! "+userDevops.Username+"can't read "+environmentProd.Name+" environment")
	canDevopsWriteProd, _ := CanUserWriteOnEnvironment(fakeRepo, userDevops, project, environmentProd)
	assert.False(t, canDevopsWriteProd, "Oops! "+userDevops.Username+"can't write "+environmentProd.Name+" environment")
	canDevopsInviteProd, _ := CanUserInviteOnEnvironment(fakeRepo, userDevops, project, environmentProd)
	assert.False(t, canDevopsInviteProd, "Oops! "+userDevops.Username+" can't invite to "+environmentProd.Name+" environment")

	// Admin user
	canAdminReadProd, _ := CanUserReadEnvironment(fakeRepo, userAdmin, project, environmentProd)
	assert.True(t, canAdminReadProd, "Oops! "+userAdmin.Username+" user should be able to read "+environmentProd.Name+" environment")
	canAdminWriteProd, _ := CanUserWriteOnEnvironment(fakeRepo, userAdmin, project, environmentProd)
	assert.True(t, canAdminWriteProd, "Oops! "+userAdmin.Username+" user should be able to write "+environmentProd.Name+" environment")
	canAdminInviteProd, _ := CanUserInviteOnEnvironment(fakeRepo, userAdmin, project, environmentProd)
	assert.True(t, canAdminInviteProd, "Oops! "+userAdmin.Username+" should be able to "+environmentProd.Name+" environment")

}
