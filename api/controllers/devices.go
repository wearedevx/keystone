package controllers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

// Returns a List of Roles
func GetDevices(params router.Params, _ io.ReadCloser, Repo repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	status = http.StatusOK
	var result = models.GetDevicesResponse{
		Devices: []models.Device{},
	}

	if err = Repo.GetDevices(user.ID, &result.Devices).Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}

		fmt.Println(result.Devices)
		return &result, status, err
	}

	return &result, status, err
}

func DeleteDevice(params router.Params, _ io.ReadCloser, Repo repo.IRepo, u models.User) (_ router.Serde, status int, err error) {
	status = http.StatusNoContent
	var result = &models.RemoveDeviceResponse{Success: true}

	var deviceUID = params.Get("uid").(string)

	if err = Repo.RevokeDevice(u.ID, deviceUID).Err(); err != nil {
		result.Error = err.Error()
		result.Success = false
		return result, http.StatusConflict, nil
	}

	return nil, status, err
}
