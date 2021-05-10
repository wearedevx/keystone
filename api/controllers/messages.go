package controllers

import (
	"io"
	"net/http"

	. "github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"

	"github.com/wearedevx/keystone/api/internal/router"
	. "github.com/wearedevx/keystone/api/internal/utils"
)

func GetMessagesFromProjectByUser(params router.Params, _ io.ReadCloser, Repo repo.Repo, user User) (router.Serde, int, error) {
	var status = http.StatusOK
	var projectID = params.Get("projectID").(string)

	var result = GetMessageByEnvironmentResponse{
		Environments: map[string]GetMessageResponse{},
	}

	var environments []Environment
	runner := NewRunner([]RunnerAction{
		NewAction(func() error {
			environments = Repo.GetEnvironmentsByProjectUUID(projectID)
			return Repo.Err()
		}),
		NewAction(func() error {
			for _, environment := range environments {
				result.Environments[environment.Name] = GetMessageResponse{}
				curr := result.Environments[environment.Name]
				Repo.GetMessagesForUserOnEnvironment(user, environment, &curr.Message)
				curr.VersionID = environment.VersionID
				result.Environments[environment.Name] = curr
			}
			return Repo.Err()
		}),
	}).Run()

	err := runner.Error()

	return &result, status, err
}
