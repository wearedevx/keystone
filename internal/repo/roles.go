package repo

import (
	. "github.com/wearedevx/keystone/internal/models"
)

func (r *Repo) GetRoles(roles *[]Role) *Repo {
	if r.err != nil {
		return r
	}

	db := r.GetDb()
	db.Model(Role{}).Find(roles)

	return r
}
