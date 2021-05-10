package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	. "github.com/wearedevx/keystone/api/internal/utils"
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

func PostProject(_ router.Params, body io.ReadCloser, Repo repo.Repo, user User) (router.Serde, int, error) {
	var status int = http.StatusOK
	var err error

	project := &Project{
		User:   user,
		UserID: user.ID,
	}

	runner := NewRunner([]RunnerAction{
		NewAction(func() error {
			return project.Deserialize(body)
		}),
		NewAction(func() error {
			Repo.GetOrCreateProject(project, user)
			return Repo.Err()
		}).SetStatusSuccess(201),
	})

	if err = runner.Run().Error(); err != nil {
		return project, status, err
	}

	status = runner.Status()
	err = runner.Error()

	return project, status, err

}

func GetProjectsPublicKeys(params router.Params, _ io.ReadCloser, Repo repo.Repo, user User) (router.Serde, int, error) {
	var status int = http.StatusOK
	var err error

	var project Project
	var projectID = params.Get("projectID").(string)
	var result projectsPublicKeys

	Repo.GetProjectByUUID(projectID, &project).
		ProjectLoadUsers(&project)

	if Repo.Err() != nil && Repo.Err() == repo.ErrorNotFound {
		return &result, 404, nil
	}

	if Repo.Err() != nil {
		return &result, 500, nil
	}

	for _, member := range project.Members {
		result.keys = append(result.keys, UserPublicKey{
			UserID:    member.User.UserID,
			PublicKey: member.User.PublicKey,
		})
	}

	return &result, status, err
}

func GetProjectsMembers(params router.Params, _ io.ReadCloser, Repo repo.Repo, user User) (router.Serde, int, error) {
	var status int = http.StatusOK
	var err error

	var project Project
	var projectID = params.Get("projectID").(string)
	var result GetMembersResponse

	runner := NewRunner([]RunnerAction{
		NewAction(func() error {
			Repo.GetProjectByUUID(projectID, &project)
			Repo.ProjectGetMembers(&project, &result.Members)

			return Repo.Err()
		}).SetStatusSuccess(200),
	}).Run()

	status = runner.Status()
	err = runner.Error()

	return &result, status, err
}

func PostProjectsMembers(params router.Params, body io.ReadCloser, Repo repo.Repo, user User) (router.Serde, int, error) {
	var status int = http.StatusOK
	var err error

	var project Project
	var projectID = params.Get("projectID").(string)
	input := models.AddMembersPayload{}
	result := models.AddMembersResponse{Success: true, Error: ""}
	err = input.Deserialize(body)

	if err != nil {
		return &result, 500, nil
	}

	Repo.GetProjectByUUID(projectID, &project).
		ProjectLoadUsers(&project)

	if Repo.Err() != nil && Repo.Err() == repo.ErrorNotFound {
		return &result, 404, nil
	}

	if Repo.Err() != nil {
		return &result, 500, nil
	}

	canInvite, err := rights.CanUserInviteUsers(Repo, &user, &project, input.Members)

	if err != nil {
		return &result, 500, nil
	}

	if !canInvite {
		return &result, 401, nil
	}

	Repo.ProjectAddMembers(project, input.Members)

	if Repo.Err() != nil {
		return &result, 500, nil
	}

	return &result, status, err
}

func DeleteProjectsMembers(params router.Params, body io.ReadCloser, Repo repo.Repo, user User) (router.Serde, int, error) {
	var status int = http.StatusOK
	var err error

	var project Project
	var projectID = params.Get("projectID").(string)
	input := models.RemoveMembersPayload{}
	result := models.RemoveMembersResponse{}
	err = input.Deserialize(body)

	if err != nil {
		return &result, 500, err
	}

	runner := NewRunner([]RunnerAction{
		NewAction(func() error {
			Repo.GetProjectByUUID(projectID, &project)

			return Repo.Err()
		}).
			SetStatusError(404),
		NewAction(func() error {
			Repo.ProjectRemoveMembers(project, input.Members)

			return Repo.Err()
		}).
			SetStatusSuccess(204).
			SetStatusError(500),
	}).Run()

	status = runner.Status()
	err = runner.Error()

	if err != nil {
		result.Success = false
		result.Error = err.Error()
	}

	return &result, status, err
}
