package controllers

import (
	"errors"
	"io"
	"net/http"

	apierrors "github.com/wearedevx/keystone/api/internal/errors"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/notification"
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

	targetUsers := []models.User{}
	result := &models.GetInviteResponse{}

	status = http.StatusOK
	log := models.ActivityLog{
		UserID: &user.ID,
		Action: "PostInvite",
	}

	var project = models.Project{Name: payload.ProjectName}
	if err = Repo.GetProject(&project).Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
			err = apierrors.ErrorFailedToGetResource(err)
			goto done
		}

		status = http.StatusInternalServerError
		err = apierrors.ErrorFailedToGetResource(err)
		goto done
	}

	// check if user exist
	if err = Repo.GetUserByEmail(payload.Email, &targetUsers).Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			err = notification.SendInvitationEmail(user, payload)
			if err != nil {
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
