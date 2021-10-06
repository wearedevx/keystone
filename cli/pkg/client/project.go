package client

import (
	"fmt"

	"github.com/wearedevx/keystone/api/pkg/models"
)

type Project struct {
	id string
	r  requester
}

func (p *Project) Init(name string, organizationID uint) (models.Project, error) {
	var project models.Project

	payload := models.Project{
		Name:           name,
		OrganizationID: organizationID,
	}

	err := p.r.post("/projects", payload, &project, nil)

	return project, err
}

func (p *Project) GetAllMembers() ([]models.ProjectMember, error) {
	var err error
	var result models.GetMembersResponse

	err = p.r.get("/projects/"+p.id+"/members", &result, nil)

	return result.Members, err
}

// Add members to a project
//
func (p *Project) AddMembers(memberRoles map[string]models.Role) error {
	var result models.AddMembersResponse
	var err error

	payload := models.AddMembersPayload{
		Members: make([]models.MemberRole, 0),
	}

	for memberID, role := range memberRoles {
		payload.Members = append(payload.Members, models.MemberRole{MemberID: memberID, RoleID: role.ID})
	}

	err = p.r.post("/projects/"+p.id+"/members", payload, &result, nil)

	if !result.Success && result.Error != "" {
		err = fmt.Errorf(result.Error)
	}

	return err
}

// Removes members from the project (from all environments)
func (p *Project) RemoveMembers(members []string) error {
	var result models.RemoveMembersResponse
	var err error

	payload := models.RemoveMembersPayload{
		Members: members,
	}

	err = p.r.del("/projects/"+p.id+"/members/", payload, &result, nil)

	if !result.Success && result.Error != "" {
		err = fmt.Errorf(result.Error)
	}

	return err
}

// Changes the role of a member
// memberId should have the form <username>@<github|gitlab>
func (p *Project) SetMemberRole(memberId string, role string) (err error) {
	payload := models.SetMemberRolePayload{
		MemberID: memberId,
		RoleName: role,
	}

	err = p.r.put("/projects/"+p.id+"/members/role", payload, nil, nil)

	return err
}

func (p *Project) GetAccessibleEnvironments() ([]models.Environment, error) {
	var result models.GetEnvironmentsResponse

	err := p.r.get("/projects/"+p.id+"/environments", &result, nil)

	return result.Environments, err
}

// Destroys the project, its environments, environments versions,
// project members, messages, etc. Permanently
func (p *Project) Destroy() (err error) {
	err = p.r.del("/projects/"+p.id, nil, nil, nil)

	return err
}

func (p *Project) GetProjectsOrganization() (models.Organization, error) {
	var result models.Organization

	err := p.r.get("/projects/"+p.id+"/organization", &result, nil)

	return result, err
}

func (p *Project) GetRoles() ([]models.Role, error) {
	var err error
	var result models.GetRolesResponse

	err = p.r.get("/projects/"+p.id+"roles", &result, nil)

	return result.Roles, err
}

func (p *Project) GetAll() ([]models.Project, error) {
	var err error
	var result models.GetProjectsResponse

	err = p.r.get("/projects", &result, nil)

	return result.Projects, err
}

func (p *Project) GetLogs(options *models.GetLogsOptions) ([]models.ActivityLogLite, error) {
	var err error
	var result models.GetActivityLogResponse

	err = p.r.post("/projects/"+p.id+"/activity-logs", options, &result, nil)

	return result.Logs, err
}
