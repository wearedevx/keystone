package client

import (
	"errors"
	"log"

	"github.com/wearedevx/keystone/api/pkg/models"
)

type Devices struct {
	log *log.Logger
	r   requester
}

// GetAll method fetches and returns a list of all the user's devices
func (c *Devices) GetAll() ([]models.Device, error) {
	var err error
	var result models.GetDevicesResponse

	err = c.r.get("/devices", &result, nil)

	return result.Devices, err
}

// Revoke method revokes access to the device with UID `uid`
func (c *Devices) Revoke(uid string) error {
	var err error

	var result models.RemoveDeviceResponse
	err = c.r.del("/devices/"+uid, nil, &result, nil)

	if len(result.Error) > 0 {
		return errors.New(result.Error)
	}
	return err
}
