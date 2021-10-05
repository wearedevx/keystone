package controllers

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
	"gorm.io/gorm"
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

	orga.UserID = user.ID

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

func GetOrganizationProjects(params router.Params, body io.ReadCloser, Repo repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	var orgaID = params.Get("orgaID").(string)

	var result = models.GetProjectsResponse{
		Projects: []models.Project{},
	}

	u64, err := strconv.ParseUint(orgaID, 10, 0)
	orga := models.Organization{ID: uint(u64)}

	if err != nil {
		status = http.StatusInternalServerError
		return &result, status, err
	}

	if err = Repo.GetOrganizationProjects(&orga, &result.Projects).Err(); err != nil {
		status = http.StatusInternalServerError
		return &result, status, err
	}

	for index, project := range result.Projects {
		projectMember := models.ProjectMember{
			UserID:    user.ID,
			ProjectID: project.ID,
		}

		err := Repo.GetProjectMember(&projectMember).Err()

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Remove project if user is not a member
				result.Projects = append(result.Projects[:index], result.Projects[index+1:]...)
			} else {
				status = http.StatusInternalServerError
				return &result, status, err
			}
		}
	}

	return &result, status, err
}

func GetOrganizationMembers(params router.Params, body io.ReadCloser, Repo repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	status = http.StatusOK

	var orgaID = params.Get("orgaID").(string)

	var result = models.GetMembersResponse{
		Members: []models.ProjectMember{},
	}

	u64, err := strconv.ParseUint(orgaID, 10, 0)
	orga := models.Organization{ID: uint(u64)}

	if err != nil {
		status = http.StatusInternalServerError
		return &result, status, err
	}

	if err = Repo.GetOrganizationMembers(orga.ID, &result.Members).Err(); err != nil {
		status = http.StatusInternalServerError
		return &result, status, err
	}

	return &result, status, err
}
