package client

import (
	"github.com/wearedevx/keystone/cli/internal/config"
	"github.com/wearedevx/keystone/cli/internal/errors"
)

type KeystoneClientImpl struct {
	r requester
}

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

func (client *KeystoneClientImpl) Roles() *Roles {
	return &Roles{
		r: client.r,
	}
}

func (client *KeystoneClientImpl) Project(projectId string) *Project {
	return &Project{
		r:  client.r,
		id: projectId,
	}
}

func (client *KeystoneClientImpl) Users() *Users {
	return &Users{
		r: client.r,
	}
}

func (client *KeystoneClientImpl) Messages() *Messages {
	return &Messages{
		r: client.r,
	}
}
func (client *KeystoneClientImpl) Devices() *Devices {
	return &Devices{
		r: client.r,
	}
}
