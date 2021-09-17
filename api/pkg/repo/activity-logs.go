package repo

import (
	"fmt"

	"github.com/wearedevx/keystone/api/pkg/models"
	"gorm.io/gorm/clause"
)

func (r *Repo) SaveActivityLog(al *models.ActivityLog) IRepo {
	r.err = r.GetDb().
		Omit(clause.Associations).
		Create(al).
		Error
	fmt.Printf("r.err: %+v\n", r.err)
	fmt.Printf("al: %+v\n", al.ID)

	return r
}
