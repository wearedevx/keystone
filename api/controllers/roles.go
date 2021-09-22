package controllers

import (
	"errors"
	"io"
	"net/http"

	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

// Returns a List of Roles
func GetRoles(params router.Params, _ io.ReadCloser, Repo repo.IRepo, user models.User) (_ router.Serde, status int, err error) {
	status = http.StatusOK
	log := models.ActivityLog{
		UserID: &user.ID,
		Action: "GetRoles",
	}

	var result = models.GetRolesResponse{
		Roles: []models.Role{},
	}

	projectID := params.Get("projectID").(string)

	if projectID == "" {
		if err = Repo.GetRoles(&result.Roles).Err(); err != nil {
			if errors.Is(err, repo.ErrorNotFound) {
				status = http.StatusNotFound
			} else {
				status = http.StatusInternalServerError
			}

			goto done
		}
	} else {
		isPaid, err := Repo.IsProjectOrganizationPaid(projectID)
		if err != nil {
			return &result, status, err
		}

		if isPaid {
			if err = Repo.GetRoles(&result.Roles).Err(); err != nil {
				status = http.StatusInternalServerError
				return &result, status, err
			}
		} else {
			adminRole := models.Role{Name: "admin"}

			if err = Repo.GetRole(&adminRole).Err(); err != nil {
				status = http.StatusInternalServerError
				goto done
			}

			result.Roles = []models.Role{adminRole}
		}
	}

done:
	return &result, status, log.SetError(err)
}
