package repo

import (
	"github.com/wearedevx/keystone/api/pkg/models"
)

func (r *Repo) GetRoles(roles *[]models.Role) IRepo {
	if r.err != nil {
		return r
	}

	r.err = r.GetDb().Find(roles).Error

	return r
}

func (repo *Repo) CreateRole(role *models.Role) IRepo {
	repo.err = repo.GetDb().Create(role).Error
	return repo
}

func (repo *Repo) GetRole(role *models.Role) IRepo {
	if repo.Err() != nil {
		return repo
	}

	repo.err = repo.GetDb().
		Where(*role).
		First(role).
		Error

	return repo
}

func (repo *Repo) GetOrCreateRole(role *models.Role) IRepo {
	if repo.err != nil {
		return repo
	}

	repo.err = repo.GetDb().Where(&role).FirstOrCreate(&role).Error

	return repo
}

func (r *Repo) GetInvitableRoles(role models.Role, roles *[]models.Role) IRepo {
	r.err = r.GetDb().Model(&models.Role{}).
		Joins("left join roles_environment_types on roles_environment_types.role_id = roles.id").
		Where("roles_environment_types.role_id = ? and roles_environment_types.invite = true", role.ID).
		Find(roles).Error

	return r
}

func (r *Repo) GetRolesMemberCanInvite(
	projectMember models.ProjectMember,
	roles *[]models.Role,
) IRepo {
	if r.Err() != nil {
		return r
	}

	if projectMember.Role.CanAddMember {
		r.GetChildrenRoles(projectMember.Role, roles)
	}

	return r
}

func (r *Repo) GetChildrenRoles(role models.Role, roles *[]models.Role) IRepo {
	if r.Err() != nil {
		return r
	}

	q := roleQueue{}
	q.push(role)

	allRoles := make([]models.Role, 0)
	r.GetRoles(&allRoles)

	if r.Err() != nil {
		return r
	}

	var toAdd []models.Role

	var current models.Role
	for ok := true; ok; current, ok = q.pop() {
		// FIXME: how come we can have empty structs here?
		if current.ID == 0 {
			continue
		}
		toAdd = make([]models.Role, 0)

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
	roles []models.Role
}

func (r *roleQueue) push(role models.Role) {
	for _, existing := range r.roles {
		if existing.ID == role.ID {
			return
		}
	}

	r.roles = append(r.roles, role)
}

func (r *roleQueue) pushMany(roles []models.Role) {
	for _, role := range roles {
		r.push(role)
	}
}

func (r *roleQueue) pop() (models.Role, bool) {
	var res models.Role
	var found bool

	if len(r.roles) >= 1 {
		role, tail := r.roles[0], r.roles[1:]

		r.roles = tail

		res = role
		found = true
	}

	return res, found
}
