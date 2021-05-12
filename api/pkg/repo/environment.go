package repo

import (
	. "github.com/wearedevx/keystone/api/pkg/models"
)

func (repo *Repo) CreateEnvironment(environment *Environment) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDb().Create(&environment).Error

	return repo
}

func (repo *Repo) GetEnvironment(environment *Environment) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDb().
		Preload("EnvironmentType").
		Where(*environment).
		First(&environment).
		Error

	return repo
}

func (repo *Repo) GetOrCreateEnvironment(environment *Environment) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDb().
		Preload("EnvironmentType").
		Where(*environment).
		FirstOrCreate(environment).
		Error

	return repo
}
