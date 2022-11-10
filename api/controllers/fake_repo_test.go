package controllers

import (
	"github.com/wearedevx/keystone/api/internal/emailer"
	"github.com/wearedevx/keystone/api/internal/redis"
	"github.com/wearedevx/keystone/api/pkg/message"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
	"gorm.io/gorm"
)

type crashers map[string]error

type funcCall struct {
	funcName string
	args     []interface{}
	result   []interface{}
}

type fakeRepo struct {
	repo.Repo
	err      error
	called   []string
	crashers crashers
	messages message.MessageService
}

var noCrashers = map[string]error{}

func newFakeRepo(crashers map[string]error) *fakeRepo {
	return &fakeRepo{
		Repo:     *repo.NewRepo(),
		crashers: crashers,
		called:   []string{},
		messages: newFakeMessageService(crashers),
	}
}

func (f *fakeRepo) Err() error {
	f.called = append(f.called, "Err")
	if f.err != nil {
		return f.err
	}

	return f.Repo.Err()
}

func (f *fakeRepo) ClearErr() repo.IRepo {
	f.called = append(f.called, "ClearErr")
	f.err = nil
	f.Repo.ClearErr()

	return f
}

func (f *fakeRepo) MessageService() message.MessageService {
	f.called = append(f.called, "MessageService")
	return f.messages
}

func (f *fakeRepo) WriteMessage(u models.User, m models.Message) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "WriteMessage")
	if e, ok := f.crashers["WriteMessage"]; ok {
		f.err = e
		return f
	}
	f.Repo.WriteMessage(u, m)
	return f
}

func (f *fakeRepo) RemoveOldMessageForRecipent(deviceID uint, environmentID string) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "RemoveOldMessageForRecipient")
	if e, ok := f.crashers["RemoveOldMessageForRecipent"]; ok {
		f.err = e
		return f
	}
	f.Repo.RemoveOldMessageForRecipient(deviceID, environmentID)

	return f
}

func (f *fakeRepo) DeleteMessage(id uint, recipientID uint) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "DeleteMessage")
	if e, ok := f.crashers["DeleteMessage"]; ok {
		f.err = e
		return f
	}
	f.Repo.DeleteMessage(id, recipientID)
	return f
}

func (f *fakeRepo) GetProject(project *models.Project) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetProject")
	if e, ok := f.crashers["GetProject"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetProject(project)
	return f
}

func (f *fakeRepo) GetDeviceByUserID(id uint, device *models.Device) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetDeviceByUserId")
	if e, ok := f.crashers["GetDeviceByUserID"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetDeviceByUserID(id, device)
	return f
}

func (f *fakeRepo) GetDevice(device *models.Device) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetDevice")
	if e, ok := f.crashers["GetDevice"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetDevice(device)
	return f
}

func (f *fakeRepo) GetEnvironmentsByProjectUUID(id string, environments *[]models.Environment) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetEnvironmentsByProjectUUID")
	if e, ok := f.crashers["GetEnvironmentsByProjectUUID"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetEnvironmentsByProjectUUID(id, environments)
	return f
}

func (f *fakeRepo) GetProjectMember(
	projectMember *models.ProjectMember,
) repo.IRepo {
	if f.err != nil {
		return f
	}
	if e, ok := f.crashers["GetProjectMember"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetProjectMember(projectMember)
	return f
}

func (f *fakeRepo) GetRolesEnvironmentType(environmentType *models.RolesEnvironmentType) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetRolesEnvironmentType")
	if e, ok := f.crashers["GetRolesEnvironmentType"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetRolesEnvironmentType(environmentType)
	return f
}

func (f *fakeRepo) GetMessagesForUserOnEnvironment(device models.Device, environment models.Environment, message *models.Message) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetMessagesForUserOnEnvironment")
	if e, ok := f.crashers["GetMessagesForUserOnEnvironment"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetMessagesForUserOnEnvironment(device, environment, message)
	return f
}

func (f *fakeRepo) GetOrganizationMembers(orgID uint, members *[]models.ProjectMember) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetOrgnaizationMembers")
	if e, ok := f.crashers["GetOrganizationMembers"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetOrganizationMembers(orgID, members)
	return f
}

func (f *fakeRepo) SetNewVersionID(environment *models.Environment) error {
	if f.err != nil {
		return f.err
	}
	f.called = append(f.called, "SetNewVersionID")
	if e, ok := f.crashers["SetNewVersionID"]; ok {
		return e
	}
	return f.Repo.SetNewVersionID(environment)
}

func (f *fakeRepo) CreateEnvironment(environment *models.Environment) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "CreateEnvironment")
	if e, ok := f.crashers["CreateEnvironment"]; ok {
		f.err = e
		return f
	}
	f.Repo.CreateEnvironment(environment)
	return f
}

func (f *fakeRepo) CreateEnvironmentType(environmentType *models.EnvironmentType) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "CreateEnvironmentType")
	if e, ok := f.crashers["CreateEnvironmentType"]; ok {
		f.err = e
		return f
	}
	f.Repo.CreateEnvironmentType(environmentType)
	return f
}

func (f *fakeRepo) CreateLoginRequest() models.LoginRequest {
	if f.err != nil {
		return models.LoginRequest{}
	}
	f.called = append(f.called, "CreateLoginRequest")
	if e, ok := f.crashers["CreateLoginRequest"]; ok {
		f.err = e
		return models.LoginRequest{}
	}
	return f.Repo.CreateLoginRequest()
}

func (f *fakeRepo) CreateProjectMember(projectMember *models.ProjectMember, role *models.Role) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "CreateProjectMember")
	if e, ok := f.crashers["CreateProjectMember"]; ok {
		f.err = e
		return f
	}
	f.Repo.CreateProjectMember(projectMember, role)
	return f
}

func (f *fakeRepo) CreateRole(role *models.Role) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "CreateRole")
	if e, ok := f.crashers["CreateRole"]; ok {
		f.err = e
		return f
	}
	f.Repo.CreateRole(role)
	return f
}

func (f *fakeRepo) CreateRoleEnvironmentType(rolesEnvironmentType *models.RolesEnvironmentType) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "CreateRoleEnvironmentType")
	if e, ok := f.crashers["CreateRoleEnvironmentType"]; ok {
		f.err = e
		return f
	}
	f.Repo.CreateRoleEnvironmentType(rolesEnvironmentType)
	return f
}

func (f *fakeRepo) DeleteLoginRequest(str string) bool {
	if f.err != nil {
		return false
	}
	f.called = append(f.called, "DeleteLoginRequest")
	if e, ok := f.crashers["DeleteLoginRequest"]; ok {
		f.err = e
		return false
	}
	return f.Repo.DeleteLoginRequest(str)
}

func (f *fakeRepo) DeleteAllProjectMembers(project *models.Project) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "DeleteAllProjectMembers")
	if e, ok := f.crashers["DeleteAllProjectMembers"]; ok {
		f.err = e
		return f
	}
	f.Repo.DeleteAllProjectMembers(project)
	return f
}

func (f *fakeRepo) DeleteExpiredMessages() repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "DeleteExpiredMessages")
	if e, ok := f.crashers["DeleteExpiredMessages"]; ok {
		f.err = e
		return f
	}
	f.Repo.DeleteExpiredMessages()
	return f
}

func (f *fakeRepo) GetGroupedMessagesWillExpireByUser(groupedMessageUser *map[uint]emailer.GroupedMessagesUser) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetGroupedMessagesWillExpireByUser")
	if e, ok := f.crashers["GetGroupedMessagesWillExpireByUser"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetGroupedMessagesWillExpireByUser(groupedMessageUser)
	return f
}

func (f *fakeRepo) DeleteProject(project *models.Project) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "DeleteProject")
	if e, ok := f.crashers["DeleteProject"]; ok {
		f.err = e
		return f
	}
	f.Repo.DeleteProject(project)
	return f
}

func (f *fakeRepo) DeleteProjectsEnvironments(project *models.Project) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "DeleteProjectsEnvironments")
	if e, ok := f.crashers["DeleteProjectsEnvironments"]; ok {
		f.err = e
		return f
	}
	f.Repo.DeleteProjectsEnvironments(project)
	return f
}

func (f *fakeRepo) FindUsers(userIDs []string, users *map[string]models.User, notFounds *[]string) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "FindUsers")
	if e, ok := f.crashers["FindUsers"]; ok {
		f.err = e
		return f
	}
	f.Repo.FindUsers(userIDs, users, notFounds)
	return f
}

func (f *fakeRepo) GetActivityLogs(projectID string, options models.GetLogsOptions, logs *[]models.ActivityLog) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetActivityLogs")
	if e, ok := f.crashers["GetActivityLogs"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetActivityLogs(projectID, options, logs)
	return f
}

func (f *fakeRepo) GetChildrenRoles(role models.Role, roles *[]models.Role) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetChildrenRoles")
	if e, ok := f.crashers["GetChildrenRoles"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetChildrenRoles(role, roles)
	return f
}

func (f *fakeRepo) GetDB() *gorm.DB {
	f.called = append(f.called, "GetDB")
	if e, ok := f.crashers["GetDB"]; ok {
		f.err = e
		return nil
	}
	return f.Repo.GetDB()
}

func (f *fakeRepo) GetEnvironment(environment *models.Environment) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetEnvironment")
	if e, ok := f.crashers["GetEnvironment"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetEnvironment(environment)
	return f
}

func (f *fakeRepo) GetEnvironmentPublicKeys(envID string, publicKeys *models.PublicKeys) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetEnvironmentPublicKeys")
	if e, ok := f.crashers["GetEnvironmentPublicKeys"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetEnvironmentPublicKeys(envID, publicKeys)
	return f
}

func (f *fakeRepo) GetEnvironmentType(environmentType *models.EnvironmentType) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetEnvironmentType")
	if e, ok := f.crashers["GetEnvironmentType"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetEnvironmentType(environmentType)
	return f
}

func (f *fakeRepo) GetInvitableRoles(role models.Role, roles *[]models.Role) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetInvitableRoles")
	if e, ok := f.crashers["GetInvitableRoles"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetInvitableRoles(role, roles)
	return f
}

func (f *fakeRepo) GetLoginRequest(loginRequest string) (models.LoginRequest, bool) {
	if f.err != nil {
		return models.LoginRequest{}, false
	}
	f.called = append(f.called, "GetLoginRequest")
	if e, ok := f.crashers["GetLoginRequest"]; ok {
		f.err = e
		return models.LoginRequest{}, false
	}
	return f.Repo.GetLoginRequest(loginRequest)
}

func (f *fakeRepo) GetMessage(message *models.Message) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetMessage")
	if e, ok := f.crashers["GetMessage"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetMessage(message)
	return f
}

func (f *fakeRepo) GetOrCreateEnvironment(environment *models.Environment) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetOrCreateEnvironment")
	if e, ok := f.crashers["GetOrCreateEnvironment"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetOrCreateEnvironment(environment)
	return f
}

func (f *fakeRepo) GetOrCreateEnvironmentType(environmentType *models.EnvironmentType) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetOrCreateEnvironmentType")
	if e, ok := f.crashers["GetOrCreateEnvironmentType"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetOrCreateEnvironmentType(environmentType)
	return f
}

func (f *fakeRepo) GetOrCreateProject(project *models.Project) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetOrCreateProject")
	if e, ok := f.crashers["GetOrCreateProject"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetOrCreateProject(project)
	return f
}

func (f *fakeRepo) GetOrCreateProjectMember(projectMember *models.ProjectMember, iRepo string) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetOrCreateProjectMember")
	if e, ok := f.crashers["GetOrCreateProjectMember"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetOrCreateProjectMember(projectMember, iRepo)
	return f
}

func (f *fakeRepo) GetOrCreateRole(role *models.Role) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetOrCreateRole")
	if e, ok := f.crashers["GetOrCreateRole"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetOrCreateRole(role)
	return f
}

func (f *fakeRepo) GetOrCreateRoleEnvType(rolesEnvironmentType *models.RolesEnvironmentType) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetOrCreateRoleEnvType")
	if e, ok := f.crashers["GetOrCreateRoleEnvType"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetOrCreateRoleEnvType(rolesEnvironmentType)
	return f
}

func (f *fakeRepo) GetOrCreateUser(user *models.User) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetOrCreateUser")
	if e, ok := f.crashers["GetOrCreateUser"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetOrCreateUser(user)
	return f
}

func (f *fakeRepo) GetProjectByUUID(projectID string, project *models.Project) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetProjectByUUID")
	if e, ok := f.crashers["GetProjectByUUID"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetProjectByUUID(projectID, project)
	return f
}

func (f *fakeRepo) GetProjectsOrganization(organizationID string, organization *models.Organization) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetProjectsOrganization")
	if e, ok := f.crashers["GetProjectsOrganization"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetProjectsOrganization(organizationID, organization)
	return f
}

func (f *fakeRepo) OrganizationCountMembers(organization *models.Organization, count *int64) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "OrganizationCountMembers")
	if e, ok := f.crashers["OrganizationCountMembers"]; ok {
		f.err = e
		return f
	}
	f.Repo.OrganizationCountMembers(organization, count)
	return f
}

func (f *fakeRepo) GetRole(role *models.Role) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetRole")
	if e, ok := f.crashers["GetRole"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetRole(role)
	return f
}

func (f *fakeRepo) GetRoles(role *[]models.Role) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetRoles")
	if e, ok := f.crashers["GetRoles"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetRoles(role)
	return f
}

func (f *fakeRepo) GetRolesMemberCanInvite(projectMember models.ProjectMember, roles *[]models.Role) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetRolesMemberCanInvite")
	if e, ok := f.crashers["GetRolesMemberCanInvite"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetRolesMemberCanInvite(projectMember, roles)
	return f
}

func (f *fakeRepo) GetUser(user *models.User) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetUser")
	if e, ok := f.crashers["GetUser"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetUser(user)
	return f
}

func (f *fakeRepo) GetUserByEmail(userID string, users *[]models.User) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetUserByEmail")
	if e, ok := f.crashers["GetUserByEmail"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetUserByEmail(userID, users)
	return f
}

func (f *fakeRepo) IsMemberOfProject(project *models.Project, projectMember *models.ProjectMember) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "IsMemberOfProject")
	if e, ok := f.crashers["IsMemberOfProject"]; ok {
		f.err = e
		return f
	}
	f.Repo.IsMemberOfProject(project, projectMember)
	return f
}

func (f *fakeRepo) ListProjectMembers(userIDList []string, projectMember *[]models.ProjectMember) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "ListProjectMembers")
	if e, ok := f.crashers["ListProjectMembers"]; ok {
		f.err = e
		return f
	}
	f.Repo.ListProjectMembers(userIDList, projectMember)
	return f
}

func (f *fakeRepo) ProjectAddMembers(project models.Project, memberRole []models.MemberRole, user models.User) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "ProjectAddMembers")
	if e, ok := f.crashers["ProjectAddMembers"]; ok {
		f.err = e
		return f
	}
	f.Repo.ProjectAddMembers(project, memberRole, user)
	return f
}

func (f *fakeRepo) UsersInMemberRoles(mers []models.MemberRole) (map[string]models.User, []string) {
	if f.err != nil {
		return nil, nil
	}
	f.called = append(f.called, "UsersInMemberRoles")
	if e, ok := f.crashers["UsersInMemberRoles"]; ok {
		f.err = e
		return map[string]models.User{}, []string{}
	}
	return f.Repo.UsersInMemberRoles(mers)
}

func (f *fakeRepo) SetNewlyCreatedDevice(flag bool, deviceID uint, userID uint) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "SetNewlyCreatedDevice")
	if e, ok := f.crashers["SetNewlyCreatedDevice"]; ok {
		f.err = e
		return f
	}
	f.Repo.SetNewlyCreatedDevice(flag, deviceID, userID)
	return f
}

func (f *fakeRepo) ProjectGetAdmins(project *models.Project, members *[]models.ProjectMember) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "ProjectGetAdmins")
	if e, ok := f.crashers["ProjectGetAdmins"]; ok {
		f.err = e
		return f
	}
	f.Repo.ProjectGetAdmins(project, members)
	return f
}

func (f *fakeRepo) ProjectIsMemberAdmin(project *models.Project, member *models.ProjectMember) bool {
	if f.err != nil {
		return false
	}
	f.called = append(f.called, "ProjectIsMemberAdmin")
	if e, ok := f.crashers["ProjectIsMemberAdmin"]; ok {
		f.err = e
		return false
	}
	return f.Repo.ProjectIsMemberAdmin(project, member)
}

func (f *fakeRepo) ProjectGetMembers(project *models.Project, projectMember *[]models.ProjectMember) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "ProjectGetMembers")
	if e, ok := f.crashers["ProjectGetMembers"]; ok {
		f.err = e
		return f
	}
	f.Repo.ProjectGetMembers(project, projectMember)
	return f
}

func (f *fakeRepo) ProjectLoadUsers(project *models.Project) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "ProjectLoadUsers")
	if e, ok := f.crashers["ProjectLoadUsers"]; ok {
		f.err = e
		return f
	}
	f.Repo.ProjectLoadUsers(project)
	return f
}

func (f *fakeRepo) ProjectRemoveMembers(project models.Project, memberIDs []string) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "ProjectRemoveMembers")
	if e, ok := f.crashers["ProjectRemoveMembers"]; ok {
		f.err = e
		return f
	}
	f.Repo.ProjectRemoveMembers(project, memberIDs)
	return f
}

func (f *fakeRepo) ProjectSetRoleForUser(project models.Project, user models.User, role models.Role) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "ProjectSetRoleForUser")
	if e, ok := f.crashers["ProjectSetRoleForUser"]; ok {
		f.err = e
		return f
	}
	f.Repo.ProjectSetRoleForUser(project, user, role)
	return f
}

func (f *fakeRepo) CheckMembersAreInProject(project models.Project, str []string) ([]string, error) {
	if f.err != nil {
		return nil, f.err
	}
	f.called = append(f.called, "CheckMembersAreInProject")
	if e, ok := f.crashers["CheckMembersAreInProject"]; ok {
		f.err = e
		return []string{}, e
	}
	return f.Repo.CheckMembersAreInProject(project, str)
}

func (f *fakeRepo) RemoveOldMessageForRecipient(userID uint, environmentID string) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "RemoveOldMessageForRecipient")
	if e, ok := f.crashers["RemoveOldMessageForRecipient"]; ok {
		f.err = e
		return f
	}
	f.Repo.RemoveOldMessageForRecipient(userID, environmentID)
	return f
}

func (f *fakeRepo) SaveActivityLog(al *models.ActivityLog) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "SaveActivityLog")
	if e, ok := f.crashers["SaveActivityLog"]; ok {
		f.err = e
		return f
	}
	f.Repo.SaveActivityLog(al)
	return f
}

func (f *fakeRepo) SetLoginRequestCode(code string, otherCode string) models.LoginRequest {
	if f.err != nil {
		return models.LoginRequest{}
	}
	f.called = append(f.called, "SetLoginRequestCode")
	if e, ok := f.crashers["SetLoginRequestCode"]; ok {
		f.err = e
		return models.LoginRequest{}
	}
	return f.Repo.SetLoginRequestCode(code, otherCode)
}

func (f *fakeRepo) GetDevices(userID uint, devices *[]models.Device) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetDevices")
	if e, ok := f.crashers["GetDevices"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetDevices(userID, devices)
	return f
}

func (f *fakeRepo) GetNewlyCreatedDevices(device *[]models.Device) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetNewlyCreatedDevices")
	if e, ok := f.crashers["GetNewlyCreatedDevices"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetNewlyCreatedDevices(device)
	return f
}

func (f *fakeRepo) UpdateDeviceLastUsedAt(deviceUID string) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "UpdateDeviceLastUsedAt")
	if e, ok := f.crashers["UpdateDeviceLastUsedAt"]; ok {
		f.err = e
		return f
	}
	f.Repo.UpdateDeviceLastUsedAt(deviceUID)
	return f
}

func (f *fakeRepo) RevokeDevice(userID uint, deviceUID string) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "RevokeDevice")
	if e, ok := f.crashers["RevokeDevice"]; ok {
		f.err = e
		return f
	}
	f.Repo.RevokeDevice(userID, deviceUID)
	return f
}

func (f *fakeRepo) GetAdminsFromUserProjects(userID uint, adminProjectsMap *map[string][]string) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetAdminsFromUserProjects")
	if e, ok := f.crashers["GetAdminsFromUserProjects"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetAdminsFromUserProjects(userID, adminProjectsMap)
	return f
}

func (f *fakeRepo) CreateOrganization(orga *models.Organization) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "CreateOrganization")
	if e, ok := f.crashers["CreateOrganization"]; ok {
		f.err = e
		return f
	}
	f.Repo.CreateOrganization(orga)
	return f
}

func (f *fakeRepo) UpdateOrganization(orga *models.Organization) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "UpdateOrganization")
	if e, ok := f.crashers["UpdateOrganization"]; ok {
		f.err = e
		return f
	}
	f.Repo.UpdateOrganization(orga)
	return f
}

func (f *fakeRepo) OrganizationSetCustomer(organization *models.Organization, customer string) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "OrganizationSetCustomer")
	if e, ok := f.crashers["OrganizationSetCustomer"]; ok {
		f.err = e
		return f
	}
	f.Repo.OrganizationSetCustomer(organization, customer)
	return f
}

func (f *fakeRepo) OrganizationSetSubscription(organization *models.Organization, subscription string) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "OrganizationSetSubscription")
	if e, ok := f.crashers["OrganizationSetSubscription"]; ok {
		f.err = e
		return f
	}
	f.Repo.OrganizationSetSubscription(organization, subscription)
	return f
}

func (f *fakeRepo) GetOrganization(orga *models.Organization) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetOrganization")
	if e, ok := f.crashers["GetOrganization"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetOrganization(orga)
	return f
}

func (f *fakeRepo) GetOrganizations(userID uint, result *[]models.Organization) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetOrganizations")
	if e, ok := f.crashers["GetOrganizations"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetOrganizations(userID, result)
	return f
}

func (f *fakeRepo) GetOwnedOrganizations(userID uint, result *[]models.Organization) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetOwnedOrganizations")
	if e, ok := f.crashers["GetOwnedOrganizations"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetOwnedOrganizations(userID, result)
	return f
}

func (f *fakeRepo) GetOwnedOrganizationByName(userID uint, name string, orgas *[]models.Organization) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetOwnedOrganizationByName")
	if e, ok := f.crashers["GetOwnedOrganizationByName"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetOwnedOrganizationByName(userID, name, orgas)
	return f
}

func (f *fakeRepo) GetOrganizationByName(userID uint, name string, orga *[]models.Organization) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetOrganizationByName")
	if e, ok := f.crashers["GetOrganizationByName"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetOrganizationByName(userID, name, orga)
	return f
}

func (f *fakeRepo) GetOrganizationProjects(organization *models.Organization, project *[]models.Project) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetOrganizationProjects")
	if e, ok := f.crashers["GetOrganizationProjects"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetOrganizationProjects(organization, project)
	return f
}

func (f *fakeRepo) IsUserOwnerOfOrga(user *models.User, organization *models.Organization) (bool, error) {
	if f.err != nil {
		return false, f.err
	}
	f.called = append(f.called, "IsUserOwnerOfOrga")
	if e, ok := f.crashers["IsUserOwnerOfOrga"]; ok {
		f.err = e
		return false, e
	}
	return f.Repo.IsUserOwnerOfOrga(user, organization)
}

func (f *fakeRepo) IsProjectOrganizationPaid(str string) (bool, error) {
	if f.err != nil {
		return false, f.err
	}
	f.called = append(f.called, "IsProjectOrganizationPaid")
	if e, ok := f.crashers["IsProjectOrganizationPaid"]; ok {
		f.err = e
		return false, e
	}
	return f.Repo.IsProjectOrganizationPaid(str)
}

func (f *fakeRepo) CreateCheckoutSession(checkoutSession *models.CheckoutSession) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "CreateCheckoutSession")
	if e, ok := f.crashers["CreateCheckoutSession"]; ok {
		f.err = e
		return f
	}
	f.Repo.CreateCheckoutSession(checkoutSession)
	return f
}

func (f *fakeRepo) GetCheckoutSession(id string, checkoutSession *models.CheckoutSession) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetCheckoutSession")
	if e, ok := f.crashers["GetCheckoutSession"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetCheckoutSession(id, checkoutSession)
	return f
}

func (f *fakeRepo) UpdateCheckoutSession(checkoutSession *models.CheckoutSession) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "UpdateCheckoutSession")
	if e, ok := f.crashers["UpdateCheckoutSession"]; ok {
		f.err = e
		return f
	}
	f.Repo.UpdateCheckoutSession(checkoutSession)
	return f
}

func (f *fakeRepo) DeleteCheckoutSession(checkoutSession *models.CheckoutSession) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "DeleteCheckoutSession")
	if e, ok := f.crashers["DeleteCheckoutSession"]; ok {
		f.err = e
		return f
	}
	f.Repo.DeleteCheckoutSession(checkoutSession)
	return f
}

func (f *fakeRepo) OrganizationSetPaid(organization *models.Organization, paid bool) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "OrganizationSetPaid")
	if e, ok := f.crashers["OrganizationSetPaid"]; ok {
		f.err = e
		return f
	}
	f.Repo.OrganizationSetPaid(organization, paid)
	return f
}

func (f *fakeRepo) GetUserProjects(userID uint, projects *[]models.Project) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "GetUserProjects")
	if e, ok := f.crashers["GetUserProjects"]; ok {
		f.err = e
		return f
	}
	f.Repo.GetUserProjects(userID, projects)
	return f
}

func (f *fakeRepo) FindUserWithRefreshToken(refreshToken string, user *models.User) repo.IRepo {
	if f.err != nil {
		return f
	}
	f.called = append(f.called, "FindUserWithRefreshToken")
	if e, ok := f.crashers["FindUserWithRefreshToken"]; ok {
		f.err = e
		return f
	}
	f.Repo.FindUserWithRefreshToken(refreshToken, user)
	return f
}

// +-------------------------

type fakeMessageService struct {
	message.MessageService
	err      error
	crashers map[string]error
	redis    redis.IRedis
}

func newFakeMessageService(crashers map[string]error) message.MessageService {
	return &fakeMessageService{
		MessageService: message.NewMessageService(),
		err:            nil,
		crashers:       crashers,
		redis:          redis.NewRedis(),
	}
}

func (f *fakeMessageService) GetMessageByUUID(uuid string) ([]byte, error) {
	if e, ok := f.crashers["MessageService.GetMessageByUUID"]; ok {
		return nil, e
	}

	return f.MessageService.GetMessageByUUID(uuid)
}

func (f *fakeMessageService) WriteMessageWithUUID(uuid string, message []byte) error {
	if e, ok := f.crashers["MessageService.WriteMessageWithUUID"]; ok {
		return e
	}

	return f.MessageService.WriteMessageWithUUID(uuid, message)
}
