package rights

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wearedevx/keystone/api/pkg/repo"
	. "github.com/wearedevx/keystone/internal/models"
)

type FakeRepo struct{}

func getRoleByEnvAndRole(environment *Environment, role *Role) RolesEnvironmentType {
	switch {
	case environment.Name == "dev" && role.Name == "dev":
		return RolesEnvironmentType{
			Read:   true,
			Write:  true,
			Invite: false,
		}
	case environment.Name == "staging" && role.Name == "dev":
		return RolesEnvironmentType{
			Read:   false,
			Write:  false,
			Invite: false,
		}
	case environment.Name == "prod" && role.Name == "dev":
		return RolesEnvironmentType{
			Read:   false,
			Write:  false,
			Invite: false,
		}

	case environment.Name == "dev" && role.Name == "devops":
		return RolesEnvironmentType{
			Read:   true,
			Write:  true,
			Invite: true,
		}
	case environment.Name == "staging" && role.Name == "devops":
		return RolesEnvironmentType{
			Read:   true,
			Write:  true,
			Invite: false,
		}
	case environment.Name == "prod" && role.Name == "devops":
		return RolesEnvironmentType{
			Read:   false,
			Write:  false,
			Invite: false,
		}

	case environment.Name == "dev" && role.Name == "admin":
		return RolesEnvironmentType{
			Read:   true,
			Write:  true,
			Invite: true,
		}
	case environment.Name == "staging" && role.Name == "admin":
		return RolesEnvironmentType{
			Read:   true,
			Write:  true,
			Invite: true,
		}
	case environment.Name == "prod" && role.Name == "admin":
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

func (fakeRepo *FakeRepo) GetRolesEnvironmentType(environment *Environment, role *Role) (*RolesEnvironmentType, error) {
	rolesEnvironmentType := getRoleByEnvAndRole(environment, role)
	return &rolesEnvironmentType, nil
}

func (fakeRepo *FakeRepo) GetProjectMember(user *User, project *Project) (ProjectMember, error) {
	role := getRoleByUsername(user.Username)
	projectMember := ProjectMember{
		Role: role,
	}
	return projectMember, nil
}
func (fakeRepo *FakeRepo) GetInvitableRoles(role Role, roles []*Role) *repo.Repo {
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
