package client

import (
	"errors"

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

func (c *Devices) Revoke(name string) error {
	var err error

	var result models.RemoveDeviceResponse
	err = c.r.del("/devices/"+name, nil, &result, nil)

	if result.Success == false {
		return errors.New(result.Error)
	}
	return err
}
