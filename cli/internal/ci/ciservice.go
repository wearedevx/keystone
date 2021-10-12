package ci

import (
	"errors"
	"fmt"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

type CiServiceType string

const (
	StubCI   CiServiceType = "stub-ci"
	GithubCI CiServiceType = "github-ci"
)

var availableServices map[CiServiceType]string

var (
	ErrorMissinCiInformation error = errors.New("missing CI information")
	ErrorNoCIServiceWithName       = errors.New("no ci service with that name")
	ErrorUnknownServiceType        = errors.New("unknown service type")
	ErrorInvalidServiceType        = errors.New("invalid service type")
	ErrorNoCIServices              = errors.New("no ci services")
)

type CiService interface {
	Name() string
	Type() CiServiceType
	Setup() CiService
	GetOptions() map[string]string

	PushSecret(message models.MessagePayload, environment string) CiService
	CleanSecret(environment string) CiService
	CheckSetup() CiService
	Error() error
	PrintSuccess(string)
	// // Finish(pkey []byte) (models.User, string, error)
	// GetKeys() ServicesKeys
	// SetKeys(ServicesKeys) error
	// GetApiKey() ApiKey
	// SetApiKey(ApiKey)
}

func init() {
	availableServices = map[CiServiceType]string{
		GithubCI: "GitHub CI",
	}
}

func GetCiService(serviceName string, ctx *core.Context, apiUrl string) (CiService, error) {
	var c CiService
	var err error
	var service keystonefile.CiService
	var found bool = false

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

	switch CiServiceType(service.Type) {
	case GithubCI:
		c = GitHubCi(ctx, serviceName, apiUrl)

	default:
		err = fmt.Errorf(
			"No service type %s: %w",
			service.Type,
			ErrorUnknownServiceType,
		)
	}

	return c, err
}

func PickCiService(name string, ctx *core.Context, apiUrl string) (CiService, error) {
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
	default:
		return nil, ErrorInvalidServiceType
	}

}

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

func FindCiServiceWithName(ctx *core.Context, name string) (service keystonefile.CiService, found bool) {
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

func AddCiService(ctx *core.Context, service CiService) (err error) {
	if ctx.Err() != nil {
		return nil
	}

	if err = new(keystonefile.KeystoneFile).
		Load(ctx.Wd).
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

func RemoveCiService(ctx *core.Context, serviceName string) (err error) {
	if ctx.Err() != nil {
		return nil
	}

	if err = new(keystonefile.KeystoneFile).
		Load(ctx.Wd).
		RemoveCiService(serviceName).
		Save().
		Err(); err != nil {
		return err
	}

	return nil
}
