package controllers

import (
	"fmt"
	"io"
	"strings"

	"github.com/wearedevx/keystone/functions/ksapi/routes"
	. "github.com/wearedevx/keystone/internal/models"
	"github.com/wearedevx/keystone/internal/repo"
	. "github.com/wearedevx/keystone/internal/utils"
)

func DoUsersExist(params routes.Params, body io.ReadCloser, Repo repo.Repo, user User) (routes.Serde, int, error) {
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
