package repo

import "github.com/wearedevx/keystone/api/pkg/models"

func (r *Repo) SaveActivityLog(al *models.ActivityLog) IRepo {
	r.err = r.GetDb().
		Omit("User").
		Omit("Project").
		Omit("Environment").
		Create(al).
		Error

	return r
}
