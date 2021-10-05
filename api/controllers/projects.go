package controllers

import (
	"errors"
	"io"
	"net/http"

	"github.com/wearedevx/keystone/api/pkg/models"

	"github.com/wearedevx/keystone/api/internal/payment"
	"github.com/wearedevx/keystone/api/internal/rights"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func PostProject(_ router.Params, body io.ReadCloser, Repo repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	status = http.StatusOK
	log := models.ActivityLog{
		UserID: &user.ID,
		Action: "PostProject",
	}

	project := models.Project{}
	orga := models.Organization{}

	if err = project.Deserialize(body); err != nil {
		status = http.StatusBadRequest
		goto done
	}

	project.UserID = user.ID

	if err = Repo.GetOrCreateProject(&project).Err(); err != nil {
		status = http.StatusInternalServerError
		goto done
	}

	// Add organization's owner to project as admin
	if err := Repo.GetProjectsOrganization(project.UUID, &orga).Err(); err != nil {
		return &project, http.StatusInternalServerError, err
	}

	if user.ID != orga.UserID {
		role := models.Role{
			Name: "admin",
		}

		if err := Repo.GetRole(&role).Err(); err != nil {
			return &project, http.StatusInternalServerError, err
		}

		orgaOwner := models.ProjectMember{
			ProjectID: project.ID,
			UserID:    orga.UserID,
			RoleID:    role.ID,
		}

		Repo.GetDb().Save(&orgaOwner)
	}

	project.User = user
	project.UserID = user.ID

	if project.ID != 0 {
		log.ProjectID = &project.ID
	}

done:
	return &project, status, log.SetError(err)
}

func GetProjects(
	params router.Params,
	_ io.ReadCloser,
	Repo repo.IRepo,
	user models.User,
) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	var result models.GetProjectsResponse

	Repo.GetUserProjects(user.ID, &result.Projects)

	if err = Repo.Err(); err != nil {
		return &result, http.StatusInternalServerError, err
	}

	return &result, status, err
}

func GetProjectsMembers(
	params router.Params,
	_ io.ReadCloser,
	Repo repo.IRepo,
	user models.User,
) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	var project models.Project
	var member models.ProjectMember
	var projectID = params.Get("projectID").(string)
	var result models.GetMembersResponse
	var log = models.ActivityLog{
		UserID: &user.ID,
		Action: "GetProjectMembers",
	}

	if projectID == "" {
		status = http.StatusBadRequest
		goto done
	}

	member.UserID = user.ID

	Repo.GetProjectByUUID(projectID, &project).
		IsMemberOfProject(&project, &member).
		ProjectGetMembers(&project, &result.Members)

	log.ProjectID = &project.ID

	if err = Repo.Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
	}

done:
	return &result, status, log.SetError(err)
}

func PostProjectsMembers(
	params router.Params,
	body io.ReadCloser,
	Repo repo.IRepo,
	user models.User,
) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	var can bool
	var isPaid bool
	var areInProjects []string
	var project models.Project
	var projectID = params.Get("projectID").(string)

	members := make([]string, 0)
	input := models.AddMembersPayload{}
	result := models.AddMembersResponse{}
	organization := models.Organization{}
	log := models.ActivityLog{
		UserID: &user.ID,
		Action: "PostProjectMembers",
	}

	err = input.Deserialize(body)

	if projectID == "" || err != nil {
		status = http.StatusBadRequest
		err = nil
		goto done
	}

	if err = Repo.
		GetProjectsOrganization(projectID, &organization).
		Err(); err != nil {
		status = http.StatusBadRequest
		goto done
	}

	isPaid = organization.Paid

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

		goto done
	}

	log.ProjectID = &project.ID

	for _, member := range input.Members {
		members = append(members, member.MemberID)
	}

	areInProjects, err = Repo.CheckMembersAreInProject(project, members)

	if len(areInProjects) > 0 {
		status = http.StatusConflict
		result.Error = "user already in project"

		goto done
	}

	can, err = checkUserCanAddMembers(Repo, user, project, input.Members)

	if err != nil {
		status = http.StatusInternalServerError
		result.Error = err.Error()

		goto done
	}

	if can {
		var seats int64
		if err = Repo.
			ProjectAddMembers(project, input.Members, user).
			OrganizationCountMembers(&organization, &seats).
			Err(); err != nil {
			status = http.StatusInternalServerError
			result.Error = err.Error()
		} else {
			result.Success = true

			p := payment.NewStripePayment()
			err = p.UpdateSubscription(
				payment.SubscriptionID(organization.SubscriptionID),
				seats,
			)
		}
	} else {
		status = http.StatusForbidden
		result.Error = "operation not allowed"
	}

done:
	return &result, status, log.SetError(err)
}

func DeleteProjectsMembers(params router.Params, body io.ReadCloser, Repo repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	status = http.StatusOK
	log := models.ActivityLog{
		UserID: &user.ID,
		Action: "DeleteProjectsMembers",
	}

	var project models.Project
	var organization models.Organization
	var projectID = params.Get("projectID").(string)
	input := models.RemoveMembersPayload{}
	result := models.RemoveMembersResponse{}
	var can, userIsAdmin bool
	var areInProjects []string
	var seats int64

	err = input.Deserialize(body)

	if projectID == "" || err != nil {
		status = http.StatusBadRequest
		goto done
	}

	if err = Repo.
		GetProjectByUUID(projectID, &project).
		GetProjectsOrganization(projectID, &organization).
		OrganizationCountMembers(&organization, &seats).
		Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
			result.Error = "No such project"
		} else {
			status = http.StatusInternalServerError
			result.Error = err.Error()
		}
		goto done
	}

	log.ProjectID = &project.ID

	// Prevent users that are not admin on the project from deleting it
	userIsAdmin = Repo.ProjectIsMemberAdmin(&project, &models.ProjectMember{UserID: user.ID})
	if err = Repo.Err(); err != nil || !userIsAdmin {
		status = http.StatusNotFound
		goto done
	}

	areInProjects, err = Repo.CheckMembersAreInProject(project, input.Members)

	if len(areInProjects) != len(input.Members) {
		status = http.StatusConflict
		result.Error = "user not in project"

		goto done
	}

	can, err = checkUserCanRemoveMembers(Repo, user, project, input.Members)
	if err != nil {
		status = http.StatusInternalServerError
		result.Error = err.Error()

		goto done
	}

	if can {
		if err = Repo.
			ProjectRemoveMembers(project, input.Members).
			Err(); err != nil {
			status = http.StatusInternalServerError
			result.Error = err.Error()
		} else {
			result.Success = true

			p := payment.NewStripePayment()
			err = p.UpdateSubscription(
				payment.SubscriptionID(organization.SubscriptionID),
				seats,
			)
		}
	} else {
		status = http.StatusForbidden
		result.Error = "operation not allowed"
	}

done:
	return &result, status, log.SetError(err)
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

		isMemberOwnerOfOrga, isMemberOwnerOfOrgaErr := rights.IsUserOwnerOfOrga(Repo, m.UserID, project)
		if isMemberOwnerOfOrgaErr != nil {
			can = false
			break
		}
		if isMemberOwnerOfOrga {
			can = false
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
	var can bool
	var log = models.ActivityLog{
		UserID: &user.ID,
		Action: "GetAccessibleEnvironments",
	}

	if err = Repo.GetProjectByUUID(projectID, &project).
		GetEnvironmentsByProjectUUID(project.UUID, &environments).
		Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
		goto done
	}

	if project.ID != 0 {
		log.ProjectID = &project.ID
	}

	for _, environment := range environments {
		log.Environment = environment

		can, err = rights.CanUserWriteOnEnvironment(Repo, user.ID, project.ID, &environment)
		if err != nil {
			status = http.StatusNotFound
			goto done
		}

		if can {
			result.Environments = append(result.Environments, environment)
		}
	}

done:
	return &result, status, log.SetError(err)
}

func DeleteProject(params router.Params, _ io.ReadCloser, Repo repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	status = http.StatusNoContent
	log := models.ActivityLog{
		UserID: &user.ID,
		Action: "DeleteProject",
	}

	var projectId = params.Get("projectID").(string)
	var project models.Project

	Repo.
		GetProjectByUUID(projectId, &project).
		DeleteAllProjectMembers(&project).
		DeleteProjectsEnvironments(&project).
		DeleteProject(&project)

	log.ProjectID = &project.ID

	if err = Repo.Err(); err != nil {
		status = http.StatusInternalServerError
	}

	return nil, status, log.SetError(err)
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
