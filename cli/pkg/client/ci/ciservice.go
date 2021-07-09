package ci

import (
	"fmt"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/pkg/core"
)

type ServicesKeys map[string]string

type ApiKey string

type CiService interface {
	Name() string
	PushSecret(models.MessagePayload) error
	// Finish(pkey []byte) (models.User, string, error)
	GetKeys() ServicesKeys
	SetKeys(ServicesKeys) error
	GetApiKey() ApiKey
	SetApiKey(ApiKey)
	InitClient() CiService
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
