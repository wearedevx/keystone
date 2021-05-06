package controllers

import (
	"io"
	"net/http"

	. "github.com/wearedevx/keystone/api/internal/utils"
	. "github.com/wearedevx/keystone/api/pkg/models"

	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func PostAddVariable(params router.Params, body io.ReadCloser, Repo repo.Repo, user User) (router.Serde, int, error) {
	projectID := params.Get("projectID").(string)

	var status int = http.StatusOK
	var err error

	var project Project
	input := AddVariablePayload{}
	err = input.Deserialize(body)

	if err != nil {
		return nil, 400, err
	}

	runner := NewRunner([]RunnerAction{
		NewAction(func() error {
			var secret Secret

			Repo.GetProjectByUUID(projectID, &project)
			Repo.GetSecretByName(input.VarName, &secret)

			// for _, uev := range input.UserEnvValue {
			// 	if environment, ok := Repo.GetEnvironmentByProjectIDAndName(project, uev.Environment); ok {
			// 		if user, ok := Repo.GetUser(uev.UserID); ok {
			// 			// Repo.EnvironmentSetVariableForUser(environment, secret, user, uev.Value)
			// 		}

			// 	}
			// }

			return Repo.Err()
		}),
	})

	err = runner.Error()
	status = runner.Status()

	return nil, status, err
}

func PutSetVariable(params router.Params, body io.ReadCloser, Repo repo.Repo, user User) (router.Serde, int, error) {
	projectID := params.Get("projectID").(string)
	// environmentName := params.Get("environment").(string)

	var status = http.StatusOK
	var project Project
	input := SetVariablePayload{}
	input.Deserialize(body)

	runner := NewRunner([]RunnerAction{
		NewAction(func() error {
			var secret Secret
			Repo.GetProjectByUUID(projectID, &project)

			Repo.GetSecretByName(input.VarName, &secret)

			// for _, uv := range body.UserValue {
			// 	if environment, ok := Repo.GetEnvironmentByProjectIDAndName(project, environmentName); ok {
			// 		if user, ok := Repo.GetUser(uv.UserID); ok {
			// 			Repo.EnvironmentSetVariableForUser(environment, secret, user, uv.Value)
			// 		}
			// 	}
			// }

			return Repo.Err()
		}),
	})

	status = runner.Status()
	err := runner.Error()

	return nil, status, err
}
