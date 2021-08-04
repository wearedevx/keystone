package controllers

import (
	"io"
	"net/http"

	"github.com/wearedevx/keystone/api/internal/emailer"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func PostInvite(_ router.Params, body io.ReadCloser, _ repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	payload := models.InvitePayload{}
	payload.Deserialize(body)

	senderEmail := user.Email
	targetEmail := payload.Email

	e, err := emailer.InviteMail(senderEmail, payload.ProjectName)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if err = e.Send([]string{targetEmail}); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil
}
