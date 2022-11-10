package ci

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/xanzy/go-gitlab"

	"github.com/wearedevx/keystone/cli/internal/config"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

var (
	gitlabClientID     string
	gitlabClientSecret string
)

const GitlabCI CiServiceType = "gitlab-ci"

type GitlabOptions struct {
	BaseURL string `yaml:"base_url"`
	Project string `yaml:"project"`
}

const (
	OptionKeyBaseURL = "base_url"
	OptionKeyAPIKey  = "api_key"
	OptionKeyProject = "project"
)

const (
	SlotSize = 1024
	NSlots   = 5
)

type gitlabCiService struct {
	log           *log.Logger
	err           error
	name          string
	apiURL        string
	ctx           *core.Context
	apiKey        ApiKey
	client        *gitlab.Client
	options       GitlabOptions
	environment   string
	fileVariables []string
}

// configServiceName function returns the key to find the ApiKey
// in the configuration file
func configServiceName(baseURL string) string {
	domain := strings.TrimPrefix(baseURL, "https://")
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.ReplaceAll(domain, ".", "#")

	return fmt.Sprintf("%s_%s", GitlabCI, domain)
}

// GitLabCi function return a `CiService` that works with the GitLab API
func GitLabCi(ctx *core.Context, name string, apiURL string) CiService {
	kf := keystonefile.KeystoneFile{}
	kf.Load(ctx.Wd)

	savedService := kf.GetCiService(name)

	apiKey := config.GetServiceApiKey(
		configServiceName(savedService.Options[OptionKeyBaseURL]),
	)

	ciService := &gitlabCiService{
		err:    nil,
		log:    log.New(log.Writer(), "[GitlabCi] ", 0),
		name:   name,
		apiURL: apiURL,
		ctx:    ctx,
		apiKey: ApiKey(apiKey),
		client: &gitlab.Client{},
		options: GitlabOptions{
			BaseURL: savedService.Options[OptionKeyBaseURL],
			Project: savedService.Options[OptionKeyProject],
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
	baseURL := g.options.BaseURL

	return configServiceName(baseURL)
}

// Setup method starts th ci service setu process, asking
// the user information through prompts
func (g *gitlabCiService) Setup() CiService {
	if g.err != nil {
		return g
	}

	g.options.BaseURL = g.askForBaseURL()
	g.apiKey = ApiKey(g.askForPersonalAccessToken())
	g.options.Project = g.askForProjectName()

	config.SetServiceApiKey(
		configServiceName(g.options.BaseURL),
		string(g.apiKey),
	)
	config.Write()

	return g
}

// GetOptions method returns the service options
func (g *gitlabCiService) GetOptions() map[string]string {
	return map[string]string{
		OptionKeyBaseURL: g.options.BaseURL,
		OptionKeyProject: g.options.Project,
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
	g.log.Printf(
		"Sending secrets to %s/%s on environment %s\n",
		g.options.BaseURL,
		g.options.Project,
		environment,
	)

	g.initClient().
		createEnvironment().
		sendEnvironmentSecrets().
		sendEnvironmentFiles()

	return g
}

// Adds the environment scope option to gitlab requests
// to set secrets for specific environments
func (g *gitlabCiService) environmentScopeOption() func(*retryablehttp.Request) error {
	return func(req *retryablehttp.Request) error {
		query := req.URL.Query()
		query.Add("filter[environment_scope]", g.environment)

		req.URL.RawQuery = query.Encode()

		return nil
	}
}

func (g *gitlabCiService) hasVariable(key string) bool {
	variable, _, err := g.client.ProjectVariables.GetVariable(
		g.options.Project,
		key,
		&gitlab.GetProjectVariableOptions{
			Filter: &gitlab.VariableFilter{
				EnvironmentScope: g.environment,
			},
		},
	)
	if err != nil {
		g.log.Printf(
			"[Warning] an error occurred getting variable %s: %v\n",
			key,
			err,
		)
	}

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

	g.log.Printf("Creating new variable %s, with value %s\n", key, value)

	_, _, err := g.client.ProjectVariables.CreateVariable(
		g.options.Project,
		options,
	)
	if err != nil {
		g.err = err
	}
}

func (g *gitlabCiService) updateVariable(key string, value string) {
	g.log.Printf("Updating variable %s, with value %s\n", key, value)

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
	g.log.Printf("Deleting variable %sn", key)

	_, err := g.client.ProjectVariables.RemoveVariable(
		g.options.Project,
		key,
		&gitlab.RemoveProjectVariableOptions{
			Filter: &gitlab.VariableFilter{
				EnvironmentScope: g.environment,
			},
		},
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

	g.log.Printf(
		"Current setup for %s:\n\tBaseURL: %s,\n\t,Project: %s\n\t,ApiKey: %s\n",
		g.Name(),
		g.options.BaseURL,
		g.options.Project,
		g.apiKey,
	)

	if len(g.options.BaseURL) == 0 ||
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

func (g *gitlabCiService) askForBaseURL() string {
	original := g.options.BaseURL

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
		g.log.Printf("Creating new environment %s on remote\n", g.environment)
		_, _, err := g.client.Environments.CreateEnvironment(
			g.options.Project,
			&gitlab.CreateEnvironmentOptions{
				Name: gitlab.String(g.environment),
			},
		)
		if err != nil {
			g.err = err
		}
	} else {
		g.log.Printf("Environment %s already exists on remote\n", g.environment)
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

		g.log.Printf(
			"Sending secret %s (environemt: %s, value: %s)",
			key,
			value,
			g.environment,
		)

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

		g.log.Printf(
			"Sending file %s (key: %s, environment: %s)\n",
			file.Path,
			key,
			g.environment,
		)

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
