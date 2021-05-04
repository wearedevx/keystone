package client

import (
	"bytes"
	"fmt"

	"github.com/wearedevx/keystone/internal/crypto"
	. "github.com/wearedevx/keystone/internal/models"
)

type SKeystoneClient struct {
	r requester
}

func NewKeystoneClient(userID string, jwtToken string) KeystoneClient {
	return &SKeystoneClient{
		r: newRequester(userID, jwtToken),
	}
}

func (client *SKeystoneClient) Roles() *Roles {
	return &Roles{
		r: client.r,
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

func (client *SKeystoneClient) CheckUsersExist(userIds []string) (CheckMembersResponse, error) {
	var err error
	var result CheckMembersResponse

	payload := CheckMembersPayload{
		MemberIDs: userIds,
	}
	err = client.r.post("/users/exist", payload, &result)

	return result, err
}

func (client *SKeystoneClient) ProjectMembers(projectId string) ([]ProjectMember, error) {
	var err error
	var result GetMembersResponse

	err = client.r.get("/projects/"+projectId+"/members", &result)

	return result.Members, err
}

func (client *SKeystoneClient) ProjectAddMembers(projectId string, memberRoles map[string]Role) error {
	var result AddMembersResponse
	var err error

	payload := AddMembersPayload{
		Members: make([]MemberRole, 0),
	}

	for memberID, role := range memberRoles {
		payload.Members = append(payload.Members, MemberRole{MemberID: memberID, RoleID: role.ID})
	}

	err = client.r.post("/projects/"+projectId+"/members", payload, &result)

	if !result.Success && result.Error != "" {
		err = fmt.Errorf(result.Error)
	}

	return err
}

func (client *SKeystoneClient) ProjectRemoveMembers(projectId string, members []string) error {
	var result RemoveMembersResponse
	var err error

	payload := RemoveMembersPayload{
		Members: members,
	}

	err = client.r.del("/projects/"+projectId+"/members/", payload, &result)

	if !result.Success && result.Error != "" {
		err = fmt.Errorf(result.Error)
	}

	return err
}

func (client *SKeystoneClient) MemberSetRole(projectId string, role string) {
}

func (client *SKeystoneClient) GetUsersKeys(projectId string) ([]UserPublicKey, error) {
	var err error
	var result struct {
		keys []UserPublicKey
	}

	err = client.r.get("/projects/"+projectId+"/public-keys", &result)

	return result.keys, err
}

// Adds a variable to all environments in a project.
// It encrypts it for all users assciated with the project
// using their publick key,
// and sends that to the server
func (client *SKeystoneClient) AddVariable(projectId string, name string, valueMap map[string]string) error {
	var err error
	var payload AddVariablePayload

	uk, err := client.GetUsersKeys(projectId)

	payload.VarName = name

	for _, u := range uk {

		for environment, value := range valueMap {
			var uev struct {
				userID      string
				environment string
				value       bytes.Buffer
			}

			uev.userID = u.UserID
			uev.environment = environment

			crypto.EncryptForPublicKey(u.PublicKey, bytes.NewBufferString(value), &uev.value)

			// payload.UserEnvValue = append(payload.UserEnvValue, uev)
		}
	}

	err = client.r.post("/projects/"+projectId+"/variables", payload, nil)

	return err
}

// Updates a variable value for the specified environment.
// The variable must already exist.
// It encrypts it for all users associated with the project,
// and with reading rights on the environment
func (client *SKeystoneClient) SetVariable(projectId string, environment string, name string, value string) error {
	var err error
	var payload SetVariablePayload

	uk, err := client.GetUsersKeys(projectId)

	payload.VarName = name

	for _, u := range uk {

		var uv struct {
			userID string
			value  bytes.Buffer
		}

		uv.userID = u.UserID

		crypto.EncryptForPublicKey(u.PublicKey, bytes.NewBufferString(value), &uv.value)

		// payload.UserValue = append(payload.UserValue, uv)
	}

	err = client.r.put("/projects/"+projectId+"/"+environment+"/variables", payload, nil)

	return err
}
