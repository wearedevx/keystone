package repo

import (
	"errors"
	"regexp"

	"github.com/wearedevx/keystone/api/pkg/models"
	"gorm.io/gorm"
)

func (r *Repo) CreateOrganization(orga *models.Organization) IRepo {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9\.\-\_]{1,}$`, orga.Name)
	if !matched {
		r.err = errors.New("Incorrect organization name. Organization name must be alphanumeric with ., -, _")
		return r
	}

	if err := r.GetDb().Where("name = ?", orga.Name).First(&orga).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			r.err = r.GetDb().Create(orga).Error
			return r
		}
		r.err = err
		return r
	} else {
		r.err = errors.New("Organization name already taken. Choose another one.")
		return r
	}
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
