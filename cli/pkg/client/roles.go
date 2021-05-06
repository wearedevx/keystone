package client

import (
	"github.com/wearedevx/keystone/api/pkg/models"
)

type Roles struct {
	r requester
}

func (c *Roles) GetAll() ([]models.Role, error) {
	var err error
	var result models.GetRolesResponse

	err = c.r.get("/roles", &result)

	return result.Roles, err
}
