package controllers

import (
	"errors"
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
		PublicKeys: []models.PublicKey{},
	}

	if err = Repo.GetPublicKeys(user.ID, &result.PublicKeys).Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}

		return &result, status, err
	}

	return &result, status, err
}
