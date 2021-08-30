package controllers

import (
	"errors"
	"io"
	"net/http"

	"github.com/wearedevx/keystone/api/internal/emailer"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func PostInvite(
	_ router.Params,
	body io.ReadCloser,
	Repo repo.IRepo,
	user models.User,
) (_ router.Serde, status int, err error) {
	payload := models.InvitePayload{}
	if err = payload.Deserialize(body); err != nil {
		return nil, http.StatusBadRequest, err
	}

	targetEmail := payload.Email

	targetUsers := []models.User{}
	result := &models.GetInviteResponse{}

	// check if user exist
	if err = Repo.GetUserByEmail(targetEmail, &targetUsers).Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {

			e, err := emailer.InviteMail(user, payload.ProjectName)
			if err != nil {
				return result, http.StatusInternalServerError, err
			}

			if err = e.Send([]string{targetEmail}); err != nil {
				return result, http.StatusInternalServerError, err
			}
		}

	} else {

		for _, user := range targetUsers {
			result.UserUIDs = append(result.UserUIDs, user.UserID)
		}
	}

	return result, http.StatusOK, nil
}
