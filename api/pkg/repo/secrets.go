package repo

import (
	. "github.com/wearedevx/keystone/api/pkg/models"
)

func (repo *Repo) CreateSecret(secret *Secret) {
	if repo.err != nil {
		return
	}

	repo.err = repo.GetDb().Model(secret).Create(secret).Error

	return
}

func (repo *Repo) GetSecretByName(name string, secret *Secret) {
	if repo.err != nil {
		return
	}

	repo.err = repo.GetDb().Where("name = ?", name).FirstOrCreate(secret).Error
}
