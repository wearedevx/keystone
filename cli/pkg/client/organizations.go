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

func (c *Organizations) CreateOrganization(organizationName string, private bool) (models.Organization, error) {
	var err error
	var result models.Organization
	payload := models.Organization{Name: organizationName, Private: private}

	err = c.r.post("/organizations", &payload, &result, nil)

	return result, err
}

func (c *Organizations) UpdateOrganization(organization models.Organization) (models.Organization, error) {
	var err error
	var result models.Organization

	err = c.r.put("/organizations", &organization, &result, nil)

	return result, err
}
