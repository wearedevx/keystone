package repo

import (
	"github.com/wearedevx/keystone/api/pkg/models"
	"gorm.io/gorm"
)

// IRepo The Repo interface
type IRepo interface {
	CreateEnvironment(*models.Environment) IRepo
	CreateEnvironmentType(*models.EnvironmentType) IRepo
	CreateLoginRequest() models.LoginRequest
	CreateProjectMember(*models.ProjectMember, *models.Role) IRepo
	CreateRole(*models.Role) IRepo
	CreateRoleEnvironmentType(*models.RolesEnvironmentType) IRepo
	CreateSecret(*models.Secret)
	DeleteLoginRequest(string) bool
	Err() error
	FindUsers([]string) (map[string]models.User, []string)
	GetDb() *gorm.DB
	GetEnvironment(*models.Environment) IRepo
	GetEnvironmentType(*models.EnvironmentType) IRepo
	GetInvitableRoles(models.Role, *[]models.Role) IRepo
	GetLoginRequest(string) (models.LoginRequest, bool)
	GetOrCreateEnvironment(*models.Environment) IRepo
	GetOrCreateEnvironmentType(*models.EnvironmentType) IRepo
	GetOrCreateProject(*models.Project) IRepo
	GetOrCreateProjectMember(*models.ProjectMember, string) IRepo
	GetOrCreateRole(*models.Role) IRepo
	GetOrCreateRoleEnvType(*models.RolesEnvironmentType) IRepo
	GetOrCreateUser(*models.User) IRepo
	GetProject(*models.Project) IRepo
	GetProjectByUUID(string, *models.Project) IRepo
	GetProjectMember(*models.ProjectMember) IRepo
	GetRoleByID(uint, *models.Role) IRepo
	GetRoleByName(string, *models.Role) IRepo
	GetRoles(*[]models.Role) IRepo
	GetRolesEnvironmentType(*models.RolesEnvironmentType) IRepo
	GetSecretByName(string, *models.Secret)
	GetUser(*models.User) IRepo
	ProjectAddMembers(models.Project, []models.MemberRole) IRepo
	ProjectGetMembers(*models.Project, *[]models.ProjectMember) IRepo
	ProjectLoadUsers(*models.Project) IRepo
	ProjectRemoveMembers(models.Project, []string) IRepo
	ProjectSetRoleForUser(models.Project, models.User, models.Role) IRepo
	SetLoginRequestCode(string, string) models.LoginRequest
}
