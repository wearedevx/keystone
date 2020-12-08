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

// Initilize a project with `name` and a "default" environment
func (client *SKeystoneClient) InitProject(name string) (Project, error) {
	var project Project

	payload := Project{
		Name: name,
	}

	err := client.r.post("/projects", payload, &project)

	return project, err
}

type UserPublicKey struct {
	UserID    string `json:"user"`
	PublicKey string `json:"public_key"`
}

func (client *SKeystoneClient) GetUsersKeys(projectId string) ([]UserPublicKey, error) {
	var err error
	var result struct {
		keys []UserPublicKey
	}

	err = client.r.get("/projects/public-keys", &result)

	return result.keys, err
}

// Adds a variable to all environments in a project.
// It encrypts it for all users assciated with the project
// using their publick key,
// and sends that to the server
func (client *SKeystoneClient) AddVariable(projectId string, name string, valueMap map[string]string) error {

	return nil
}

// Updates a variable value for the specified environment.
// The variable must already exist.
// It encrypts it for all users associated with the project,
// and with reading rights on the environment
func (client *SKeystoneClient) SetVariable(projectId string, environment string, name string, value string) error {

	return nil
}
