package ci

import (
	"context"
	"encoding/base64"

	"github.com/google/go-github/v32/github"
	"github.com/jamesruan/sodium"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/config"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/pkg/core"
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

	savedService := kf.GetCiService("github")

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

	// There should go the prompts for keys and such
	// as those are all github specifics

	return g
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

	_, err = g.client.Actions.CreateOrUpdateRepoSecret(
		context.Background(),
		g.servicesKeys["Owner"],
		g.servicesKeys["Project"],
		encryptedSecret,
	)

	if err != nil {
		return err
	}

	return nil
}

func (g gitHubCiService) initClient() CiService {
	context := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: string(g.apiKey)},
	)
	tc := oauth2.NewClient(context, ts)

	client := github.NewClient(tc)
	g.client = client

	return g
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
