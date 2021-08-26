package controllers

import (
	"errors"
	"io"
	"net/http"

	"github.com/wearedevx/keystone/api/pkg/models"

	"github.com/wearedevx/keystone/api/internal/rights"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func GetEnvironmentPublicKeys(params router.Params, _ io.ReadCloser, Repo repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	status = http.StatusOK
	result := models.PublicKeys{
		Keys: make([]models.UserPublicKeys, 0),
	}

	var envID = params.Get("envID").(string)

	// - fetch the environment to check rights
	environment := models.Environment{EnvironmentID: envID}
	if err = Repo.GetEnvironment(&environment).
		Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}

		return &result, status, err
	}

	// - check user has access to that environment
	can, err := rights.CanUserReadEnvironment(Repo, user.ID, environment.ProjectID, &environment)
	if err != nil {
		return &result, http.StatusInternalServerError, err
	}

	if !can {
		return &result, http.StatusForbidden, err
	}

	// - do the work
	if err = Repo.GetEnvironmentPublicKeys(envID, &result).
		Err(); err != nil {
		return &result, http.StatusInternalServerError, err
	}

	return &result, status, err
}
