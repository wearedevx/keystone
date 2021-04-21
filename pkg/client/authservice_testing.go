// +build test

package client

import (
	"context"

	"github.com/wearedevx/keystone/internal/models"
)

type AuthService interface {
	Name() string
	Start() (string, error)
	WaitForExternalLogin() error
	CheckAccount(account map[string]string) (bool, error)
	Finish(pkey []byte) (models.User, string, error)
}

func GetAuthService(serviceName string, ctx context.Context) (AuthService, error) {
	return new(dummyAuthService), nil
}
