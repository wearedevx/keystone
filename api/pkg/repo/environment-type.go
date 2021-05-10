package repo

import (
	. "github.com/wearedevx/keystone/api/pkg/models"
)

func (repo *Repo) CreateEnvironmentType(envType *EnvironmentType) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDb().Create(&envType).Error

	return repo
}

func (repo *Repo) GetEnvironmentType(envType *EnvironmentType) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDb().First(&envType).Error

	return repo
}

func (repo *Repo) GetOrCreateEnvironmentType(envType *EnvironmentType) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDb().FirstOrCreate(envType).Error

	return repo
}
