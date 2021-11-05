package ci

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/archive"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/pkg/core"
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
	ApiKey  string `yaml:"api_key"`
	Project string `yaml:"project"`
}

const (
	OPTION_KEY_BASE_URL = "base_url"
	OPTION_KEY_API_KEY  = "api_key"
	OPTION_KEY_PROJECT  = "project"
)

type gitLabCi struct {
	err                 error
	name                string
	apiUrl              string
	ctx                 *core.Context
	apiKey              ApiKey
	client              *gitlab.Client
	options             GitlabOptions
	lastEnvironmentSent string
}

func GitLabCi(ctx *core.Context, name string, apiUrl string) CiService {
	kf := keystonefile.KeystoneFile{}
	kf.Load(ctx.Wd)

	savedService := kf.GetCiService(name)

	ciService := &gitLabCi{
		err:    nil,
		name:   name,
		apiUrl: apiUrl,
		ctx:    ctx,
		apiKey: "",
		client: &gitlab.Client{},
		options: GitlabOptions{
			BaseUrl: savedService.Options[OPTION_KEY_BASE_URL],
			ApiKey:  savedService.Options[OPTION_KEY_API_KEY],
			Project: savedService.Options[OPTION_KEY_PROJECT],
		},
		lastEnvironmentSent: "",
	}

	return ciService
}

// Name method returns the name of the service
func (g *gitLabCi) Name() string { return g.name }

func (g *gitLabCi) Usage() string {
	return fmt.Sprintf(
		`To use them in your pipeline, add the following job in your gitlab-ci.yml:

keystone:
    script: |-
        apt-get install snap
        snap install keystone-cli
        
        echo $KEYSTONE_%s_SLOT > keystone.tar.gz
        tar -xzf keystone.tar.gz         

        eval "$(ks source --env %s)"


You should run this job prior to any other job on your pipline.
`,
		strings.ToUpper(g.lastEnvironmentSent),
		g.lastEnvironmentSent,
	)
}

// Type method returns the type of the service
func (g *gitLabCi) Type() CiServiceType { return GitlabCI }

// Setup method starts th ci service setu process, asking
// the user information through prompts
func (g *gitLabCi) Setup() CiService {
	if g.err != nil {
		return g
	}

	g.options.BaseUrl = g.askForBaseUrl()
	g.options.ApiKey = g.askForPersonalAccessToken()
	g.options.Project = g.askForProjectName()

	return g
}

// GetOptions method returns the service options
func (g *gitLabCi) GetOptions() map[string]string {
	return map[string]string{
		OPTION_KEY_BASE_URL: g.options.BaseUrl,
		OPTION_KEY_API_KEY:  g.options.ApiKey,
		OPTION_KEY_PROJECT:  g.options.Project,
	}
}

// PushSecret method sends a "Message" (that's a completed encrypted environment)
// to GitLab as one project variable
func (g *gitLabCi) PushSecret(
	message models.MessagePayload,
	environment string,
) CiService {
	if g.err != nil {
		return g
	}

	g.initClient()
	if g.err != nil {
		return g
	}

	archive := g.getArchiveBuffer(environment)
	if g.err != nil {
		return g
	}

	vartype := gitlab.EnvVariableType
	key := fmt.Sprintf("KEYSTONE_%s_SLOT", strings.ToUpper(environment))
	protected, masked := true, true
	envScope := "*"

	_, _, err := g.client.ProjectVariables.CreateVariable(
		g.options.Project,
		&gitlab.CreateProjectVariableOptions{
			Key:              &key,
			Value:            &archive,
			VariableType:     &vartype,
			Protected:        &protected,
			Masked:           &masked,
			EnvironmentScope: &envScope,
		},
	)
	if err != nil {
		g.err = err
	}

	g.lastEnvironmentSent = environment

	return g
}

// CleanSecret method remove all the secrets for the given environment
// from the CI service
func (g *gitLabCi) CleanSecret(environment string) CiService {
	if g.err != nil {
		return g
	}

	g.initClient()
	if g.err != nil {
		return g
	}

	key := fmt.Sprintf("KEYSTONE_%s_SLOT", strings.ToUpper(environment))

	_, err := g.client.ProjectVariables.RemoveVariable(
		g.options.Project,
		key,
	)
	if err != nil {
		g.err = err
	}

	return g
}

// CheckSetup method verifies the user submitted information is valid
func (g *gitLabCi) CheckSetup() CiService {
	if g.err != nil {
		return g
	}

	if len(g.options.BaseUrl) == 0 ||
		len(g.options.ApiKey) == 0 ||
		len(g.options.Project) == 0 {
		g.err = ErrorMissingCiInformation
	}

	return g
}

// Error method returns the last error encountered
func (g *gitLabCi) Error() error {
	return g.err
}

func (g *gitLabCi) askForBaseUrl() string {
	original := g.options.BaseUrl

	if original == "" {
		original = "gitlab.com"
	}

	b := prompts.StringInput("Gitlab base URL:", original)

	return b
}

func (g *gitLabCi) askForPersonalAccessToken() string {
	t := prompts.StringInput("Gitlab Personal Access Token:", g.options.ApiKey)
	return t
}

func (g *gitLabCi) askForProjectName() string {
	p := prompts.StringInput(
		"Gitlab project (namespace/project_path):",
		g.options.Project,
	)
	return p
}

func (g *gitLabCi) initClient() {
	g.client, g.err = gitlab.NewClient(g.options.ApiKey)
}

func (g *gitLabCi) getFileList(environmentName string) []utils.FileInfo {
	if g.err != nil {
		return nil
	}

	fileList := make([]utils.FileInfo, 0)
	source := g.ctx.DotKeystonePath()
	prefix := filepath.Join("cache", environmentName)

	err := utils.DirWalk(source,
		func(info utils.FileInfo) error {
			if strings.HasPrefix(info.Path, prefix) {
				fileList = append(fileList, info)
			}

			return nil
		})
	if err != nil {
		g.err = err
		return nil
	}

	return fileList
}

func (g *gitLabCi) getArchiveBuffer(environmentName string) string {
	if g.err != nil {
		return ""
	}

	fileList := g.getFileList(environmentName)
	if g.err != nil {
		return ""
	}

	buffer, err := archive.TarFileList(fileList)
	if err != nil {
		g.err = err
		return ""
	}

	gzipBuffer, err := archive.Gzip(buffer)
	if err != nil {
		g.err = err
		return ""
	}

	sb := new(strings.Builder)
	_, err = io.Copy(sb, gzipBuffer)
	if err != nil {
		g.err = err
		return ""
	}

	return sb.String()
}
