package rights

import (
	"math/rand"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"github.com/wearedevx/keystone/api/internal/emailer"
	"github.com/wearedevx/keystone/api/pkg/message"
	"github.com/wearedevx/keystone/api/pkg/models"
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
	fakeProjects              []Project
	fakeOrganizations         []Organization
)

const (
	DEV     int = 0
	LEAD        = 1
	DEVOPS      = 2
	ADMIN       = 3
	NOTHING     = 4
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

	fakeRoles[LEAD].Parent = &fakeRoles[DEVOPS]
	fakeRoles[DEV].Parent = &fakeRoles[NOTHING]
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

	fakeProjects = make([]Project, 0)
	fakeOrganizations = make([]Organization, 0)
}

type FakeRepo struct {
	called []string
	err    error
}

func newFakeRepo() *FakeRepo {
	f := new(FakeRepo)
	f.called = make([]string, 0)

	return f
}

func (f *FakeRepo) CreateEnvironment(_ *Environment) IRepo {
	f.called = append(f.called, "CreateEnvironment")
	return f
}

func (f *FakeRepo) CreateEnvironmentType(_ *EnvironmentType) IRepo {
	f.called = append(f.called, "CreateEnvironmentType")
	return f
}

func (f *FakeRepo) CreateLoginRequest() LoginRequest {
	f.called = append(f.called, "CreateLoginRequest")
	return LoginRequest{}
}

func (f *FakeRepo) CreateProjectMember(_ *ProjectMember, _ *Role) IRepo {
	f.called = append(f.called, "CreateProjectMember")
	return f
}

func (f *FakeRepo) CreateRole(_ *Role) IRepo {
	f.called = append(f.called, "CreateRole")
	return f
}

func (f *FakeRepo) CreateRoleEnvironmentType(_ *RolesEnvironmentType) IRepo {
	f.called = append(f.called, "CreateRoleEnvironmentType")
	return f
}

func (f *FakeRepo) CreateSecret(_ *Secret) {
	f.called = append(f.called, "CreateSecret")
}

func (f *FakeRepo) DeleteLoginRequest(_ string) bool {
	f.called = append(f.called, "DeleteLoginRequest")
	return false
}

func (f *FakeRepo) DeleteMessage(_ uint, _ uint) IRepo {
	f.called = append(f.called, "DeleteMessage")
	return f
}

func (f *FakeRepo) DeleteExpiredMessages() IRepo {
	f.called = append(f.called, "DeleteExpiredMessages")
	return f
}

func (f *FakeRepo) GetGroupedMessagesWillExpireByUser(
	groupedMessageUser *map[uint]emailer.GroupedMessagesUser,
) IRepo {
	f.called = append(f.called, "GetGroupedMessagesWillExpireByUser")
	return f
}

func (f *FakeRepo) FindUsers(
	_ []string,
	_ *map[string]User,
	_ *[]string,
) IRepo {
	f.called = append(f.called, "FindUsers")
	return f
}

func (f *FakeRepo) GetDb() *gorm.DB {
	f.called = append(f.called, "GetDb")
	return nil
}

func (f *FakeRepo) GetEnvironment(_ *Environment) IRepo {
	f.called = append(f.called, "GetEnvironment")
	return f
}

func (f *FakeRepo) GetEnvironmentsByProjectUUID(
	_ string,
	_ *[]Environment,
) IRepo {
	f.called = append(f.called, "GetEnvironmentsByProjectUUID")
	return f
}

func (f *FakeRepo) GetEnvironmentPublicKeys(_ string, _ *PublicKeys) IRepo {
	f.called = append(f.called, "GetEnvironmentPublicKeys")
	return f
}

func (f *FakeRepo) GetEnvironmentType(_ *EnvironmentType) IRepo {
	f.called = append(f.called, "GetEnvironmentType")
	return f
}

func (f *FakeRepo) GetLoginRequest(_ string) (LoginRequest, bool) {
	f.called = append(f.called, "GetLoginRequest")
	return LoginRequest{}, false
}

func (f *FakeRepo) GetMessagesForUserOnEnvironment(
	_ Device,
	_ Environment,
	_ *Message,
) IRepo {
	f.called = append(f.called, "GetMessagesForUserOnEnvironment")
	return f
}

func (f *FakeRepo) GetOrCreateEnvironment(_ *Environment) IRepo {
	f.called = append(f.called, "GetOrCreateEnvironment")
	return f
}

func (f *FakeRepo) GetOrCreateEnvironmentType(_ *EnvironmentType) IRepo {
	f.called = append(f.called, "GetOrCreateEnvironmentType")
	return f
}

func (f *FakeRepo) GetOrCreateProject(_ *Project) IRepo {
	f.called = append(f.called, "GetOrCreateProject")
	return f
}

func (f *FakeRepo) GetOrCreateProjectMember(_ *ProjectMember, _ string) IRepo {
	f.called = append(f.called, "GetOrCreateProjectMember")
	return f
}

func (f *FakeRepo) GetOrCreateRole(_ *Role) IRepo {
	f.called = append(f.called, "GetOrCreateRole")
	return f
}

func (f *FakeRepo) GetOrCreateRoleEnvType(_ *RolesEnvironmentType) IRepo {
	f.called = append(f.called, "GetOrCreateRoleEnvType")
	return f
}

func (f *FakeRepo) GetOrCreateUser(_ *User) IRepo {
	f.called = append(f.called, "GetOrCreateUser")
	return f
}

func (f *FakeRepo) GetProject(_ *Project) IRepo {
	f.called = append(f.called, "GetProject")
	return f
}

func (f *FakeRepo) GetProjectByUUID(_ string, _ *Project) IRepo {
	f.called = append(f.called, "GetProjectByUUID")
	return f
}

func (f *FakeRepo) GetRole(_ *Role) IRepo {
	f.called = append(f.called, "GetRole")
	return f
}

func (f *FakeRepo) GetRoles(_ *[]Role) IRepo {
	f.called = append(f.called, "GetRoles")
	return f
}

func (f *FakeRepo) GetRolesMemberCanInvite(
	projectMember ProjectMember,
	roles *[]Role,
) IRepo {
	f.called = append(f.called, "GetRolesMemberCanInvite")
	return f
}

func (f *FakeRepo) GetSecretByName(_ string, _ *Secret) {
	f.called = append(f.called, "GetSecretByName")
}

func (f *FakeRepo) GetUser(_ *User) IRepo {
	f.called = append(f.called, "GetUser")
	return f
}

func (f *FakeRepo) ListProjectMembers(
	userIDList []string,
	projectMember *[]ProjectMember,
) IRepo {
	f.called = append(f.called, "ListProjectMembers")
	return f
}

func (f *FakeRepo) ProjectAddMembers(_ Project, _ []MemberRole, _ User) IRepo {
	f.called = append(f.called, "ProjectAddMembers")
	return f
}

func (f *FakeRepo) UsersInMemberRoles(
	mers []MemberRole,
) (map[string]User, []string) {
	f.called = append(f.called, "UsersInMemberRoles")
	return map[string]User{}, []string{}
}

func (f *FakeRepo) ProjectGetMembers(_ *Project, _ *[]ProjectMember) IRepo {
	f.called = append(f.called, "ProjectGetMembers")
	return f
}

func (f *FakeRepo) ProjectLoadUsers(_ *Project) IRepo {
	f.called = append(f.called, "ProjectLoadUsers")
	return f
}

func (f *FakeRepo) ProjectRemoveMembers(_ Project, _ []string) IRepo {
	f.called = append(f.called, "ProjectRemoveMembers")
	return f
}

func (f *FakeRepo) ProjectSetRoleForUser(_ Project, _ User, _ Role) IRepo {
	f.called = append(f.called, "ProjectSetRoleForUser")
	return f
}

func (f *FakeRepo) RemoveOldMessageForRecipient(_ uint, _ string) IRepo {
	f.called = append(f.called, "RemoveOldMessageForRecipient")
	return f
}

func (f *FakeRepo) SetLoginRequestCode(_ string, _ string) LoginRequest {
	f.called = append(f.called, "SetLoginRequestCode")
	return LoginRequest{}
}

func (f *FakeRepo) SetNewVersionID(_ *Environment) error {
	f.called = append(f.called, "SetNewVersionID")
	return f.err
}

func (f *FakeRepo) WriteMessage(_ User, _ Message) IRepo {
	f.called = append(f.called, "WriteMessage")
	return f
}

func (f *FakeRepo) CheckMembersAreInProject(
	_ Project,
	_ []string,
) ([]string, error) {
	f.called = append(f.called, "{")
	return []string{}, nil
}

func (f *FakeRepo) DeleteAllProjectMembers(project *Project) IRepo {
	f.called = append(f.called, "DeleteAllProjectMembers")
	return f
}

func (f *FakeRepo) DeleteProject(project *Project) IRepo {
	f.called = append(f.called, "DeleteProject")
	return f
}

func (f *FakeRepo) DeleteProjectsEnvironments(project *Project) IRepo {
	f.called = append(f.called, "DeleteProjectsEnvironments")
	return f
}

func (f *FakeRepo) GetActivityLogs(
	projectID string,
	options GetLogsOptions,
	logs *[]ActivityLog,
) IRepo {
	f.called = append(f.called, "GetActivityLogs")
	return f
}

func (f *FakeRepo) GetMessage(message *Message) IRepo {
	f.called = append(f.called, "GetMessage")
	return f
}

func (f *FakeRepo) GetProjectsOrganization(
	projectuuid string,
	orga *Organization,
) IRepo {
	var project *Project

	for _, p := range fakeProjects {
		if p.UUID == projectuuid {
			project = &p
			break
		}
	}

	if project == nil {
		f.err = repo.ErrorNotFound
		return f
	}

	found := false
	for _, o := range fakeOrganizations {
		if o.ID == project.OrganizationID {
			*orga = o
			found = true
			break
		}
	}

	if !found {
		f.err = repo.ErrorNotFound
	}

	return f
}

func (f *FakeRepo) SetNewlyCreatedDevice(
	flag bool,
	deviceID uint,
	userID uint,
) repo.IRepo {
	f.called = append(f.called, "SetNewlyCreatedDevice")
	return f
}

func (f *FakeRepo) OrganizationCountMembers(_ *Organization, _ *int64) IRepo {
	f.called = append(f.called, "OrganizationCountMembers")
	return f
}

func (f *FakeRepo) GetUserByEmail(_ string, _ *[]User) IRepo {
	f.called = append(f.called, "GetUserByEmail")
	return f
}

func (f *FakeRepo) IsMemberOfProject(_ *Project, _ *ProjectMember) IRepo {
	f.called = append(f.called, "IsMemberOfProject")
	return f
}

func (f *FakeRepo) MessageService() *message.MessageService {
	f.called = append(f.called, "MessageService")
	return nil
}

func (f *FakeRepo) ProjectGetAdmins(
	project *Project,
	members *[]ProjectMember,
) IRepo {
	f.called = append(f.called, "ProjectGetAdmins")
	return f
}

func (f *FakeRepo) ProjectIsMemberAdmin(
	project *Project,
	member *ProjectMember,
) bool {
	f.called = append(f.called, "ProjectIsMemberAdmin")
	return false
}

func (f *FakeRepo) SaveActivityLog(al *ActivityLog) IRepo {
	f.called = append(f.called, "SaveActivityLog")
	return f
}

func (f *FakeRepo) GetDevices(_ uint, _ *[]Device) IRepo {
	f.called = append(f.called, "GetDevices")
	return f
}

func (f *FakeRepo) GetNewlyCreatedDevices(_ *[]Device) IRepo {
	f.called = append(f.called, "GetNewlyCreatedDevices")
	return f
}

func (f *FakeRepo) GetDevice(device *Device) IRepo {
	f.called = append(f.called, "GetDevice")
	return f
}

func (f *FakeRepo) GetDeviceByUserID(userID uint, device *Device) IRepo {
	f.called = append(f.called, "GetDeviceByUserID")
	return f
}

func (f *FakeRepo) UpdateDeviceLastUsedAt(deviceUID string) IRepo {
	f.called = append(f.called, "UpdateDeviceLastUsedAt")
	return f
}

func (f *FakeRepo) RevokeDevice(userID uint, deviceUID string) IRepo {
	f.called = append(f.called, "RevokeDevice")
	return f
}

func (f *FakeRepo) GetAdminsFromUserProjects(
	userID uint,
	adminProjectsMap *map[string][]string,
) IRepo {
	f.called = append(f.called, "GetAdminsFromUserProjects")
	return f
}

func (f *FakeRepo) CreateOrganization(orga *Organization) IRepo {
	f.called = append(f.called, "CreateOrganization")
	return f
}

func (f *FakeRepo) UpdateOrganization(orga *Organization) IRepo {
	f.called = append(f.called, "UpdateOrganization")
	return f
}

func (f *FakeRepo) OrganizationSetCustomer(
	organization *Organization,
	customer string,
) IRepo {
	f.called = append(f.called, "OrganizationSetCustomer")
	return f
}

func (f *FakeRepo) OrganizationSetSubscription(
	organization *Organization,
	subscription string,
) IRepo {
	f.called = append(f.called, "OrganizationSetSubscription")
	return f
}

func (f *FakeRepo) GetOrganization(orga *Organization) IRepo {
	f.called = append(f.called, "GetOrganization")
	return f
}

func (f *FakeRepo) GetOrganizations(userID uint, result *[]Organization) IRepo {
	f.called = append(f.called, "GetOrganizations")
	return f
}

func (f *FakeRepo) GetOwnedOrganizations(
	userID uint,
	result *[]Organization,
) IRepo {
	f.called = append(f.called, "GetOwnedOrganizations")
	return f
}

func (f *FakeRepo) GetOwnedOrganizationByName(
	userID uint,
	name string,
	orgas *[]Organization,
) IRepo {
	f.called = append(f.called, "GetOwnedOrganizations")
	return f
}

func (f *FakeRepo) GetOrganizationByName(
	userID uint,
	name string,
	orga *[]Organization,
) IRepo {
	f.called = append(f.called, "GetOrganizationByName")
	return f
}

func (f *FakeRepo) GetOrganizationProjects(
	_ *Organization,
	_ *[]Project,
) IRepo {
	f.called = append(f.called, "GetOrganizationProjects")
	return f
}

func (f *FakeRepo) GetOrganizationMembers(
	orgaID uint,
	result *[]ProjectMember,
) IRepo {
	{
		f.called = append(f.called, "GetOrganizationMembers")

		if orgaID == 0 {
			f.err = repo.ErrorNotFound
			return f
		}

		found := false
		for _, org := range fakeOrganizations {
			if org.ID == orgaID {
				found = true
				break
			}
		}

		if !found {
			f.err = repo.ErrorNotFound
			return f
		}

		*result = make([]ProjectMember, 0)
		for _, p := range fakeProjects {
			if p.OrganizationID == orgaID {
				*result = append(*result, p.Members...)
			}
		}

		if len(*result) == 0 {
			f.err = repo.ErrorNotFound
		}

		return f
	}
}

func (f *FakeRepo) IsUserOwnerOfOrga(_ *User, _ *Organization) (bool, error) {
	f.called = append(f.called, "IsUserOwnerOfOrga")
	return false, f.err
}

func (f *FakeRepo) IsProjectOrganizationPaid(
	projectUUID string,
) (paid bool, _ error) {
	f.called = append(f.called, "IsProjectOrganizationPaid")

	if projectUUID == "" {
		f.err = repo.ErrorNotFound
	} else {
		for _, p := range fakeProjects {
			if p.UUID == projectUUID {
				paid = p.Organization.Paid
				f.err = nil
				break
			}
		}
	}

	return paid, f.err
}

func (f *FakeRepo) CreateCheckoutSession(_ *CheckoutSession) IRepo {
	f.called = append(f.called, "CreateCheckoutSession")
	return f
}

func (f *FakeRepo) GetCheckoutSession(_ string, _ *CheckoutSession) IRepo {
	f.called = append(f.called, "GetCheckoutSession")
	return f
}

func (f *FakeRepo) UpdateCheckoutSession(_ *CheckoutSession) IRepo {
	f.called = append(f.called, "UpdateCheckoutSession")
	return f
}

func (f *FakeRepo) DeleteCheckoutSession(_ *CheckoutSession) IRepo {
	f.called = append(f.called, "DeleteCheckoutSession")
	return f
}

func (f *FakeRepo) OrganizationSetPaid(
	organization *Organization,
	paid bool,
) IRepo {
	f.called = append(f.called, "OrganizationSetPaid")
	return f
}

func (f *FakeRepo) GetUserProjects(userID uint, projects *[]Project) IRepo {
	f.called = append(f.called, "GetUserProjects")
	return f
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
	*role = fakeRoles[NOTHING]
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
	if rolesEnvironmentType.EnvironmentTypeID == 0 {
		fakeRepo.err = repo.ErrorNotFound
	}

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
	if projectMember.UserID == 0 {
		fakeRepo.err = repo.ErrorNotFound
		return fakeRepo
	}

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
		*roles = []Role{fakeRoles[DEV], fakeRoles[LEAD], fakeRoles[DEVOPS]}
	case 3:
		*roles = []Role{fakeRoles[DEV], fakeRoles[LEAD]}

	case 2:
		*roles = []Role{fakeRoles[DEV]}

	default:
		*roles = []Role{}
		fakeRepo.err = repo.ErrorNotFound
	}

	return fakeRepo
}

type rw struct {
	r   bool
	w   bool
	err bool
}

func TestCanUserHasRightEnvironment(t *testing.T) {
	fakeRepo := newFakeRepo()
	project := &Project{}

	userDev := &User{ID: 1, Username: "dev", UserID: "dev"}
	userLeadDev := &User{ID: 2, Username: "lead", UserID: "lead"}
	userDevops := &User{ID: 3, Username: "devops", UserID: "devops"}
	userAdmin := &User{ID: 4, Username: "admin", UserID: "admin"}
	notUser := &User{ID: 0, Username: "---", UserID: "---"}

	environmentDev := &Environment{
		ID:                1,
		Name:              "dev",
		EnvironmentTypeID: 1,
		EnvironmentType: EnvironmentType{
			ID:   1,
			Name: "dev",
		},
	}
	environmentStaging := &Environment{
		ID:                2,
		Name:              "staging",
		EnvironmentTypeID: 2,
		EnvironmentType: EnvironmentType{
			ID:   2,
			Name: "staging",
		},
	}
	environmentProd := &Environment{
		ID:                3,
		Name:              "prod",
		EnvironmentTypeID: 3,
		EnvironmentType: EnvironmentType{
			ID:   3,
			Name: "prod",
		},
	}
	environmentNot := &Environment{
		ID:                4,
		Name:              "---",
		EnvironmentTypeID: 0,
		EnvironmentType:   EnvironmentType{},
	}

	users := map[string]*User{
		"developer":      userDev,
		"lead developer": userLeadDev,
		"devops":         userDevops,
		"admin":          userAdmin,
		"---":            notUser,
	}

	environments := map[string]*Environment{
		"dev":     environmentDev,
		"staging": environmentStaging,
		"prod":    environmentProd,
		"---":     environmentNot,
	}

	rightsMatrix := map[string]map[string]rw{
		"developer": {
			//         read, write, err
			"dev":     {true, true, false},
			"staging": {false, false, false},
			"prod":    {false, false, false},
			"---":     {false, false, true},
		},
		"lead developer": {
			//         read, write, err
			"dev":     {true, true, false},
			"staging": {false, false, false},
			"prod":    {false, false, false},
			"---":     {false, false, true},
		},
		"devops": {
			//         read, write, err
			"dev":     {true, true, false},
			"staging": {true, true, false},
			"prod":    {true, true, false},
			"---":     {false, false, true},
		},
		"admin": {
			//         read, write, err
			"dev":     {true, true, false},
			"staging": {true, true, false},
			"prod":    {true, true, false},
			"---":     {false, false, true},
		},
		"---": {
			//         read, write, err
			"dev":     {false, false, true},
			"staging": {false, false, true},
			"prod":    {false, false, true},
			"---":     {false, false, true},
		},
	}

	for name, user := range users {
		for envName, environment := range environments {
			want := rightsMatrix[name][envName]
			fakeRepo.err = nil

			canRead, err := CanUserReadEnvironment(
				fakeRepo,
				user.ID,
				project.ID,
				environment,
			)

			assert.Equal(
				t,
				want.err,
				(err != nil),
				"Oops! Error while trying can read for user %s on environment %s: %v",
				name,
				envName,
				err,
			)

			canWrite, err := CanUserWriteOnEnvironment(
				fakeRepo,
				user.ID,
				project.ID,
				environment,
			)

			assert.Equal(
				t,
				want.err,
				(fakeRepo.Err() != nil),
				"Oops! Error whilr trying can write for user %s on environment %s: %v",
				name,
				envName,
				err,
			)

			assert.Equal(
				t,
				want.r,
				canRead,
				"Oops! User %s has unexpected read rights on %s environment",
				name,
				envName,
			)

			assert.Equal(
				t,
				want.w,
				canWrite,
				"Oops! User %s has unexpected write rights on %s environment",
				name,
				envName,
			)
		}
	}
}

func TestCanRoleAddRole(t *testing.T) {
	type args struct {
		Repo         repo.IRepo
		role         Role
		roleToInvite Role
	}

	tests := []struct {
		name    string
		args    args
		wantCan bool
		wantErr bool
	}{
		{
			name: "admin can add admin",
			args: args{
				Repo:         newFakeRepo(),
				role:         fakeRoles[ADMIN],
				roleToInvite: fakeRoles[ADMIN],
			},
			wantCan: true,
			wantErr: false,
		},
		{
			name: "admin can add devops",
			args: args{
				Repo:         newFakeRepo(),
				role:         fakeRoles[ADMIN],
				roleToInvite: fakeRoles[DEVOPS],
			},
			wantCan: true,
			wantErr: false,
		},
		{
			name: "admin can add lead-dev",
			args: args{
				Repo:         newFakeRepo(),
				role:         fakeRoles[ADMIN],
				roleToInvite: fakeRoles[LEAD],
			},
			wantCan: true,
			wantErr: false,
		},
		{
			name: "admin can add dev",
			args: args{
				Repo:         newFakeRepo(),
				role:         fakeRoles[ADMIN],
				roleToInvite: fakeRoles[DEV],
			},
			wantCan: true,
			wantErr: false,
		},
		{
			name: "devops can NOT add admin",
			args: args{
				Repo:         newFakeRepo(),
				role:         fakeRoles[DEVOPS],
				roleToInvite: fakeRoles[ADMIN],
			},
			wantCan: false,
			wantErr: false,
		},
		{
			name: "devops can add devops",
			args: args{
				Repo:         newFakeRepo(),
				role:         fakeRoles[DEVOPS],
				roleToInvite: fakeRoles[DEVOPS],
			},
			wantCan: true,
			wantErr: false,
		},
		{
			name: "devops can add lead-dev",
			args: args{
				Repo:         newFakeRepo(),
				role:         fakeRoles[DEVOPS],
				roleToInvite: fakeRoles[LEAD],
			},
			wantCan: true,
			wantErr: false,
		},
		{
			name: "devops can add dev",
			args: args{
				Repo:         newFakeRepo(),
				role:         fakeRoles[DEVOPS],
				roleToInvite: fakeRoles[DEV],
			},
			wantCan: true,
			wantErr: false,
		},
		{
			name: "lead-dev can NOT add admin",
			args: args{
				Repo:         newFakeRepo(),
				role:         fakeRoles[LEAD],
				roleToInvite: fakeRoles[ADMIN],
			},
			wantCan: false,
			wantErr: false,
		},
		{
			name: "lead-dev can NOT add devops",
			args: args{
				Repo:         newFakeRepo(),
				role:         fakeRoles[LEAD],
				roleToInvite: fakeRoles[DEVOPS],
			},
			wantCan: false,
			wantErr: false,
		},
		{
			name: "lead-dev can add lead-dev",
			args: args{
				Repo:         newFakeRepo(),
				role:         fakeRoles[LEAD],
				roleToInvite: fakeRoles[LEAD],
			},
			wantCan: true,
			wantErr: false,
		},
		{
			name: "lead-dev can add dev",
			args: args{
				Repo:         newFakeRepo(),
				role:         fakeRoles[LEAD],
				roleToInvite: fakeRoles[DEV],
			},
			wantCan: true,
			wantErr: false,
		},
		{
			name: "dev can NOT add admin",
			args: args{
				Repo:         newFakeRepo(),
				role:         fakeRoles[DEV],
				roleToInvite: fakeRoles[ADMIN],
			},
			wantCan: false,
			wantErr: false,
		},
		{
			name: "dev can NOT add devops",
			args: args{
				Repo:         newFakeRepo(),
				role:         fakeRoles[DEV],
				roleToInvite: fakeRoles[DEVOPS],
			},
			wantCan: false,
			wantErr: false,
		},
		{
			name: "dev can NOT add lead-dev",
			args: args{
				Repo:         newFakeRepo(),
				role:         fakeRoles[DEV],
				roleToInvite: fakeRoles[LEAD],
			},
			wantCan: false,
			wantErr: false,
		},
		{
			name: "dev can NOT add dev",
			args: args{
				Repo:         newFakeRepo(),
				role:         fakeRoles[DEV],
				roleToInvite: fakeRoles[DEV],
			},
			wantCan: false,
			wantErr: false,
		},
		{
			name: "error on bad role",
			args: args{
				Repo:         newFakeRepo(),
				role:         Role{ID: 12384, CanAddMember: true},
				roleToInvite: Role{ID: 25892},
			},
			wantCan: false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			can, err := CanRoleAddRole(
				tt.args.Repo,
				tt.args.role,
				tt.args.roleToInvite,
			)

			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Error CanRoleAddRole() err = %v, want %v",
					err,
					tt.wantErr,
				)
			}

			if can != tt.wantCan {
				t.Errorf(
					"Error CanRoleAddRole() can = %v, want %v",
					can,
					tt.wantCan,
				)
			}
		})
	}
}

func TestUserCanSetMemberRole(t *testing.T) {
	organization := Organization{}
	faker.FakeData(&organization)

	project := Project{
		ID:             12,
		UUID:           "8E733BFA-FFC7-412D-91AB-D1A9C3210A56",
		OrganizationID: organization.ID,
		Organization:   organization,
	}

	userDev := &User{ID: 1, Username: "dev", UserID: "dev"}
	userLeadDev := &User{ID: 2, Username: "lead", UserID: "lead"}
	userDevops := &User{ID: 3, Username: "devops", UserID: "devops"}
	userAdmin := &User{ID: 4, Username: "admin", UserID: "admin"}
	userAdminNotOwner := &User{ID: 5, Username: "notowner", UserID: "notowner"}
	userNot := &User{ID: 0, Username: "---", UserID: "---"}
	userBadRole := &User{ID: 135, Username: "badrole", UserID: "badrole"}

	organization.UserID = userAdmin.ID
	organization.User = *userAdmin

	fakeOrganizations = append(fakeOrganizations, organization)
	fakeProjects = append(fakeProjects, project)

	users := map[string]*User{
		"dev":      userDev,
		"lead":     userLeadDev,
		"devops":   userDevops,
		"admin":    userAdmin,
		"notowner": userAdminNotOwner,
		"---":      userNot,
		"bad":      userBadRole,
	}

	rightsMatrix := map[string]map[*User]map[int]struct {
		can bool
		err bool
	}{
		"dev": {
			// can dev set the role of a developer to
			userDev: {
				DEV:     {can: false, err: false},
				LEAD:    {can: false, err: false},
				DEVOPS:  {can: false, err: false},
				ADMIN:   {can: false, err: false},
				NOTHING: {can: false, err: false},
			},
			userLeadDev: {
				DEV:     {can: false, err: false},
				LEAD:    {can: false, err: false},
				DEVOPS:  {can: false, err: false},
				ADMIN:   {can: false, err: false},
				NOTHING: {can: false, err: false},
			},
			userDevops: {
				DEV:     {can: false, err: false},
				LEAD:    {can: false, err: false},
				DEVOPS:  {can: false, err: false},
				ADMIN:   {can: false, err: false},
				NOTHING: {can: false, err: false},
			},
			userAdmin: {
				DEV:     {can: false, err: false},
				LEAD:    {can: false, err: false},
				DEVOPS:  {can: false, err: false},
				ADMIN:   {can: false, err: false},
				NOTHING: {can: false, err: false},
			},
			userAdminNotOwner: {
				DEV:     {can: false, err: false},
				LEAD:    {can: false, err: false},
				DEVOPS:  {can: false, err: false},
				ADMIN:   {can: false, err: false},
				NOTHING: {can: false, err: false},
			},
			userNot: {
				DEV:     {can: false, err: true},
				LEAD:    {can: false, err: true},
				DEVOPS:  {can: false, err: true},
				ADMIN:   {can: false, err: true},
				NOTHING: {can: false, err: true},
			},
			userBadRole: {
				DEV:     {can: false, err: false},
				LEAD:    {can: false, err: false},
				DEVOPS:  {can: false, err: false},
				ADMIN:   {can: false, err: false},
				NOTHING: {can: false, err: false},
			},
		},
		"lead": {
			userDev: {
				DEV:     {can: true, err: false},
				LEAD:    {can: true, err: false},
				DEVOPS:  {can: false, err: false},
				ADMIN:   {can: false, err: false},
				NOTHING: {can: false, err: false},
			},
			userLeadDev: {
				DEV:     {can: true, err: false},
				LEAD:    {can: true, err: false},
				DEVOPS:  {can: false, err: false},
				ADMIN:   {can: false, err: false},
				NOTHING: {can: false, err: false},
			},
			userDevops: {
				DEV:     {can: false, err: false},
				LEAD:    {can: false, err: false},
				DEVOPS:  {can: false, err: false},
				ADMIN:   {can: false, err: false},
				NOTHING: {can: false, err: false},
			},
			userAdmin: {
				DEV:     {can: false, err: false},
				LEAD:    {can: false, err: false},
				DEVOPS:  {can: false, err: false},
				ADMIN:   {can: false, err: false},
				NOTHING: {can: false, err: false},
			},
			userAdminNotOwner: {
				DEV:     {can: false, err: false},
				LEAD:    {can: false, err: false},
				DEVOPS:  {can: false, err: false},
				ADMIN:   {can: false, err: false},
				NOTHING: {can: false, err: false},
			},
			userNot: {
				DEV:     {can: false, err: true},
				LEAD:    {can: false, err: true},
				DEVOPS:  {can: false, err: true},
				ADMIN:   {can: false, err: true},
				NOTHING: {can: false, err: true},
			},
			userBadRole: {
				DEV:     {can: false, err: false},
				LEAD:    {can: false, err: false},
				DEVOPS:  {can: false, err: false},
				ADMIN:   {can: false, err: false},
				NOTHING: {can: false, err: false},
			},
		},
		"devops": {
			userDev: {
				DEV:     {can: true, err: false},
				LEAD:    {can: true, err: false},
				DEVOPS:  {can: true, err: false},
				ADMIN:   {can: false, err: false},
				NOTHING: {can: false, err: false},
			},
			userLeadDev: {
				DEV:     {can: true, err: false},
				LEAD:    {can: true, err: false},
				DEVOPS:  {can: true, err: false},
				ADMIN:   {can: false, err: false},
				NOTHING: {can: false, err: false},
			},
			userDevops: {
				DEV:     {can: true, err: false},
				LEAD:    {can: true, err: false},
				DEVOPS:  {can: true, err: false},
				ADMIN:   {can: false, err: false},
				NOTHING: {can: false, err: false},
			},
			userAdmin: {
				DEV:     {can: false, err: false},
				LEAD:    {can: false, err: false},
				DEVOPS:  {can: false, err: false},
				ADMIN:   {can: false, err: false},
				NOTHING: {can: false, err: false},
			},
			userAdminNotOwner: {
				DEV:     {can: false, err: false},
				LEAD:    {can: false, err: false},
				DEVOPS:  {can: false, err: false},
				ADMIN:   {can: false, err: false},
				NOTHING: {can: false, err: false},
			},
			userNot: {
				DEV:     {can: false, err: true},
				LEAD:    {can: false, err: true},
				DEVOPS:  {can: false, err: true},
				ADMIN:   {can: false, err: true},
				NOTHING: {can: false, err: true},
			},
			userBadRole: {
				DEV:     {can: false, err: false},
				LEAD:    {can: false, err: false},
				DEVOPS:  {can: false, err: false},
				ADMIN:   {can: false, err: false},
				NOTHING: {can: false, err: false},
			},
		},
		"admin": {
			userDev: {
				DEV:     {can: true, err: false},
				LEAD:    {can: true, err: false},
				DEVOPS:  {can: true, err: false},
				ADMIN:   {can: true, err: false},
				NOTHING: {can: false, err: false},
			},
			userLeadDev: {
				DEV:     {can: true, err: false},
				LEAD:    {can: true, err: false},
				DEVOPS:  {can: true, err: false},
				ADMIN:   {can: true, err: false},
				NOTHING: {can: false, err: false},
			},
			userDevops: {
				DEV:     {can: true, err: false},
				LEAD:    {can: true, err: false},
				DEVOPS:  {can: true, err: false},
				ADMIN:   {can: true, err: false},
				NOTHING: {can: false, err: false},
			},
			// userAdmin is the owner, of the organization owning the
			// project, so their role cannot be changed
			userAdmin: {
				DEV:     {can: false, err: false},
				LEAD:    {can: false, err: false},
				DEVOPS:  {can: false, err: false},
				ADMIN:   {can: false, err: false},
				NOTHING: {can: false, err: false},
			},
			userAdminNotOwner: {
				DEV:     {can: false, err: false},
				LEAD:    {can: false, err: false},
				DEVOPS:  {can: false, err: false},
				ADMIN:   {can: false, err: false},
				NOTHING: {can: false, err: false},
			},
			userNot: {
				DEV:     {can: false, err: true},
				LEAD:    {can: false, err: true},
				DEVOPS:  {can: false, err: true},
				ADMIN:   {can: false, err: true},
				NOTHING: {can: false, err: true},
			},
			userBadRole: {
				DEV:     {can: false, err: false},
				LEAD:    {can: false, err: false},
				DEVOPS:  {can: false, err: false},
				ADMIN:   {can: false, err: false},
				NOTHING: {can: false, err: false},
			},
		},
		"---": {
			userDev: {
				DEV:     {can: false, err: true},
				LEAD:    {can: false, err: true},
				DEVOPS:  {can: false, err: true},
				ADMIN:   {can: false, err: true},
				NOTHING: {can: false, err: true},
			},
			userLeadDev: {
				DEV:     {can: false, err: true},
				LEAD:    {can: false, err: true},
				DEVOPS:  {can: false, err: true},
				ADMIN:   {can: false, err: true},
				NOTHING: {can: false, err: true},
			},
			userDevops: {
				DEV:     {can: false, err: true},
				LEAD:    {can: false, err: true},
				DEVOPS:  {can: false, err: true},
				ADMIN:   {can: false, err: true},
				NOTHING: {can: false, err: true},
			},
			userAdminNotOwner: {
				DEV:     {can: false, err: true},
				LEAD:    {can: false, err: true},
				DEVOPS:  {can: false, err: true},
				ADMIN:   {can: false, err: true},
				NOTHING: {can: false, err: true},
			},
			userAdmin: {
				DEV:     {can: false, err: true},
				LEAD:    {can: false, err: true},
				DEVOPS:  {can: false, err: true},
				ADMIN:   {can: false, err: true},
				NOTHING: {can: false, err: true},
			},
			userNot: {
				DEV:     {can: false, err: true},
				LEAD:    {can: false, err: true},
				DEVOPS:  {can: false, err: true},
				ADMIN:   {can: false, err: true},
				NOTHING: {can: false, err: true},
			},
			userBadRole: {
				DEV:     {can: false, err: true},
				LEAD:    {can: false, err: true},
				DEVOPS:  {can: false, err: true},
				ADMIN:   {can: false, err: true},
				NOTHING: {can: false, err: true},
			},
		},
		"bad": {
			userDev: {
				DEV:     {can: false, err: true},
				LEAD:    {can: false, err: true},
				DEVOPS:  {can: false, err: true},
				ADMIN:   {can: false, err: true},
				NOTHING: {can: false, err: true},
			},
			userLeadDev: {
				DEV:     {can: false, err: true},
				LEAD:    {can: false, err: true},
				DEVOPS:  {can: false, err: true},
				ADMIN:   {can: false, err: true},
				NOTHING: {can: false, err: true},
			},
			userDevops: {
				DEV:     {can: false, err: true},
				LEAD:    {can: false, err: true},
				DEVOPS:  {can: false, err: true},
				ADMIN:   {can: false, err: true},
				NOTHING: {can: false, err: true},
			},
			userAdmin: {
				DEV:     {can: false, err: true},
				LEAD:    {can: false, err: true},
				DEVOPS:  {can: false, err: true},
				ADMIN:   {can: false, err: true},
				NOTHING: {can: false, err: true},
			},
			userAdminNotOwner: {
				DEV:     {can: false, err: true},
				LEAD:    {can: false, err: true},
				DEVOPS:  {can: false, err: true},
				ADMIN:   {can: false, err: true},
				NOTHING: {can: true, err: false},
			},
			userNot: {
				DEV:     {can: false, err: true},
				LEAD:    {can: false, err: true},
				DEVOPS:  {can: false, err: true},
				ADMIN:   {can: false, err: true},
				NOTHING: {can: false, err: true},
			},
			userBadRole: {
				DEV:     {can: false, err: true},
				LEAD:    {can: false, err: true},
				DEVOPS:  {can: false, err: true},
				ADMIN:   {can: false, err: true},
				NOTHING: {can: true, err: false},
			},
		},
	}

	for name, user := range users {
		for otherName, otherUser := range users {
			spec := rightsMatrix[name][otherUser]

			for i, want := range spec {
				fakeRepo := newFakeRepo()
				targetRole := fakeRoles[i]

				can, err := CanUserSetMemberRole(
					fakeRepo,
					*user,
					*otherUser,
					targetRole,
					project,
				)

				if (err != nil) != want.err {
					t.Errorf(
						"Error CanUserSetMemberRole(%s, %s, to %s) err = %v, want %v",
						name,
						otherName,
						targetRole.Name,
						err,
						want.err,
					)
				}

				if can != want.can {
					t.Errorf(
						"Error CanUserSetMemberRole(%s, %s, to %s) can = %v, want %v",
						name,
						otherName,
						targetRole.Name,
						can,
						want.can,
					)
				}
			}
		}
	}
}

func TestCanUserAddMemberWithRole(t *testing.T) {
	project := Project{}

	userNotFound := &User{ID: 0, Username: "---", UserID: "---"}
	userDev := &User{ID: 1, Username: "dev", UserID: "dev"}
	userLeadDev := &User{ID: 2, Username: "lead", UserID: "lead"}
	userDevops := &User{ID: 3, Username: "devops", UserID: "devops"}
	userAdmin := &User{ID: 4, Username: "admin", UserID: "admin"}

	users := map[string]*User{
		"dev":    userDev,
		"lead":   userLeadDev,
		"devops": userDevops,
		"admin":  userAdmin,
		"---":    userNotFound,
	}

	type want struct {
		can bool
		err bool
	}

	rightsMatrix := map[string]map[*User]want{
		"dev": {
			userDev:     {can: false, err: false},
			userLeadDev: {can: false, err: false},
			userDevops:  {can: false, err: false},
			userAdmin:   {can: false, err: false},
		},
		"lead": {
			userDev:     {can: true, err: false},
			userLeadDev: {can: true, err: false},
			userDevops:  {can: false, err: false},
			userAdmin:   {can: false, err: false},
		},
		"devops": {
			userDev:     {can: true, err: false},
			userLeadDev: {can: true, err: false},
			userDevops:  {can: true, err: false},
			userAdmin:   {can: false, err: false},
		},
		"admin": {
			userDev:     {can: true, err: false},
			userLeadDev: {can: true, err: false},
			userDevops:  {can: true, err: false},
			userAdmin:   {can: true, err: false},
		},
		"---": {
			userDev:      {can: false, err: true},
			userLeadDev:  {can: false, err: true},
			userDevops:   {can: false, err: true},
			userAdmin:    {can: false, err: true},
			userNotFound: {can: false, err: true},
		},
	}

	for name, user := range users {
		for otherName, otherUser := range users {
			want := rightsMatrix[name][otherUser]

			role := getRoleByUserID(otherUser.ID)
			fakeRepo := newFakeRepo()

			can, err := CanUserAddMemberWithRole(
				fakeRepo,
				*user,
				role,
				project,
			)

			if (err != nil) != want.err {
				t.Errorf(
					"Error CanUserAddMemberWithRole(%s, %s to %s) err = %v, want %v",
					name,
					otherName,
					role.Name,
					err,
					want.err,
				)
			}

			if can != want.can {
				t.Errorf(
					"Error CanUserAddMemberWithRole(%s, %s to %s) can = %v, want %v",
					name,
					otherName,
					role.Name,
					can,
					want.can,
				)
			}
		}
	}
}

func fakeManyOrgs(orgs []*Organization) {
	for _, org := range orgs {
		o := Organization{}
		err := faker.FakeData(&o)
		if err != nil {
			panic(err)
		}
		o.ID = uint(rand.Intn(900)) + 100
		*org = o
	}
}

func fakeManyProjects(projects []*Project) {
	for _, project := range projects {
		p := Project{}
		err := faker.FakeData(&p)
		if err != nil {
			panic(err)
		}
		*project = p
	}
}

func TestHasOrganizationNotPaidAndHasNonAdmin(t *testing.T) {
	freeFailingOrg := Organization{}
	freeOKOrg := Organization{}
	paidOrg := Organization{}
	fakeManyOrgs([]*Organization{&freeFailingOrg, &freeOKOrg, &paidOrg})
	paidOrg.Paid = true

	fakeOrganizations = append(
		fakeOrganizations,
		freeFailingOrg,
		freeOKOrg,
		paidOrg,
	)

	badFreeProject := Project{}
	okFreeProject := Project{}
	okPaidProject := Project{}
	projectBadOrg := Project{}
	fakeManyProjects([]*Project{
		&badFreeProject,
		&okFreeProject,
		&okPaidProject,
		&projectBadOrg,
	})
	badFreeProject.OrganizationID = freeFailingOrg.ID
	badFreeProject.Organization = freeFailingOrg
	okFreeProject.OrganizationID = freeOKOrg.ID
	okFreeProject.Organization = freeOKOrg
	okPaidProject.OrganizationID = paidOrg.ID
	okPaidProject.Organization = paidOrg

	badFreeProject.Members = []ProjectMember{
		{
			ProjectID: badFreeProject.ID,
			Role:      fakeRoles[ADMIN],
		},
		{
			ProjectID: badFreeProject.ID,
			Role:      fakeRoles[DEV],
		},
	}

	okFreeProject.Members = []ProjectMember{
		{
			ProjectID: okFreeProject.ID,
			Role:      fakeRoles[ADMIN],
		},
		{
			ProjectID: okFreeProject.ID,
			Role:      fakeRoles[ADMIN],
		},
	}

	okPaidProject.Members = []ProjectMember{
		{
			ProjectID: okPaidProject.ID,
			Role:      fakeRoles[ADMIN],
		},
		{
			ProjectID: okPaidProject.ID,
			Role:      fakeRoles[DEV],
		},
	}

	fakeProjects = append(
		fakeProjects,
		badFreeProject,
		okFreeProject,
		okPaidProject,
		projectBadOrg,
	)

	type args struct {
		Repo    repo.IRepo
		project models.Project
	}
	tests := []struct {
		name    string
		args    args
		wantHas bool
		wantErr bool
	}{
		{
			name: "works and returns true (free orga, non admin member)",
			args: args{
				Repo:    newFakeRepo(),
				project: badFreeProject,
			},
			wantHas: true,
			wantErr: false,
		},
		{
			name: "works and returns false (paid orga, non admin members)",
			args: args{
				Repo:    newFakeRepo(),
				project: okPaidProject,
			},
			wantHas: false,
			wantErr: false,
		},
		{
			name: "works and returns false (free orga, only admin members)",
			args: args{
				Repo:    newFakeRepo(),
				project: okFreeProject,
			},
			wantHas: false,
			wantErr: false,
		},
		{
			name: "organization does not exists",
			args: args{
				Repo:    newFakeRepo(),
				project: projectBadOrg,
			},
			wantHas: false,
			wantErr: true,
		},
		{
			name: "project does not exists",
			args: args{
				Repo:    newFakeRepo(),
				project: Project{},
			},
			wantHas: false,
			wantErr: true,
		},
		// {
		// 	name: "fails getting members",
		// 	args: args{
		// 		Repo:    newFakeRepo(),
		// 		project: Project{},
		// 	},
		// 	wantHas: false,
		// 	wantErr: false,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHas, err := HasOrganizationNotPaidAndHasNonAdmin(
				tt.args.Repo,
				tt.args.project,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"HasOrganizationNotPaidAndHasNonAdmin() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if gotHas != tt.wantHas {
				t.Errorf(
					"HasOrganizationNotPaidAndHasNonAdmin() = %v, want %v",
					gotHas,
					tt.wantHas,
				)
			}
		})
	}
}
