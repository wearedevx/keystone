package rights

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wearedevx/keystone/api/internal/emailer"
	"github.com/wearedevx/keystone/api/pkg/message"
	. "github.com/wearedevx/keystone/api/pkg/models"
	. "github.com/wearedevx/keystone/api/pkg/repo"
	"gorm.io/gorm"
)

var (
	fakeRoles    []Role
	fakeUserRole map[string]string
)

func initFakeRoles() {
	fakeRoles = []Role{
		{
			ID:           1,
			Name:         "developer",
			ParentID:     2,
			CanAddMember: false,
		},
		{
			ID:           2,
			Name:         "lead developer",
			ParentID:     3,
			CanAddMember: true,
		},
		{
			ID:           3,
			Name:         "devops",
			ParentID:     4,
			CanAddMember: true,
		},
		{
			ID:           4,
			Name:         "admin",
			CanAddMember: true,
		},
		{
			ID:           5,
			Name:         "nothing",
			CanAddMember: true,
		},
	}

	fakeRoles[1].Parent = &fakeRoles[3]
	fakeRoles[0].Parent = &fakeRoles[1]
}

func initFakeUserRoles() {
	fakeUserRole = map[string]string{
		"dev":    "developer",
		"lead":   "lead developer",
		"devops": "devops",
		"admin":  "admin",
	}
}

// Let’s setup some fixtures
func init() {
	initFakeRoles()
	initFakeUserRoles()
}

type FakeRepo struct{}

func (f *FakeRepo) CreateEnvironment(_ *Environment) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) CreateEnvironmentType(_ *EnvironmentType) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) CreateLoginRequest() LoginRequest {
	panic("not implemented")
}

func (f *FakeRepo) CreateProjectMember(_ *ProjectMember, _ *Role) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) CreateRole(_ *Role) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) CreateRoleEnvironmentType(_ *RolesEnvironmentType) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) CreateSecret(_ *Secret) {
	panic("not implemented")
}

func (f *FakeRepo) DeleteLoginRequest(_ string) bool {
	panic("not implemented")
}

func (f *FakeRepo) DeleteMessage(_ uint, _ uint) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) DeleteExpiredMessages() IRepo {
	panic("not implemented")
}

func (repo *FakeRepo) GetGroupedMessagesWillExpireByUser(
	groupedMessageUser *map[uint]emailer.GroupedMessagesUser,
) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) FindUsers(
	_ []string,
	_ *map[string]User,
	_ *[]string,
) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetDb() *gorm.DB {
	panic("not implemented")
}

func (f *FakeRepo) GetEnvironment(_ *Environment) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetEnvironmentsByProjectUUID(
	_ string,
	_ *[]Environment,
) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetEnvironmentPublicKeys(_ string, _ *PublicKeys) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetEnvironmentType(_ *EnvironmentType) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetLoginRequest(_ string) (LoginRequest, bool) {
	panic("not implemented")
}

func (f *FakeRepo) GetMessagesForUserOnEnvironment(
	_ Device,
	_ Environment,
	_ *Message,
) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOrCreateEnvironment(_ *Environment) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOrCreateEnvironmentType(_ *EnvironmentType) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOrCreateProject(_ *Project) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOrCreateProjectMember(_ *ProjectMember, _ string) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOrCreateRole(_ *Role) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOrCreateRoleEnvType(_ *RolesEnvironmentType) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOrCreateUser(_ *User) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetProject(_ *Project) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetProjectByUUID(_ string, _ *Project) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetRole(_ *Role) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetRoles(_ *[]Role) IRepo {
	panic("not implemented")
}

func (r *FakeRepo) GetRolesMemberCanInvite(
	projectMember ProjectMember,
	roles *[]Role,
) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetSecretByName(_ string, _ *Secret) {
	panic("not implemented")
}

func (f *FakeRepo) GetUser(_ *User) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) ListProjectMembers(
	userIDList []string,
	projectMember *[]ProjectMember,
) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) ProjectAddMembers(_ Project, _ []MemberRole, _ User) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) ProjectGetMembers(_ *Project, _ *[]ProjectMember) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) ProjectLoadUsers(_ *Project) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) ProjectRemoveMembers(_ Project, _ []string) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) ProjectSetRoleForUser(_ Project, _ User, _ Role) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) RemoveOldMessageForRecipient(_ uint, _ string) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) SetLoginRequestCode(_ string, _ string) LoginRequest {
	panic("not implemented")
}

func (f *FakeRepo) SetNewVersionID(_ *Environment) error {
	panic("not implemented")
}

func (f *FakeRepo) WriteMessage(_ User, _ Message) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) CheckMembersAreInProject(_ Project, _ []string) ([]string, error) {
	panic("not implemented")
}

func (f *FakeRepo) DeleteAllProjectMembers(project *Project) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) DeleteProject(project *Project) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) DeleteProjectsEnvironments(project *Project) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetActivityLogs(projectID string, options GetLogsOptions, logs *[]ActivityLog) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetMessage(message *Message) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetProjectsOrganization(_ string, _ *Organization) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) OrganizationCountMembers(_ *Organization, _ *int64) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetUserByEmail(_ string, _ *[]User) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) IsMemberOfProject(_ *Project, _ *ProjectMember) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) MessageService() *message.MessageService {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) ProjectGetAdmins(project *Project, members *[]ProjectMember) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) ProjectIsMemberAdmin(project *Project, member *ProjectMember) bool {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) SaveActivityLog(al *ActivityLog) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetDevices(_ uint, _ *[]Device) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetDevice(device *Device) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetDeviceByUserID(userID uint, device *Device) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) UpdateDeviceLastUsedAt(deviceUID string) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) RevokeDevice(userID uint, deviceUID string) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetAdminsFromUserProjects(userID uint, adminProjectsMap *map[string][]string) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) CreateOrganization(orga *Organization) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) UpdateOrganization(orga *Organization) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) OrganizationSetCustomer(organization *Organization, customer string) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) OrganizationSetSubscription(organization *Organization, subscription string) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetOrganization(orga *Organization) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetOrganizations(userID uint, result *[]Organization) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetOwnedOrganizationByName(userID uint, name string, orgas *[]Organization) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetOrganizationByName(userID uint, name string, orga *[]Organization) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetOrganizationProjects(_ *Organization, _ *[]Project) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetOrganizationMembers(orgaID uint, result *[]ProjectMember) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) IsUserOwnerOfOrga(_ *User, _ *Organization) (bool, error) {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) IsProjectOrganizationPaid(_ string) (bool, error) {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) CreateCheckoutSession(_ *CheckoutSession) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetCheckoutSession(_ string, _ *CheckoutSession) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) UpdateCheckoutSession(_ *CheckoutSession) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) DeleteCheckoutSession(_ *CheckoutSession) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) OrganizationSetPaid(organization *Organization, paid bool) IRepo {
	panic("not implemented") // TODO: Implement
}

func (f *FakeRepo) GetUserProjects(userID uint, projects *[]Project) IRepo {
	panic("not implemented") // TODO: Implement
}

func getRoleByEnvironmentTypeAndRole(
	environmentType *EnvironmentType,
	role *Role,
) RolesEnvironmentType {
	switch {
	case environmentType.Name == "dev" && (role.Name == "developer" || role.Name == "lead developer"):
		return RolesEnvironmentType{
			Read:  true,
			Write: true,
		}
	case environmentType.Name == "staging" && (role.Name == "developer" || role.Name == "lead developer"):
		return RolesEnvironmentType{
			Read:  false,
			Write: false,
		}
	case environmentType.Name == "prod" && (role.Name == "developer" || role.Name == "lead developer"):
		return RolesEnvironmentType{
			Read:  false,
			Write: false,
		}

	case environmentType.Name == "dev" && role.Name == "devops":
		return RolesEnvironmentType{
			Read:  true,
			Write: true,
		}
	case environmentType.Name == "staging" && role.Name == "devops":
		return RolesEnvironmentType{
			Read:  true,
			Write: true,
		}
	case environmentType.Name == "prod" && role.Name == "devops":
		return RolesEnvironmentType{
			Read:  false,
			Write: false,
		}

	case environmentType.Name == "dev" && role.Name == "admin":
		return RolesEnvironmentType{
			Read:  true,
			Write: true,
		}
	case environmentType.Name == "staging" && role.Name == "admin":
		return RolesEnvironmentType{
			Read:  true,
			Write: true,
		}
	case environmentType.Name == "prod" && role.Name == "admin":
		return RolesEnvironmentType{
			Read:  true,
			Write: true,
		}

	default:
		return RolesEnvironmentType{
			Read:  false,
			Write: false,
		}
	}
}

func findRole(role *Role) {
	for _, r := range fakeRoles {
		if r.Name == role.Name {
			*role = r
			return
		}
	}

	// role not found ?
	// role with the "nothing" name
	*role = fakeRoles[4]
}

func getRoleByUsername(userName string) (role Role) {
	roleName, ok := fakeUserRole[userName]
	if !ok {
		roleName = "nothing"
	}

	role = Role{Name: roleName}
	findRole(&role)

	return role
}

func (fakeRepo *FakeRepo) Err() error {
	return nil
}

func (fakeRepo *FakeRepo) GetRolesEnvironmentType(
	rolesEnvironmentType *RolesEnvironmentType,
) IRepo {
	*rolesEnvironmentType = getRoleByEnvironmentTypeAndRole(
		&rolesEnvironmentType.EnvironmentType,
		&rolesEnvironmentType.Role,
	)

	return fakeRepo
}

func (fakeRepo *FakeRepo) GetProjectMember(projectMember *ProjectMember) IRepo {
	role := getRoleByUsername(projectMember.User.Username)
	*projectMember = ProjectMember{
		Role: role,
	}

	return fakeRepo
}

func (fakeRepo *FakeRepo) GetInvitableRoles(role Role, roles *[]Role) IRepo {
	// return fakeRepo
	// TODO
	return fakeRepo
}

func (fakeRepo *FakeRepo) GetChildrenRoles(role Role, roles *[]Role) IRepo {
	switch role.ID {
	case 3:
		*roles = []Role{fakeRoles[0], fakeRoles[1]}

	case 2:
		*roles = []Role{fakeRoles[0]}

	default:
		*roles = []Role{}
	}

	return fakeRepo
}

type rw struct {
	r bool
	w bool
}

func TestCanUserHasRightEnvironment(t *testing.T) {
	fakeRepo := &FakeRepo{}
	project := &Project{}

	userDev := &User{Username: "dev"}
	userLeadDev := &User{Username: "lead"}
	userDevops := &User{Username: "devops"}
	userAdmin := &User{Username: "admin"}

	environmentDev := &Environment{Name: "dev"}
	environmentStaging := &Environment{Name: "staging"}
	environmentProd := &Environment{Name: "prod"}

	users := map[string]*User{
		"dev":    userDev,
		"lead":   userLeadDev,
		"devops": userDevops,
		"admin":  userAdmin,
	}

	environments := map[string]*Environment{
		"dev":     environmentDev,
		"staging": environmentStaging,
		"prod":    environmentProd,
	}

	rightsMatrix := map[string]map[string]rw{
		"dev": {
			"dev":     {true, true},
			"staging": {false, false},
			"prod":    {false, false},
		},
		"lead": {
			"dev":     {true, true},
			"staging": {false, false},
			"prod":    {false, false},
		},
		"devops": {
			"dev":     {true, true},
			"staging": {true, true},
			"prod":    {false, false},
		},
		"admin": {
			"dev":     {true, true},
			"staging": {true, true},
			"prod":    {true, true},
		},
	}

	for name, user := range users {
		for envName, environment := range environments {
			expectation := rightsMatrix[name][envName]

			canRead, _ := CanUserReadEnvironment(
				fakeRepo,
				user.ID,
				project.ID,
				environment,
			)
			canWrite, _ := CanUserWriteOnEnvironment(
				fakeRepo,
				user.ID,
				project.ID,
				environment,
			)

			assert.Equal(
				t,
				expectation.r,
				canRead,
				"Oops! User %s has unexpected read rights on %s environment",
				name,
				envName,
			)
			assert.Equal(
				t,
				expectation.w,
				canWrite,
				"Oops! User %s has unexpected write rights on %s environment",
				name,
				envName,
			)
		}
	}
}

func TestUserCanSetMemberRole(t *testing.T) {
	fakeRepo := FakeRepo{}
	project := Project{}

	userDev := User{Username: "dev"}
	userLeadDev := User{Username: "lead"}
	userDevops := User{Username: "devops"}
	userAdmin := User{Username: "admin"}

	users := map[string]User{
		"dev":    userDev,
		"lead":   userLeadDev,
		"devops": userDevops,
		"admin":  userAdmin,
	}

	rightsMatrix := map[string]map[string]bool{
		"dev": {
			"dev":    false,
			"lead":   false,
			"devops": false,
			"admin":  false,
		},
		"lead": {
			"dev":    true,
			"lead":   true,
			"devops": false,
			"admin":  false,
		},
		"devops": {
			"dev":    true,
			"lead":   true,
			"devops": true,
			"admin":  false,
		},
		"admin": {
			"dev":    true,
			"lead":   true,
			"devops": true,
			"admin":  true,
		},
	}

	for name, user := range users {
		for otherName, otherUser := range users {
			expectation := rightsMatrix[name][otherName]
			role := getRoleByUsername(otherName)

			canSetMemberRole, _ := CanUserSetMemberRole(
				&fakeRepo,
				user,
				otherUser,
				role,
				project,
			)

			assert.Equal(
				t,
				expectation,
				canSetMemberRole,
				"Oops! User %s has unexpected role setting rights on user %s",
				name,
				otherName,
			)
		}
	}
}
