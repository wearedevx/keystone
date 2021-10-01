package controllers

import (
	"errors"
	"io"
	"net/http"

	"github.com/wearedevx/keystone/api/pkg/models"

	"github.com/wearedevx/keystone/api/internal/rights"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func PostProject(_ router.Params, body io.ReadCloser, Repo repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	project := models.Project{}

	if err = project.Deserialize(body); err != nil {
		return &project, http.StatusBadRequest, err
	}

	project.UserID = user.ID

	if err = Repo.GetOrCreateProject(&project).Err(); err != nil {
		return &project, http.StatusInternalServerError, err
	}

	project.User = user
	project.UserID = user.ID

	return &project, status, err

}

func GetProjects(params router.Params, _ io.ReadCloser, Repo repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	var result models.GetProjectsResponse

	Repo.GetUserProjects(user.ID, &result.Projects)

	if err = Repo.Err(); err != nil {
		return &result, http.StatusInternalServerError, err
	}

	return &result, status, err
}

func GetProjectsMembers(params router.Params, _ io.ReadCloser, Repo repo.IRepo, _ models.User) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	var project models.Project
	var projectID = params.Get("projectID").(string)
	var result models.GetMembersResponse

	if projectID == "" {
		return &result, http.StatusBadRequest, nil
	}

	Repo.GetProjectByUUID(projectID, &project).
		ProjectGetMembers(&project, &result.Members)

	if err = Repo.Err(); err != nil {
		return &result, http.StatusInternalServerError, err
	}

	return &result, status, err
}

func PostProjectsMembers(params router.Params, body io.ReadCloser, Repo repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	var project models.Project
	var projectID = params.Get("projectID").(string)

	input := models.AddMembersPayload{}
	result := models.AddMembersResponse{}
	err = input.Deserialize(body)

	if projectID == "" || err != nil {
		return &result, http.StatusBadRequest, nil
	}

	isPaid, err := Repo.IsProjectOrganizationPaid(projectID)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	for _, member := range input.Members {
		role := models.Role{ID: member.RoleID}
		if err := Repo.GetRole(&role).Err(); err == nil {
			if role.Name != "admin" && !isPaid {
				return nil, http.StatusInternalServerError, errors.New("You are not allowed to set role other than admin for free organization")
			}
		}
	}

	if err = Repo.GetProjectByUUID(projectID, &project).Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
		}

		return &result, status, err
	}

	members := make([]string, 0)
	for _, member := range input.Members {
		members = append(members, member.MemberID)
	}

	areInProjects, err := Repo.CheckMembersAreInProject(project, members)

	if len(areInProjects) > 0 {
		status = http.StatusConflict
		result.Error = "user already in project"

		return &result, status, err
	}

	can, err := checkUserCanAddMembers(Repo, user, project, input.Members)

	if err != nil {
		status = http.StatusInternalServerError
		result.Error = err.Error()

		return &result, status, err
	}

	if can {
		err = Repo.ProjectAddMembers(project, input.Members, user).Err()

		if err != nil {
			status = http.StatusInternalServerError
			result.Error = err.Error()
		} else {
			result.Success = true
		}
	} else {
		status = http.StatusForbidden
		result.Error = "operation not allowed"
	}

	return &result, status, err
}

func DeleteProjectsMembers(params router.Params, body io.ReadCloser, Repo repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	var project models.Project
	var projectID = params.Get("projectID").(string)
	input := models.RemoveMembersPayload{}
	result := models.RemoveMembersResponse{}
	err = input.Deserialize(body)

	if projectID == "" || err != nil {
		return &result, http.StatusBadRequest, err
	}

	if err = Repo.GetProjectByUUID(projectID, &project).Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
			result.Error = "No such project"
		} else {
			status = http.StatusInternalServerError
			result.Error = err.Error()
		}
	}

	areInProjects, err := Repo.CheckMembersAreInProject(project, input.Members)

	if len(areInProjects) != len(input.Members) {
		status = http.StatusConflict
		result.Error = "user not in project"

		return &result, status, err
	}

	can, err := checkUserCanRemoveMembers(Repo, user, project, input.Members)
	if err != nil {
		status = http.StatusInternalServerError
		result.Error = err.Error()

		return &result, status, err
	}

	if can {
		err = Repo.
			ProjectRemoveMembers(project, input.Members).
			Err()

		if err != nil {
			status = http.StatusInternalServerError
			result.Error = err.Error()
		} else {
			result.Success = true
		}
	} else {
		status = http.StatusForbidden
		result.Error = "operation not allowed"
	}

	return &result, status, err
}

// checkUserCanAddMembers checks wether a user can add all the members in `members` to `project`
// Returns false if at least one of the members cannot be added
func checkUserCanAddMembers(Repo repo.IRepo, user models.User, project models.Project, members []models.MemberRole) (can bool, err error) {
	can = true

	projectMember := models.ProjectMember{
		UserID:    user.ID,
		ProjectID: project.ID,
	}

	if err = Repo.GetProjectMember(&projectMember).Err(); err != nil {
		return false, err
	}

	for _, m := range members {
		role := models.Role{ID: m.RoleID}

		can, err = rights.CanRoleAddRole(Repo, projectMember.Role, role)

		if err != nil {
			can = false
			break
		}

		if !can {
			break
		}
	}

	return can, err
}

// checkUserCanRemoveMembers checks wether a user can remove all the members in `members` from `project`
// Returns false if at least one of the members cannot be removed
func checkUserCanRemoveMembers(Repo repo.IRepo, user models.User, project models.Project, members []string) (can bool, err error) {
	can = true

	projectMember := models.ProjectMember{
		UserID:    user.ID,
		ProjectID: project.ID,
	}

	projectMembers := []models.ProjectMember{}

	if err = Repo.
		GetProjectMember(&projectMember).
		ListProjectMembers(members, &projectMembers).
		Err(); err != nil {

		return false, err
	}

	for _, m := range projectMembers {
		can, err = rights.CanRoleAddRole(Repo, projectMember.Role, m.Role)

		if err != nil {
			can = false
			break
		}

		if !can {
			break
		}
	}

	return can, err
}

func GetAccessibleEnvironments(params router.Params, _ io.ReadCloser, Repo repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	result := models.GetEnvironmentsResponse{
		Environments: make([]models.Environment, 0),
	}
	var environments []models.Environment
	var projectID = params.Get("projectID").(string)
	var project models.Project

	Repo.GetProjectByUUID(projectID, &project)

	Repo.GetEnvironmentsByProjectUUID(project.UUID, &environments)

	for _, environment := range environments {
		can, _ := rights.CanUserWriteOnEnvironment(Repo, user.ID, project.ID, &environment)
		if can {
			result.Environments = append(result.Environments, environment)
		}
	}
	return &result, status, err
}

func DeleteProject(params router.Params, _ io.ReadCloser, Repo repo.IRepo, _ models.User) (_ router.Serde, status int, err error) {
	status = http.StatusNoContent

	var projectId = params.Get("projectID").(string)

	var project models.Project

	Repo.
		GetProjectByUUID(projectId, &project).
		DeleteAllProjectMembers(&project).
		DeleteProjectsEnvironments(&project).
		DeleteProject(&project)

	if err = Repo.Err(); err != nil {
		status = http.StatusInternalServerError
	}

	return nil, status, err
}

func GetProjectsOrganization(params router.Params, _ io.ReadCloser, Repo repo.IRepo, _ models.User) (_ router.Serde, status int, err error) {

	result := models.Organization{}
	var projectId = params.Get("projectID").(string)
	var organization models.Organization

	Repo.GetProjectsOrganization(projectId, &organization)
	if err != nil {
		return nil, status, err
	}
	result = organization

	return &result, status, err
}
