package client

import (
	"github.com/wearedevx/keystone/cli/internal/config"
	"github.com/wearedevx/keystone/cli/internal/errors"
)

type KeystoneClientImpl struct {
	r requester
}

// NewKeystoneClient function returns a new instance of KeysotneClient
func NewKeystoneClient() (KeystoneClient, *errors.Error) {
	account, index := config.GetCurrentAccount()
	token := config.GetAuthToken()

	if index < 0 {
		return nil, errors.MustBeLoggedIn(nil)
	}

	return &KeystoneClientImpl{
		r: newRequester(account.UserID, token),
	}, nil
}

// Roles method returns an interface to interact with roles
func (client *KeystoneClientImpl) Roles() *Roles {
	return &Roles{
		r: client.r,
	}
}

// Project method returns an interface to interact with projects
func (client *KeystoneClientImpl) Project(projectId string) *Project {
	return &Project{
		r:  client.r,
		id: projectId,
	}
}

// Users method returns an interface to interact with users/members
func (client *KeystoneClientImpl) Users() *Users {
	return &Users{
		r: client.r,
	}
}

// Messages method returns an interface to interact with messages
func (client *KeystoneClientImpl) Messages() *Messages {
	return &Messages{
		r: client.r,
	}
}

// Devices method returns an interface to interact with devices
func (client *KeystoneClientImpl) Devices() *Devices {
	return &Devices{
		r: client.r,
	}
}

// Organizations method returns an interfaec to interact with organizations
func (client *KeystoneClientImpl) Organizations() *Organizations {
	return &Organizations{
		r: client.r,
	}
}

// Logs method returns an interface to interact with logs
func (client *KeystoneClientImpl) Logs() *Logs {
	return &Logs{
		r: client.r,
	}
}
