package ci

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/google/go-github/v40/github"
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

const GithubCI CiServiceType = "github-ci"

type ServicesKeys map[string]string

type ApiKey string

type gitHubCiService struct {
	err          error
	name         string
	apiUrl       string
	ctx          *core.Context
	servicesKeys ServicesKeys
	apiKey       ApiKey
	environment  string
	client       *github.Client
	pk           *github.PublicKey
	repo         *github.Repository
	sentFiles    []string
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
	return ui.RenderTemplate(
		"github-ci-usage",
		`Secrets will be available in the {{ .Environment }} environment.
To make use of files, add the following step in your jobs (it should go right after the Checkout step):

      - name: Load Secrets
        uses: wearedevx/keystone-action@v2
        with:
          files: |-
            {{ range .Files }}${{"{{"}} secrets.{{ . }} {{"}}"}}
            {{ end }}
`,
		map[string]interface{}{
			"Environment": g.environment,
			"Files":       g.sentFiles,
		},
	)
}

// Type method returns the type of the service
func (g *gitHubCiService) Type() string { return string(GithubCI) }

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
// Secrets will be sent as secrets to the GitHub repository.
// Files will receive a different treatment:
//   - an uppercase snakecase name will be derived from the path
//   - the content will be base64 encoded
//   - the actual path is prepended to the content, separated by a #
// TODO: the # separator might be problematic ?
func (g *gitHubCiService) PushSecret(
	message models.MessagePayload,
	environment string,
) CiService {
	if g.err != nil {
		return g
	}

	g.environment = environment

	return g.
		initClient().
		createEnvironment().
		getGithubPublicKey().
		sendEnvironmentSecrets().
		sendEnvironmentFiles()
}

// CleanSecret method remove all the secrets for the given environment
// from the CI service
func (g *gitHubCiService) CleanSecret(environment string) CiService {
	if g.err != nil {
		return g
	}

	g.environment = environment

	g.initClient().
		cleanSecrets().
		cleanFiles()

	return g
}

func (g *gitHubCiService) initClient() *gitHubCiService {
	context := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: string(g.apiKey)},
	)
	tc := oauth2.NewClient(context, ts)

	client := github.NewClient(tc)

	g.client = client

	options := g.getOptions()
	g.repo, _, g.err = g.client.Repositories.Get(
		context,
		options["Owner"],
		options["Project"],
	)

	return g
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

func (g *gitHubCiService) getGithubPublicKey() *gitHubCiService {
	if g.err != nil {
		return g
	}

	publicKey, resp, err := g.client.Actions.GetEnvPublicKey(
		context.Background(),
		int(*g.repo.ID),
		g.environment,
	)
	switch {
	case resp.StatusCode == 403:
		g.err = ErrorGithubCIPermissionDenied
		return g

	case resp.StatusCode == 404:
		g.err = ErrorGithubCINoSuchRepository
		return g

	case err != nil:
		g.err = err
		return g
	}

	g.pk = publicKey

	return g
}

func (g *gitHubCiService) encryptSecret(
	key, value string,
	secret *github.EncryptedSecret,
) *gitHubCiService {
	publicKey, err := base64.StdEncoding.DecodeString(g.pk.GetKey())
	if err != nil {
		g.err = err
		return g
	}

	bpk := sodium.BoxPublicKey{
		Bytes: sodium.Bytes(publicKey),
	}

	encryptedValue := sodium.Bytes(value).SealedBox(bpk)
	base64data := base64.StdEncoding.EncodeToString(encryptedValue)

	*secret = github.EncryptedSecret{
		Name:           key,
		KeyID:          g.pk.GetKeyID(),
		EncryptedValue: base64data,
	}

	return g
}

func (g *gitHubCiService) sendSecret(
	secret *github.EncryptedSecret,
	environment string,
) *gitHubCiService {
	resp, err := g.client.Actions.CreateOrUpdateEnvSecret(
		context.Background(),
		int(*g.repo.ID),
		environment,
		secret,
	)
	switch {
	case resp.StatusCode == 401:
		g.err = ErrorGithubCIPermissionDenied
		return g

	case err != nil:
		g.err = err
		return g
	}

	return g
}

func (g *gitHubCiService) createEnvironment() *gitHubCiService {
	g.client.Repositories.CreateUpdateEnvironment(
		context.Background(),
		g.servicesKeys["Owner"],
		g.servicesKeys["Project"],
		g.environment,
		&github.CreateUpdateEnvironment{},
	)

	return g
}

func (g *gitHubCiService) sendEnvironmentSecrets() *gitHubCiService {
	if g.err != nil {
		return g
	}

	secrets := g.ctx.ListSecrets()

	for _, secret := range secrets {
		value, ok := secret.Values[core.EnvironmentName(g.environment)]
		if !ok && secret.Required {
			// TODO: have a better error handling
			g.err = fmt.Errorf("missing required secret: %s", string(value))
			break
		}

		encryptedSecret := github.EncryptedSecret{}
		if g.
			encryptSecret(secret.Name, string(value), &encryptedSecret).
			sendSecret(&encryptedSecret, g.environment).
			err != nil {
			break
		}

	}

	return g
}

// TODO: error out if file is required and not present or empty
func (g *gitHubCiService) sendEnvironmentFiles() *gitHubCiService {
	if g.err != nil {
		return g
	}

	files := g.ctx.ListFiles()

	g.sentFiles = make([]string, 0)

	filecachepath := g.ctx.CachedEnvironmentFilesPath(g.environment)
	for _, file := range files {
		fullpath := path.Join(filecachepath, file.Path)
		freader, err := os.Open(fullpath)
		if err != nil {
			g.err = err
			break
		}

		contents, err := base64encode(freader)
		if err != nil {
			break
		}

		// file var name, uses the path to create an uppercase snakecase name
		key := pathToVarname(file.Path)
		// prepend the path, separated with a #
		value := fmt.Sprintf("%s#%s", file.Path, contents)

		encryptedSecret := github.EncryptedSecret{}
		if g.
			encryptSecret(key, value, &encryptedSecret).
			sendSecret(&encryptedSecret, g.environment).
			err != nil {
			break
		}

		g.sentFiles = append(g.sentFiles, key)
	}

	return g
}

// hasSecret method returns true if the secret exists for the given environment
func (g *gitHubCiService) hasSecret(secret string) bool {
	if g.err != nil {
		return false
	}

	s, _, err := g.client.Actions.GetEnvSecret(
		context.Background(),
		int(*g.repo.ID),
		g.environment,
		secret,
	)

	switch {
	case s != nil:
		return true

	case err != nil && strings.Contains(err.Error(), "404"):
		return false

	default:
		g.err = err
		return false
	}
}

// deleteSecret method removes one secret from the remote
func (g *gitHubCiService) deleteSecret(secret string) *gitHubCiService {
	if g.err != nil {
		return g
	}

	_, err := g.client.Actions.DeleteEnvSecret(
		context.Background(),
		int(*g.repo.ID),
		g.environment,
		secret,
	)
	g.err = err

	return g
}

// cleanSecrets method removes all the secrets on the remote for the given
// environment
func (g *gitHubCiService) cleanSecrets() *gitHubCiService {
	if g.err != nil {
		return g
	}

	secrets := g.ctx.ListSecrets()
	for _, secret := range secrets {
		key := secret.Name

		if g.hasSecret(key) {
			g.deleteSecret(key)
		}
	}

	return g
}

// cleanFiles method removes all the files secrets on the remote for the given
// envrionment
func (g *gitHubCiService) cleanFiles() *gitHubCiService {
	if g.err != nil {
		return g
	}

	files := g.ctx.ListFiles()
	for _, file := range files {
		key := pathToVarname(file.Path)

		if g.hasSecret(key) {
			g.deleteSecret(key)
		}
	}

	return g
}
