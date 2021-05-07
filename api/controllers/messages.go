package controllers

import (
	"io"
	"net/http"

	"github.com/wearedevx/keystone/api/routes"
	_ "github.com/wearedevx/keystone/api/routes"
	. "github.com/wearedevx/keystone/internal/models"
	"github.com/wearedevx/keystone/internal/repo"

	. "github.com/wearedevx/keystone/internal/utils"
)

func GetMessagesFromEnvironmentByUser(params routes.Params, _ io.ReadCloser, Repo repo.Repo, user User) (routes.Serde, int, error) {
	var status = http.StatusOK
	var environmentID = params.Get("environmentID").(string)
	// var versionID = params.Get("versionid").(string)

	var result = GetMessagesResponse{
		Messages: []Message{},
	}

	var environment Environment
	runner := NewRunner([]RunnerAction{
		NewAction(func() error {
			environment = Repo.GetEnvironment(environmentID)
			return Repo.Err()
		}),
		NewAction(func() error {
			Repo.GetMessagesForUserOnEnvironment(&user, &environment, &result.Messages)

			return Repo.Err()
		}),
	}).Run()

	status = runner.Status()
	err := runner.Error()

	return &result, status, err
}
