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

func (repo *Repo) GetRoleByName(name string, role *Role) *Repo {
	if repo.Err() != nil {
		return repo
	}

	repo.err = repo.GetDb().Where("name = ?", name).First(&role).Error

	return repo
}

func (repo *Repo) GetRoleByID(id uint, role *Role) *Repo {
	if repo.Err() != nil {
		return repo
	}

	repo.err = repo.GetDb().First(role, id).Error

	return repo
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
