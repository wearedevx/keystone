package repo

import (
	"github.com/wearedevx/keystone/api/pkg/models"
	"gorm.io/gorm/clause"
)

func (r *Repo) SaveActivityLog(al *models.ActivityLog) IRepo {
	r.err = r.GetDb().
		Omit(clause.Associations).
		Create(al).
		Error

	return r
}
