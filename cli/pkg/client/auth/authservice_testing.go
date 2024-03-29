// +build test

package auth

import (
	"github.com/wearedevx/keystone/api/pkg/models"
)

type AuthService interface {
	Name() string
	Start() (string, error)
	WaitForExternalLogin() error
	CheckAccount(account map[string]string) (bool, error)
	Finish(pkey []byte, device string, deviceUID string) (models.User, string, error)
}

func GetAuthService(serviceName string, apiUrl string) (AuthService, error) {
	a := new(dummyAuthService)
	a.apiUrl = apiUrl

	return a, nil
}
