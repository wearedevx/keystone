package client

import (
	"fmt"

	. "github.com/wearedevx/keystone/api/pkg/models"
)

type Messages struct {
	r requester
}

func (client *Messages) GetMessages(projectID string) (GetMessageByEnvironmentResponse, error) {
	var err error
	var result = GetMessageByEnvironmentResponse{
		map[string]GetMessageResponse{},
	}

	err = client.r.get("/messages/"+projectID, &result, nil)
	fmt.Println(err)
	fmt.Println(result)

	return result, err
}

func (client *Messages) SaveMessages(MessageByEnvironments GetMessageByEnvironmentResponse) (GetMessageByEnvironmentResponse, error) {
	for _, environment := range MessageByEnvironments.Environments {
		fmt.Println(environment)
	}

	return MessageByEnvironments, nil
}
