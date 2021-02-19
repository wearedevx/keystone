package repo

import (
	. "github.com/wearedevx/keystone/internal/models"
)

func (repo *Repo) CreateSecret(secret *Secret) {
	if repo.err != nil {
		return
	}

	repo.err = repo.db.Model(secret).Create(secret).Error

	return
}

func (repo *Repo) GetSecretByName(name string, secret *Secret) {
	if repo.err != nil {
		return
	}

	repo.err = repo.db.Where("name = ?", name).FirstOrCreate(secret).Error
}
