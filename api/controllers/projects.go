package controllers

import (
	"errors"
	"io"
	"net/http"

	. "github.com/wearedevx/keystone/api/pkg/models"

	"github.com/wearedevx/keystone/api/internal/rights"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

// type PublicKeys struct {
// 	Keys []UserPublicKey
// }

// func (p *PublicKeys) Deserialize(in io.Reader) error {
// 	return json.NewDecoder(in).Decode(p)
// }

// func (p *PublicKeys) Serialize(out *string) (err error) {
// 	var sb strings.Builder

// 	err = json.NewEncoder(&sb).Encode(p)

// 	*out = sb.String()

// 	return err
// }

func PostProject(_ router.Params, body io.ReadCloser, Repo repo.Repo, user User) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	project := &Project{
		User:   user,
		UserID: user.ID,
	}

	if err = project.Deserialize(body); err != nil {
		return project, http.StatusBadRequest, err
	}

	if err = Repo.GetOrCreateProject(project).Err(); err != nil {
		return project, http.StatusInternalServerError, err
	}

	return project, status, err

}

func GetProjectsPublicKeys(params router.Params, _ io.ReadCloser, Repo repo.Repo, _ User) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	var project Project
	var projectID string = params.Get("projectID").(string)
	var result PublicKeys

	if projectID == "" {
		return &result, http.StatusBadRequest, nil
	}

	Repo.
		GetProjectByUUID(projectID, &project).
		ProjectLoadUsers(&project)

	if err = Repo.Err(); err != nil {
		return &result, http.StatusInternalServerError, err
	}

	for _, member := range project.Members {
		result.Keys = append(result.Keys, UserPublicKey{
			UserID:    member.User.UserID,
			PublicKey: member.User.PublicKey,
		})
	}

	return &result, status, err
}

func GetProjectsMembers(params router.Params, _ io.ReadCloser, Repo repo.Repo, _ User) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	var project Project
	var projectID = params.Get("projectID").(string)
	var result GetMembersResponse

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

func PostProjectsMembers(params router.Params, body io.ReadCloser, Repo repo.Repo, user User) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	var project Project
	var projectID = params.Get("projectID").(string)

	input := AddMembersPayload{}
	result := AddMembersResponse{}
	err = input.Deserialize(body)

	if projectID == "" || err != nil {
		return &result, http.StatusBadRequest, nil
	}

	if err = Repo.GetProjectByUUID(projectID, &project).Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
		}

		return &result, status, err
	}

	can, err := checkUserCanAddMembers(&Repo, user, project, input.Members)
	if err != nil {
		status = http.StatusInternalServerError
		result.Error = err.Error()

		return &result, status, err
	}

	if can {
		err = Repo.ProjectAddMembers(project, input.Members).Err()

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

func DeleteProjectsMembers(params router.Params, body io.ReadCloser, Repo repo.Repo, user User) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	var project Project
	var projectID = params.Get("projectID").(string)
	input := RemoveMembersPayload{}
	result := RemoveMembersResponse{}
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

	can, err := checkUserCanRemoveMembers(&Repo, user, project, input.Members)
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
func checkUserCanAddMembers(Repo *repo.Repo, user User, project Project, members []MemberRole) (can bool, err error) {
	can = true

	projectMember := ProjectMember{
		UserID:    user.ID,
		ProjectID: project.ID,
	}

	if err = Repo.GetProjectMember(&projectMember).Err(); err != nil {
		return false, err
	}

	for _, m := range members {
		role := Role{ID: m.RoleID}

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
func checkUserCanRemoveMembers(Repo repo.IRepo, user User, project Project, members []string) (can bool, err error) {
	can = true

	projectMember := ProjectMember{
		UserID:    user.ID,
		ProjectID: project.ID,
	}

	projectMembers := []ProjectMember{}

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
