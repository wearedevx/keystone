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
func DoUsersExist(_ router.Params, body io.ReadCloser, Repo repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	status = http.StatusOK
	response := &models.CheckMembersResponse{}
	payload := &models.CheckMembersPayload{}
	log := models.ActivityLog{
		UserID: user.ID,
		Action: "DoUsersExist",
	}

	// input check
	err = payload.Deserialize(body)

	users := make(map[string]models.User)
	notFounds := make([]string, 0)

	if err != nil {
		goto done
	}

	// actual work
	Repo.FindUsers(payload.MemberIDs, &users, &notFounds)

	if err = Repo.Err(); err != nil {
		status = http.StatusInternalServerError
		goto done
	}

	if len(notFounds) != 0 {
		response.Error = fmt.Sprintf("%s do not exists", strings.Join(notFounds, ", "))
		response.Success = false

		status = http.StatusNotFound
	}

done:
	return response, status, log.SetError(err)
}

// PutMembersSetRole sets the role for a given project member
func PutMembersSetRole(params router.Params, body io.ReadCloser, Repo repo.IRepo, user models.User) (response router.Serde, status int, err error) {
	status = http.StatusOK
	payload := &models.SetMemberRolePayload{}
	project := models.Project{}
	member := models.User{}
	role := models.Role{}
	can := false

	// input check
	var projectID = params.Get("projectID").(string)

	log := models.ActivityLog{
		UserID: user.ID,
		Action: "PutMemberSetRole",
	}

	err = payload.Deserialize(body)
	if err != nil {
		status = http.StatusInternalServerError
		goto done
	}

	// actual work
	member.UserID = payload.MemberID
	role.Name = payload.RoleName

	if err = Repo.GetProjectByUUID(projectID, &project).
		GetUser(&member).
		GetRole(&role).
		Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}

		goto done
	}

	log.ProjectID = project.ID

	can, err = rights.CanUserSetMemberRole(Repo, user, member, role, project)
	if err != nil {
		status = http.StatusInternalServerError
		goto done
	}

	if !can {
		status = http.StatusForbidden
		goto done
	}

	if err = Repo.ProjectSetRoleForUser(project, member, role).Err(); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
	}

done:
	return response, status, log.SetError(err)
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
