package ci

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
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

const (
	GithubCINbSlots    = 5
	GithubCISLotLength = 64 * 1024
)

type ServicesKeys map[string]string

type ApiKey string

type gitHubCiService struct {
	err          error
	log          *log.Logger
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
		log:    log.New(log.Writer(), "[GitHubCI]", 0),
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
		getGithubPublicKey().
		sendSlot(message)
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

	g.client = github.NewClient(tc)

	options := g.getOptions()
	g.repo, _, g.err = g.client.Repositories.Get(
		context,
		options["Owner"],
		options["Project"],
	)

	g.log.Printf(
		"Initializing GitHub CI client for project %s/%s with PAT %s\n",
		options["Owner"],
		options["Project"],
		string(g.apiKey),
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
	g.log.Printf("Got PAT from config: %s\n", apiKey)

	return ApiKey(apiKey)
}

func (g *gitHubCiService) setApiKey(apiKey ApiKey) {
	g.apiKey = apiKey
	g.log.Printf("Set PAT: %s\n", apiKey)

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
		g.log.Printf("Existing Repo in service options: %sv\n", serviceUrl)
	}

	urlIsValid := false

	for !urlIsValid {
		serviceUrl = prompts.StringInput(
			"GitHub repository URL",
			serviceUrl,
		)

		// url.URL will say the url is invalid if it ends with a slash ?
		serviceUrl = strings.TrimSuffix(serviceUrl, "/")
		g.log.Printf("User input service url: %s\n", serviceUrl)

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
			g.log.Printf("parts (length: %d): %+v, %s\n", len(parts), parts, u.Hostname())
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

	log.Printf("got service options from URL: %+v\n", serviceOptions)

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
		g.log.Printf("Permission Denied on getGithubPublicKey: %d, %v\n", resp.StatusCode, err)
		g.err = ErrorGithubCIPermissionDenied
		return g

	case resp.StatusCode == 404:
		g.log.Printf("No Such Repository on getGithubPublicKey: %d, %v\n", resp.StatusCode, err)
		g.err = ErrorGithubCINoSuchRepository
		return g

	case err != nil:
		g.log.Printf("Error on getGithubPublicKey: %d, %v\n", resp.StatusCode, err)
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

func (g *gitHubCiService) sendSlot(message models.MessagePayload) *gitHubCiService {
	slots, err := makeSlots(message, GithubCINbSlots, GithubCISLotLength)
	if err != nil {
		g.err = err
		return g
	}

	for i, slot := range slots {
		s := slotName(g.environment, i)
		encryptedSecret := github.EncryptedSecret{}

		g.
			encryptSecret(s, slot, &encryptedSecret).
			sendRepoSecret(&encryptedSecret)

		if g.err != nil {
			return g
		}
	}

	return g
}

func (g *gitHubCiService) sendEnvironmentSecret(
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
	case resp.StatusCode == 401 || resp.StatusCode == 403:
		g.log.Printf("Permission Denied on sendEnvironmentSecret: %d, %v\n", resp.StatusCode, err)
		g.err = ErrorGithubCIPermissionDenied
		return g

	case err != nil:
		g.log.Printf("Error on sendEnvironmentSecret: %d, %v\n", resp.StatusCode, err)
		g.err = err
		return g
	}

	return g
}

func (g *gitHubCiService) sendRepoSecret(
	secret *github.EncryptedSecret,
) *gitHubCiService {
	options := g.getOptions()
	owner := options["Owner"]
	repo := options["Project"]

	resp, err := g.client.Actions.CreateOrUpdateRepoSecret(
		context.Background(),
		owner,
		repo,
		secret,
	)
	switch {
	case resp.StatusCode == 401 || resp.StatusCode == 403:
		g.log.Printf("Permission Denied on sendRepoSecret: %d, %v\n", resp.StatusCode, err)
		g.err = ErrorGithubCIPermissionDenied
		return g

	case err != nil:
		g.log.Printf("Error on sendRepoSecret: %d, %v\n", resp.StatusCode, err)
		g.err = err
		return g
	}

	return g
}

func (g *gitHubCiService) createEnvironment() *gitHubCiService {
	_, resp, err := g.client.Repositories.CreateUpdateEnvironment(
		context.Background(),
		g.servicesKeys["Owner"],
		g.servicesKeys["Project"],
		g.environment,
		&github.CreateUpdateEnvironment{},
	)

	switch {
	case resp.StatusCode == 401 || resp.StatusCode == 403:
		g.log.Printf("Permission Denied on createEnvironment: %d, %v\n", resp.StatusCode, err)
		g.err = ErrorGithubCIPermissionDenied
		return g

	case err != nil:
		g.log.Printf("Error on createEnvironment: %d, %v\n", resp.StatusCode, err)
		g.err = err
		return g
	}

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
			sendEnvironmentSecret(&encryptedSecret, g.environment).
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
		g.log.Printf("Sending file with path: %s\n", file.Path)

		fullpath := path.Join(filecachepath, file.Path)
		freader, err := os.Open(fullpath)
		if err != nil {
			g.log.Printf("Error opening %s\n", fullpath)
			g.err = err
			break
		}

		contents, err := base64encode(freader)
		if err != nil {
			g.log.Printf("Error base64 encoding file %s\n", fullpath)
			break
		}

		// file var name, uses the path to create an uppercase snakecase name
		key := pathToVarname(file.Path)
		// prepend the path, separated with a #
		value := fmt.Sprintf("%s#%s", file.Path, contents)

		encryptedSecret := github.EncryptedSecret{}
		if g.
			encryptSecret(key, value, &encryptedSecret).
			sendEnvironmentSecret(&encryptedSecret, g.environment).
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

	s, resp, err := g.client.Actions.GetEnvSecret(
		context.Background(),
		int(*g.repo.ID),
		g.environment,
		secret,
	)

	switch {
	case s != nil:
		return true

	case err != nil && resp.StatusCode == 404:
		g.log.Printf("Secret %s not found\n", secret)
		return false

	default:
		g.log.Printf("Error trying to figure out if %s exists %v\n", secret, err)
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

	if err != nil {
		g.log.Printf("Error trying to delete secret %s %v\n", secret, err)
	}

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
