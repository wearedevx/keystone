package repo

import (
	"errors"
	"regexp"

	"github.com/wearedevx/keystone/api/pkg/models"
	"gorm.io/gorm"
)

func matchOrganizationName(name string) error {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9\.\-\_\@]{1,}$`, name)
	if !matched {
		return ErrorBadName
	}
	return nil
}

func (r *Repo) GetOrganizationByName(
	userID uint,
	name string,
	orgas *[]models.Organization,
) IRepo {
	if r.err != nil {
		return r
	}

	r.err = r.GetDb().
		Preload("User").
		Joins("left join projects p on p.organization_id = organizations.id").
		Joins("left join project_members pm on pm.project_id = p.id").
		Joins("left join users u on u.id = pm.user_id").
		Where("(pm.user_id = ? and organizations.private = false) or organizations.user_id = ?", userID, userID).
		Where("organizations.name = ?", name).
		Distinct("organizations.id").
		Find(&orgas).
		Error

	if len(*orgas) == 0 {
		r.err = ErrorNotFound
	}

	return r
}

func (r *Repo) GetOwnedOrganizationByName(
	userID uint,
	name string,
	orgas *[]models.Organization,
) IRepo {
	if r.err != nil {
		return r
	}

	r.err = r.GetDb().
		Preload("User").
		Where("organizations.user_id = ?", userID).
		Where("organizations.name = ?", name).
		Find(&orgas).
		Error

	if len(*orgas) == 0 {
		r.err = gorm.ErrRecordNotFound
	}

	return r
}

func (r *Repo) GetOrganizationProjects(
	orga *models.Organization,
	projects *[]models.Project,
) IRepo {
	if err := r.GetDb().Where("organization_id = ?", orga.ID).Find(&projects).Error; err != nil {
		r.err = err
		return r
	}
	return r
}

func (r *Repo) CreateOrganization(orga *models.Organization) IRepo {
	if err := matchOrganizationName(orga.Name); err != nil {
		r.err = err
		return r
	}

	if err := r.GetDb().Where("name = ?", orga.Name).First(&orga).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			r.err = r.GetDb().Omit("paid").Create(&orga).Error
			return r
		}
		r.err = err
		return r
	} else {
		r.err = ErrorNameTaken
		return r
	}
}

func (r *Repo) GetOrganization(orga *models.Organization) IRepo {
	if r.Err() != nil {
		return r
	}

	r.err = r.GetDb().
		Where(*orga).
		Preload("User").
		First(&orga).
		Error

	return r
}

func (r *Repo) UpdateOrganization(orga *models.Organization) IRepo {
	if err := matchOrganizationName(orga.Name); err != nil {
		r.err = err
		return r
	}

	foundOrga := models.Organization{}
	if err := r.GetDb().
		Where("name = ? and id != ?", orga.Name, orga.ID).
		First(&foundOrga).
		Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			r.err = r.GetDb().Omit("paid", "user_id").Save(&orga).Error
			return r
		}
		r.err = err
		return r
	} else {
		r.err = ErrorNameTaken
		return r
	}
}

func (r *Repo) OrganizationSetCustomer(
	organization *models.Organization,
	customer string,
) IRepo {
	if r.err != nil {
		return r
	}

	if organization.CustomerID == customer {
		return r
	}

	r.err = r.GetDb().
		Model(&models.Organization{}).
		Where("id = ?", organization.ID).
		Update("customer_id", customer).
		Error

	return r
}

func (r *Repo) OrganizationSetSubscription(
	organization *models.Organization,
	subscription string,
) IRepo {
	if r.err != nil {
		return nil
	}

	r.err = r.GetDb().
		Model(&models.Organization{}).
		Where("id = ?", organization.ID).
		Update("subscription_id", subscription).
		Error

	return r
}

func (r *Repo) OrganizationSetPaid(
	organization *models.Organization,
	paid bool,
) IRepo {
	if r.err != nil {
		return r
	}

	r.err = r.GetDb().
		Model(&models.Organization{}).
		Where(organization).
		Update("paid", paid).
		Error

	return r
}

func (r *Repo) GetOrganizations(
	userID uint,
	orgas *[]models.Organization,
) IRepo {
	if r.Err() != nil {
		return r
	}

	err := r.GetDb().
		Preload("User").
		Joins("left join projects p on p.organization_id = organizations.id").
		Joins("left join project_members pm on pm.project_id = p.id").
		Joins("left join users u on u.id = pm.user_id").
		Where("(pm.user_id = ? and organizations.private = false) or organizations.user_id = ?", userID, userID).
		Group("organizations.name").Group("organizations.id").
		Find(&orgas).Error
	if err != nil {
		r.err = err
		return r
	}

	return r
}

func (r *Repo) GetOwnedOrganizations(
	userID uint,
	orgas *[]models.Organization,
) IRepo {
	if r.err != nil {
		return r
	}

	r.err = r.GetDb().
		Preload("User").
		Where("organizations.user_id = ?", userID).
		Find(&orgas).
		Error

	return r
}

func (r *Repo) OrganizationCountMembers(
	organization *models.Organization,
	count *int64,
) IRepo {
	if r.Err() != nil {
		return r
	}

	r.err = r.GetDb().
		Table("project_members").
		Joins("inner join projects on projects.id = project_id").
		Joins("inner join organizations on organizations.id = projects.organization_id").
		Where("organizations.id = ?", organization.ID).
		Select("project_members.user_id").
		Group("project_members.user_id").
		Count(count).
		Error

	return r
}

func (r *Repo) IsUserOwnerOfOrga(
	user *models.User,
	orga *models.Organization,
) (isOwner bool, err error) {
	foundOrga := models.Organization{}

	if err = r.GetDb().
		Where("user_id = ?", user.ID).
		Where("id=?", orga.ID).
		First(&foundOrga).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func setRoleToAdminForAllProjectsFromOrga(orga *models.Organization) error {
	return nil
}

func (r *Repo) GetOrganizationMembers(
	orgaID uint,
	result *[]models.ProjectMember,
) IRepo {
	if r.Err() != nil {
		return r
	}

	err := r.GetDb().
		Preload("User").
		Preload("Role").
		Joins("left join projects p on p.id = project_members.project_id").
		Joins("left join organizations o on o.id = p.organization_id").
		Where("o.id = ?", orgaID).
		Group("project_members.user_id").Group("project_members.id").
		Find(result).Error
	if err != nil {
		r.err = err
		return r
	}

	return r
}
