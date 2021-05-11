package client

import (
	"fmt"
	"strconv"

	. "github.com/wearedevx/keystone/api/pkg/models"
)

type Messages struct {
	r requester
}

type GenericResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func (client *Messages) GetMessages(projectID string) (GetMessageByEnvironmentResponse, error) {
	var err error
	var result = GetMessageByEnvironmentResponse{
		map[string]GetMessageResponse{},
	}

	err = client.r.get("/projects/"+projectID+"/messages/", &result, nil)
	fmt.Println(err)

	return result, err
}

func (client *Messages) DeleteMessage(messageID uint) (GenericResponse, error) {
	var err error

	var result GenericResponse

	var stringMessageID = strconv.FormatUint(uint64(messageID), 10)

	err = client.r.del("/messages/"+stringMessageID, nil, &result, nil)

	return result, err
}
