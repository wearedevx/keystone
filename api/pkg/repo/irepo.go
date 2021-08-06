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
	DeleteLoginRequest(string) bool
	DeleteMessage(messageID uint, userID uint) IRepo
	DeleteExpiredMessages() IRepo
	Err() error
	FindUsers(userIDs []string, users *map[string]models.User, notFounds *[]string) IRepo
	GetChildrenRoles(role models.Role, roles *[]models.Role) IRepo
	GetDb() *gorm.DB
	GetEnvironment(*models.Environment) IRepo
	GetEnvironmentPublicKeys(envID string, publicKeys *models.PublicKeys) IRepo
	GetEnvironmentType(*models.EnvironmentType) IRepo
	GetEnvironmentsByProjectUUID(projectUUID string, foundEnvironments *[]models.Environment) IRepo
	GetInvitableRoles(models.Role, *[]models.Role) IRepo
	GetLoginRequest(string) (models.LoginRequest, bool)
	GetMessagesForUserOnEnvironment(user models.User, environment models.Environment, message *models.Message) IRepo
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
	GetRole(*models.Role) IRepo
	GetRoles(*[]models.Role) IRepo
	GetRolesEnvironmentType(*models.RolesEnvironmentType) IRepo
	GetRolesMemberCanInvite(projectMember models.ProjectMember, roles *[]models.Role) IRepo
	GetUser(*models.User) IRepo
	GetUserByEmail(string, *[]models.User) IRepo
	ListProjectMembers(userIDList []string, projectMember *[]models.ProjectMember) IRepo
	ProjectAddMembers(models.Project, []models.MemberRole, models.User) IRepo
	ProjectGetMembers(*models.Project, *[]models.ProjectMember) IRepo
	ProjectLoadUsers(*models.Project) IRepo
	ProjectRemoveMembers(models.Project, []string) IRepo
	ProjectSetRoleForUser(models.Project, models.User, models.Role) IRepo
	CheckMembersAreInProject(models.Project, []string) ([]string, error)
	RemoveOldMessageForRecipient(userID uint, environmentID string) IRepo
	SetLoginRequestCode(string, string) models.LoginRequest
	SetNewVersionID(environment *models.Environment) error
	WriteMessage(user models.User, message models.Message) IRepo
}
