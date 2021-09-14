package controllers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/wearedevx/keystone/api/internal/activitylog"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

// Returns a List of Roles
func GetDevices(params router.Params, _ io.ReadCloser, Repo repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	e := activitylog.Context{
		UserID: user.ID,
		Action: "GetDevices",
	}

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
	}

	return &result, status, e.IntoError(err)
}

func DeleteDevice(params router.Params, _ io.ReadCloser, Repo repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	status = http.StatusNoContent
	e := activitylog.Context{
		UserID: user.ID,
		Action: "DeleteDevice",
	}

	var result = &models.RemoveDeviceResponse{Success: true}

	var deviceUID = params.Get("uid").(string)

	if err = Repo.RevokeDevice(user.ID, deviceName).Err(); err != nil {
		result.Error = err.Error()
		result.Success = false
		status = http.StatusConflict
	} else {
		result = nil
	}

	return result, status, e.IntoError(err)
}
