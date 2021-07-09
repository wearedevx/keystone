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

type gitHubCiService struct {
	apiUrl       string
	ctx          core.Context
	kf           keystonefile.KeystoneFile
	conf         *oauth2.Config
	servicesKeys ServicesKeys
	apiKey       ApiKey
	client       *github.Client
}

func GitHubCi(ctx core.Context, apiUrl string) CiService {
	kf := keystonefile.KeystoneFile{}
	kf.Load(ctx.Wd)

	savedService := kf.GetCiService("github")

	ciService := &gitHubCiService{
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

func (g gitHubCiService) Name() string { return "github" }

// PushSecret sends a "Message" (that's a complete encrypted environment)
// to GitHub as one repository Secret
func (g gitHubCiService) PushSecret(message models.MessagePayload) error {
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
		Name:           "keystone_slot_1",
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

func (g gitHubCiService) GetKeys() ServicesKeys {
	return g.servicesKeys
}

func (g gitHubCiService) SetKeys(servicesKeys ServicesKeys) error {
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

func (g gitHubCiService) GetApiKey() ApiKey {
	apiKey := config.GetServiceApiKey(g.Name())
	return ApiKey(apiKey)
}

func (g gitHubCiService) SetApiKey(apiKey ApiKey) {
	g.apiKey = apiKey
	config.SetServiceApiKey(g.Name(), string(apiKey))
	config.Write()
}

func (g gitHubCiService) InitClient() CiService {
	context := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: string(g.apiKey)},
	)
	tc := oauth2.NewClient(context, ts)

	client := github.NewClient(tc)
	g.client = client

	return g
}
