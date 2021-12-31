package client

import (
	"fmt"
	"log"

	"github.com/wearedevx/keystone/api/pkg/models"
)

type Roles struct {
	log *log.Logger
	r   requester
}

// GetAll method returns all the roles available in project
// identifiled by `projectID`
func (c *Roles) GetAll(projectID string) ([]models.Role, error) {
	var err error
	var result models.GetRolesResponse

	err = c.r.get(fmt.Sprintf("/roles/%s", projectID), &result, nil)

	return result.Roles, err
}
