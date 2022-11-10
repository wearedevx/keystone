//go:build !test
// +build !test

package auth

import (
	"context"
	"fmt"

	"github.com/wearedevx/keystone/api/pkg/models"
)

type AuthService interface {
	Name() string
	Start() (string, error)
	WaitForExternalLogin() error
	CheckAccount(account map[string]string) (bool, error)
	Finish(
		pkey []byte,
		device string,
		deviceUID string,
	) (models.User, string, string, error)
}

// GetAuthService function returns the instance of an AuthService implementor
// based on `serviceName`
func GetAuthService(
	serviceName string,
	apiURL string,
) (AuthService, error) {
	var c AuthService
	var err error
	ctx := context.Background()

	switch serviceName {
	case "github":
		c = GitHubAuth(ctx, apiURL)

	case "gitlab":
		c = GitlabAuth(ctx, apiURL)

	default:
		err = fmt.Errorf("unknown service name %s", serviceName)
	}

	return c, err
}
