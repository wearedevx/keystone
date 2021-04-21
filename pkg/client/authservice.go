// +build !test

package client

import (
	"context"
	"fmt"

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
	var c AuthService
	var err error

	switch serviceName {
	case "GitHub":
		c = GitHubAuth(ctx)

	case "GitLab":
		c = GitLabAuth(ctx)

	default:
		err = fmt.Errorf("Unknown service name %s", serviceName)
	}

	return c, err
}
