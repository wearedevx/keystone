package controllers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/wearedevx/keystone/api/internal/rights"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
	"gorm.io/gorm"
)

// DoUsersExist checks if users exists in the Keystone database
// This is not project dependant, it checks all users in the whole world
func DoUsersExist(_ router.Params, body io.ReadCloser, Repo repo.IRepo, _ models.User) (_ router.Serde, status int, err error) {
	status = http.StatusBadRequest
	response := &models.CheckMembersResponse{}
	payload := &models.CheckMembersPayload{}

	// input check
	err = payload.Deserialize(body)

	if err != nil {
		return response, status, err
	}

	// actual work

	users := make(map[string]models.User)
	notFounds := make([]string, 0)

	Repo.FindUsers(payload.MemberIDs, &users, &notFounds)

	if Repo.Err() != nil {
		return response, http.StatusInternalServerError, Repo.Err()
	}

	if len(notFounds) != 0 {
		response.Error = fmt.Sprintf("%s do not exists", strings.Join(notFounds, ", "))
		response.Success = false

		status = http.StatusNotFound
	}

	return response, status, err
}

// PutMembersSetRole sets the role for a given project member
func PutMembersSetRole(params router.Params, body io.ReadCloser, Repo repo.IRepo, user models.User) (response router.Serde, status int, err error) {
	status = http.StatusOK
	payload := &models.SetMemberRolePayload{}
	project := models.Project{}

	// input check
	var projectID = params.Get("projectID").(string)

	err = payload.Deserialize(body)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	// actual work
	member := models.User{
		UserID: payload.MemberID,
	}

	role := models.Role{Name: payload.RoleName}

	if err = Repo.GetProjectByUUID(projectID, &project).
		GetUser(&member).
		GetRole(&role).
		Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}

		return response, status, err
	}

	isPaid, err := Repo.IsProjectOrganizationPaid(projectID)

	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	if role.Name != "admin" && !isPaid {
		return response, http.StatusForbidden, errors.New("This feature is not allowed for free organization")
	}
	can, err := rights.CanUserSetMemberRole(Repo, user, member, role, project)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	if !can {
		return response, http.StatusForbidden, errors.New("operation not allowed")
	}

	if err = Repo.ProjectSetRoleForUser(project, member, role).Err(); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
	}

	return response, status, err
}

func checkUserCanChangeMember(Repo repo.IRepo, user models.User, project models.Project, other models.User) (can bool, err error) {
	projectMember := models.ProjectMember{
		UserID:    user.ID,
		ProjectID: project.ID,
	}
	otherProjectMember := models.ProjectMember{
		UserID:    other.ID,
		ProjectID: project.ID,
	}

	if err = Repo.
		GetProjectMember(&projectMember).
		GetProjectMember(&otherProjectMember).
		Err(); err != nil {

		return false, err
	}

	return rights.CanRoleAddRole(Repo, projectMember.Role, otherProjectMember.Role)
}
