package client

import (
	. "github.com/wearedevx/keystone/internal/models"
)

type Roles struct {
	r requester
}

func (c *Roles) GetAll() ([]Role, error) {
	var err error
	var result GetRolesResponse

	err = c.r.get("/roles", &result, nil)

	return result.Roles, err
}
