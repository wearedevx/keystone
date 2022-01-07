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
	"github.com/wearedevx/keystone/api/pkg/notification"
	"github.com/wearedevx/keystone/api/pkg/repo"

	"github.com/wearedevx/keystone/api/internal/constants"
	"github.com/wearedevx/keystone/api/internal/emailer"
	apierrors "github.com/wearedevx/keystone/api/internal/errors"
	"github.com/wearedevx/keystone/api/internal/redis"
	"github.com/wearedevx/keystone/api/internal/rights"
	"github.com/wearedevx/keystone/api/internal/router"
)

var Redis *redis.Redis

type GenericResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
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

func GetMessagesFromProjectByUser(
	params router.Params,
	_ io.ReadCloser,
	Repo repo.IRepo,
	user models.User,
) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	projectID := params.Get("projectID")
	deviceUID := params.Get("device")
	project := models.Project{UUID: projectID}
	device := models.Device{}
	var environments []models.Environment
	result := models.GetMessageByEnvironmentResponse{
		Environments: map[string]models.GetMessageResponse{},
	}
	log := models.ActivityLog{
		UserID: &user.ID,
		Action: "GetMessagesFromProjectByUser",
	}

	// Get publicKey by device name to send message to current user device
	device.UID = deviceUID

	if err = Repo.
		GetProject(&project).
		GetDeviceByUserID(user.ID, &device).
		GetEnvironmentsByProjectUUID(projectID, &environments).
		Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
			goto done
		}

		status = http.StatusBadRequest
		err = apierrors.ErrorBadRequest(err)
		goto done
	}

	log.ProjectID = &project.ID

	for _, environment := range environments {
		// - rights check
		var can bool
		log.Environment = environment

		can, err = rights.
			CanUserReadEnvironment(
				Repo,
				user.ID,
				project.ID,
				&environment,
			)
		if err != nil {
			status = http.StatusNotFound
			err = apierrors.ErrorFailedToGetPermission(err)
			goto done
		}

		if can {
			curr := models.GetMessageResponse{}
			if err = Repo.
				GetMessagesForUserOnEnvironment(
					device,
					environment,
					&curr.Message,
				).Err(); err != nil {
				status = http.StatusBadRequest
				err = apierrors.ErrorFailedToGetResource(err)
				goto done
			}

			curr.Environment = environment
			curr.Message.Payload, err = Repo.
				MessageService().
				GetMessageByUuid(curr.Message.Uuid)

			if err != nil {
				fmt.Printf(
					"Error getting message payload from redis: %+v\n",
					err,
				)
				// If there is no matching message, the database entry
				// should be deleted
				if err = Repo.
					DeleteMessage(
						curr.Message.ID,
						curr.Message.RecipientID,
					).
					Err(); err != nil {
					fmt.Printf(
						"Error deleting message without a payload: %+v\n",
						err,
					)
					Repo.ClearErr()
				}
			} else {
				result.Environments[environment.Name] = curr
			}
		}
	}

done:
	return &result, status, log.SetError(err)
}

// WriteMessages writes messages to users
func WriteMessages(
	_ router.Params,
	body io.ReadCloser,
	Repo repo.IRepo,
	user models.User,
) (_ router.Serde, status int, err error) {
	status = http.StatusOK
	response := &models.GetEnvironmentsResponse{}
	log := models.ActivityLog{
		UserID: &user.ID,
		Action: "WriteMessages",
	}
	senderDevice := models.Device{}
	var has bool
	var canRead, canWrite bool
	var errCanRead, errCanWrite error

	payload := &models.MessagesToWritePayload{}
	if err = payload.Deserialize(body); err != nil {
		status = http.StatusBadRequest
		response = nil
		err = apierrors.ErrorBadRequest(err)
		goto done
	}

	for _, clientMessage := range payload.Messages {
		if len(clientMessage.Payload) == 0 {
			status = http.StatusBadRequest
			err = apierrors.ErrorEmptyPayload()
			response = nil
			goto done
		}
		environment := models.Environment{
			EnvironmentID: clientMessage.EnvironmentID,
		}
		if err = Repo.
			GetEnvironment(&environment).
			Err(); err != nil {
			status = http.StatusNotFound
			response = nil
			goto done
		}

		// - gather information for the checks
		projectMember := models.ProjectMember{
			UserID:    clientMessage.RecipientID,
			ProjectID: environment.ProjectID,
		}

		if err = Repo.
			GetProjectMember(&projectMember).
			Err(); err != nil {
			status = http.StatusNotFound
			response = nil
			goto done
		}

		log.ProjectID = &projectMember.ProjectID
		log.EnvironmentID = &environment.ID

		// If organization has not paid and there is no admin in the project,
		// messages cannot be written
		has, err = rights.
			HasOrganizationNotPaidAndHasNonAdmin(
				Repo,
				environment.Project,
			)

		if err != nil {
			status = http.StatusInternalServerError
			err = apierrors.ErrorOrganizationWithoutAnAdmin(err)

			goto done
		}

		if has {
			status = http.StatusForbidden
			err = apierrors.ErrorNeedsUpgrade()

			goto done
		}

		// - check if user has rights to write on environment
		canWrite, errCanWrite = rights.
			CanUserWriteOnEnvironment(
				Repo,
				user.ID,
				environment.Project.ID,
				&environment,
			)

		// - check recipient exists with read rights.
		canRead, errCanRead = rights.
			CanUserReadEnvironment(Repo,
				projectMember.UserID,
				projectMember.ProjectID,
				&environment,
			)

		if errCanWrite != nil || errCanRead != nil {
			status = http.StatusNotFound
			err = apierrors.ErrorFailedToGetPermission(err)
			goto done
		}

		if !canRead || !canWrite {
			continue
		}

		// If ok, remove potential old messages for recipient.
		if err = Repo.
			RemoveOldMessageForRecipient(
				clientMessage.RecipientDeviceID,
				clientMessage.EnvironmentID,
			).
			Err(); err != nil {
			status = http.StatusInternalServerError
			err = apierrors.ErrorFailedToDeleteResource(err)
			goto done
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
			err = apierrors.ErrorNoDevice()
			response = nil

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
			status = http.StatusInternalServerError
			err = apierrors.ErrorFailedToWriteMessage(err)
			response = nil
			break
		}

		if err = Repo.
			MessageService().
			WriteMessageWithUuid(
				messageToWrite.Uuid,
				clientMessage.Payload,
			); err != nil {
			status = http.StatusInternalServerError
			err = apierrors.ErrorFailedToWriteMessage(err)
			response = nil

			goto done
		}

		if clientMessage.UpdateEnvironmentVersion {
			// Change environment version id.
			err = Repo.SetNewVersionID(&environment)

			if err != nil {
				if errors.Is(err, repo.ErrorNotFound) {
					status = http.StatusNotFound
				} else {
					status = http.StatusInternalServerError
					err = apierrors.ErrorFailedToSetEnvironmentVersion(err)
				}
				response = nil

				goto done
			}
		}

		// Change environment version id.
		response.Environments = append(
			response.Environments,
			environment,
		)
	}

done:
	return response, status, log.SetError(err)
}

func DeleteMessage(
	params router.Params,
	_ io.ReadCloser,
	Repo repo.IRepo,
	user models.User,
) (_ router.Serde, status int, err error) {
	status = http.StatusNoContent
	response := &GenericResponse{}
	var message models.Message
	response.Success = true
	log := models.ActivityLog{
		UserID: &user.ID,
		Action: "DeleteMessage",
	}

	messageID := params.Get("messageID")

	id, err := strconv.ParseUint(messageID, 10, 64)
	if err != nil {
		status = http.StatusBadRequest
		response.Success = false
		response.Error = err.Error()

		goto done
	}

	message.ID = uint(id)

	if err = Repo.
		GetMessage(&message).
		DeleteMessage(uint(id), user.ID).Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
			goto done
		}
		response.Error = err.Error()
		response.Success = false
		err = apierrors.ErrorFailedToDeleteResource(err)
		status = http.StatusInternalServerError

		goto done
	}

	if err = Repo.
		MessageService().
		DeleteMessageWithUuid(message.Uuid); err != nil {
		fmt.Printf("Error deleting message on redis: %+v\n", err)
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
func DeleteExpiredMessages(
	w http.ResponseWriter,
	r *http.Request,
	_ httprouter.Params,
) {
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
		fmt.Printf("[ERROR] Removing expired messages: %+v\n", err)
		http.Error(
			w,
			"Internal Server Error",
			http.StatusInternalServerError,
		)
		return
	}

	http.Error(w, "OK", http.StatusOK)
}

func AlertMessagesWillExpire(
	w http.ResponseWriter,
	r *http.Request,
	_ httprouter.Params,
) {
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
		groupedMessageUser := make(
			map[uint]emailer.GroupedMessagesUser,
		)

		Repo.GetGroupedMessagesWillExpireByUser(&groupedMessageUser)

		notification.SendExpireMessageToUsers(groupedMessageUser, &errors)

		return Repo.Err()
	})

	if err != nil {
		fmt.Printf("[ERROR] Sending Expiration Emails: %+v\n", err)
		http.Error(
			w,
			"Internal Server Error",
			http.StatusInternalServerError,
		)
		return
	} else if len(errors) > 0 {
		fmt.Printf("[ERROR] Sending Expiration Emails: %+v\n", errors)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Error(w, "OK", http.StatusOK)
}

func init() {
	Redis = new(redis.Redis)
}
