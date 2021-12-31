package rights

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wearedevx/keystone/api/internal/emailer"
	"github.com/wearedevx/keystone/api/pkg/message"
	. "github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
	. "github.com/wearedevx/keystone/api/pkg/repo"
	"gorm.io/gorm"
)

var (
	fakeRoles                 []Role
	fakeUserRole              map[uint]string
	fakeEnvironmentTypes      []EnvironmentType
	fakeRolesEnvironmentTypes []RolesEnvironmentType
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
	fakeUserRole = map[uint]string{
		1: "developer",
		2: "lead developer",
		3: "devops",
		4: "admin",
	}
}

func initFakeEnvironmentTypes() {
	fakeEnvironmentTypes = []EnvironmentType{
		{
			ID:   1,
			Name: "dev",
		},
		{
			ID:   2,
			Name: "staging",
		},
		{
			ID:   3,
			Name: "prod",
		},
	}
}

func initFakeRolesEnvironmentTypes() {
	fakeRolesEnvironmentTypes = []RolesEnvironmentType{}

	matrix := [][]struct {
		read  bool
		write bool
	}{
		{ // role dev
			{true, true},   // dev
			{false, false}, // staging
			{false, false}, // prod
		},
		{ // role lead dev
			{true, true},   // dev
			{false, false}, // staging
			{false, false}, // prod
		},
		{ // role devops
			{true, true}, // dev
			{true, true}, // staging
			{true, true}, // prod
		},
		{ // role admin
			{true, true}, // dev
			{true, true}, // staging
			{true, true}, // prod
		},
	}

	for n, line := range matrix {
		for m, r := range line {
			fakeRolesEnvironmentTypes = append(fakeRolesEnvironmentTypes,
				RolesEnvironmentType{
					ID:                uint(len(fakeRolesEnvironmentTypes) + 1),
					RoleID:            fakeRoles[n].ID,
					Role:              fakeRoles[n],
					EnvironmentTypeID: fakeEnvironmentTypes[m].ID,
					EnvironmentType:   fakeEnvironmentTypes[m],
					Name:              "",
					Read:              r.read,
					Write:             r.write,
				},
			)
		}
	}
}

// Let’s setup some fixtures
func init() {
	initFakeRoles()
	initFakeUserRoles()
	initFakeEnvironmentTypes()
	initFakeRolesEnvironmentTypes()
}

type FakeRepo struct {
	err error
}

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
func (f *FakeRepo) UsersInMemberRoles(mers []MemberRole) (map[string]User, []string) {
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

func (f *FakeRepo) CheckMembersAreInProject(
	_ Project,
	_ []string,
) ([]string, error) {
	panic("not implemented")
}

func (f *FakeRepo) DeleteAllProjectMembers(project *Project) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) DeleteProject(project *Project) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) DeleteProjectsEnvironments(project *Project) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetActivityLogs(
	projectID string,
	options GetLogsOptions,
	logs *[]ActivityLog,
) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetMessage(message *Message) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetProjectsOrganization(
	projectuuid string,
	orga *Organization,
) IRepo {
	*orga = Organization{
		ID:             0,
		Name:           "organization-namel",
		Paid:           false,
		Private:        false,
		CustomerID:     "",
		SubscriptionID: "",
		UserID:         0,
		User: User{
			ID:            0,
			AccountType:   "",
			UserID:        "",
			ExtID:         "",
			Username:      "",
			Fullname:      "",
			Email:         "",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Devices:       []Device{},
			Organizations: []Organization{},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return f
}

func (r *FakeRepo) SetNewlyCreatedDevice(flag bool, deviceID uint, userID uint) repo.IRepo {
	panic("not implemented")
}
func (f *FakeRepo) OrganizationCountMembers(_ *Organization, _ *int64) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetUserByEmail(_ string, _ *[]User) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) IsMemberOfProject(_ *Project, _ *ProjectMember) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) MessageService() *message.MessageService {
	panic("not implemented")
}

func (f *FakeRepo) ProjectGetAdmins(
	project *Project,
	members *[]ProjectMember,
) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) ProjectIsMemberAdmin(
	project *Project,
	member *ProjectMember,
) bool {
	panic("not implemented")
}

func (f *FakeRepo) SaveActivityLog(al *ActivityLog) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetDevices(_ uint, _ *[]Device) IRepo {
	panic("not implemented")
}
func (f *FakeRepo) GetNewlyCreatedDevices(_ *[]Device) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetDevice(device *Device) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetDeviceByUserID(userID uint, device *Device) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) UpdateDeviceLastUsedAt(deviceUID string) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) RevokeDevice(userID uint, deviceUID string) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetAdminsFromUserProjects(
	userID uint,
	adminProjectsMap *map[string][]string,
) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) CreateOrganization(orga *Organization) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) UpdateOrganization(orga *Organization) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) OrganizationSetCustomer(
	organization *Organization,
	customer string,
) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) OrganizationSetSubscription(
	organization *Organization,
	subscription string,
) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOrganization(orga *Organization) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOrganizations(userID uint, result *[]Organization) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOwnedOrganizations(userID uint, result *[]Organization) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOwnedOrganizationByName(
	userID uint,
	name string,
	orgas *[]Organization,
) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOrganizationByName(
	userID uint,
	name string,
	orga *[]Organization,
) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOrganizationProjects(
	_ *Organization,
	_ *[]Project,
) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOrganizationMembers(
	orgaID uint,
	result *[]ProjectMember,
) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) IsUserOwnerOfOrga(_ *User, _ *Organization) (bool, error) {
	panic("not implemented")
}

func (f *FakeRepo) IsProjectOrganizationPaid(_ string) (bool, error) {
	panic("not implemented")
}

func (f *FakeRepo) CreateCheckoutSession(_ *CheckoutSession) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetCheckoutSession(_ string, _ *CheckoutSession) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) UpdateCheckoutSession(_ *CheckoutSession) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) DeleteCheckoutSession(_ *CheckoutSession) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) OrganizationSetPaid(
	organization *Organization,
	paid bool,
) IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetUserProjects(userID uint, projects *[]Project) IRepo {
	panic("not implemented")
}

func getRoleByEnvironmentTypeAndRole(
	environmentTypeID uint,
	roleID uint,
) RolesEnvironmentType {
	for _, re := range fakeRolesEnvironmentTypes {
		if re.RoleID == roleID && re.EnvironmentTypeID == environmentTypeID {
			return re
		}
	}

	return RolesEnvironmentType{}
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

func getRoleByUserID(userID uint) (role Role) {
	roleName, ok := fakeUserRole[userID]
	if !ok {
		roleName = "nothing"
	}

	role = Role{Name: roleName}
	findRole(&role)

	return role
}

func (fakeRepo *FakeRepo) Err() error {
	return fakeRepo.err
}

func (fakeRepo *FakeRepo) GetRolesEnvironmentType(
	rolesEnvironmentType *RolesEnvironmentType,
) IRepo {
	*rolesEnvironmentType = getRoleByEnvironmentTypeAndRole(
		rolesEnvironmentType.EnvironmentTypeID,
		rolesEnvironmentType.RoleID,
	)
	if rolesEnvironmentType.ID == 0 {
		fakeRepo.err = repo.ErrorNotFound
	}

	return fakeRepo
}

func (fakeRepo *FakeRepo) GetProjectMember(projectMember *ProjectMember) IRepo {
	role := getRoleByUserID(projectMember.UserID)
	projectMember.RoleID = role.ID
	projectMember.Role = role

	return fakeRepo
}

func (fakeRepo *FakeRepo) GetInvitableRoles(role Role, roles *[]Role) IRepo {
	// return fakeRepo
	// TODO
	return fakeRepo
}

func (fakeRepo *FakeRepo) GetChildrenRoles(role Role, roles *[]Role) IRepo {
	switch role.ID {
	case 4:
		*roles = []Role{fakeRoles[0], fakeRoles[1], fakeRoles[2]}
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

	userDev := &User{ID: 1, Username: "dev", UserID: "dev"}
	userLeadDev := &User{ID: 2, Username: "lead", UserID: "lead"}
	userDevops := &User{ID: 3, Username: "devops", UserID: "devops"}
	userAdmin := &User{ID: 4, Username: "admin", UserID: "admin"}

	environmentDev := &Environment{
		Name:              "dev",
		EnvironmentTypeID: 1,
		EnvironmentType: EnvironmentType{
			ID:   1,
			Name: "dev",
		},
	}
	environmentStaging := &Environment{
		Name:              "staging",
		EnvironmentTypeID: 2,
		EnvironmentType: EnvironmentType{
			ID:   2,
			Name: "staging",
		},
	}
	environmentProd := &Environment{
		Name:              "prod",
		EnvironmentTypeID: 3,
		EnvironmentType: EnvironmentType{
			ID:   3,
			Name: "prod",
		},
	}

	users := map[string]*User{
		"developer":      userDev,
		"lead developer": userLeadDev,
		"devops":         userDevops,
		"admin":          userAdmin,
	}

	environments := map[string]*Environment{
		"dev":     environmentDev,
		"staging": environmentStaging,
		"prod":    environmentProd,
	}

	rightsMatrix := map[string]map[string]rw{
		"developer": {
			"dev":     {true, true},
			"staging": {false, false},
			"prod":    {false, false},
		},
		"lead developer": {
			"dev":     {true, true},
			"staging": {false, false},
			"prod":    {false, false},
		},
		"devops": {
			"dev":     {true, true},
			"staging": {true, true},
			"prod":    {true, true},
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

	userDev := &User{ID: 1, Username: "dev", UserID: "dev"}
	userLeadDev := &User{ID: 2, Username: "lead", UserID: "lead"}
	userDevops := &User{ID: 3, Username: "devops", UserID: "devops"}
	userAdmin := &User{ID: 4, Username: "admin", UserID: "admin"}

	users := map[string]*User{
		"dev":    userDev,
		"lead":   userLeadDev,
		"devops": userDevops,
		"admin":  userAdmin,
	}

	rightsMatrix := map[string]map[*User]bool{
		"dev": {
			userDev:     false,
			userLeadDev: false,
			userDevops:  false,
			userAdmin:   false,
		},
		"lead": {
			userDev:     true,
			userLeadDev: true,
			userDevops:  false,
			userAdmin:   false,
		},
		"devops": {
			userDev:     true,
			userLeadDev: true,
			userDevops:  true,
			userAdmin:   false,
		},
		"admin": {
			userDev:     true,
			userLeadDev: true,
			userDevops:  true,
			userAdmin:   true,
		},
	}

	for name, user := range users {
		for otherName, otherUser := range users {
			expectation := rightsMatrix[name][otherUser]
			role := getRoleByUserID(otherUser.ID)

			canSetMemberRole, _ := CanUserSetMemberRole(
				&fakeRepo,
				*user,
				*otherUser,
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
