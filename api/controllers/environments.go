package controllers

import (
	"io"
	"net/http"

	"github.com/wearedevx/keystone/api/pkg/models"

	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func GetEnvironmentPublicKeys(params router.Params, _ io.ReadCloser, Repo repo.Repo, _ models.User) (router.Serde, int, error) {
	var status int = http.StatusOK
	var err error

	// TODO
	// Check user rights to read

	var envID = params.Get("envID").(string)

	result := models.PublicKeys{
		Keys: make([]models.UserPublicKey, 0),
	}

	Repo.GetEnvironmentPublicKeys(envID, &result)

	if Repo.Err() != nil {
		status = http.StatusInternalServerError
		err = Repo.Err()
	}

	// status = runner.Status()

	return &result, status, err
}
