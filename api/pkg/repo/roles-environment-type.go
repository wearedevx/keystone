package repo

import (
	. "github.com/wearedevx/keystone/internal/models"
)

func (repo *Repo) CreateRoleEnvironmentType(rolesEnvironmentType *RolesEnvironmentType) *RolesEnvironmentType {
	// rolesEnvironmentTypeToCreate := RolesEnvironmentType{
	// 	RoleID:            rolesEnvironmentTypeToCreate.roleID,
	// 	EnvironmentTypeID: environmentTypeID,
	// }
	repo.err = repo.GetDb().Create(&rolesEnvironmentType).Error
	return rolesEnvironmentType
}

func (repo *Repo) GetRoleEnvTypeByRoleAndEnvironment(role *Role, environmentType *EnvironmentType) (*RolesEnvironmentType, bool) {
	var foundRoleEnvType *RolesEnvironmentType

	// if repo.err != nil {
	// 	return foundRoleEnvType, false
	// }

	repo.err = repo.GetDb().Where("roleID = ? and environmentID = ?", role.ID, environmentType.ID).First(&foundRoleEnvType).Error

	return foundRoleEnvType, repo.err == nil
}

func (repo *Repo) GetOrCreateRoleEnvType(rolesEnvironmentType *RolesEnvironmentType) *RolesEnvironmentType {
	if rolesEnvironmentType, ok := repo.GetRoleEnvTypeByRoleAndEnvironment(&rolesEnvironmentType.Role, &rolesEnvironmentType.EnvironmentType); ok {
		return rolesEnvironmentType
	}

	return repo.CreateRoleEnvironmentType(rolesEnvironmentType)
}

func (repo *Repo) GetRolesEnvironmentType(environment *Environment, role *Role) (*RolesEnvironmentType, error) {
	rolesEnvironmentType := RolesEnvironmentType{}

	repo.err = db.Where(
		"environment_type_id = ? and role_id = ?",
		environment.ID, role.ID,
	).First(&rolesEnvironmentType).Error

	return &rolesEnvironmentType, repo.err
}
