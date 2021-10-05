package repo

import (
	. "github.com/wearedevx/keystone/api/pkg/models"
)

func (r *Repo) GetRoles(roles *[]Role) IRepo {
	if r.err != nil {
		return r
	}

	r.err = r.GetDb().Find(roles).Error

	return r
}

func (repo *Repo) CreateRole(role *Role) IRepo {
	repo.err = repo.GetDb().Create(role).Error
	return repo
}

func (repo *Repo) GetRole(role *Role) IRepo {
	if repo.Err() != nil {
		return repo
	}

	repo.err = repo.GetDb().
		Where(*role).
		First(role).
		Error

	return repo
}

func (repo *Repo) GetOrCreateRole(role *Role) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDb().Where(&role).FirstOrCreate(&role).Error

	return repo
}

func (r *Repo) GetInvitableRoles(role Role, roles *[]Role) IRepo {
	r.err = r.GetDb().Model(&Role{}).
		Joins("left join roles_environment_types on roles_environment_types.role_id = roles.id").
		Where("roles_environment_types.role_id = ? and roles_environment_types.invite = true", role.ID).
		Find(roles).Error

	return r
}

func (r *Repo) GetRolesMemberCanInvite(projectMember ProjectMember, roles *[]Role) IRepo {
	if r.Err() != nil {
		return r
	}

	if projectMember.Role.CanAddMember {
		r.GetChildrenRoles(projectMember.Role, roles)
	}

	return r
}

func (r *Repo) GetChildrenRoles(role Role, roles *[]Role) IRepo {
	if r.Err() != nil {
		return r
	}

	q := roleQueue{}
	q.push(role)

	allRoles := make([]Role, 0)
	r.GetRoles(&allRoles)

	if r.Err() != nil {
		return r
	}

	var toAdd []Role

	var current Role
	for ok := true; ok; current, ok = q.pop() {
		toAdd = make([]Role, 0)

		for _, dbRole := range allRoles {
			if dbRole.ParentID == current.ID {
				toAdd = append(toAdd, dbRole)
			}
		}

		if len(toAdd) > 0 {
			*roles = append(*roles, toAdd...)
			q.pushMany(toAdd)
		}
	}

	return r
}

type roleQueue struct {
	roles []Role
}

func (r *roleQueue) push(role Role) {
	for _, existing := range r.roles {
		if existing.ID == role.ID {
			return
		}
	}

	r.roles = append(r.roles, role)
}

func (r *roleQueue) pushMany(roles []Role) {
	for _, role := range roles {
		r.push(role)
	}
}

func (r *roleQueue) pop() (Role, bool) {
	var res Role
	var found bool

	if len(r.roles) >= 1 {
		role, tail := r.roles[0], r.roles[1:]

		r.roles = tail

		res = role
		found = true
	}

	return res, found
}
