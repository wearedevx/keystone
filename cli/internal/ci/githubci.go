package ci

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-github/v32/github"
	"github.com/jamesruan/sodium"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/config"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
	"github.com/wearedevx/keystone/cli/ui/prompts"
	"golang.org/x/oauth2"
)

var (
	githubClientId     string
	githubClientSecret string
)

type ServicesKeys map[string]string

type ApiKey string

type gitHubCiService struct {
	err          error
	name         string
	apiUrl       string
	ctx          *core.Context
	servicesKeys ServicesKeys
	apiKey       ApiKey
	client       *github.Client
}

var (
	ErrorGithubCIPermissionDenied error = errors.New(
		"you don't have rights to send secrets to the repo. Please ensure your personal access token has access to \"repo\" scope",
	)
	ErrorGithubCINoSuchRepository = errors.New(
		"you are trying to send secret to a repository that doesn't exist. Please make sure repo's name and owner is correct",
	)
	ErrorGithubCITooLarge = errors.New(
		"secrets and files are too large to send to CI",
	)
)

// GitHubCi function return a `CiService` that works with the GitHub API
func GitHubCi(ctx *core.Context, name string, apiUrl string) CiService {
	kf := keystonefile.KeystoneFile{}
	kf.Load(ctx.Wd)

	savedService := kf.GetCiService(name)

	ciService := &gitHubCiService{
		err:    nil,
		name:   name,
		apiUrl: apiUrl,
		ctx:    ctx,
		servicesKeys: ServicesKeys{
			"Owner":   savedService.Options["Owner"],
			"Project": savedService.Options["Project"],
		},
		apiKey: ApiKey(config.GetServiceApiKey(string(GithubCI))),
	}

	return ciService
}

// Name method returns the name of the service
func (g *gitHubCiService) Name() string { return g.name }

func (g *gitHubCiService) Usage() string {
	return `See https://github.com/wearedevx/keystone-action to use them.`
}

// Type method returns the type of the service
func (g *gitHubCiService) Type() CiServiceType { return GithubCI }

// GetOptions method returns the service options
func (g *gitHubCiService) GetOptions() map[string]string {
	return g.servicesKeys
}

// Setup method starts the ci service setup process, asking
// the user information through prompts
func (g *gitHubCiService) Setup() CiService {
	if g.err != nil {
		return g
	}

	// These are the prompts for keys and such
	// as those are all github specifics
	g.askForRepoUrl()
	g.askForApiKey()

	return g
}

// CheckSetup method verifies the user submitted information is valid
func (g *gitHubCiService) CheckSetup() CiService {
	if g.err != nil {
		return g
	}

	if len(g.servicesKeys["Owner"]) == 0 ||
		len(g.servicesKeys["Project"]) == 0 ||
		len(g.getApiKey()) == 0 {
		g.err = ErrorMissingCiInformation
	}

	return g
}

// PushSecret sends a "Message" (that's a complete encrypted environment)
// to GitHub as one repository Secret
func (g *gitHubCiService) PushSecret(
	message models.MessagePayload,
	environment string,
) CiService {
	if g.err != nil {
		return g
	}

	g.initClient()

	var payload string

	g.err = message.Serialize(&payload)
	if g.err != nil {
		return g
	}

	publicKey, resp, err := g.client.Actions.GetRepoPublicKey(
		context.Background(),
		g.servicesKeys["Owner"],
		g.servicesKeys["Project"],
	)

	if resp.StatusCode == 403 {
		g.err = ErrorGithubCIPermissionDenied
		return g
	}

	if resp.StatusCode == 404 {
		g.err = ErrorGithubCINoSuchRepository
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
			Name: fmt.Sprintf(
				"KEYSTONE_%s_SLOT_%o",
				strings.ToUpper(environment),
				i+1,
			),
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
			g.err = ErrorGithubCIPermissionDenied
			continue
		}

		if err != nil {
			g.err = err
			continue
		}
	}

	return g
}

// CleanSecret method remove all the secrets for the given environment
// from the CI service
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
		if strings.Contains(err.Error(), "404") {
			g.err = ErrorNoSecretsForEnvironment
		} else {
			g.err = err
		}
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

func (g *gitHubCiService) getOptions() ServicesKeys {
	return g.servicesKeys
}

func (g *gitHubCiService) setKeys(servicesKeys ServicesKeys) CiService {
	var service keystonefile.CiService

	g.servicesKeys = servicesKeys
	service.Name = g.Name()
	service.Type = string(GithubCI)
	service.Options = g.servicesKeys

	// Write the local keystone.yaml changes
	new(keystonefile.KeystoneFile).
		Load(g.ctx.Wd).
		AddCiService(service).
		Save()

	return g
}

func (g *gitHubCiService) getApiKey() ApiKey {
	apiKey := config.GetServiceApiKey(string(g.Type()))
	return ApiKey(apiKey)
}

func (g *gitHubCiService) setApiKey(apiKey ApiKey) {
	g.apiKey = apiKey
	config.SetServiceApiKey(string(g.Type()), string(apiKey))
	config.Write()
}

func (g *gitHubCiService) askForRepoUrl() CiService {
	// serviceName := availableServices[GithubCI]
	serviceOptions := g.getOptions()
	owner := serviceOptions["Owner"]
	project := serviceOptions["Project"]

	serviceUrl := ""

	if serviceOptions["Owner"] != "" && serviceOptions["Project"] != "" {
		serviceUrl = "https://github.com/" + serviceOptions["Owner"] + "/" + serviceOptions["Project"]
	}

	urlIsValid := false

	for !urlIsValid {
		serviceUrl = prompts.StringInput(
			"GitHub repository URL",
			serviceUrl,
		)

		u, err := new(url.URL).Parse(serviceUrl)
		if err != nil {
			ui.Print(ui.RenderTemplate(
				"malformed-url",
				`{{ "Warning" | yellow }} The url is malformed
This caused by: {{ .Cause }}`,
				map[string]string{
					"Cause": err.Error(),
				},
			),
			)
			continue
		}

		p := u.EscapedPath()
		p = strings.TrimPrefix(p, "/")
		parts := strings.Split(p, "/")

		if (len(parts) != 2) || (u.Hostname() != "github.com") {
			ui.Print(ui.RenderTemplate(
				"not-a-github-url",
				`{{ "Warning" | yellow }} This is not a valid github URL`,
				map[string]string{},
			),
			)
			continue
		}

		owner = parts[0]
		project = parts[1]

		if len(owner) == 0 || len(project) == 0 {
			ui.Print(ui.RenderTemplate(
				"not-a-valid-repo",
				`{{ "Warning" | yellow }} This is not a valid repository URL`,
				map[string]string{},
			),
			)
			continue
		}

		urlIsValid = true
	}

	serviceOptions["Owner"] = owner
	serviceOptions["Project"] = project

	g.setKeys(serviceOptions)

	return g
}

func (g *gitHubCiService) askForApiKey() CiService {
	serviceName := g.Name()
	apiKey := g.getApiKey()

	fmt.Println(
		"Personal access token can be generated here: https://github.com/settings/tokens/new\nIt should have access to \"repo\" scope.",
	)
	apiKey = ApiKey(
		prompts.StringInput(serviceName+" Access Token", string(apiKey)),
	)

	g.setApiKey(apiKey)

	return g
}

// Error method returns the last error encountered
func (g *gitHubCiService) Error() error {
	return g.err
}

func (g *gitHubCiService) sliceMessageInParts(
	message string,
) ([]string, error) {
	slots := make([]string, 5)

	// Add spaces to message to make it divisible by 5 (number of slots)
	for len(message)%5 != 0 {
		message += " "
	}
	slotSize := (len(message) / 5)

	var err error

	// 64Kb is maximum size for a slot in github
	if slotSize*(4/3) > 64000 { // base64 encoding make 4 bytes out of 3
		err = ErrorGithubCITooLarge
	}

	slots[0] = message[0:slotSize]
	slots[1] = message[slotSize : slotSize*2]
	slots[2] = message[slotSize*2 : slotSize*3]
	slots[3] = message[slotSize*3 : slotSize*4]
	slots[4] = message[slotSize*4 : slotSize*5]

	return slots, err
}
