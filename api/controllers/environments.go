package controllers

import (
	"fmt"
	"io"
	"net/http"

	. "github.com/wearedevx/keystone/api/internal/utils"
	. "github.com/wearedevx/keystone/api/pkg/models"

	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func GetEnvironmentPublicKeys(params router.Params, _ io.ReadCloser, Repo repo.Repo, user User) (router.Serde, int, error) {
	var status int = http.StatusOK
	var err error

	// TODO
	// Check user rights to read

	var envID = params.Get("envID").(string)
	var result PublicKeys

	runner := NewRunner([]RunnerAction{
		NewAction(func() error {

			result = PublicKeys{
				Keys: make([]UserPublicKey, 0),
			}

			Repo.GetEnvironmentPublicKeys(envID, &result)

			return Repo.Err()
		}).SetStatusSuccess(200),
	}).Run()

	status = runner.Status()
	err = runner.Error()

	fmt.Println("api ~ environments.go ~ result", result)
	return &result, status, err
}
