package ci

import (
	"context"
	"fmt"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/config"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"golang.org/x/oauth2"
)

var githubClientId string
var githubClientSecret string

type gitHubCiService struct {
	apiUrl       string
	ctx          core.Context
	kf           keystonefile.KeystoneFile
	conf         *oauth2.Config
	servicesKeys ServicesKeys
	apiKey       ApiKey
}

func GitHubCi(ctx core.Context, apiUrl string) CiService {
	kf := keystonefile.KeystoneFile{}
	kf.Load(ctx.Wd)

	savedService := kf.GetCiService("github")

	return &gitHubCiService{
		apiUrl: apiUrl,
		ctx:    ctx,
		kf:     kf,
		servicesKeys: ServicesKeys{
			"Owner":   savedService.Keys["Owner"],
			"Project": savedService.Keys["Project"],
		},
		apiKey: ApiKey(config.GetServiceApiKey("github")),
	}
}

func (g gitHubCiService) Name() string { return "github" }

func (g gitHubCiService) PushSecret(ctx context.Context, message models.MessagePayload) error {
	token := g.GetApiKey()

	fmt.Println(token)
	fmt.Println(message)

	// ts := oauth2.StaticTokenSource(token)
	// tc := oauth2.NewClient(g.ctx, ts)

	// g.token = token
	// g.client = github.NewClient(tc)

	return nil
}

func (g gitHubCiService) GetKeys() ServicesKeys {
	return g.servicesKeys
}

func (g gitHubCiService) SetKeys(servicesKeys ServicesKeys) error {
	g.servicesKeys = servicesKeys
	var service keystonefile.CiService
	service.Name = g.Name()
	service.Keys = g.servicesKeys
	file := g.kf.SetCiService(service)
	if file.Err() != nil {
		return file.Err()
	}
	return nil
}

func (g gitHubCiService) GetApiKey() ApiKey {
	apiKey := config.GetServiceApiKey(g.Name())
	return ApiKey(apiKey)
}

func (g gitHubCiService) SetApiKey(apiKey ApiKey) {
	g.apiKey = apiKey
	config.SetServiceApiKey(g.Name(), string(apiKey))
}
