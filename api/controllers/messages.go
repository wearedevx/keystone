package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	. "github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"

	"github.com/wearedevx/keystone/api/internal/router"
	. "github.com/wearedevx/keystone/api/internal/utils"
)

type GenericResponse struct {
	Success bool  `json:"success"`
	Error   error `json:"error"`
}

func (gr *GenericResponse) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(gr)
}

func (gr *GenericResponse) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(gr)

	*out = sb.String()

	return err
}

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

func WriteMessages(params router.Params, body io.ReadCloser, Repo repo.Repo, user User) (router.Serde, int, error) {
	var status = http.StatusOK
	response := &GenericResponse{}

	// Check if user is a member of project
	// var projectID = params.Get("projectID").(string)

	payload := &repo.MessagesPayload{}
	payload.Deserialize(body)

	for _, message := range payload.Messages {
		Repo.WriteMessage(user, message)
	}

	// For each message,
	// - check if user can write to environment
	// - check recipient exist
	// - If yes, write message

	return response, status, nil
}

func DeleteMessage(params router.Params, body io.ReadCloser, Repo repo.Repo, user User) (router.Serde, int, error) {

	var status = http.StatusNoContent
	response := &GenericResponse{}
	response.Success = true

	var messageID = params.Get("messageID").(string)

	id, err := strconv.Atoi(messageID)

	if err != nil {
		response.Success = false
		response.Error = err
		return response, status, nil
	}

	err = Repo.DeleteMessage(id)

	if err != nil {
		response.Error = err
		response.Success = false
	}

	return response, status, nil
}
