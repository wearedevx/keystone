package client

import (
	"github.com/wearedevx/keystone/api/pkg/models"
)

type Logs struct {
	r requester
}

func (c *Logs) GetAll() ([]models.ActivityLogLite, error) {
	var err error
	var result models.GetActivityLogResponse

	err = c.r.get("/devices", &result, nil)

	return result.Logs, err
}
