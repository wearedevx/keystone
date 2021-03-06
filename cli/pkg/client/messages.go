package client

import (
	"strconv"

	"github.com/wearedevx/keystone/api/pkg/models"
)

type Messages struct {
	r requester
}

type GenericResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func (client *Messages) GetMessages(projectID string) (models.GetMessageByEnvironmentResponse, error) {
	var err error
	var result = models.GetMessageByEnvironmentResponse{
		Environments: map[string]models.GetMessageResponse{},
	}

	err = client.r.get("/projects/"+projectID+"/messages/", &result, nil)

	return result, err
}

func (client *Messages) DeleteMessage(messageID uint) (GenericResponse, error) {
	var err error

	var result GenericResponse

	var stringMessageID = strconv.FormatUint(uint64(messageID), 10)

	err = client.r.del("/messages/"+stringMessageID, nil, &result, nil)

	return result, err
}

func (client *Messages) SendMessages(messages models.MessagesToWritePayload) (models.GetEnvironmentsResponse, error) {
	var result models.GetEnvironmentsResponse

	err := client.r.post("/messages", messages, &result, nil)

	return result, err
}
