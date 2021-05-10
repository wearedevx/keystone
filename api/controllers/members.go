package controllers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/wearedevx/keystone/api/internal/router"
	. "github.com/wearedevx/keystone/api/internal/utils"
	. "github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
	"gorm.io/gorm"
)

func DoUsersExist(params router.Params, body io.ReadCloser, Repo repo.Repo, user User) (router.Serde, int, error) {
	var err error
	status := 500
	response := &CheckMembersResponse{}

	payload := &CheckMembersPayload{}

	runner := NewRunner([]RunnerAction{
		NewAction(func() error {
			return payload.Deserialize(body)
		}),
		NewAction(func() error {
			_, notFound := Repo.FindUsers(payload.MemberIDs)

			if len(notFound) != 0 {
				response.Error = fmt.Sprintf("%s do not exists", strings.Join(notFound, ", "))
				response.Success = false

				status = 404
			}

			return Repo.Err()
		}).SetStatusSuccess(200),
	})

	err = runner.Run().Error()
	return response, status, err
}

func PutMembersSetRole(params router.Params, body io.ReadCloser, Repo repo.Repo, user User) (response router.Serde, status int, err error) {
	status = http.StatusOK
	payload := &SetMemberRolePayload{}
	project := Project{}
	role := Role{}

	var projectID = params.Get("projectID").(string)

	err = payload.Deserialize(body)
	if err != nil {
		return response, http.StatusInternalServerError, err
	}

	member := User{
		UserID: payload.MemberID,
	}

	Repo.GetProjectByUUID(projectID, &project).
		GetUser(&member).
		GetRoleByName(payload.RoleName, &role).
		ProjectSetRoleForUser(project, member, role)

	err = Repo.Err()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
	}

	return response, status, err
}
