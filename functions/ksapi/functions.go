// Package p contains an HTTP Cloud Function.
package ksusers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/julienschmidt/httprouter"

	"github.com/wearedevx/keystone/functions/ksapi/routes"
	log "github.com/wearedevx/keystone/internal/cloudlogger"
	"github.com/wearedevx/keystone/internal/crypto"
	. "github.com/wearedevx/keystone/internal/jwt"
	. "github.com/wearedevx/keystone/internal/models"
	"github.com/wearedevx/keystone/internal/repo"
	. "github.com/wearedevx/keystone/internal/utils"
)

// postUser Gets or Creates a user
func postUser(w http.ResponseWriter, r *http.Request, _params httprouter.Params) {
	var status int = http.StatusOK
	var responseBody bytes.Buffer
	var err error

	Repo := new(repo.Repo)
	var user *User = &User{}
	var serializedUser string

	runner := NewRunner([]RunnerAction{
		NewAction(func() error {
			Repo.Connect()

			return Repo.Err()
		}),
		NewAction(func() error {
			return user.Deserialize(r.Body)
		}).SetStatusError(http.StatusBadRequest),
		NewAction(func() error {
			Repo.GetOrCreateUser(user)

			return Repo.Err()
		}),
		NewAction(func() error {
			return user.Serialize(&serializedUser)
		}),
		NewAction(func() error {
			in := bytes.NewBufferString(serializedUser)
			_, e := crypto.EncryptForUser(user, in, &responseBody)

			return e
		}),
	})

	if err = runner.Run().Error(); err != nil {
		log.Error(r, err.Error())
		http.Error(w, err.Error(), status)
		return
	}

	status = runner.Status()

	if responseBody.Len() > 0 {
		w.Header().Add("Content-Type", "application/octet-stream")
		w.Header().Add("Content-Length", strconv.Itoa(responseBody.Len()))
		w.Write(responseBody.Bytes())
	}

	w.WriteHeader(status)
}

// getUser gets a user
func getUser(_ routes.Params, _ io.ReadCloser, _ repo.Repo, user User) (routes.Serde, int, error) {
	return &user, http.StatusOK, nil
}

func postProject(_ routes.Params, body io.ReadCloser, Repo repo.Repo, user User) (routes.Serde, int, error) {
	var status int = http.StatusOK
	var err error

	project := &Project{}
	var environment Environment

	runner := NewRunner([]RunnerAction{
		NewAction(func() error {
			return project.Deserialize(body)
		}),
		NewAction(func() error {
			Repo.GetOrCreateProject(project, user)

			environment = Repo.GetOrCreateEnvironment(*project, "default")

			return Repo.Err()
		}).SetStatusSuccess(201),
	})

	status = runner.Status()
	err = runner.Error()

	return project, status, err

}

type projectsPublicKeys struct {
	keys []UserPublicKey
}

func (p *projectsPublicKeys) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(p)
}

func (p *projectsPublicKeys) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(p)

	*out = sb.String()

	return err
}

func getProjectsPublicKeys(params routes.Params, _ io.ReadCloser, Repo repo.Repo, user User) (routes.Serde, int, error) {
	var status int = http.StatusOK
	var err error

	var project Project
	var projectID = params.Get("id").(string)
	var result projectsPublicKeys

	runner := NewRunner([]RunnerAction{
		NewAction(func() error {
			Repo.GetProjectByUUID(projectID, &project)
			Repo.ProjectLoadUsers(&project)

			// for _, user := range project.Users {
			// 	result.keys = append(result.keys, UserPublicKey{
			// 		UserID:    user.UserID,
			// 		PublicKey: user.PublicKey,
			// 	})
			// }

			return Repo.Err()
		}).SetStatusSuccess(200),
	})

	status = runner.Status()
	err = runner.Error()

	return &result, status, err
}

func postAddVariable(params routes.Params, body io.ReadCloser, Repo repo.Repo, user User) (routes.Serde, int, error) {
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

			for _, uev := range input.UserEnvValue {
				if environment, ok := Repo.GetEnvironmentByProjectIDAndName(project, uev.Environment); ok {
					if user, ok := Repo.GetUser(uev.UserID); ok {
						Repo.EnvironmentSetVariableForUser(environment, secret, user, uev.Value)
					}

				}
			}

			return Repo.Err()
		}),
	})

	err = runner.Error()
	status = runner.Status()

	return nil, status, err
}

func putSetVariable(params routes.Params, body io.ReadCloser, Repo repo.Repo, user User) (routes.Serde, int, error) {
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

// Auth shows the code to copy paste into the cli
func UserService(w http.ResponseWriter, r *http.Request) {
	router := httprouter.New()

	router.POST("/", postUser)
	router.GET("/", routes.AuthedHandler(getUser))

	router.POST("/projects", routes.AuthedHandler(postProject))

	router.GET("/projects/:id/public-keys", routes.AuthedHandler(getProjectsPublicKeys))

	router.POST("/projects/:projectID/variables", routes.AuthedHandler(postAddVariable))
	router.PUT("/projects/:projectID/:environment/variables", routes.AuthedHandler(putSetVariable))

	router.ServeHTTP(w, r)
}
