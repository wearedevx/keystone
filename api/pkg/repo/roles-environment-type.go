package repo

import (
	"github.com/wearedevx/keystone/api/pkg/models"
)

func (repo *Repo) CreateRoleEnvironmentType(
	rolesEnvironmentType *models.RolesEnvironmentType,
) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDB().Create(&rolesEnvironmentType).Error

	return repo
}

func (repo *Repo) GetOrCreateRoleEnvType(
	ret *models.RolesEnvironmentType,
) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDB().
		Where(
			"role_id = ? and environment_type_id = ?",
			ret.RoleID,
			ret.EnvironmentTypeID,
		).
		FirstOrCreate(ret).
		Error

	return repo
}

func (repo *Repo) GetRolesEnvironmentType(
	rolesEnvironmentType *models.RolesEnvironmentType,
) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDB().
		Where(*rolesEnvironmentType).
		First(&rolesEnvironmentType).
		Error

	return repo
}
