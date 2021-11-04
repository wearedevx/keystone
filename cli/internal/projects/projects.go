package projects

import (
	"errors"
	"path"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

var ErrorAlreadyInKeystoneProject = errors.New("already in keystone project")

type ProjectService struct {
	err    error
	ctx    *core.Context
	ksfile *keystonefile.KeystoneFile
	cli    client.KeystoneClient
}

// NewProjectService function returns a new instance of ProjectService
func NewProjectService(
	ctx *core.Context,
	currentfolder, projectName string,
) *ProjectService {
	s := new(ProjectService)
	cli, err := client.NewKeystoneClient()
	if err != nil {
		s.err = err
		return s
	}

	s.ctx = ctx
	s.ksfile = new(keystonefile.KeystoneFile)
	s.cli = cli

	s.load(currentfolder, projectName)

	return s
}

// Err method returns the last error that occurred in the service
func (s *ProjectService) Err() error {
	return s.err
}

// loads the keystone file, but sets the ProjectService error if we
// already are in a Keystone project
func (s *ProjectService) load(currentfolder, projectName string) {
	if keystonefile.ExistsKeystoneFile(currentfolder) {
		s.ksfile.Load(currentfolder)

		// If there is already a keystone file around here,
		// inform the user they are in a keystone project
		if s.ksfile.ProjectId != "" && s.ksfile.ProjectName != projectName {
			// check if .keystone directory too
			if utils.DirExists(path.Join(s.ctx.Wd, ".keystone")) {
				s.err = ErrorAlreadyInKeystoneProject
			}
		}
	} else {
		s.ksfile = keystonefile.NewKeystoneFile(
			currentfolder,
			models.Project{},
		)
	}
}

// GetOrCreate method gets or create a project in organization
func (s *ProjectService) GetOrCreate(
	project *models.Project,
	organizationName string,
) *ProjectService {
	if s.ksfile.ProjectId == "" {
		s.Create(project, organizationName)
	} else {
		s.Get(project)
	}

	return s
}

// Create method creates a project in an organization
func (s *ProjectService) Create(
	project *models.Project,
	organizationName string,
) *ProjectService {
	if s.err != nil {
		return s
	}

	var organizationID uint
	s.pickOrganization(organizationName, &organizationID)
	if s.err != nil {
		return s
	}

	// Remote Project Creation
	*project, s.err = s.cli.Project("").Init(project.Name, organizationID)

	// Handle invalid token
	if s.err != nil {
		return s
	}

	// Update the ksfile
	// So that it keeps secrets and files
	// if the file exited without a project-id
	s.ksfile.ProjectId = project.UUID
	s.ksfile.ProjectName = project.Name

	if s.err = s.ksfile.Save().Err(); s.err != nil {
		return s
	}

	return s
}

// Get method gets the project information from the API
func (s *ProjectService) Get(project *models.Project) *ProjectService {
	if s.err != nil {
		return s
	}

	project.UUID = s.ksfile.ProjectId
	project.Name = s.ksfile.ProjectName

	// But environment data is still on the server
	environments, err := s.cli.Project(project.UUID).GetAccessibleEnvironments()
	if err != nil {
		s.err = err
		return s
	}

	project.Environments = environments

	return s
}

// prompts the user the pick the organizaiton the project is going to belong to
func (s *ProjectService) pickOrganization(
	organizationName string,
	organizationID *uint,
) *ProjectService {
	if s.err != nil {
		return s
	}

	organizations, err := s.cli.Organizations().GetAll()
	if err != nil {
		s.err = err
		return s
	}

	orga := models.Organization{}

	if organizationName == "" {
		orga = prompts.OrganizationsSelect(organizations)
		*organizationID = orga.ID
	} else {
		for _, o := range organizations {
			if organizationName == o.Name {
				orga = o
			}
		}

		if orga.ID == 0 {
			s.err = errors.New("organization not found")
			return s
		}

		*organizationID = orga.ID
	}

	return s
}
