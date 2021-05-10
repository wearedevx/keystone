package client

import (
	"fmt"

	. "github.com/wearedevx/keystone/api/pkg/models"
)

type Messages struct {
	r requester
}

func (client *Messages) GetMessages(projectID string, environmentVersion string) (GetMessagesByEnvironmentResponse, error) {
	var err error
	var result = GetMessagesByEnvironmentResponse{
		map[string]GetMessagesResponse{},
	}

	err = client.r.get("/messages/"+projectID, &result, nil)
	fmt.Println(err)
	fmt.Println(result)

	return result, err
}
