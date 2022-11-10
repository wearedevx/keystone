package repo

import (
	"github.com/wearedevx/keystone/api/pkg/models"
)

func (repo *Repo) CreateEnvironmentType(envType *models.EnvironmentType) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDB().Create(&envType).Error

	return repo
}

func (repo *Repo) GetEnvironmentType(envType *models.EnvironmentType) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDB().Where(*envType).First(&envType).Error

	return repo
}

func (repo *Repo) GetOrCreateEnvironmentType(envType *models.EnvironmentType) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDB().Where(*envType).FirstOrCreate(envType).Error

	return repo
}
