package repo

import (
	"fmt"

	. "github.com/wearedevx/keystone/api/pkg/models"
)

func (r *Repo) GetRoles(roles *[]Role) *Repo {
	if r.err != nil {
		return r
	}

	db := r.GetDb()
	r.err = db.Find(roles).Error

	return r
}

func (repo *Repo) CreateRole(role *Role) *Repo {
	repo.err = repo.GetDb().Create(&role).Error
	return repo
}

func (repo *Repo) GetRoleByName(name string, role *Role) *Repo {
	if repo.Err() != nil {
		return repo
	}

	repo.err = repo.GetDb().Where("name = ?", name).First(&role).Error

	return repo
}

func (repo *Repo) GetRoleByID(id uint, role *Role) *Repo {
	if repo.Err() != nil {
		return repo
	}

	repo.err = repo.GetDb().First(role, id).Error

	return repo
}

func (repo *Repo) GetOrCreateRole(name string, role *Role) *Repo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDb().Where(Role{Name: name}).FirstOrCreate(&role).Error

	return repo
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
