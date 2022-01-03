package environments

import (
	"errors"
	"log"
	"strings"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/config"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/spinner"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/client/auth"
	"github.com/wearedevx/keystone/cli/pkg/core"
)

var ErrorNoAccess = errors.New("has no access")

type environmentService struct {
	log    *log.Logger
	err    *kserrors.Error
	ctx    *core.Context
	client client.KeystoneClient
}

type EnvironmentService interface {
	Err() *kserrors.Error
	GetAccessibleEnvironments() []models.Environment
}

// NewEnvironmentService function return an instance of EnvironmentService
func NewEnvironmentService(ctx *core.Context) EnvironmentService {
	var ksc client.KeystoneClient
	var err *kserrors.Error

	if err = ctx.Err(); err == nil {
		ksc, err = client.NewKeystoneClient()
	}

	s := &environmentService{
		log:    log.New(log.Writer(), "[Environments] ", 0),
		err:    err,
		ctx:    ctx,
		client: ksc,
	}

	return s
}

// Err method returns the last error encountered
func (s *environmentService) Err() *kserrors.Error {
	return s.err
}

// GetAccessibleEnvironments method returns the environments the currently
// logged in user has access to.
func (s *environmentService) GetAccessibleEnvironments() []models.Environment {
	if s.err != nil {
		return []models.Environment{}
	}
	projectID := s.ctx.GetProjectID()

	sp := spinner.Spinner("")
	sp.Start()

	accessibleEnvironments, err := s.client.Project(projectID).
		GetAccessibleEnvironments()
	sp.Stop()

	if err != nil {
		if errors.Is(err, auth.ErrorUnauthorized) {
			config.Logout()
			s.ctx.SetError(kserrors.InvalidConnectionToken(err))
		} else if errors.Is(err, auth.ErrorServiceNotAvailable) {
			s.ctx.SetError(kserrors.ServiceNotAvailable(err))
		} else if strings.Contains(err.Error(), auth.ErrorDeviceNotRegistered.Error()) {
			s.ctx.SetError(kserrors.DeviceNotRegistered(err))
		} else if strings.Contains(err.Error(), "not found") {
			s.ctx.SetError(kserrors.ProjectDoesntExist("", "", err))
		} else {
			s.ctx.SetError(kserrors.UnkownError(err))
		}
	}

	return accessibleEnvironments
}
