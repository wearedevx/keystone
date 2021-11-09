package ci

import (
	"encoding/base64"
	"fmt"
	"io"
	"strings"

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
	err                 error
	name                string
	apiUrl              string
	ctx                 *core.Context
	apiKey              ApiKey
	client              *gitlab.Client
	options             GitlabOptions
	lastEnvironmentSent string
	nSlots              int
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
		lastEnvironmentSent: "",
		nSlots:              1,
	}

	return ciService
}

// Name method returns the name of the service
func (g *gitlabCiService) Name() string { return g.name }

// Usage method returns a usage string that will be displayed
// to the user
func (g *gitlabCiService) Usage() string {
	slots := make([]string, g.nSlots)

	for i := range slots {
		slots[i] = slot(g.lastEnvironmentSent, i)
	}

	return ui.RenderTemplate(
		"gitlab-ci-usage",
		`To use them in your pipeline, add the following job in your gitlab-ci.yml:

default:
  before_script:
    - |
      archive="\{{ range .Slots}}
      ${{.}}\{{end}}
      "
    - echo -n $archive | base64 -d > keystone.tar.gz
    - tar -xzf keystone.tar.gz; rm keystone.tar.gz
    - unset archive;
    - set -o allexport; source .keystone/cache/{{ .Environment }}/.env; set +o allexport
    - |
      if [ "$(ls -A .keystone/cache/{{ .Environment }}/files)" ]; then
        cp -R .keystone/cache/{{ .Environment }}/files/* ./;
      fi
`,
		map[string]interface{}{
			"Slots":       slots,
			"Environment": g.lastEnvironmentSent,
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

	var err error

	g.initClient()
	if g.err != nil {
		return g
	}

	archive, err := getArchiveBuffer(g.ctx, environment)
	if err != nil {
		g.err = err
		return g
	}

	str, err := base64encode(archive)
	if err != nil {
		return g
	}

	splits, err := splitString(str, SLOT_SIZE, N_SLOTS)
	if err != nil {
		g.err = err
		return g
	}

	nslots := 0

	for i, split := range splits {
		key := slot(environment, i)
		remoteHasVariable := g.hasVariable(key)
		hasContent := len(split) > 0

		if hasContent {
			nslots++
		}

		switch {
		case remoteHasVariable && hasContent:
			g.updateVariable(key, split)
		case remoteHasVariable && !hasContent:
			g.deleteVariable(key)
		case !remoteHasVariable:
			g.createVariable(key, split)
		}

		if err != nil {
			g.err = err
			return g
		}
	}

	g.lastEnvironmentSent = environment
	g.nSlots = nslots

	return g
}

func (g *gitlabCiService) hasVariable(key string) bool {
	variable, _, _ := g.client.ProjectVariables.GetVariable(
		g.options.Project,
		key,
	)

	return variable != nil
}

func (g *gitlabCiService) createVariable(key string, value string) {
	if len(value) == 0 {
		return
	}

	_, _, err := g.client.ProjectVariables.CreateVariable(
		g.options.Project,
		&gitlab.CreateProjectVariableOptions{
			Key:          gitlab.String(key),
			Value:        gitlab.String(value),
			VariableType: gitlab.VariableType("env_var"),
			Masked:       gitlab.Bool(true),
		},
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
			Value:        gitlab.String(value),
			VariableType: gitlab.VariableType("env_var"),
			Masked:       gitlab.Bool(true),
		},
	)
	if err != nil {
		g.err = err
	}
}

func (g *gitlabCiService) deleteVariable(key string) {
	_, err := g.client.ProjectVariables.RemoveVariable(
		g.options.Project,
		key,
	)
	if err != nil {
		g.err = err
	}
}

// CleanSecret method remove all the secrets for the given environment
// from the CI service
func (g *gitlabCiService) CleanSecret(environment string) CiService {
	if g.err != nil {
		return g
	}

	g.initClient()
	if g.err != nil {
		return g
	}

	for i := range make([]int, 5) {
		key := slot(environment, i)

		if g.hasVariable(key) {
			if g.deleteVariable(key); g.err != nil {
				return g
			}
		}
	}

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

func (g *gitlabCiService) initClient() {
	g.client, g.err = gitlab.NewClient(string(g.apiKey))
}

func slot(environmentName string, i int) string {
	return fmt.Sprintf(
		"KEYSTONE_%s_SLOT_%d",
		strings.ToUpper(environmentName),
		i+1,
	)
}

func splitString(s string, chunkSize int, nChunks int) ([]string, error) {
	if len(s) == 0 {
		return nil, nil
	}

	chunks := make([]string, nChunks)

	if chunkSize >= len(s) {
		chunks[0] = s
		return chunks, nil
	}

	c := 0
	currentLen := 0
	currentStart := 0

	for i := range s {
		if currentLen == chunkSize {
			chunks[c] = s[currentStart:i]
			currentLen = 0
			currentStart = i

			c += 1

			if c == len(chunks)-1 {
				break
			}
		}

		currentLen++
	}

	lastChunk := s[currentStart:]
	if len(lastChunk) > chunkSize {
		return nil, fmt.Errorf("keystone archive too big: %d", len(s))
	}

	chunks[c] = lastChunk

	return chunks, nil
}

func base64encode(reader io.Reader) (string, error) {
	sb := new(strings.Builder)

	_, err := io.Copy(sb, reader)
	if err != nil {
		return "", err
	}

	s := base64.StdEncoding.EncodeToString([]byte(sb.String()))

	return s, err
}
