package controllers

import (
	"errors"
	"io"
	"net/http"

	"github.com/wearedevx/keystone/api/pkg/models"

	apierrors "github.com/wearedevx/keystone/api/internal/errors"
	"github.com/wearedevx/keystone/api/internal/rights"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func GetEnvironmentPublicKeys(
	params router.Params,
	_ io.ReadCloser,
	Repo repo.IRepo,
	user models.User,
) (_ router.Serde, status int, err error) {
	status = http.StatusOK
	result := models.PublicKeys{
		Keys: make([]models.UserDevices, 0),
	}

	envID := params.Get("envID")
	var can bool

	// - fetch the environment to check rights
	environment := models.Environment{EnvironmentID: envID}

	log := models.ActivityLog{
		UserID: &user.ID,
		Action: "GetEnvironmentPublicKeys",
	}

	if err = Repo.GetEnvironment(&environment).
		Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
			err = apierrors.ErrorFailedToGetResource(err)
		}

		goto done
	}

	log.EnvironmentID = &environment.ID

	// - check user has access to that environment
	can, err = rights.CanUserReadEnvironment(
		Repo,
		user.ID,
		environment.ProjectID,
		&environment,
	)
	if err != nil {
		status = http.StatusInternalServerError
		goto done
	}

	if !can {
		err = apierrors.ErrorPermissionDenied()
		status = http.StatusForbidden
		goto done
	}

	// - do the work
	if err = Repo.GetEnvironmentPublicKeys(envID, &result).
		Err(); err != nil {
		status = http.StatusInternalServerError
		err = apierrors.ErrorFailedToGetResource(err)
		goto done
	}

done:
	return &result, status, log.SetError(err)
}
