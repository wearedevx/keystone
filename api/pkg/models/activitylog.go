package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
	"unsafe"

	"gorm.io/gorm"
)

type ActivityLog struct {
	ID            uint        `json:"id" gorm:"primaryKey"`
	UserID        *uint       `json:"user_id"`
	User          User        `json:"user"`
	ProjectID     *uint       `json:"project_id"`
	Project       Project     `json:"project"`
	EnvironmentID *uint       `json:"environment_id"`
	Environment   Environment `json:"environment"`
	Action        string      `json:"action"`
	Success       bool        `json:"success"`
	Message       string      `json:"error"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

func (pm *ActivityLog) BeforeCreate(tx *gorm.DB) (err error) {
	pm.CreatedAt = time.Now()
	pm.UpdatedAt = time.Now()

	return nil
}

func (pm *ActivityLog) BeforeUpdate(tx *gorm.DB) (err error) {
	pm.UpdatedAt = time.Now()

	return nil
}

func (pm *ActivityLog) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(pm)
}

func (pm *ActivityLog) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(pm)

	*out = sb.String()

	return err
}

func (pm ActivityLog) Error() string {
	return pm.Message
}

func (pm *ActivityLog) SetError(err error) *ActivityLog {
	if err != nil {
		message := err.Error()
		parent := errors.Unwrap(err)

		if parent != nil {
			message = message + ": " + parent.Error()
		}

		pm.Message = message
		pm.Success = false
	} else {
		pm.Success = true
	}

	return pm
}

func (pm *ActivityLog) Lite() (l ActivityLogLite) {
	l.UserID = pm.User.UserID
	l.ProjectName = pm.Project.Name
	l.EnvironmentName = pm.Environment.Name
	l.Action = pm.Action
	l.Success = pm.Success
	l.ErrorMessage = pm.Message
	l.CreatedAt = pm.CreatedAt

	return l
}

func (pm *ActivityLog) Ptr() unsafe.Pointer {
	return (unsafe.Pointer)(pm)
}

func ErrorIsActivityLog(err error) bool {
	activityLogPtrType := fmt.Sprintf("%T", &ActivityLog{})
	errType := fmt.Sprintf("%T", err)

	return activityLogPtrType == errType
}

/// API types

// A lighter version of the activity log with only information
// that is safe to display (no db identifiers)
type ActivityLogLite struct {
	UserID          string    `json:"user_id"`
	ProjectName     string    `json:"project_name"`
	EnvironmentName string    `json:"environment_name"`
	Action          string    `json:"action"`
	Success         bool      `json:"success"`
	ErrorMessage    string    `json:"error_message"`
	CreatedAt       time.Time `json:"created_at"`
}

func (pm *ActivityLogLite) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(pm)
}

func (pm *ActivityLogLite) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(pm)
	*out = sb.String()

	return err
}

type GetActivityLogResponse struct {
	Logs []ActivityLogLite
}

func (pm *GetActivityLogResponse) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(pm)
}

func (pm *GetActivityLogResponse) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(pm)
	*out = sb.String()

	return err
}

type ActionFilter string

const (
	ActionFilterGetRoles                     ActionFilter = "GetRoles"
	ActionFilterDoUsersExist                 ActionFilter = "DoUsersExist"
	ActionFilterPutMemberSetRole             ActionFilter = "PutMemberSetRole"
	ActionFilterPostProject                  ActionFilter = "PostProject"
	ActionFilterGetProjects                  ActionFilter = "GetProjects"
	ActionFilterGetProjectMembers            ActionFilter = "GetProjectMembers"
	ActionFilterPostProjectMembers           ActionFilter = "PostProjectMembers"
	ActionFilterDeleteProjectsMembers        ActionFilter = "DeleteProjectsMembers"
	ActionFilterGetAccessibleEnvironments    ActionFilter = "GetAccessibleEnvironments"
	ActionFilterDeleteProject                ActionFilter = "DeleteProject"
	ActionFilterGetProjectsOrganization      ActionFilter = "GetProjectsOrganization"
	ActionFilterGetDevices                   ActionFilter = "GetDevices"
	ActionFilterDeleteDevice                 ActionFilter = "DeleteDevice"
	ActionFilterGetOrganizations             ActionFilter = "GetOrganizations"
	ActionFilterPostOrganization             ActionFilter = "PostOrganization"
	ActionFilterGetEnvironmentPublicKeys     ActionFilter = "GetEnvironmentPublicKeys"
	ActionFilterPostInvite                   ActionFilter = "PostInvite"
	ActionFilterGetMessagesFromProjectByUser ActionFilter = "GetMessagesFromProjectByUser"
	ActionFilterWriteMessages                ActionFilter = "WriteMessages"
	ActionFilterDeleteMessage                ActionFilter = "DeleteMessage"
	ActionFilterPostSubscription             ActionFilter = "PostSubscription"
	ActionFilterGetPollSubscriptionSuccess   ActionFilter = "GetPollSubscriptionSuccess"
	ActionFilterManageSubscription           ActionFilter = "ManageSubscription"
	ActionFilterPostUser                     ActionFilter = "PostUser"
	ActionFilterPostUserToken                ActionFilter = "PostUserToken"
	ActionFilterPostLoginRequest             ActionFilter = "PostLoginRequest"
	ActionFilterGetLoginRequest              ActionFilter = "GetLoginRequest"
	ActionFilterGetUserKeys                  ActionFilter = "GetUserKeys"
)

func (af ActionFilter) Validate() (ok bool) {
	switch af {
	case ActionFilterGetRoles,
		ActionFilterDoUsersExist,
		ActionFilterPutMemberSetRole,
		ActionFilterPostProject,
		ActionFilterGetProjects,
		ActionFilterGetProjectMembers,
		ActionFilterPostProjectMembers,
		ActionFilterDeleteProjectsMembers,
		ActionFilterGetAccessibleEnvironments,
		ActionFilterDeleteProject,
		ActionFilterGetProjectsOrganization,
		ActionFilterGetDevices,
		ActionFilterDeleteDevice,
		ActionFilterGetOrganizations,
		ActionFilterPostOrganization,
		ActionFilterGetEnvironmentPublicKeys,
		ActionFilterPostInvite,
		ActionFilterGetMessagesFromProjectByUser,
		ActionFilterWriteMessages,
		ActionFilterDeleteMessage,
		ActionFilterPostSubscription,
		ActionFilterGetPollSubscriptionSuccess,
		ActionFilterManageSubscription,
		ActionFilterPostUser,
		ActionFilterPostUserToken,
		ActionFilterPostLoginRequest,
		ActionFilterGetLoginRequest,
		ActionFilterGetUserKeys:
		return true
	default:
		return false
	}
}

type EnvironmentFilter string

const (
	EnvironmentFilterDev     EnvironmentFilter = "dev"
	EnvironmentFilterStaging EnvironmentFilter = "staging"
	EnvironmentFilterProd    EnvironmentFilter = "prod"
)

func (ef EnvironmentFilter) Validate() (ok bool) {
	switch ef {
	case EnvironmentFilterDev,
		EnvironmentFilterStaging,
		EnvironmentFilterProd:
		return true
	default:
		return false
	}
}

type GetLogsOptions struct {
	Actions      []ActionFilter      `json:"actions"`
	Environments []EnvironmentFilter `json:"environments"`
	Users        []string            `json:"users"`
	Limit        uint64              `json:"limit" default:"200"`
}

func NewGetLogsOption() *GetLogsOptions {
	return &GetLogsOptions{
		Actions:      make([]ActionFilter, 0),
		Environments: make([]EnvironmentFilter, 0),
		Users:        make([]string, 0),
	}
}

func (pm *GetLogsOptions) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(pm)
}

func (pm *GetLogsOptions) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(pm)
	*out = sb.String()

	return err
}

func (o *GetLogsOptions) SetActions(actions string) *GetLogsOptions {
	if actions == "" {
		return o
	}

	for _, action := range strings.Split(actions, ",") {
		if f := ActionFilter(strings.TrimSpace(action)); f.Validate() {
			o.Actions = append(o.Actions, f)
		}
	}

	return o
}

func (o *GetLogsOptions) SetEnvironments(environments string) *GetLogsOptions {
	if environments == "" {
		return o
	}

	for _, environment := range strings.Split(environments, ",") {
		if e := EnvironmentFilter(strings.TrimSpace(environment)); e.Validate() {
			o.Environments = append(o.Environments, EnvironmentFilter(environment))
		}
	}

	return o
}

func (o *GetLogsOptions) SetUsers(users string) *GetLogsOptions {
	if users == "" {
		return o
	}

	for _, user := range strings.Split(users, ",") {
		user = strings.TrimSpace(user)
		if user != "" {
			o.Users = append(o.Users, user)
		}
	}
	return o
}

func (o *GetLogsOptions) SetLimit(limit uint64) *GetLogsOptions {
	if limit == 0 {
		return o
	}

	o.Limit = limit

	return o
}
