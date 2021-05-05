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

	return r
}

func (r *Repo) GetInvitableRoles(role Role, roles *[]Role) *Repo {

	r.err = db.Model(&Role{}).
		Joins("left join roles_environment_types on roles_environment_types.role_id = roles.id").
		Where("roles_environment_types.role_id = ? and roles_environment_types.invite = true", role.ID).
		Find(roles).Error
	fmt.Println("keystone ~ roles.go ~ r.err", r.err)

	return r
	// repo.err = repo.GetDb().Where("roleID = ? and environmentID = ?", role.ID, environmentType.ID).Find(&roles).Error
}
