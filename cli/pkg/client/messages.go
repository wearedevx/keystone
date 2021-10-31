package client

import (
	"strconv"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/config"
)

type Messages struct {
	r requester
}

type GenericResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// GetMessages method gets all messages for the users device, regarding
// the current project
func (client *Messages) GetMessages(
	projectID string,
) (models.GetMessageByEnvironmentResponse, error) {
	var err error
	result := models.GetMessageByEnvironmentResponse{
		Environments: map[string]models.GetMessageResponse{},
	}
	device := config.GetDeviceUID()
	err = client.r.get("/projects/"+projectID+"/messages/"+device, &result, nil)

	return result, err
}

// DeleteMessage method deletes the message
func (client *Messages) DeleteMessage(messageID uint) (GenericResponse, error) {
	var err error

	var result GenericResponse

	stringMessageID := strconv.FormatUint(uint64(messageID), 10)

	err = client.r.del("/messages/"+stringMessageID, nil, &result, nil)

	return result, err
}

// SendMessages method sends messages to members of the project
func (client *Messages) SendMessages(
	messages models.MessagesToWritePayload,
) (models.GetEnvironmentsResponse, error) {
	var result models.GetEnvironmentsResponse

	err := client.r.post("/messages", messages, &result, nil)

	return result, err
}
