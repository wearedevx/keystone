package client

import (
	"fmt"

	"github.com/wearedevx/keystone/api/pkg/models"
)

type Organizations struct {
	r requester
}

func (c *Organizations) GetAll() ([]models.Organization, error) {
	var err error
	var result models.GetOrganizationsResponse

	err = c.r.get("/organizations", &result, nil)
	fmt.Println(result)

	return result.Organizations, err
}
