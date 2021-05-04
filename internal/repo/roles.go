package repo

import (
	"fmt"

	. "github.com/wearedevx/keystone/internal/models"
)

func (r *Repo) GetRoles(roles *[]Role) *Repo {
	if r.err != nil {
		return r
	}

	db := r.GetDb()
	r.err = db.Find(roles).Error
	fmt.Println("roles:", roles)

	return r
}
