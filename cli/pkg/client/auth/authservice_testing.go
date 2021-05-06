// +build test

package auth

import (
	"context"

	"github.com/wearedevx/keystone/api/pkg/models"
)

type AuthService interface {
	Name() string
	Start() (string, error)
	WaitForExternalLogin() error
	CheckAccount(account map[string]string) (bool, error)
	Finish(pkey []byte) (models.User, string, error)
}

func GetAuthService(serviceName string, ctx context.Context, apiUrl string) (AuthService, error) {
	return new(dummyAuthService), nil
}
