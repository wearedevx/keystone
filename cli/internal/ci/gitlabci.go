package ci

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/config"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
	"github.com/wearedevx/keystone/cli/ui/prompts"
	"github.com/xanzy/go-gitlab"
)

var (
	gitlabClientId     string
	gitlabClientSecret string
)

const GitlabCI CiServiceType = "gitlab-ci"

type GitlabOptions struct {
	BaseUrl string `yaml:"base_url"`
	Project string `yaml:"project"`
}

const (
	OPTION_KEY_BASE_URL = "base_url"
	OPTION_KEY_API_KEY  = "api_key"
	OPTION_KEY_PROJECT  = "project"
)

const (
	SLOT_SIZE = 1024
	N_SLOTS   = 5
)

type gitlabCiService struct {
	err           error
	name          string
	apiUrl        string
	ctx           *core.Context
	apiKey        ApiKey
	client        *gitlab.Client
	options       GitlabOptions
	environment   string
	fileVariables []string
}

// configServiceName function returns the key to find the ApiKey
// in the configuration file
func configServiceName(baseUrl string) string {
	domain := strings.TrimPrefix(baseUrl, "https://")
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.ReplaceAll(domain, ".", "#")

	return fmt.Sprintf("%s_%s", GitlabCI, domain)
}

// GitLabCi function return a `CiService` that works with the GitLab API
func GitLabCi(ctx *core.Context, name string, apiUrl string) CiService {
	kf := keystonefile.KeystoneFile{}
	kf.Load(ctx.Wd)

	savedService := kf.GetCiService(name)

	apiKey := config.GetServiceApiKey(
		configServiceName(savedService.Options[OPTION_KEY_BASE_URL]),
	)

	ciService := &gitlabCiService{
		err:    nil,
		name:   name,
		apiUrl: apiUrl,
		ctx:    ctx,
		apiKey: ApiKey(apiKey),
		client: &gitlab.Client{},
		options: GitlabOptions{
			BaseUrl: savedService.Options[OPTION_KEY_BASE_URL],
			Project: savedService.Options[OPTION_KEY_PROJECT],
		},
		environment:   "",
		fileVariables: []string{},
	}

	return ciService
}

// Name method returns the name of the service
func (g *gitlabCiService) Name() string { return g.name }

// Usage method returns a usage string that will be displayed
// to the user
func (g *gitlabCiService) Usage() string {
	return ui.RenderTemplate(
		"gitlab-ci-usage",
		`To use them in your pipeline, add the following job in your gitlab-ci.yml:

default:
  before_script:
    {{ range .FileVariables }}- mkdir -p $(dirname ${ {{- . -}} %#*})
    - echo "${ {{- . -}} #*#}" | base64 -d > ${ {{- . -}} %#*}
    {{end}}
`,
		map[string]interface{}{
			"FileVariables": g.fileVariables,
			"Environment":   g.environment,
		},
	)
}

// Type method returns the type of the service
func (g *gitlabCiService) Type() string {
	baseUrl := g.options.BaseUrl

	return configServiceName(baseUrl)
}

// Setup method starts th ci service setu process, asking
// the user information through prompts
func (g *gitlabCiService) Setup() CiService {
	if g.err != nil {
		return g
	}

	g.options.BaseUrl = g.askForBaseUrl()
	g.apiKey = ApiKey(g.askForPersonalAccessToken())
	g.options.Project = g.askForProjectName()

	config.SetServiceApiKey(
		configServiceName(g.options.BaseUrl),
		string(g.apiKey),
	)
	config.Write()

	return g
}

// GetOptions method returns the service options
func (g *gitlabCiService) GetOptions() map[string]string {
	return map[string]string{
		OPTION_KEY_BASE_URL: g.options.BaseUrl,
		OPTION_KEY_PROJECT:  g.options.Project,
	}
}

// PushSecret method sends a "Message" (that's a completed encrypted environment)
// to GitLab as one project variable
func (g *gitlabCiService) PushSecret(
	message models.MessagePayload,
	environment string,
) CiService {
	if g.err != nil {
		return g
	}

	g.environment = environment

	g.initClient().
		createEnvironment().
		sendEnvironmentSecrets().
		sendEnvironmentFiles()

	return g
}

func (g *gitlabCiService) environmentScopeOption() func(*retryablehttp.Request) error {
	return func(req *retryablehttp.Request) error {
		query := req.URL.Query()
		query.Add("filter[environment_scope]", g.environment)

		req.URL.RawQuery = query.Encode()

		return nil
	}
}

func (g *gitlabCiService) hasVariable(key string) bool {
	variable, _, _ := g.client.ProjectVariables.GetVariable(
		g.options.Project,
		key,
		g.environmentScopeOption(),
	)

	return variable != nil
}

func (g *gitlabCiService) createVariable(key string, value string) {
	if len(value) == 0 {
		return
	}

	options := &gitlab.CreateProjectVariableOptions{
		Key:              gitlab.String(key),
		Value:            gitlab.String(value),
		VariableType:     gitlab.VariableType("env_var"),
		EnvironmentScope: gitlab.String(g.environment),
	}

	_, _, err := g.client.ProjectVariables.CreateVariable(
		g.options.Project,
		options,
	)
	if err != nil {
		g.err = err
	}
}

func (g *gitlabCiService) updateVariable(key string, value string) {
	_, _, err := g.client.ProjectVariables.UpdateVariable(
		g.options.Project,
		key,
		&gitlab.UpdateProjectVariableOptions{
			Value:            gitlab.String(value),
			VariableType:     gitlab.VariableType("env_var"),
			EnvironmentScope: gitlab.String(g.environment),
		},
	)
	if err != nil {
		g.err = err
	}
}

func (g *gitlabCiService) deleteVariable(key string) *gitlabCiService {
	_, err := g.client.ProjectVariables.RemoveVariable(
		g.options.Project,
		key,
		g.environmentScopeOption(),
	)
	if err != nil {
		g.err = err
	}

	return g
}

// CleanSecret method remove all the secrets for the given environment
// from the CI service
func (g *gitlabCiService) CleanSecret(environment string) CiService {
	if g.err != nil {
		return g
	}

	g.environment = environment

	g.initClient().
		cleanSecrets().
		cleanFiles()

	return g
}

// CheckSetup method verifies the user submitted information is valid
func (g *gitlabCiService) CheckSetup() CiService {
	if g.err != nil {
		return g
	}

	if len(g.options.BaseUrl) == 0 ||
		len(g.apiKey) == 0 ||
		len(g.options.Project) == 0 {
		g.err = ErrorMissingCiInformation
	}

	return g
}

// Error method returns the last error encountered
func (g *gitlabCiService) Error() error {
	return g.err
}

func (g *gitlabCiService) askForBaseUrl() string {
	original := g.options.BaseUrl

	if original == "" {
		original = "gitlab.com"
	}

	b := prompts.StringInput("Gitlab base URL:", original)

	return b
}

func (g *gitlabCiService) askForPersonalAccessToken() string {
	t := prompts.StringInput("Gitlab Personal Access Token:", string(g.apiKey))
	return t
}

func (g *gitlabCiService) askForProjectName() string {
	p := prompts.StringInput(
		"Gitlab project (namespace/project_path):",
		g.options.Project,
	)
	return p
}

func (g *gitlabCiService) initClient() *gitlabCiService {
	g.client, g.err = gitlab.NewClient(string(g.apiKey))

	return g
}

func (g *gitlabCiService) createEnvironment() *gitlabCiService {
	environments, _, err := g.client.Environments.ListEnvironments(
		g.options.Project,
		&gitlab.ListEnvironmentsOptions{
			Name: gitlab.String(g.environment),
		},
	)
	if err != nil {
		g.err = err
		return g
	}

	if len(environments) == 0 {
		_, _, err := g.client.Environments.CreateEnvironment(
			g.options.Project,
			&gitlab.CreateEnvironmentOptions{
				Name: gitlab.String(g.environment),
			},
		)
		if err != nil {
			g.err = err
		}
	}

	return g
}

func (g *gitlabCiService) createOrUpdateVariable(
	key, value string,
) *gitlabCiService {
	if g.err != nil {
		return g
	}

	if g.hasVariable(key) {
		g.updateVariable(key, value)
	} else {
		g.createVariable(key, value)
	}

	return g
}

func (g *gitlabCiService) sendEnvironmentSecrets() *gitlabCiService {
	if g.err != nil {
		return g
	}

	secrets := g.ctx.ListSecrets()

	for _, secret := range secrets {
		key := secret.Name
		value, ok := secret.Values[core.EnvironmentName(g.environment)]
		if !ok && secret.Required {
			g.err = fmt.Errorf("required secret is missing %s", key)
			break
		}

		if g.createOrUpdateVariable(key, string(value)).err != nil {
			break
		}
	}

	return g
}

func (g *gitlabCiService) sendEnvironmentFiles() *gitlabCiService {
	if g.err != nil {
		return g
	}

	files := g.ctx.ListFiles()
	filecachepath := g.ctx.CachedEnvironmentFilesPath(g.environment)

	for _, file := range files {
		fullpath := path.Join(filecachepath, file.Path)
		f, err := os.Open(fullpath)
		if err != nil {
			g.err = err
			break
		}

		contents, err := base64encode(f)
		if err != nil {
			g.err = err
			break
		}

		key := pathToVarname(file.Path)
		value := fmt.Sprintf("%s#%s", file.Path, contents)

		g.createOrUpdateVariable(key, value)

		g.fileVariables = append(g.fileVariables, key)
	}

	return g
}

func (g *gitlabCiService) cleanSecrets() *gitlabCiService {
	if g.err != nil {
		return g
	}

	secrets := g.ctx.ListSecrets()
	for _, secret := range secrets {
		key := secret.Name

		if g.hasVariable(key) {
			g.deleteVariable(key)
		}
	}
	return g
}

func (g *gitlabCiService) cleanFiles() *gitlabCiService {
	if g.err != nil {
		return g
	}

	files := g.ctx.ListFiles()
	for _, file := range files {
		key := pathToVarname(file.Path)

		if g.hasVariable(key) {
			g.deleteVariable(key)
		}

	}

	return g
}
