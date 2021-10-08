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
	var result = models.GetOrganizationsResponse{
		Organizations: []models.Organization{},
	}

	log := models.ActivityLog{
		UserID: &user.ID,
		Action: "GetOrganizations",
	}

	if err = Repo.GetOrganizations(user.ID, &result).Err(); err != nil {
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
	params router.Params,
	body io.ReadCloser,
	Repo repo.IRepo,
	user models.User,
) (_ router.Serde, status int, err error) {
	status = http.StatusOK
	log := models.ActivityLog{
		UserID: &user.ID,
		Action: "PostOrganization",
	}

	orga := models.Organization{}

	if err = orga.Deserialize(body); err != nil {
		status = http.StatusBadRequest
		err = apierrors.ErrorBadRequest(err)
		goto done
	}

	orga.UserID = user.ID

	if err = Repo.CreateOrganization(&orga).Err(); err != nil {
		switch {
		case errors.Is(err, repo.ErrorBadName):
			status = http.StatusBadRequest
			err = apierrors.ErrorBadOrganizationName()
		case errors.Is(err, repo.ErrorNameTaken):
			status = http.StatusConflict
			err = apierrors.ErrorOrganizationNameAlreadyTaken()
		default:
			status = http.StatusInternalServerError
			err = apierrors.ErrorFailedToCreateResource(err)
		}
	}

done:
	return &orga, status, log.SetError(err)
}

func UpdateOrganization(
	params router.Params,
	body io.ReadCloser,
	Repo repo.IRepo,
	user models.User,
) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	log := models.ActivityLog{
		UserID: &user.ID,
	}

	var isOwner bool
	orga := models.Organization{}

	if err = orga.Deserialize(body); err != nil {
		status = http.StatusBadRequest
		err = apierrors.ErrorBadRequest(err)
		goto done
	}

	isOwner, err = Repo.IsUserOwnerOfOrga(&user, &orga)

	switch {
	case err != nil:
		status = http.StatusBadRequest
		err = apierrors.ErrorBadRequest(err)
	case !isOwner:
		status = http.StatusForbidden
		err = apierrors.ErrorNotOrganizationOwner()
	default:
		if err = Repo.UpdateOrganization(&orga).Err(); err != nil {
			if errors.Is(err, repo.ErrorBadName) {
				status = http.StatusForbidden
				err = apierrors.ErrorBadOrganizationName()

			} else if errors.Is(err, repo.ErrorNameTaken) {
				status = http.StatusConflict
				err = apierrors.ErrorOrganizationNameAlreadyTaken()
			} else {
				status = http.StatusInternalServerError
				err = apierrors.ErrorFailedToUpdateResource(err)
			}
		}
	}

done:
	return &orga, status, log.SetError(err)
}

func GetOrganizationProjects(
	params router.Params,
	body io.ReadCloser,
	Repo repo.IRepo,
	user models.User,
) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	log := models.ActivityLog{
		UserID: &user.ID,
	}

	var orgaID = params.Get("orgaID").(string)

	var result = models.GetProjectsResponse{
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

func GetOrganizationMembers(params router.Params, body io.ReadCloser, Repo repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	log := models.ActivityLog{
		UserID: &user.ID,
	}

	var orgaID = params.Get("orgaID").(string)

	var result = models.GetMembersResponse{
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
