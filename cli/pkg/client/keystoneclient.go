package client

import (
	"log"

	"github.com/wearedevx/keystone/cli/internal/config"
	"github.com/wearedevx/keystone/cli/internal/errors"
)

type KeystoneClientImpl struct {
	r requester
}

// NewKeystoneClient function returns a new instance of KeysotneClient
func NewKeystoneClient() (KeystoneClient, *errors.Error) {
	account, index := config.GetCurrentAccount()
	token, refreshToken := config.GetAuthToken()

	if index < 0 {
		return nil, errors.MustBeLoggedIn(nil)
	}

	return &KeystoneClientImpl{
		r: newRequester(account.UserID, token, refreshToken),
	}, nil
}

// Roles method returns an interface to interact with roles
func (client *KeystoneClientImpl) Roles() *Roles {
	return &Roles{
		log: log.New(log.Writer(), "[Roles] ", 0),
		r:   client.r,
	}
}

// Project method returns an interface to interact with projects
func (client *KeystoneClientImpl) Project(projectID string) *Project {
	return &Project{
		log: log.New(log.Writer(), "[Project] ", 0),
		r:   client.r,
		id:  projectID,
	}
}

// Users method returns an interface to interact with users/members
func (client *KeystoneClientImpl) Users() *Users {
	return &Users{
		log: log.New(log.Writer(), "[Users] ", 0),
		r:   client.r,
	}
}

// Messages method returns an interface to interact with messages
func (client *KeystoneClientImpl) Messages() *Messages {
	return &Messages{
		log: log.New(log.Writer(), "[Messages] ", 0),
		r:   client.r,
	}
}

// Devices method returns an interface to interact with devices
func (client *KeystoneClientImpl) Devices() *Devices {
	return &Devices{
		log: log.New(log.Writer(), "[Devices] ", 0),
		r:   client.r,
	}
}

// Organizations method returns an interfaec to interact with organizations
func (client *KeystoneClientImpl) Organizations() *Organizations {
	return &Organizations{
		log: log.New(log.Writer(), "[Organizations] ", 0),
		r:   client.r,
	}
}
