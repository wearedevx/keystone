package controllers

import (
	"errors"
	"io"
	"net/http"

	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func GetOrganizations(params router.Params, _ io.ReadCloser, Repo repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
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

		return &result, status, err
	}

	return &result, status, err
}

func PostOrganization(params router.Params, body io.ReadCloser, Repo repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	orga := models.Organization{}

	if err = orga.Deserialize(body); err != nil {
		return &orga, http.StatusBadRequest, err
	}

	orga.OwnerID = user.ID

	if err = Repo.GetDb().Create(&orga).Error; err != nil {
		return &orga, http.StatusInternalServerError, err
	}

	return &orga, status, err
}
