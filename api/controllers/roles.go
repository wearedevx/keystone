package controllers

import (
	"io"
	"net/http"

	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/repo"
	. "github.com/wearedevx/keystone/internal/models"

	. "github.com/wearedevx/keystone/internal/utils"
)

// Returns a List of Roles
func GetRoles(params router.Params, body io.ReadCloser, Repo repo.Repo, user User) (router.Serde, int, error) {
	var status = http.StatusOK
	var result = GetRolesResponse{
		Roles: []Role{},
	}

	runner := NewRunner([]RunnerAction{
		NewAction(func() error {
			Repo.GetRoles(&result.Roles)

			return Repo.Err()
		}),
	}).Run()

	status = runner.Status()
	err := runner.Error()

	return &result, status, err
}
