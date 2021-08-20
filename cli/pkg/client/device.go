package client

import (
	"github.com/wearedevx/keystone/api/pkg/models"
)

type Devices struct {
	r requester
}

func (c *Devices) GetAll() ([]models.PublicKey, error) {
	var err error
	var result models.GetDevicesResponse

	err = c.r.get("/devices", &result, nil)

	return result.PublicKeys, err
}
