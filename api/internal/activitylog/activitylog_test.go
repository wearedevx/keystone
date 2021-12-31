package activitylog

import (
	"errors"
	"reflect"
	"testing"

	"github.com/wearedevx/keystone/api/internal/emailer"
	"github.com/wearedevx/keystone/api/pkg/message"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
	"gorm.io/gorm"
)

func TestNewActivityLogger(t *testing.T) {
	type args struct {
		repo repo.IRepo
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "creates an activity logger",
			args: args{
				repo: new(FakeRepo),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewActivityLogger(tt.args.repo).Err(); err != nil &&
				err != tt.wantErr {
				t.Errorf(
					"NewActivityLogger().Err() = %v, want %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func Test_activityLogger_Err(t *testing.T) {
	type fields struct {
		err  error
		repo repo.IRepo
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "no error in logger",
			fields: fields{
				err:  nil,
				repo: new(FakeRepo),
			},
			wantErr: false,
		},
		{
			name: "error in logger",
			fields: fields{
				err:  errors.New("some error that happened"),
				repo: new(FakeRepo),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &activityLogger{
				err:  tt.fields.err,
				repo: tt.fields.repo,
			}
			if err := logger.Err(); (err != nil) != tt.wantErr {
				t.Errorf(
					"activityLogger.Err() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func Test_activityLogger_Save(t *testing.T) {
	type fields struct {
		err  error
		repo repo.IRepo
	}
	type args struct {
		err error
	}
	fakeRepo := new(FakeRepo)
	tests := []struct {
		name   string
		fields fields
		args   args
		want   ActivityLogger
	}{
		{
			name: "saves an activity log",
			fields: fields{
				err:  nil,
				repo: fakeRepo,
			},
			args: args{},
			want: &activityLogger{
				nil,
				fakeRepo,
			},
		},
		{
			name: "does not save a plain error",
			fields: fields{
				err:  nil,
				repo: new(FakeRepo),
			},
			args: args{},
			want: &activityLogger{
				nil,
				fakeRepo,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &activityLogger{
				err:  tt.fields.err,
				repo: tt.fields.repo,
			}
			if got := logger.Save(tt.args.err); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("activityLogger.Save() = %v, want %v", got, tt.want)
			}
		})
	}
}

type FakeRepo struct{}

func (f *FakeRepo) CreateEnvironment(_ *models.Environment) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) CreateEnvironmentType(_ *models.EnvironmentType) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) CreateLoginRequest() models.LoginRequest {
	panic("not implemented")
}

func (f *FakeRepo) CreateProjectMember(
	_ *models.ProjectMember,
	_ *models.Role,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) CreateRole(_ *models.Role) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) CreateRoleEnvironmentType(
	_ *models.RolesEnvironmentType,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) DeleteLoginRequest(_ string) bool {
	panic("not implemented")
}

func (f *FakeRepo) DeleteAllProjectMembers(project *models.Project) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) DeleteExpiredMessages() repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetGroupedMessagesWillExpireByUser(
	groupedMessageUser *map[uint]emailer.GroupedMessagesUser,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) DeleteMessage(messageID uint, userID uint) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) DeleteProject(project *models.Project) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) DeleteProjectsEnvironments(
	project *models.Project,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) Err() error {
	panic("not implemented")
}

func (f *FakeRepo) FindUsers(
	userIDs []string,
	users *map[string]models.User,
	notFounds *[]string,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetActivityLogs(
	projectID string,
	options models.GetLogsOptions,
	logs *[]models.ActivityLog,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetChildrenRoles(
	role models.Role,
	roles *[]models.Role,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetDb() *gorm.DB {
	panic("not implemented")
}

func (f *FakeRepo) GetEnvironment(_ *models.Environment) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetEnvironmentPublicKeys(
	envID string,
	publicKeys *models.PublicKeys,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetEnvironmentType(_ *models.EnvironmentType) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetEnvironmentsByProjectUUID(
	projectUUID string,
	foundEnvironments *[]models.Environment,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetInvitableRoles(
	_ models.Role,
	_ *[]models.Role,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetLoginRequest(_ string) (models.LoginRequest, bool) {
	panic("not implemented")
}

func (f *FakeRepo) GetMessage(message *models.Message) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetMessagesForUserOnEnvironment(
	device models.Device,
	environment models.Environment,
	message *models.Message,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOrCreateEnvironment(_ *models.Environment) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOrCreateEnvironmentType(
	_ *models.EnvironmentType,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOrCreateProject(_ *models.Project) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOrCreateProjectMember(
	_ *models.ProjectMember,
	_ string,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOrCreateRole(_ *models.Role) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOrCreateRoleEnvType(
	_ *models.RolesEnvironmentType,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOrCreateUser(_ *models.User) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetProject(_ *models.Project) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetProjectByUUID(_ string, _ *models.Project) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetProjectMember(_ *models.ProjectMember) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) UsersInMemberRoles(mers []models.MemberRole) (map[string]models.User, []string) {
	panic("not implemented")
}

func (r *FakeRepo) SetNewlyCreatedDevice(flag bool, deviceID uint, userID uint) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetProjectsOrganization(
	_ string,
	_ *models.Organization,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) OrganizationCountMembers(
	_ *models.Organization,
	_ *int64,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetRole(_ *models.Role) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetRoles(_ *[]models.Role) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetRolesEnvironmentType(
	_ *models.RolesEnvironmentType,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetRolesMemberCanInvite(
	projectMember models.ProjectMember,
	roles *[]models.Role,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetUser(_ *models.User) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetUserByEmail(_ string, _ *[]models.User) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) IsMemberOfProject(
	_ *models.Project,
	_ *models.ProjectMember,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) ListProjectMembers(
	userIDList []string,
	projectMember *[]models.ProjectMember,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) MessageService() *message.MessageService {
	panic("not implemented")
}

func (f *FakeRepo) ProjectAddMembers(
	_ models.Project,
	_ []models.MemberRole,
	_ models.User,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) ProjectGetAdmins(
	project *models.Project,
	members *[]models.ProjectMember,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) ProjectIsMemberAdmin(
	project *models.Project,
	member *models.ProjectMember,
) bool {
	panic("not implemented")
}

func (f *FakeRepo) ProjectGetMembers(
	_ *models.Project,
	_ *[]models.ProjectMember,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) ProjectLoadUsers(_ *models.Project) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) ProjectRemoveMembers(
	_ models.Project,
	_ []string,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) ProjectSetRoleForUser(
	_ models.Project,
	_ models.User,
	_ models.Role,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) CheckMembersAreInProject(
	_ models.Project,
	_ []string,
) ([]string, error) {
	panic("not implemented")
}

func (f *FakeRepo) RemoveOldMessageForRecipient(
	userID uint,
	environmentID string,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) SaveActivityLog(al *models.ActivityLog) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) SetLoginRequestCode(_ string, _ string) models.LoginRequest {
	panic("not implemented")
}

func (f *FakeRepo) SetNewVersionID(environment *models.Environment) error {
	panic("not implemented")
}

func (f *FakeRepo) WriteMessage(
	user models.User,
	message models.Message,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetDevices(_ uint, _ *[]models.Device) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetDevice(device *models.Device) repo.IRepo {
	panic("not implemented")
}
func (f *FakeRepo) GetNewlyCreatedDevices(_ *[]models.Device) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetDeviceByUserID(
	userID uint,
	device *models.Device,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) UpdateDeviceLastUsedAt(deviceUID string) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) RevokeDevice(userID uint, deviceUID string) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetAdminsFromUserProjects(
	userID uint,
	adminProjectsMap *map[string][]string,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) CreateOrganization(orga *models.Organization) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) UpdateOrganization(orga *models.Organization) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) OrganizationSetCustomer(
	organization *models.Organization,
	customer string,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) OrganizationSetSubscription(
	organization *models.Organization,
	subscription string,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOrganization(orga *models.Organization) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOrganizations(
	userID uint,
	result *[]models.Organization,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOwnedOrganizations(
	userID uint,
	result *[]models.Organization,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOwnedOrganizationByName(
	userID uint,
	name string,
	orgas *[]models.Organization,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOrganizationByName(
	userID uint,
	name string,
	orga *[]models.Organization,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOrganizationProjects(
	_ *models.Organization,
	_ *[]models.Project,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetOrganizationMembers(
	orgaID uint,
	result *[]models.ProjectMember,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) IsUserOwnerOfOrga(
	_ *models.User,
	_ *models.Organization,
) (bool, error) {
	panic("not implemented")
}

func (f *FakeRepo) IsProjectOrganizationPaid(_ string) (bool, error) {
	panic("not implemented")
}

func (f *FakeRepo) CreateCheckoutSession(_ *models.CheckoutSession) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetCheckoutSession(
	_ string,
	_ *models.CheckoutSession,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) UpdateCheckoutSession(_ *models.CheckoutSession) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) DeleteCheckoutSession(_ *models.CheckoutSession) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) OrganizationSetPaid(
	organization *models.Organization,
	paid bool,
) repo.IRepo {
	panic("not implemented")
}

func (f *FakeRepo) GetUserProjects(
	userID uint,
	projects *[]models.Project,
) repo.IRepo {
	panic("not implemented")
}
