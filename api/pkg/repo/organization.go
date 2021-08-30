package repo

import (
	"github.com/wearedevx/keystone/api/pkg/models"
)

func (repo *Repo) CreateOrganization(orga *models.Organization) IRepo {
	repo.err = repo.GetDb().Create(orga).Error
	return repo
}
