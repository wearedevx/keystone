package repo

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/wearedevx/keystone/api/pkg/models"
	"gorm.io/gorm"
)

func matchOrganizationName(name string) error {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9\.\-\_\@]{1,}$`, name)
	if !matched {
		err := errors.New("Incorrect organization name. Organization name must be alphanumeric with ., -, _, @")
		return err
	}
	return nil
}

func (r *Repo) GetOrganizationByName(orga *models.Organization) IRepo {
	if err := r.GetDb().Where("name = ?", orga.Name).First(&orga).Error; err != nil {
		r.err = err
		return r
	}
	return r
}

func (r *Repo) GetOrganizationProjects(orga *models.Organization, projects *[]models.Project) IRepo {
	fmt.Println(orga)
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
		r.err = errors.New("Organization name already taken. Choose another one.")
		return r
	}
}

func (r *Repo) UpdateOrganization(orga *models.Organization) IRepo {
	if err := matchOrganizationName(orga.Name); err != nil {
		r.err = err
		return r
	}

	foundOrga := models.Organization{}
	if err := r.GetDb().Where("name = ? and id != ?", orga.Name, orga.ID).First(&foundOrga).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			r.err = r.GetDb().Omit("paid", "user_id").Save(&orga).Error
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

	orgas := make([]models.Organization, 0)
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

	result.Organizations = orgas

	return r
}

func (r *Repo) IsUserOwnerOfOrga(user *models.User, orga *models.Organization) (isOwner bool, err error) {
	foundOrga := models.Organization{}

	if err = r.GetDb().Where("id=?", orga.ID).First(&foundOrga).Error; err != nil {
		return false, err
	}
	if foundOrga.UserID == user.ID {
		return true, nil
	}

	return isOwner, err
}

func setRoleToAdminForAllProjectsFromOrga(orga *models.Organization) error {

	return nil
}
