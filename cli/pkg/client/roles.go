package client

import (
	"fmt"

	"github.com/wearedevx/keystone/api/pkg/models"
)

type Roles struct {
	r requester
}

func (c *Roles) GetAll(projectID string) ([]models.Role, error) {
	var err error
	var result models.GetRolesResponse

	err = c.r.get(fmt.Sprintf("/roles/%s", projectID), &result, nil)

	return result.Roles, err
}
