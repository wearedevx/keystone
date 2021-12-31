package client

import (
	"fmt"
	"log"

	"github.com/wearedevx/keystone/api/pkg/apierrors"
	"github.com/wearedevx/keystone/api/pkg/models"
)

type GetOptions bool

const (
	OWNED_ONLY GetOptions = true
	ALL_KNWON  GetOptions = false
)

type Organizations struct {
	log *log.Logger
	r   requester
}

// GetAll method returns all the organizations the user is a member of
func (c *Organizations) GetAll() ([]models.Organization, error) {
	var err error
	var result models.GetOrganizationsResponse

	err = c.r.get("/organizations", &result, nil)

	return result.Organizations, err
}

// GetByName method returns the organiziation withe the name `name`.
// Set `owned` to `OWNED_ONLY` to only search among the organizations owned
// by the user. Set it to `ALL_KNWON` to search among all the organizations
// the user is a member of.
func (c *Organizations) GetByName(
	name string,
	owned GetOptions,
) (models.Organization, error) {
	var err error
	var result models.GetOrganizationsResponse
	orga := models.Organization{}

	params := map[string]string{
		"name": name,
	}

	if owned {
		params["owned"] = "1"
	}

	c.log.Printf("Getting organizations with params %v\n", params)

	err = c.r.get("/organizations", &result, params)

	if err != nil {
		return orga, err
	}

	if len(result.Organizations) == 0 {
		return orga, apierrors.ErrorFailedToGetResource
	}

	orga = result.Organizations[0]

	return orga, nil
}

// CreateOrganization method creates a new organization
func (c *Organizations) CreateOrganization(
	organizationName string,
	private bool,
) (models.Organization, error) {
	var err error
	var result models.Organization
	payload := models.Organization{Name: organizationName, Private: private}

	var s string
	payload.Serialize(&s)
	c.log.Printf("Create organization payload %s\n", s)

	err = c.r.post("/organizations", &payload, &result, nil)

	return result, err
}

// UpdateOrganization method updates (renames) an organizaiton
func (c *Organizations) UpdateOrganization(
	organization models.Organization,
) (models.Organization, error) {
	var err error
	var result models.Organization

	var s string
	organization.Serialize(&s)
	c.log.Printf("Update organization payload %s\n", s)

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

// GetProjects method returns all the projects that belong to the
// organization
func (c *Organizations) GetProjects(
	orga models.Organization,
) ([]models.Project, error) {
	var err error
	var result models.GetProjectsResponse

	err = c.r.get(
		fmt.Sprintf("/organizations/%d/projects", orga.ID),
		&result,
		nil,
	)

	return result.Projects, err
}

// GetMembers method returns all the members having access to the organization
func (c *Organizations) GetMembers(
	orga models.Organization,
) ([]models.ProjectMember, error) {
	var err error
	var result models.GetMembersResponse

	err = c.r.get(
		fmt.Sprintf("/organizations/%d/members", orga.ID),
		&result,
		nil,
	)
	return result.Members, err
}
