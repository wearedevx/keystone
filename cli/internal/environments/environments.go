package environments

import (
	"errors"

	"github.com/wearedevx/keystone/api/pkg/models"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/spinner"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/client/auth"
	"github.com/wearedevx/keystone/cli/pkg/core"
)

type environmentService struct {
	err    *kserrors.Error
	ctx    *core.Context
	client client.KeystoneClient
}

type EnvironmentService interface {
	Err() *kserrors.Error
	GetAccessibleEnvironments() []models.Environment
}

func NewEnvironmentService(ctx *core.Context) EnvironmentService {
	var ksc client.KeystoneClient
	var err *kserrors.Error

	if err = ctx.Err(); err == nil {
		ksc, err = client.NewKeystoneClient()
	}

	s := &environmentService{
		err:    err,
		ctx:    ctx,
		client: ksc,
	}

	return s
}

func (s *environmentService) Err() *kserrors.Error {
	return s.err
}

func (s *environmentService) GetAccessibleEnvironments() []models.Environment {
	if s.err != nil {
		return []models.Environment{}
	}
	projectID := s.ctx.GetProjectID()

	sp := spinner.Spinner(" ")
	sp.Start()

	accessibleEnvironments, err := s.client.Project(projectID).GetAccessibleEnvironments()
	sp.Stop()

	if err != nil {
		if errors.Is(err, auth.ErrorUnauthorized) {
			s.ctx.SetError(kserrors.InvalidConnectionToken(err))
		} else {
			s.ctx.SetError(kserrors.UnkownError(err))
		}
	}

	return accessibleEnvironments
}
