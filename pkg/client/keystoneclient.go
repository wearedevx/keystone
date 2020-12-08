package client

import (
	. "github.com/wearedevx/keystone/internal/models"
)

type SKeystoneClient struct {
	r requester
}

func NewKeystoneClient(userID string, pk string) KeystoneClient {
	return &SKeystoneClient{
		r: newRequester(userID, pk),
	}
}

func (client *SKeystoneClient) InitProject(name string) (Project, error) {
	var project Project

	payload := Project{
		Name: name,
	}

	err := client.r.post("/projects", payload, &project)

	return project, err
}
