//go:build test
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
	Finish(pkey []byte, device string, deviceUID string) (models.User, string, string, error)
}

func GetAuthService(serviceName string, apiURL string) (AuthService, error) {
	a := new(dummyAuthService)
	a.apiURL = apiURL

	return a, nil
}
