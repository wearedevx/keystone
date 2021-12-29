package controllers

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	apierrors "github.com/wearedevx/keystone/api/internal/errors"
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
	result := models.GetOrganizationsResponse{
		Organizations: []models.Organization{},
	}

	log := models.ActivityLog{
		UserID: &user.ID,
		Action: "GetOrganizations",
	}

	organizationName := params.Get("name")
	owned := params.Get("owned") == "1"

	var method func() error

	switch {
	case organizationName != "" && owned:
		method = func() error {
			return Repo.
				GetOwnedOrganizationByName(
					user.ID,
					organizationName,
					&result.Organizations,
				).
				Err()
		}

	case organizationName != "" && !owned:
		method = func() error {
			return Repo.
				GetOrganizationByName(
					user.ID,
					organizationName,
					&result.Organizations,
				).
				Err()
		}

	case organizationName == "" && owned:
		method = func() error {
			return Repo.
				GetOwnedOrganizations(
					user.ID,
					&result.Organizations,
				).
				Err()
		}

	default:
		method = func() error {
			return Repo.
				GetOrganizations(
					user.ID,
					&result.Organizations,
				).
				Err()
		}
	}

	if err = method(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
		} else {
			err = apierrors.ErrorFailedToGetResource(err)
			status = http.StatusInternalServerError
		}
	}

	return &result, status, log.SetError(err)
}

func PostOrganization(
	_ router.Params,
	body io.ReadCloser,
	Repo repo.IRepo,
	user models.User,
) (_ router.Serde, status int, err error) {
	status = http.StatusCreated
	log := models.ActivityLog{
		UserID: &user.ID,
		Action: "PostOrganization",
	}

	orga := &models.Organization{}

	if err = orga.Deserialize(body); err != nil {
		status = http.StatusBadRequest
		err = apierrors.ErrorBadRequest(err)
		orga = nil
		goto done
	}

	orga.UserID = user.ID
	orga.Paid = false

	if err = Repo.CreateOrganization(orga).Err(); err != nil {
		switch {
		case errors.Is(err, repo.ErrorBadName):
			status = http.StatusBadRequest
			err = apierrors.ErrorBadOrganizationName()
			orga = nil
		case errors.Is(err, repo.ErrorNameTaken):
			status = http.StatusConflict
			err = apierrors.ErrorOrganizationNameAlreadyTaken()
			orga = nil
		default:
			status = http.StatusInternalServerError
			err = apierrors.ErrorFailedToCreateResource(err)
			orga = nil
		}
	}

done:
	return orga, status, log.SetError(err)
}

func UpdateOrganization(
	_ router.Params,
	body io.ReadCloser,
	Repo repo.IRepo,
	user models.User,
) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	log := models.ActivityLog{
		UserID: &user.ID,
	}

	var isOwner bool
	var orga *models.Organization = &models.Organization{}
	inputOrga := models.Organization{}

	if err = inputOrga.Deserialize(body); err != nil {
		status = http.StatusBadRequest
		err = apierrors.ErrorBadRequest(err)
		orga = nil
		goto done
	}

	orga.ID = inputOrga.ID
	if err = Repo.GetOrganization(orga).Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
			orga = nil
			goto done
		}
	}

	if orga.Name != inputOrga.Name && inputOrga.Name != "" {
		orga.Name = inputOrga.Name
	}
	orga.Private = inputOrga.Private

	isOwner, err = Repo.IsUserOwnerOfOrga(&user, orga)

	switch {
	case err != nil:
		status = http.StatusBadRequest
		err = apierrors.ErrorBadRequest(err)
	case !isOwner:
		status = http.StatusForbidden
		err = apierrors.ErrorNotOrganizationOwner()
		orga = nil
	default:
		if err = Repo.UpdateOrganization(orga).Err(); err != nil {
			if errors.Is(err, repo.ErrorBadName) {
				status = http.StatusBadRequest
				err = apierrors.ErrorBadOrganizationName()
				orga = nil

			} else if errors.Is(err, repo.ErrorNameTaken) {
				status = http.StatusConflict
				err = apierrors.ErrorOrganizationNameAlreadyTaken()
				orga = nil
			} else {
				status = http.StatusInternalServerError
				err = apierrors.ErrorFailedToUpdateResource(err)
				orga = nil
			}
		}
	}

done:
	return orga, status, log.SetError(err)
}

func GetOrganizationProjects(
	params router.Params,
	_ io.ReadCloser,
	Repo repo.IRepo,
	user models.User,
) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	log := models.ActivityLog{
		UserID: &user.ID,
	}

	orgaID := params.Get("orgaID")

	result := models.GetProjectsResponse{
		Projects: []models.Project{},
	}

	u64, err := strconv.ParseUint(orgaID, 10, 0)
	orga := models.Organization{ID: uint(u64)}

	if err != nil {
		status = http.StatusBadRequest
		err = apierrors.ErrorBadRequest(err)

		goto done
	}

	if err = Repo.
		GetOrganizationProjects(&orga, &result.Projects).
		Err(); err != nil {
		status = http.StatusInternalServerError
		err = apierrors.ErrorFailedToGetResource(err)
		goto done
	}

	for index, project := range result.Projects {
		projectMember := models.ProjectMember{
			UserID:    user.ID,
			ProjectID: project.ID,
		}

		if err = Repo.
			GetProjectMember(&projectMember).
			Err(); err != nil {
			if errors.Is(err, repo.ErrorNotFound) {
				// Remove project if user is not a member
				result.Projects = append(
					result.Projects[:index],
					result.Projects[index+1:]...,
				)
			} else {
				status = http.StatusInternalServerError
				err = apierrors.ErrorFailedToGetResource(err)
			}
		}
	}

done:
	return &result, status, log.SetError(err)
}

func GetOrganizationMembers(
	params router.Params,
	_ io.ReadCloser,
	Repo repo.IRepo,
	user models.User,
) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	log := models.ActivityLog{
		UserID: &user.ID,
	}

	orgaID := params.Get("orgaID")

	result := models.GetMembersResponse{
		Members: []models.ProjectMember{},
	}

	u64, err := strconv.ParseUint(orgaID, 10, 0)
	orga := models.Organization{ID: uint(u64)}

	if err != nil {
		status = http.StatusBadRequest
		err = apierrors.ErrorBadRequest(err)

		goto done
	}

	if err = Repo.
		GetOrganizationMembers(orga.ID, &result.Members).
		Err(); err != nil {
		status = http.StatusInternalServerError
		err = apierrors.ErrorFailedToGetResource(err)

		goto done
	}

done:
	return &result, status, log.SetError(err)
}
