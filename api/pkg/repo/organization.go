package repo

import (
	"github.com/wearedevx/keystone/api/pkg/models"
)

func (repo *Repo) CreateOrganization(orga *models.Organization) IRepo {
	repo.err = repo.GetDb().Create(orga).Error
	return repo
}

func (r *Repo) GetOrganizations(userID uint, result *models.GetOrganizationsResponse) IRepo {
	if r.Err() != nil {
		return r
	}

	rows, err := r.GetDb().
		Raw(`
	select o.name, o.id from organizations o
	left join projects p on p.organization_id = o.id
	left join project_members pm on pm.project_id = p.id
	left join users u on u.id = pm.user_id
	where (pm.user_id = ? and o.private = false) or o.owner_id = ?
	group by o.id, o.name
	`, userID, userID).
		Rows()

	if err != nil {
		r.err = err
		return r
	}
	var name string
	var id uint

	orgas := make([]models.Organization, 0)
	for rows.Next() {
		rows.Scan(&name, &id)

		orgas = append(orgas, models.Organization{
			ID:   id,
			Name: name,
		})
	}
	result.Organizations = orgas

	return r
}
