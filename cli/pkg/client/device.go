package client

import (
	"errors"

	"github.com/wearedevx/keystone/api/pkg/models"
)

type Devices struct {
	r requester
}

func (c *Devices) GetAll() ([]models.Device, error) {
	var err error
	var result models.GetDevicesResponse

	err = c.r.get("/devices", &result, nil)

	return result.Devices, err
}

func (c *Devices) Revoke(uid string) error {
	var err error

	var result models.RemoveDeviceResponse
	err = c.r.del("/devices/"+uid, nil, &result, nil)

	if len(result.Error) > 0 {
		return errors.New(result.Error)
	}
	return err
}
