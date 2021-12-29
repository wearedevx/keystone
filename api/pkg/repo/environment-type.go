package repo

import (
	"github.com/wearedevx/keystone/api/pkg/models"
)

func (repo *Repo) CreateEnvironmentType(envType *models.EnvironmentType) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDb().Create(&envType).Error

	return repo
}

func (repo *Repo) GetEnvironmentType(envType *models.EnvironmentType) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDb().Where(*envType).First(&envType).Error

	return repo
}

func (repo *Repo) GetOrCreateEnvironmentType(envType *models.EnvironmentType) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDb().Where(*envType).FirstOrCreate(envType).Error

	return repo
}
