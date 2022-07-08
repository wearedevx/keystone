package client

import (
	"fmt"
	"log"

	"github.com/wearedevx/keystone/api/pkg/models"
)

type Project struct {
	log *log.Logger
	id  string
	r   requester
}

// Init method creates a new project
func (p *Project) Init(
	name string,
	organizationID uint,
) (models.Project, error) {
	var project models.Project

	payload := models.Project{
		Name:           name,
		OrganizationID: organizationID,
	}

	p.log.Printf("Init with %+v\n", payload)

	err := p.r.post("/projects", payload, &project, nil)

	return project, err
}

// GetAllMembers method returns all members associated with the project
func (p *Project) GetAllMembers() ([]models.ProjectMember, error) {
	var err error
	var result models.GetMembersResponse

	err = p.r.get(
		fmt.Sprintf("/projects/%s/members", p.id),
		&result,
		nil,
	)

	return result.Members, err
}

// Add members to a project
func (p *Project) AddMembers(memberRoles map[string]models.Role) error {
	var result models.AddMembersResponse
	var err error

	payload := models.AddMembersPayload{
		Members: make([]models.MemberRole, 0),
	}

	for memberID, role := range memberRoles {
		payload.Members = append(
			payload.Members,
			models.MemberRole{MemberID: memberID, RoleID: role.ID},
		)
	}

	p.log.Printf("Adding Members %+v\n", payload)
	var path = fmt.Sprintf("/projects/%s/members", p.id)

	err = p.r.post(path, payload, &result, nil)

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

	p.log.Printf("Removing members %+v\n", payload)

	err = p.r.del(
		fmt.Sprintf("/projects/%s/members/", p.id),
		payload,
		&result,
		nil,
	)

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

	p.log.Printf("Set member role %+v\n", payload)

	err = p.r.put(
		fmt.Sprintf("/projects/%s/members/role", p.id),
		payload,
		nil,
		nil,
	)

	return err
}

// GetAccessibleEnvironments method returns the list of environments the
// current user is allowed access to.
func (p *Project) GetAccessibleEnvironments() ([]models.Environment, error) {
	var result models.GetEnvironmentsResponse

	err := p.r.get("/projects/"+p.id+"/environments", &result, nil)

	return result.Environments, err
}

// Destroys the project, its environments, environments versions,
// project members, messages, etc. Permanently
func (p *Project) Destroy() (err error) {
	p.log.Printf("Destroy project %s", p.id)

	err = p.r.del(fmt.Sprintf("/projects/%s", p.id), nil, nil, nil)

	return err
}

// GetProjectsOrganization method fetches the organziation this
// project belongs to
func (p *Project) GetProjectsOrganization() (models.Organization, error) {
	var result models.Organization

	err := p.r.get("/projects/"+p.id+"/organization", &result, nil)

	return result, err
}

// IsOrganizationPaid method returns true if the organization the project
// belongs to is paid
func (p *Project) IsOrganizationPaid() (bool, error) {
	organization, err := p.GetProjectsOrganization()

	if err != nil {
		return false, err
	}

	return organization.Paid, nil
}

// GetRoles method returns roles available for that organization
func (p *Project) GetRoles() ([]models.Role, error) {
	var err error
	var result models.GetRolesResponse

	err = p.r.get(
		fmt.Sprintf("/projects/%s/roles", p.id),
		&result, nil)

	return result.Roles, err
}

// GetAll method returns a list of all the projects the current user
// is associated with
func (p *Project) GetAll() ([]models.Project, error) {
	var err error
	var result models.GetProjectsResponse

	err = p.r.get("/projects", &result, nil)

	return result.Projects, err
}

// GetLogs method returns all the logs relative to the current project
func (p *Project) GetLogs(
	options *models.GetLogsOptions,
) ([]models.ActivityLogLite, error) {
	var err error
	var result models.GetActivityLogResponse

	p.log.Printf("Get logs with options: %+v\n", options)

	err = p.r.post(
		fmt.Sprintf("/projects/%s/activity-logs", p.id),
		options,
		&result,
		nil,
	)

	return result.Logs, err
}
