package ci

import (
	"errors"
	"fmt"
	"strings"

	"github.com/wearedevx/keystone/api/pkg/models"

	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

type CiServiceType string

const (
	StubCI CiServiceType = "stub-ci"
)

// A map of available service types and their display name.
// Typically used for UI selection
var availableServices map[CiServiceType]string

var (
	ErrorMissingCiInformation error = errors.New("missing CI information")
	ErrorNoCIServiceWithName        = errors.New(
		"no ci service with that name",
	)
	ErrorUnknownServiceType      = errors.New("unknown service type")
	ErrorInvalidServiceType      = errors.New("invalid service type")
	ErrorNoCIServices            = errors.New("no ci services")
	ErrorNoSecretsForEnvironment = errors.New(
		"no secrets for environment",
	)
)

type CiService interface {
	Name() string
	Type() string
	Usage() string
	Setup() CiService
	GetOptions() map[string]string

	PushSecret(message models.MessagePayload, environment string) CiService
	CleanSecret(environment string) CiService
	CheckSetup() CiService
	Error() error
}

func init() {
	availableServices = map[CiServiceType]string{
		GithubCI:  "GitHub CI",
		GitlabCI:  "Gitlab CI",
		GenericCI: "Generic CI",
	}
}

// GetCiService function returns an instance of CiService.
// The service name should be coming from config file, and will be used
// to determine the type of service to instanciate.
func GetCiService(
	serviceName string,
	ctx *core.Context,
	apiUrl string,
) (CiService, error) {
	var c CiService
	var err error
	var service keystonefile.CiService
	found := false

	services, _ := ListCiServices(ctx)

	for _, s := range services {
		if s.Name == serviceName {
			service = s
			found = true
			break
		}
	}

	if !found {
		return nil, ErrorNoCIServiceWithName
	}

	t := strings.Split(service.Type, "_")[0]

	switch CiServiceType(t) {
	case GithubCI:
		c = GitHubCi(ctx, serviceName, apiUrl)

	case GitlabCI:
		c = GitLabCi(ctx, serviceName, apiUrl)

	case GenericCI:
		c = GenericCi(ctx, serviceName)

	default:
		err = fmt.Errorf(
			"no service type %s: %w",
			service.Type,
			ErrorUnknownServiceType,
		)
	}

	return c, err
}

// Asks the user to pick a type of CI service among the available ones
// It returns a `CiService` instance ready to be setup, and an error
func PickCiService(
	name string,
	ctx *core.Context,
	apiUrl string,
) (CiService, error) {
	var err error
	if err != nil {
		return nil, err
	}

	services := make([]prompts.SelectCIServiceItem, 0)

	for typ, name := range availableServices {
		services = append(services, prompts.SelectCIServiceItem{
			Name: name,
			Type: string(typ),
		})
	}

	s := prompts.SelectCIService(services)

	if name == "" {
		name = s.Name
	}

	switch CiServiceType(s.Type) {
	case GithubCI:
		return GitHubCi(ctx, name, apiUrl), nil
	case GitlabCI:
		return GitLabCi(ctx, name, apiUrl), nil
	case GenericCI:
		return GenericCi(ctx, name), nil
	default:
		return nil, ErrorInvalidServiceType
	}
}

// Asks the user the select a CI Service Configuration
func SelectCiServiceConfiguration(
	serviceName string,
	ctx *core.Context,
	apiUrl string,
) (CiService, error) {
	var err error

	services, err := ListCiServices(ctx)
	if err != nil {
		return nil, err
	}

	if len(services) == 0 {
		return nil, ErrorNoCIServices
	}

	items := make([]string, len(services))

	for idx, service := range services {
		items[idx] = service.Name
	}

	if serviceName == "" {
		_, serviceName = prompts.Select(
			"Select a CI service configuration",
			items,
		)
	}

	return GetCiService(serviceName, ctx, apiUrl)
}

// ListCiServices returns the list of configured CI services
func ListCiServices(ctx *core.Context) (_ []keystonefile.CiService, err error) {
	if ctx.Err() != nil {
		return []keystonefile.CiService{}, nil
	}

	services := []keystonefile.CiService{}

	ksfile := new(keystonefile.KeystoneFile).Load(ctx.Wd)
	if ksfile.Err() != nil {
		return services, err
	}

	return ksfile.CiServices, nil
}

// FindCiServiceWithName returns the CI service configuration
// matching `name`
func FindCiServiceWithName(
	ctx *core.Context,
	name string,
) (service keystonefile.CiService, found bool) {
	if ctx.Err() != nil {
		return service, false
	}

	services, err := ListCiServices(ctx)
	if err != nil {
		return service, false
	}

	for _, candidate := range services {
		if candidate.Name == name {
			service = candidate
			found = true
			break
		}
	}

	return service, found
}

// AddCiService adds a CI service configuration to the keystone file
func AddCiService(ctx *core.Context, service CiService) (err error) {
	if ctx.Err() != nil {
		return nil
	}

	if err = new(keystonefile.KeystoneFile).Load(ctx.Wd).
		AddCiService(keystonefile.CiService{
			Name:    service.Name(),
			Type:    string(service.Type()),
			Options: service.GetOptions(),
		}).
		Save().
		Err(); err != nil {
		return err
	}

	return nil
}

// RemoveCiService remove the CI service configuration matching `serviceName`
// from the keystone file.
func RemoveCiService(ctx *core.Context, serviceName string) (err error) {
	if ctx.Err() != nil {
		return nil
	}

	if err = new(keystonefile.KeystoneFile).Load(ctx.Wd).
		RemoveCiService(serviceName).
		Save().
		Err(); err != nil {
		return err
	}

	return nil
}
