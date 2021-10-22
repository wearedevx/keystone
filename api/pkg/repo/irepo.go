package repo

import (
	"github.com/wearedevx/keystone/api/internal/emailer"
	"github.com/wearedevx/keystone/api/pkg/message"
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
	DeleteAllProjectMembers(project *models.Project) IRepo
	DeleteExpiredMessages() IRepo
	GetGroupedMessagesWillExpireByUser(groupedMessageUser *map[uint]emailer.GroupedMessagesUser) IRepo
	DeleteMessage(messageID uint, userID uint) IRepo
	DeleteProject(project *models.Project) IRepo
	DeleteProjectsEnvironments(project *models.Project) IRepo
	Err() error
	FindUsers(userIDs []string, users *map[string]models.User, notFounds *[]string) IRepo
	GetActivityLogs(projectID string, options models.GetLogsOptions, logs *[]models.ActivityLog) IRepo
	GetChildrenRoles(role models.Role, roles *[]models.Role) IRepo
	GetDb() *gorm.DB
	GetEnvironment(*models.Environment) IRepo
	GetEnvironmentPublicKeys(envID string, publicKeys *models.PublicKeys) IRepo
	GetEnvironmentType(*models.EnvironmentType) IRepo
	GetEnvironmentsByProjectUUID(projectUUID string, foundEnvironments *[]models.Environment) IRepo
	GetInvitableRoles(models.Role, *[]models.Role) IRepo
	GetLoginRequest(string) (models.LoginRequest, bool)
	GetMessage(message *models.Message) IRepo
	GetMessagesForUserOnEnvironment(device models.Device, environment models.Environment, message *models.Message) IRepo
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
	GetProjectsOrganization(string, *models.Organization) IRepo
	OrganizationCountMembers(*models.Organization, *int64) IRepo
	GetRole(*models.Role) IRepo
	GetRoles(*[]models.Role) IRepo
	GetRolesEnvironmentType(*models.RolesEnvironmentType) IRepo
	GetRolesMemberCanInvite(projectMember models.ProjectMember, roles *[]models.Role) IRepo
	GetUser(*models.User) IRepo
	GetUserByEmail(string, *[]models.User) IRepo
	IsMemberOfProject(*models.Project, *models.ProjectMember) IRepo
	ListProjectMembers(userIDList []string, projectMember *[]models.ProjectMember) IRepo
	MessageService() *message.MessageService
	ProjectAddMembers(models.Project, []models.MemberRole, models.User) IRepo
	ProjectGetAdmins(project *models.Project, members *[]models.ProjectMember) IRepo
	ProjectIsMemberAdmin(project *models.Project, member *models.ProjectMember) bool
	ProjectGetMembers(*models.Project, *[]models.ProjectMember) IRepo
	ProjectLoadUsers(*models.Project) IRepo
	ProjectRemoveMembers(models.Project, []string) IRepo
	ProjectSetRoleForUser(models.Project, models.User, models.Role) IRepo
	CheckMembersAreInProject(models.Project, []string) ([]string, error)
	RemoveOldMessageForRecipient(userID uint, environmentID string) IRepo
	SaveActivityLog(al *models.ActivityLog) IRepo
	SetLoginRequestCode(string, string) models.LoginRequest
	SetNewVersionID(environment *models.Environment) error
	WriteMessage(user models.User, message models.Message) IRepo
	GetDevices(uint, *[]models.Device) IRepo
	GetDevice(device *models.Device) IRepo
	GetDeviceByUserID(userID uint, device *models.Device) IRepo
	UpdateDeviceLastUsedAt(deviceUID string) IRepo
	RevokeDevice(userID uint, deviceUID string) IRepo
	GetAdminsFromUserProjects(userID uint, adminProjectsMap *map[string][]string) IRepo
	CreateOrganization(orga *models.Organization) IRepo
	UpdateOrganization(orga *models.Organization) IRepo
	OrganizationSetCustomer(organization *models.Organization, customer string) IRepo
	OrganizationSetSubscription(organization *models.Organization, subscription string) IRepo
	GetOrganization(orga *models.Organization) IRepo
	GetOrganizations(userID uint, result *models.GetOrganizationsResponse) IRepo
	GetOrganizationByName(orga *models.Organization) IRepo
	GetOrganizationProjects(*models.Organization, *[]models.Project) IRepo
	GetOrganizationMembers(orgaID uint, result *[]models.ProjectMember) IRepo
	IsUserOwnerOfOrga(*models.User, *models.Organization) (bool, error)
	IsProjectOrganizationPaid(string) (bool, error)
	CreateCheckoutSession(*models.CheckoutSession) IRepo
	GetCheckoutSession(string, *models.CheckoutSession) IRepo
	UpdateCheckoutSession(*models.CheckoutSession) IRepo
	DeleteCheckoutSession(*models.CheckoutSession) IRepo
	OrganizationSetPaid(organization *models.Organization, paid bool) IRepo
	GetUserProjects(userID uint, projects *[]models.Project) IRepo
}
