package controllers

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func GetOrganizations(
	params router.Params,
	_ io.ReadCloser,
	Repo repo.IRepo,
	user models.User,
) (_ router.Serde, status int, err error) {
	status = http.StatusOK
	var result = models.GetOrganizationsResponse{
		Organizations: []models.Organization{},
	}

	if err = Repo.GetOrganizations(user.ID, &result).Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
	}

	return &result, status, err
}

func PostOrganization(
	params router.Params,
	body io.ReadCloser,
	Repo repo.IRepo,
	user models.User,
) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	orga := models.Organization{}

	if err = orga.Deserialize(body); err != nil {
		status = http.StatusBadRequest
		goto done
	}

	orga.OwnerID = user.ID

	if err = Repo.CreateOrganization(&orga).Err(); err != nil {
		if strings.Contains(err.Error(), "Incorrect organization name") {
			status = http.StatusForbidden
		} else if strings.Contains(err.Error(), "Organization name already taken") {
			status = http.StatusConflict
		} else {
			status = http.StatusInternalServerError
		}
	}

done:
	return &orga, status, err
}

func UpdateOrganization(
	params router.Params,
	body io.ReadCloser,
	Repo repo.IRepo,
	user models.User,
) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	var isOwner bool
	orga := models.Organization{}

	if err = orga.Deserialize(body); err != nil {
		status = http.StatusBadRequest
		goto done
	}

	isOwner, err = Repo.IsUserOwnerOfOrga(&user, &orga)

	switch {
	case err != nil:
		status = http.StatusBadRequest
	case !isOwner:
		status = http.StatusForbidden
		err = errors.New("You are not the organization's owner")
	default:
		if err = Repo.UpdateOrganization(&orga).Err(); err != nil {
			if strings.Contains(err.Error(), "Incorrect organization name") {
				status = http.StatusForbidden
			} else if strings.Contains(err.Error(), "Organization name already taken") {
				status = http.StatusConflict
			} else {
				status = http.StatusInternalServerError
			}
		}
	}

done:
	return &orga, status, err
}
