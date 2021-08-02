package ci

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"

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

const NAME = "github-ci"

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

	savedService := kf.GetCiService(NAME)

	ciService := &gitHubCiService{
		err:    nil,
		apiUrl: apiUrl,
		ctx:    ctx,
		kf:     kf,
		servicesKeys: ServicesKeys{
			"Owner":   savedService.Keys["Owner"],
			"Project": savedService.Keys["Project"],
		},
		apiKey: ApiKey(config.GetServiceApiKey(NAME)),
	}

	return ciService
}

func (g *gitHubCiService) Name() string { return NAME }

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
func (g *gitHubCiService) PushSecret(message models.MessagePayload, environment string) CiService {
	if g.err != nil {
		return g
	}

	g.initClient()

	var payload string

	message.Serialize(&payload)
	publicKey, resp, err := g.client.Actions.GetRepoPublicKey(
		context.Background(),
		g.servicesKeys["Owner"],
		g.servicesKeys["Project"],
	)

	if resp.StatusCode == 403 {
		g.err = errors.New("You don't have rights to send secrets to the repo. Please ensure your personal access token has access to \"repo\" scope.")
		return g
	}

	if resp.StatusCode == 404 {
		g.err = errors.New("You are trying to send secret to a repository that doesn't exist. Please make sure repo's name and owner is correct.")
		return g
	}

	if err != nil {
		g.err = err
		return g
	}

	data, err := base64.StdEncoding.DecodeString(publicKey.GetKey())

	if err != nil {
		g.err = err
		return g
	}

	boxPK := sodium.BoxPublicKey{
		Bytes: sodium.Bytes(data),
	}

	slots, err := g.sliceMessageInParts(payload)

	if err != nil {
		g.err = err
		return g
	}

	for i, slot := range slots {
		encryptedValue := sodium.Bytes(slot).SealedBox(boxPK)

		base64data := base64.StdEncoding.EncodeToString(encryptedValue)

		encryptedSecret := &github.EncryptedSecret{
			Name:           fmt.Sprintf("KEYSTONE_%s_SLOT_%o", strings.ToUpper(environment), i+1),
			KeyID:          publicKey.GetKeyID(),
			EncryptedValue: base64data,
		}

		resp, err := g.client.Actions.CreateOrUpdateRepoSecret(
			context.Background(),
			g.servicesKeys["Owner"],
			g.servicesKeys["Project"],
			encryptedSecret,
		)

		if resp.StatusCode == 401 {
			g.err = errors.New("You don't have rights to send secrets to the repo. Please ensure your personal access token has access to \"repo\" scope.")
			continue
		}

		if err != nil {
			g.err = err
			continue
		}
	}

	return g
}

func (g *gitHubCiService) CleanSecret(environment string) CiService {
	if g.err != nil {
		return g
	}

	g.initClient()

	_, err := g.client.Actions.DeleteRepoSecret(
		context.Background(),
		g.servicesKeys["Owner"],
		g.servicesKeys["Project"],
		fmt.Sprintf("KEYSTONE_%s_SLOT_1", strings.ToUpper(environment)),
	)

	if err != nil {
		g.err = err
	}

	return g
}

func (g *gitHubCiService) initClient() {
	context := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: string(g.apiKey)},
	)
	tc := oauth2.NewClient(context, ts)

	client := github.NewClient(tc)

	g.client = client
}

func (g *gitHubCiService) getKeys() ServicesKeys {
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

func (g *gitHubCiService) getApiKey() ApiKey {
	apiKey := config.GetServiceApiKey(g.Name())
	return ApiKey(apiKey)
}

func (g *gitHubCiService) setApiKey(apiKey ApiKey) {
	g.apiKey = apiKey
	config.SetServiceApiKey(g.Name(), string(apiKey))
	config.Write()
}

func (g *gitHubCiService) askForKeys() {
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

func (g *gitHubCiService) askForApiKey() {
	serviceName := g.Name()

	fmt.Println("Personal access token can be generated here: https://github.com/settings/tokens/new\nIt should have access to \"repo\" scope.")
	p := promptui.Prompt{
		Label:   serviceName + "'s Access token",
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

func (g *gitHubCiService) Error() error {
	return g.err
}

func (g *gitHubCiService) sliceMessageInParts(message string) ([]string, error) {
	slots := make([]string, 5)

	// Add spaces to message to make it divisible by 5 (number of slots)
	for len(message)%5 != 0 {
		message += " "
	}
	slotSize := (len(message) / 5)

	var err error

	// 64Kb is maximum size for a slot in github
	if slotSize*(4/3) > 64000 { // base64 encoding make 4 bytes out of 3
		err = errors.New("Secrets and files are too large to send to CI")
	}

	slots[0] = message[0:slotSize]
	slots[1] = message[slotSize : slotSize*2]
	slots[2] = message[slotSize*2 : slotSize*3]
	slots[3] = message[slotSize*3 : slotSize*4]
	slots[4] = message[slotSize*4 : slotSize*5]

	return slots, err
}

func (g *gitHubCiService) PrintSuccess(environment string) {
	ui.PrintSuccess(fmt.Sprintf("Secrets successfully sent to CI service, environment %s. See https://github.com/wearedevx/keystone-action to use them.", environment))
}
