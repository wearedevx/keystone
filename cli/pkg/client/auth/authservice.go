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
	) (models.User, string, error)
}

// GetAuthService function returns the instance of an AuthService implementor
// based on `serviceName`
func GetAuthService(
	serviceName string,
	apiUrl string,
) (AuthService, error) {
	var c AuthService
	var err error
	ctx := context.Background()

	switch serviceName {
	case "github":
		c = GitHubAuth(ctx, apiUrl)

	case "gitlab":
		c = GitlabAuth(ctx, apiUrl)

	default:
		err = fmt.Errorf("unknown service name %s", serviceName)
	}

	return c, err
}
