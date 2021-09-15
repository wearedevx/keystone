package client

import (
	"github.com/wearedevx/keystone/api/pkg/models"
)

type Organizations struct {
	r requester
}

func (c *Organizations) GetAll() ([]models.Organization, error) {
	var err error
	var result models.GetOrganizationsResponse

	err = c.r.get("/organizations", &result, nil)

	return result.Organizations, err
}

func (c *Organizations) CreateOrganization(organizationName string) (models.Organization, error) {
	var err error
	var result models.Organization
	payload := models.Organization{Name: organizationName}

	err = c.r.post("/organizations", &payload, &result, nil)

	return result, err
}
