package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"

	"github.com/wearedevx/keystone/api/internal/rights"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/internal/utils"
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

func GetMessagesFromProjectByUser(params router.Params, _ io.ReadCloser, Repo repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	status = http.StatusOK
	var projectID = params.Get("projectID").(string)
	var device = params.Get("device").(string)
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

		return &response, http.StatusInternalServerError, err
	}

	// Get publicKey by device name to send message to current user device
	publicKey := models.PublicKey{
		Device: device,
		UserID: user.ID,
	}

	if err = Repo.GetPublicKey(&publicKey).Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			response.Error = err
			return &response, http.StatusNotFound, err
		}

		return &response, http.StatusInternalServerError, err
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
		// - rights check
		can, err := rights.CanUserReadEnvironment(Repo, user.ID, project.ID, &environment)
		if err != nil {
			response.Error = err
			return &response, http.StatusInternalServerError, err
		}

		if can {
			curr := models.GetMessageResponse{}
			if err = Repo.GetMessagesForUserOnEnvironment(publicKey, environment, &curr.Message).Err(); err != nil {
				response.Error = Repo.Err()
				response.Success = false
				return &response, http.StatusBadRequest, nil
			}

			curr.Environment = environment
			result.Environments[environment.Name] = curr
		}

	}

	return &result, status, nil
}

// WriteMessages writes messages to users
func WriteMessages(_ router.Params, body io.ReadCloser, Repo repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	status = http.StatusOK
	response := &models.GetEnvironmentsResponse{}

	payload := &models.MessagesToWritePayload{}
	payload.Deserialize(body)

	for _, message := range payload.Messages {
		// - gather information for the checks
		projectMember := models.ProjectMember{
			UserID: message.RecipientID,
		}
		environment := &models.Environment{
			EnvironmentID: message.EnvironmentID,
		}

		if err = Repo.
			GetProjectMember(&projectMember).
			GetEnvironment(environment).
			Err(); err != nil {
			return response, http.StatusNotFound, err
		}

		// - check if user has rights to write on environment
		can, err := rights.CanUserWriteOnEnvironment(Repo, user.ID, environment.Project.ID, environment)

		if err != nil {
			return response, http.StatusInternalServerError, err
		}

		if !can {
			continue
		}

		// - check recipient exists with read rights.
		can, err = rights.CanUserReadEnvironment(Repo, projectMember.UserID, projectMember.ProjectID, environment)
		if err != nil {
			return response, http.StatusInternalServerError, err
		}

		if !can {
			continue
		}

		// If ok, remove potential old messages for recipient.
		if err = Repo.RemoveOldMessageForRecipient(message.PublicKeyID, message.EnvironmentID).Err(); err != nil {
			fmt.Printf("err: %+v\n", err)
			break
		}

		messageToWrite := &models.Message{
			RecipientID:   message.RecipientID,
			Payload:       message.Payload,
			EnvironmentID: message.EnvironmentID,
			SenderID:      user.ID,
			PublicKeyID:   message.PublicKeyID,
		}

		if err = Repo.WriteMessage(user, *messageToWrite).Err(); err != nil {
			fmt.Printf("err: %+v\n", err)
			break
		}

		if message.UpdateEnvironmentVersion {
			// Change environment version id.
			err = Repo.SetNewVersionID(environment)

			if err != nil {
				return response, http.StatusInternalServerError, err
			}
		}

		// Change environment version id.
		response.Environments = append(response.Environments, *environment)
	}

	return response, status, nil
}

func DeleteMessage(params router.Params, _ io.ReadCloser, Repo repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	status = http.StatusNoContent
	response := &GenericResponse{}
	response.Success = true

	var messageID = params.Get("messageID").(string)

	id, err := strconv.Atoi(messageID)

	if err != nil {
		response.Success = false
		response.Error = err
		return response, status, nil
	}

	if err = Repo.DeleteMessage(uint(id), user.ID).Err(); err != nil {
		response.Error = err
		response.Success = false
	}

	return response, status, nil
}

func getTTLToken(r *http.Request) (string, bool) {
	headers := r.Header["X-Ks-Ttl"]

	if len(headers) > 0 {
		if t := headers[0]; t != "" {
			return t, true
		}

		return "", false
	}

	return "", false
}

// Delete every message older than a week
func DeleteExpiredMessages(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Check caller with some sort of token
	token, ok := getTTLToken(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	actual := utils.GetEnv("X_KS_TTL", "")
	if actual == "" {
		fmt.Println("Missing TTL authorization key")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if token == "" || token != actual {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Actual work
	err := repo.Transaction(func(Repo repo.IRepo) error {
		Repo.DeleteExpiredMessages()

		return Repo.Err()
	})

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	http.Error(w, "OK", http.StatusOK)
}
