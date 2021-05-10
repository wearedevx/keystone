package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/wearedevx/keystone/api/pkg/models"
	. "github.com/wearedevx/keystone/api/pkg/models"

	"github.com/wearedevx/keystone/api/internal/rights"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

type projectsPublicKeys struct {
	keys []UserPublicKey
}

func (p *projectsPublicKeys) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(p)
}

func (p *projectsPublicKeys) Serialize(out *string) error {
	var sb strings.Builder
	var err error

	err = json.NewEncoder(&sb).Encode(p)

	*out = sb.String()

	return err
}

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

func GetProjectsPublicKeys(params router.Params, _ io.ReadCloser, Repo repo.Repo, user User) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	var project Project
	var projectID string = params.Get("projectID").(string)
	var result projectsPublicKeys

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
		result.keys = append(result.keys, UserPublicKey{
			UserID:    member.User.UserID,
			PublicKey: member.User.PublicKey,
		})
	}

	return &result, status, err
}

func GetProjectsMembers(params router.Params, _ io.ReadCloser, Repo repo.Repo, user User) (_ router.Serde, status int, err error) {
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

	input := models.AddMembersPayload{}
	result := models.AddMembersResponse{Success: true, Error: ""}
	err = input.Deserialize(body)

	if projectID == "" || err != nil {
		return &result, http.StatusBadRequest, nil
	}

	// Check if user can invite
	if err = Repo.GetProjectByUUID(projectID, &project).Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
		}

		return &result, status, err
	}

	// Need to change parameter to models.AddMembersPayload type
	can, err := rights.CanUserInviteRole(&Repo, &user, &project, &Role{ID: input.Members[0].RoleID})

	if can {
		err = Repo.ProjectAddMembers(project, input.Members).Err()
	}

	if err != nil {
		status = http.StatusInternalServerError
		result.Success = false
		result.Error = err.Error()
	}

	return &result, status, err
}

func DeleteProjectsMembers(params router.Params, body io.ReadCloser, Repo repo.Repo, user User) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	var project Project
	var projectID = params.Get("projectID").(string)
	input := models.RemoveMembersPayload{}
	result := models.RemoveMembersResponse{}
	err = input.Deserialize(body)

	if projectID == "" || err != nil {
		return &result, http.StatusBadRequest, err
	}

	err = Repo.
		GetProjectByUUID(projectID, &project).
		ProjectRemoveMembers(project, input.Members).
		Err()

	if err != nil {
		status = http.StatusInternalServerError
		result.Success = false
		result.Error = err.Error()
	}

	return &result, status, err
}
