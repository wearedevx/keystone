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
func GetRoles(_ router.Params, _ io.ReadCloser, Repo repo.Repo, _ models.User) (_ router.Serde, status int, err error) {
	status = http.StatusOK
	var result = models.GetRolesResponse{
		Roles: []models.Role{},
	}

	if err = Repo.GetRoles(&result.Roles).Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
	}

	return &result, status, err
}
