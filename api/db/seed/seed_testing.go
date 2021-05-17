// +build test

package seed

import (
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func SeedRoles() error {

	Repo := new(repo.Repo)

	devRole := models.Role{Name: "developer"}
	leadDevRole := models.Role{Name: "lead-dev", CanAddMember: true}
	devopsRole := models.Role{Name: "devops", CanAddMember: true}
	adminRole := models.Role{Name: "admin", CanAddMember: true}

	if err := Repo.GetOrCreateRole(&adminRole).Err(); err != nil {
		return err
	}

	devopsRole.ParentID = adminRole.ID
	if err := Repo.GetOrCreateRole(&devopsRole).Err(); err != nil {
		return err
	}

	leadDevRole.ParentID = devopsRole.ID
	if err := Repo.GetOrCreateRole(&leadDevRole).Err(); err != nil {
		return err
	}

	devRole.ParentID = leadDevRole.ID
	if err := Repo.GetOrCreateRole(&devRole).Err(); err != nil {
		return err
	}

	devEnvironmentType := models.EnvironmentType{Name: "dev"}
	stagingEnvironmentType := models.EnvironmentType{Name: "staging"}
	prodEnvironmentType := models.EnvironmentType{Name: "prod"}

	if err := Repo.
		GetOrCreateEnvironmentType(&devEnvironmentType).
		GetOrCreateEnvironmentType(&stagingEnvironmentType).
		GetOrCreateEnvironmentType(&prodEnvironmentType).
		Err(); err != nil {
		return err
	}

	// developer
	devRoleDevEnvType := models.RolesEnvironmentType{
		RoleID:            devRole.ID,
		EnvironmentTypeID: devEnvironmentType.ID,
		Read:              true,
		Write:             true,
	}
	devRoleStagingEnvType := models.RolesEnvironmentType{
		RoleID:            devRole.ID,
		EnvironmentTypeID: stagingEnvironmentType.ID,
		Read:              false,
		Write:             false,
	}
	devRoleProdEnvType := models.RolesEnvironmentType{
		RoleID:            devRole.ID,
		EnvironmentTypeID: prodEnvironmentType.ID,
		Read:              false,
		Write:             false,
	}

	// lead
	leadDevRoleDevEnvType := models.RolesEnvironmentType{
		RoleID:            leadDevRole.ID,
		EnvironmentTypeID: devEnvironmentType.ID,
		Read:              true,
		Write:             true,
	}
	leadDevRoleStagingEnvType := models.RolesEnvironmentType{
		RoleID:            leadDevRole.ID,
		EnvironmentTypeID: stagingEnvironmentType.ID,
		Read:              false,
		Write:             false,
	}
	leadDevRoleProdEnvType := models.RolesEnvironmentType{
		RoleID:            leadDevRole.ID,
		EnvironmentTypeID: prodEnvironmentType.ID,
		Read:              false,
		Write:             false,
	}

	// devops
	devopsRoleDevEnvType := models.RolesEnvironmentType{
		RoleID:            devopsRole.ID,
		EnvironmentTypeID: devEnvironmentType.ID,
		Read:              true,
		Write:             true,
	}
	devopsRoleStagingEnvType := models.RolesEnvironmentType{
		RoleID:            devopsRole.ID,
		EnvironmentTypeID: stagingEnvironmentType.ID,
		Read:              true,
		Write:             true,
	}
	devopsRoleProdEnvType := models.RolesEnvironmentType{
		RoleID:            devopsRole.ID,
		EnvironmentTypeID: prodEnvironmentType.ID,
		Read:              true,
		Write:             true,
	}

	// admin
	adminRoleDevEnvType := models.RolesEnvironmentType{
		RoleID:            adminRole.ID,
		EnvironmentTypeID: devEnvironmentType.ID,
		Read:              true,
		Write:             true,
	}
	adminRoleStagingEnvType := models.RolesEnvironmentType{
		RoleID:            adminRole.ID,
		EnvironmentTypeID: stagingEnvironmentType.ID,
		Read:              true,
		Write:             true,
	}
	adminRoleProdEnvType := models.RolesEnvironmentType{
		RoleID:            adminRole.ID,
		EnvironmentTypeID: prodEnvironmentType.ID,
		Read:              true,
		Write:             true,
	}

	if err := Repo.
		GetOrCreateRoleEnvType(&devRoleDevEnvType).
		GetOrCreateRoleEnvType(&devRoleStagingEnvType).
		GetOrCreateRoleEnvType(&devRoleProdEnvType).
		GetOrCreateRoleEnvType(&leadDevRoleDevEnvType).
		GetOrCreateRoleEnvType(&leadDevRoleStagingEnvType).
		GetOrCreateRoleEnvType(&leadDevRoleProdEnvType).
		GetOrCreateRoleEnvType(&devopsRoleDevEnvType).
		GetOrCreateRoleEnvType(&devopsRoleStagingEnvType).
		GetOrCreateRoleEnvType(&devopsRoleProdEnvType).
		GetOrCreateRoleEnvType(&adminRoleDevEnvType).
		GetOrCreateRoleEnvType(&adminRoleStagingEnvType).
		GetOrCreateRoleEnvType(&adminRoleProdEnvType).
		Err(); err != nil {
		return err
	}

	return nil
}
