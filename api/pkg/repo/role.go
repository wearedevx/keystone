package repo

import (
	. "github.com/wearedevx/keystone/api/pkg/models"
)

func (repo *Repo) CreateRole(name string) Role {
	role := Role{
		Name: name,
	}
	repo.err = repo.GetDb().Create(&role).Error
	return role
}

func (repo *Repo) GetRoleByName(name string) (Role, bool) {
	var foundRole Role

	// if repo.err != nil {
	// 	return foundRole, false
	// }

	repo.err = repo.GetDb().Where("name = ?", name).First(&foundRole).Error

	return foundRole, repo.err == nil
}

func (repo *Repo) GetOrCreateRole(name string) Role {
	if repo.err != nil {
		return Role{}
	}

	if env, ok := repo.GetRoleByName(name); ok {
		return env
	}

	return repo.CreateRole(name)
}
