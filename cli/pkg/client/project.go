package client

import (
	"fmt"

	"github.com/wearedevx/keystone/api/pkg/models"
)

type Project struct {
	id string
	r  requester
}

func (p *Project) Init(name string) (models.Project, error) {
	var project models.Project

	payload := models.Project{
		Name: name,
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
