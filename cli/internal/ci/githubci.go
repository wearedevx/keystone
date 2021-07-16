package ci

import (
	"context"
	"encoding/base64"
	"errors"
	"os"

	"github.com/google/go-github/v32/github"
	"github.com/jamesruan/sodium"
	"github.com/manifoldco/promptui"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/config"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
	"golang.org/x/oauth2"
)

var githubClientId string
var githubClientSecret string

type ServicesKeys map[string]string

type ApiKey string

type gitHubCiService struct {
	err          error
	apiUrl       string
	ctx          core.Context
	kf           keystonefile.KeystoneFile
	servicesKeys ServicesKeys
	apiKey       ApiKey
	client       *github.Client
}

func GitHubCi(ctx core.Context, apiUrl string) CiService {
	kf := keystonefile.KeystoneFile{}
	kf.Load(ctx.Wd)

	savedService := kf.GetCiService("github-ci")

	ciService := &gitHubCiService{
		err:    nil,
		apiUrl: apiUrl,
		ctx:    ctx,
		kf:     kf,
		servicesKeys: ServicesKeys{
			"Owner":   savedService.Keys["Owner"],
			"Project": savedService.Keys["Project"],
		},
		apiKey: ApiKey(config.GetServiceApiKey("github")),
	}

	return ciService
}

func (g gitHubCiService) Name() string { return "github-ci" }

func (g *gitHubCiService) Setup() CiService {
	if g.err != nil {
		return g
	}
	g.askForKeys()
	g.askForApiKey()

	// There should go the prompts for keys and such
	// as those are all github specifics

	return g
}

func (g *gitHubCiService) CheckSetup() {
	if len(g.servicesKeys["Owner"]) == 0 || len(g.servicesKeys["Project"]) == 0 || len(g.getApiKey()) == 0 {
		g.err = errors.New("There is missing information in CI service settings.\nUse $ ks ci setup")
	}
}

// PushSecret sends a "Message" (that's a complete encrypted environment)
// to GitHub as one repository Secret
func (g *gitHubCiService) PushSecret(message models.MessagePayload) CiService {
	if g.err != nil {
		return g
	}

	var payload string

	message.Serialize(&payload)
	publicKey, _, err := g.client.Actions.GetRepoPublicKey(
		context.Background(),
		g.servicesKeys["Owner"],
		g.servicesKeys["Project"],
	)
	data, err := base64.StdEncoding.DecodeString(publicKey.GetKey())

	if err != nil {
		panic(err)
	}

	boxPK := sodium.BoxPublicKey{
		Bytes: sodium.Bytes(data),
	}

	encryptedValue := sodium.Bytes(payload).SealedBox(boxPK)
	base64data := base64.StdEncoding.EncodeToString(encryptedValue)

	encryptedSecret := &github.EncryptedSecret{
		Name:           "KEYSTONE_SLOT_1",
		KeyID:          publicKey.GetKeyID(),
		EncryptedValue: base64data,
	}

	g.initClient()
	_, err = g.client.Actions.CreateOrUpdateRepoSecret(
		context.Background(),
		g.servicesKeys["Owner"],
		g.servicesKeys["Project"],
		encryptedSecret,
	)

	if err != nil {
		g.err = err
	}

	return g
}

func (g gitHubCiService) initClient() {
	context := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: string(g.apiKey)},
	)
	tc := oauth2.NewClient(context, ts)

	client := github.NewClient(tc)
	g.client = client
}

func (g gitHubCiService) getKeys() ServicesKeys {
	return g.servicesKeys
}

func (g *gitHubCiService) setKeys(servicesKeys ServicesKeys) error {
	var service keystonefile.CiService

	g.servicesKeys = servicesKeys
	service.Name = g.Name()
	service.Keys = g.servicesKeys
	file := g.kf.SetCiService(service)

	if file.Err() != nil {
		return file.Err()
	}

	return nil
}

func (g gitHubCiService) getApiKey() ApiKey {
	apiKey := config.GetServiceApiKey(g.Name())
	return ApiKey(apiKey)
}

func (g gitHubCiService) setApiKey(apiKey ApiKey) {
	g.apiKey = apiKey
	config.SetServiceApiKey(g.Name(), string(apiKey))
	config.Write()
}

func (g gitHubCiService) askForKeys() {
	serviceName := g.Name()
	servicesKeys := g.getKeys()
	for key, value := range servicesKeys {
		p := promptui.Prompt{
			Label:   serviceName + "'s " + key,
			Default: value,
		}
		result, err := p.Run()

		// Handle user cancelation
		// or prompt error
		if err != nil {
			if err.Error() != "^C" {
				ui.PrintError(err.Error())
				os.Exit(1)
			}
			os.Exit(0)
		}
		servicesKeys[key] = result
	}

	err := g.setKeys(servicesKeys)

	if err != nil {
		ui.PrintError(err.Error())
	}
}

func (g gitHubCiService) askForApiKey() {
	serviceName := g.Name()

	p := promptui.Prompt{
		Label:   serviceName + "'s Api key",
		Default: string(g.getApiKey()),
	}

	result, err := p.Run()

	// Handle user cancelation
	// or prompt error
	if err != nil {
		if err.Error() != "^C" {
			ui.PrintError(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	g.setApiKey(ApiKey(result))

	if err != nil {
		ui.PrintError(err.Error())
	}

}

func (g gitHubCiService) Error() error {
	return g.err
}
