package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	. "github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
	"gorm.io/gorm"

	"github.com/wearedevx/keystone/api/internal/router"
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
	response := &GenericResponse{}

	var result = GetMessageByEnvironmentResponse{
		Environments: map[string]GetMessageResponse{},
	}

	var environments *[]Environment
	Repo.GetEnvironmentsByProjectUUID(projectID, environments)

	if Repo.Err() != nil {
		response.Error = Repo.Err()
		response.Success = false
		return response, 400, nil
	}

	for _, environment := range *environments {
		result.Environments[environment.Name] = GetMessageResponse{}
		curr := result.Environments[environment.Name]
		Repo.GetMessagesForUserOnEnvironment(user, environment, &curr.Message)
		curr.VersionID = environment.VersionID
		result.Environments[environment.Name] = curr

		if Repo.Err() != nil {
			response.Error = Repo.Err()
			response.Success = false
			return response, 400, nil
		}
	}

	return &result, status, nil
}

func WriteMessages(params router.Params, body io.ReadCloser, Repo repo.Repo, user User) (router.Serde, int, error) {
	var status = http.StatusOK
	response := &GenericResponse{}

	var envID = params.Get("envID").(string)

	// u64, err := strconv.ParseUint(envID, 10, 64)

	// if err != nil {
	// 	fmt.Println("api ~ messages.go ~ err", err)

	// 	response.Error = err
	// 	response.Success = false
	// 	return response, 400, nil
	// }

	environment := &Environment{
		EnvironmentID: envID,
	}

	// Create transaction
	Repo.GetDb().Transaction(func(tx *gorm.DB) error {

		// TODO
		// Check if user can write on env

		payload := &repo.MessagesPayload{}
		payload.Deserialize(body)
		fmt.Println("api ~ messages.go ~ payload", payload)

		var err error

		for _, message := range payload.Messages {
			// TODO
			// - check recipient exists with read rights.

			// If ok, remove potential old messages for recipient.
			Repo.WriteMessage(user, message)

			if Repo.Err() != nil {
				err = Repo.Err()
				break
			}
		}

		if err != nil {
			fmt.Println("api ~ messages.go ~ err", err)
			return err
		}

		// Change environment version id.
		err = Repo.SetNewVersionID(environment)

		if err != nil {
			fmt.Println("api ~ messages.go ~ err", err)
			return err
		}

		// Return nil commit transaction.
		return nil
	})

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
