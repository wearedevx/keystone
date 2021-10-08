package controllers

import (
	"errors"
	"io"
	"net/http"

	apierrors "github.com/wearedevx/keystone/api/internal/errors"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

// Returns a List of Roles
func GetDevices(_ router.Params, _ io.ReadCloser, Repo repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	log := models.ActivityLog{
		UserID: &user.ID,
		Action: "GetDevices",
	}

	status = http.StatusOK
	var result = models.GetDevicesResponse{
		Devices: []models.Device{},
	}

	if err = Repo.GetDevices(user.ID, &result.Devices).Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
			err = apierrors.ErrorNoDevice()
		} else {
			status = http.StatusInternalServerError
			err = apierrors.ErrorFailedToGetResource(err)
		}
	}

	return &result, status, log.SetError(err)
}

func DeleteDevice(params router.Params, _ io.ReadCloser, Repo repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	status = http.StatusNoContent
	log := models.ActivityLog{
		UserID: &user.ID,
		Action: "DeleteDevice",
	}

	var result = &models.RemoveDeviceResponse{Success: true}

	var deviceUID = params.Get("uid").(string)

	if err = Repo.RevokeDevice(user.ID, deviceUID).Err(); err != nil {
		result.Error = err.Error()
		result.Success = false

		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
			err = apierrors.ErrorNoDevice()
		} else {
			status = http.StatusConflict
			err = apierrors.ErrorFailedToDeleteResource(err)
		}
	} else {
		result = nil
	}

	return result, status, log.SetError(err)
}
