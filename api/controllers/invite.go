package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/wearedevx/keystone/api/internal/emailer"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

type InvitePayload struct {
	Email       string
	ProjectName string
}

func (pm *InvitePayload) Deserialize(in io.Reader) error {
	return json.NewDecoder(in).Decode(pm)
}

func (pm *InvitePayload) Serialize(out *string) (err error) {
	var sb strings.Builder

	err = json.NewEncoder(&sb).Encode(pm)

	*out = sb.String()

	return err
}

func PostInvite(_ router.Params, body io.ReadCloser, _ repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	payload := InvitePayload{}
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
