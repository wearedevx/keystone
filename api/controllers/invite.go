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

	status = http.StatusOK
	log := models.ActivityLog{
		UserID: user.ID,
		Action: "PostInvite",
	}

	// check if user exist
	if err = Repo.GetUserByEmail(targetEmail, &targetUsers).Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {

			e, err := emailer.InviteMail(user, payload.ProjectName)
			if err != nil {
				status = http.StatusInternalServerError
				goto done
			}

			if err = e.Send([]string{targetEmail}); err != nil {
				status = http.StatusInternalServerError
				goto done
			}
		}

	} else {
		for _, user := range targetUsers {
			result.UserUIDs = append(result.UserUIDs, user.UserID)
		}
	}

done:
	return result, status, log.SetError(err)
}
