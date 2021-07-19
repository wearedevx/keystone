package ci

import (
	"fmt"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/pkg/core"
)

type CiService interface {
	Name() string
	Setup() CiService

	PushSecret(message models.MessagePayload, environment string) CiService
	CleanSecret(environment string) CiService
	CheckSetup()
	Error() error
	// // Finish(pkey []byte) (models.User, string, error)
	// GetKeys() ServicesKeys
	// SetKeys(ServicesKeys) error
	// GetApiKey() ApiKey
	// SetApiKey(ApiKey)
}

func GetCiService(serviceName string, ctx core.Context, apiUrl string) (CiService, error) {
	var c CiService
	var err error

	switch serviceName {
	case "github":
		c = GitHubCi(ctx, apiUrl)

	default:
		err = fmt.Errorf("Unknown service name %s", serviceName)
	}

	return c, err
}
