package client

import (
	"fmt"
	"strconv"

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

func (c *Organizations) CreateOrganization(
	organizationName string,
	private bool,
) (models.Organization, error) {
	var err error
	var result models.Organization
	payload := models.Organization{Name: organizationName, Private: private}

	err = c.r.post("/organizations", &payload, &result, nil)

	return result, err
}

func (c *Organizations) UpdateOrganization(
	organization models.Organization,
) (models.Organization, error) {
	var err error
	var result models.Organization

	err = c.r.put("/organizations", &organization, &result, nil)

	return result, err
}

// GetUpgradeUrl returns an URL to a service
// that creates a new subscription for the named organization
// To be called to turn a non-paid organization into a paid one
func (c *Organizations) GetUpgradeUrl(
	organizationName string,
) (url string, err error) {
	var result models.StartSubscriptionResponse

	err = c.r.post(
		fmt.Sprintf("/organization/%s/upgrade", organizationName),
		map[string]string{},
		&result,
		nil,
	)

	if err != nil {
		return "", err
	}

	url = result.Url

	return url, nil
}

// GetManagementUrl returns an URL to a service
// that manages an existing subscription for the named organization
// To be called to update payment method or cancel an existing
// plan
func (c *Organizations) GetManagementUrl(
	organizationName string,
) (url string, err error) {
	var result models.ManageSubscriptionResponse

	err = c.r.post(
		fmt.Sprintf("/organization/%s/manage", organizationName),
		map[string]string{},
		&result,
		nil,
	)

	if err != nil {
		return "", err
	}

	url = result.Url

	return url, nil
}

func (c *Organizations) GetProjects(
	orga models.Organization,
) ([]models.Project, error) {
	var err error
	var result models.GetProjectsResponse
	orgaIDString := strconv.FormatUint(uint64(orga.ID), 10)

	err = c.r.get("/organizations/"+orgaIDString+"/projects", &result, nil)

	return result.Projects, err
}

func (c *Organizations) GetMembers(
	orga models.Organization,
) ([]models.ProjectMember, error) {
	var err error
	var result models.GetMembersResponse
	orgaIDString := strconv.FormatUint(uint64(orga.ID), 10)

	err = c.r.get("/organizations/"+orgaIDString+"/members", &result, nil)

	return result.Members, err
}
