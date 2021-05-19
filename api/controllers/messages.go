package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
	"gorm.io/gorm"

	"github.com/wearedevx/keystone/api/internal/rights"
	"github.com/wearedevx/keystone/api/internal/router"
)

type GenericResponse struct {
	Success bool  `json:"success"`
	Error   error `json:"error"`
}

func (gr *GenericResponse) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(gr)
}

func (gr *GenericResponse) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(gr)

	*out = sb.String()

	return err
}

func GetMessagesFromProjectByUser(params router.Params, _ io.ReadCloser, Repo repo.Repo, user models.User) (_ router.Serde, status int, err error) {
	status = http.StatusOK
	var projectID = params.Get("projectID").(string)
	response := GenericResponse{
		Success: false,
	}

	project := models.Project{
		UUID: projectID,
	}
	if err = Repo.GetProject(&project).Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			response.Error = err
			return &response, http.StatusNotFound, err
		}

		return &response, http.StatusInsufficientStorage, err
	}

	var result = models.GetMessageByEnvironmentResponse{
		Environments: map[string]models.GetMessageResponse{},
	}

	var environments []models.Environment
	if err = Repo.GetEnvironmentsByProjectUUID(projectID, &environments).Err(); err != nil {
		response.Error = err
		return &response, http.StatusBadRequest, err
	}

	for _, environment := range environments {
		result.Environments[environment.Name] = models.GetMessageResponse{}
		curr := result.Environments[environment.Name]

		// - rights check
		can, err := rights.CanUserReadEnvironment(&Repo, user.ID, project.ID, &environment)
		if err != nil {
			response.Error = err
			return &response, http.StatusInternalServerError, err
		}

		if can {
			if err = Repo.GetMessagesForUserOnEnvironment(user, environment, &curr.Message).Err(); err != nil {
				response.Error = Repo.Err()
				response.Success = false
				return &response, http.StatusBadRequest, nil
			}
		}

		curr.VersionID = environment.VersionID
		result.Environments[environment.Name] = curr
	}

	return &result, status, nil
}

// WriteMessages writes messages to users
// TODO: on the client side, each message should be associated with the target EnvironmentID,
// 		 and therefore, there is no need to pass envID in the url, query or body in the HTTP query
func WriteMessages(params router.Params, body io.ReadCloser, Repo repo.Repo, user models.User) (_ router.Serde, status int, err error) {
	status = http.StatusOK
	response := &GenericResponse{}

	// Create transaction
	// TODO: @kévin ? Qu’est-ce qu’on fait du `tx` ?
	Repo.GetDb().Transaction(func(tx *gorm.DB) (err error) {
		payload := &repo.MessagesPayload{}
		payload.Deserialize(body)

		for _, message := range payload.Messages {
			// - gather information for the checks
			projectMember := models.ProjectMember{
				ID: message.RecipientID,
			}
			environment := models.Environment{
				EnvironmentID: message.EnvironmentID,
			}

			err = Repo.
				GetProjectMember(&projectMember).
				GetEnvironment(&environment).
				Err()
			if err != nil {
				break
			}

			// - check if user has rights to write on environment
			can, err := rights.CanUserWriteOnEnvironment(&Repo, user.ID, environment.Project.ID, &environment)

			if err != nil {
				return err
			}

			if !can {
				continue
			}

			// - check recipient exists with read rights.
			can, err = rights.CanUserReadEnvironment(&Repo, projectMember.UserID, projectMember.ProjectID, &environment)
			if err != nil {
				return err
			}

			if !can {
				continue
			}

			// If ok, remove potential old messages for recipient.
			if err = Repo.WriteMessage(user, message).Err(); err != nil {
				break
			}

			// Change environment version id.
			err = Repo.SetNewVersionID(&environment)

			if err != nil {
				fmt.Println("api ~ messages.go ~ err", err)
				return err
			}
		}

		if err != nil {
			fmt.Println("api ~ messages.go ~ err", err)
			return err
		}

		// Return nil commit transaction.
		return nil
	})

	return response, status, nil
}

func DeleteMessage(params router.Params, _ io.ReadCloser, Repo repo.Repo, user models.User) (router.Serde, int, error) {

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

	err = Repo.DeleteMessage(id, user.ID)

	if err != nil {
		response.Error = err
		response.Success = false
	}

	return response, status, nil
}
