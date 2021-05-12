package repo

import (
	. "github.com/wearedevx/keystone/api/pkg/models"
)

func (repo *Repo) CreateRoleEnvironmentType(rolesEnvironmentType *RolesEnvironmentType) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDb().Create(&rolesEnvironmentType).Error

	return repo
}

func (repo *Repo) GetOrCreateRoleEnvType(rolesEnvironmentType *RolesEnvironmentType) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDb().
		Where(*rolesEnvironmentType).
		FirstOrCreate(rolesEnvironmentType).
		Error

	return repo
}

func (repo *Repo) GetRolesEnvironmentType(rolesEnvironmentType *RolesEnvironmentType) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = db.
		Where(*rolesEnvironmentType).
		First(&rolesEnvironmentType).
		Error

	return repo
}
