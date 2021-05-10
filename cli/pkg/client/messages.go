package client

import (
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

	err = client.r.get("/projects/"+projectID+"/messages/", &result, nil)

	return result, err
}
