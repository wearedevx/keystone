// +build !test

package auth

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

func GetAuthService(serviceName string, ctx context.Context, apiUrl string) (AuthService, error) {
	var c AuthService
	var err error

	switch serviceName {
	case "github":
		c = GitHubAuth(ctx, apiUrl)

	case "gitlab":
		c = GitLabAuth(ctx, apiUrl)

	default:
		err = fmt.Errorf("Unknown service name %s", serviceName)
	}

	return c, err
}
