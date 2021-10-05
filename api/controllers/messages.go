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
	uuid "github.com/satori/go.uuid"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"

	"github.com/wearedevx/keystone/api/internal/constants"
	"github.com/wearedevx/keystone/api/internal/emailer"
	"github.com/wearedevx/keystone/api/internal/redis"
	"github.com/wearedevx/keystone/api/internal/rights"
	"github.com/wearedevx/keystone/api/internal/router"
)

var Redis *redis.Redis

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
	response := GenericResponse{
		Success: false,
	}

	var projectID = params.Get("projectID").(string)
	var deviceUID = params.Get("device").(string)
	var project = models.Project{UUID: projectID}
	var publicKey = models.Device{}
	var environments []models.Environment
	var result = models.GetMessageByEnvironmentResponse{
		Environments: map[string]models.GetMessageResponse{},
	}
	var log = models.ActivityLog{
		UserID: &user.ID,
		Action: "GetMessagesFromProjectByUser",
	}

	if err = Repo.GetProject(&project).Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			response.Error = err
			status = http.StatusNotFound
			goto done
		}

		status = http.StatusInternalServerError
		goto done
	}

	log.ProjectID = &project.ID

	// Get publicKey by device name to send message to current user device
	publicKey.UID = deviceUID

	if err = Repo.GetDeviceByUserID(user.ID, &publicKey).Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			response.Error = err
			status = http.StatusNotFound
			goto done
		}

		status = http.StatusAlreadyReported
		goto done
	}

	if err = Repo.GetEnvironmentsByProjectUUID(projectID, &environments).Err(); err != nil {
		response.Error = err
		status = http.StatusBadRequest
		goto done
	}

	for _, environment := range environments {
		// - rights check
		log.Environment = environment

		can, err := rights.CanUserReadEnvironment(Repo, user.ID, project.ID, &environment)
		if err != nil {
			response.Error = err
			status = http.StatusNotFound
			goto done
		}

		if can {
			curr := models.GetMessageResponse{}
			err = Repo.GetMessagesForUserOnEnvironment(publicKey, environment, &curr.Message).Err()

			if err != nil {
				response.Error = Repo.Err()
				response.Success = false
				status = http.StatusBadRequest
				goto done
			}

			curr.Environment = environment
			curr.Message.Payload, err = Repo.MessageService().GetMessageByUuid(curr.Message.Uuid)

			if err != nil {
				// fmt.Println("api ~ messages.go ~ err", err)
				response.Error = err
				return &response, http.StatusNotFound, err
			} else {
				result.Environments[environment.Name] = curr
			}
		}
	}

done:
	return &result, status, log.SetError(err)
}

// WriteMessages writes messages to users
func WriteMessages(_ router.Params, body io.ReadCloser, Repo repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	status = http.StatusOK
	response := &models.GetEnvironmentsResponse{}
	log := models.ActivityLog{
		UserID: &user.ID,
		Action: "WriteMessages",
	}
	senderDevice := models.Device{}

	payload := &models.MessagesToWritePayload{}
	if err = payload.Deserialize(body); err != nil {
		status = http.StatusBadRequest
		response = nil
		goto done
	}

	for _, clientMessage := range payload.Messages {
		// - gather information for the checks
		projectMember := models.ProjectMember{
			UserID: clientMessage.RecipientID,
		}
		environment := models.Environment{
			EnvironmentID: clientMessage.EnvironmentID,
		}

		if err = Repo.
			GetProjectMember(&projectMember).
			GetEnvironment(&environment).
			Err(); err != nil {
			status = http.StatusNotFound
			goto done
		}

		log.ProjectID = &projectMember.ProjectID
		log.EnvironmentID = &environment.ID
		// If organization has not paid and there is non admin in the project, messages cannot be written
		has, err := rights.HasOrganizationNotPaidAndHasNonAdmin(&repo.Repo{}, environment.Project)
		if err != nil {
			return response, http.StatusInternalServerError, err
		}
		if has {
			return response, http.StatusInternalServerError, errors.New("not paid")
		}

		// - check if user has rights to write on environment
		can, err := rights.CanUserWriteOnEnvironment(Repo, user.ID, environment.Project.ID, &environment)

		if err != nil {
			status = http.StatusNotFound
			goto done
		}

		if !can {
			continue
		}

		// - check recipient exists with read rights.
		can, err = rights.CanUserReadEnvironment(Repo, projectMember.UserID, projectMember.ProjectID, &environment)
		if err != nil {
			status = http.StatusNotFound
			goto done
		}

		if !can {
			continue
		}

		// If ok, remove potential old messages for recipient.
		if err = Repo.RemoveOldMessageForRecipient(clientMessage.RecipientDeviceID, clientMessage.EnvironmentID).Err(); err != nil {
			fmt.Printf("err: %+v\n", err)
			break
		}

		senderDevice = models.Device{
			UID: clientMessage.SenderDeviceUID,
		}

		if err = Repo.GetDevice(&senderDevice).Err(); err != nil {
			if errors.Is(err, repo.ErrorNotFound) {
				status = http.StatusNotFound
			} else {
				status = http.StatusInternalServerError
			}

			goto done
		}

		messageToWrite := &models.Message{
			RecipientID:       clientMessage.RecipientID,
			Uuid:              uuid.NewV4().String(),
			EnvironmentID:     clientMessage.EnvironmentID,
			SenderID:          user.ID,
			RecipientDeviceID: clientMessage.RecipientDeviceID,
			SenderDeviceID:    senderDevice.ID,
		}

		if err = Repo.WriteMessage(user, *messageToWrite).Err(); err != nil {
			fmt.Printf("WriteMessage error: %+v\n", err)
			break
		}

		Repo.
			MessageService().
			WriteMessageWithUuid(
				messageToWrite.Uuid,
				clientMessage.Payload,
			)

		if clientMessage.UpdateEnvironmentVersion {
			// Change environment version id.
			err = Repo.SetNewVersionID(&environment)

			if err != nil {
				if errors.Is(err, repo.ErrorNotFound) {
					status = http.StatusNotFound
				} else {
					status = http.StatusInternalServerError
				}

				goto done
			}
		}

		// Change environment version id.
		response.Environments = append(response.Environments, environment)
	}

done:
	return response, status, log.SetError(err)
}

func DeleteMessage(params router.Params, _ io.ReadCloser, Repo repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	status = http.StatusNoContent
	response := &GenericResponse{}
	response.Success = true
	log := models.ActivityLog{
		UserID: &user.ID,
		Action: "DeleteMessage",
	}

	var messageID = params.Get("messageID").(string)

	id, err := strconv.Atoi(messageID)

	if err != nil {
		response.Success = false
		response.Error = err
		err = nil

		goto done
	}

	if err = Repo.
		DeleteMessage(uint(id), user.ID).Err(); err != nil {
		response.Error = err
		response.Success = false
	}

done:
	return response, status, log.SetError(err)
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

	actual := constants.KsTTLHeader
	if actual == "" {
		actual = "a very long secret header"
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

func AlertMessagesWillExpire(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Check caller with some sort of token
	token, ok := getTTLToken(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	actual := constants.KsTTLHeader
	if actual == "" {
		actual = "a very long secret header"
	}

	if token == "" || token != actual {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	errors := []error{}

	// Actual work
	err := repo.Transaction(func(Repo repo.IRepo) error {
		groupedMessageUser := make(map[uint]emailer.GroupedMessagesUser)

		Repo.GetGroupedMessagesWillExpireByUser(&groupedMessageUser)

		// For each recipients, send message.
		for _, groupedMessagesUser := range groupedMessageUser {
			email, err := emailer.MessageWillExpireMail(5, groupedMessagesUser.Projects)
			if err != nil {
				errors = append(errors, err)
			} else if err = email.Send([]string{groupedMessagesUser.Recipient.Email}); err != nil {
				errors = append(errors, err)
			}
		}

		return Repo.Err()
	})

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	} else if len(errors) > 0 {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	http.Error(w, "OK", http.StatusOK)
}

func init() {
	Redis = new(redis.Redis)
}
