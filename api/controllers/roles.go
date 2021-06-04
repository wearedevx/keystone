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
	var result = models.GetRolesResponse{
		Roles: []models.Role{},
	}

	projectID := params.Get("add-to-project").(string)

	if projectID == "" {
		if err = Repo.GetRoles(&result.Roles).Err(); err != nil {
			if errors.Is(err, repo.ErrorNotFound) {
				status = http.StatusNotFound
			} else {
				status = http.StatusInternalServerError
			}

			return &result, status, err
		}
	} else {
		project := models.Project{UUID: projectID}
		projectMember := models.ProjectMember{
			UserID:    user.ID,
			ProjectID: project.ID,
		}
		roles := []models.Role{}

		if err = Repo.
			GetProject(&project).
			GetProjectMember(&projectMember).
			GetRolesMemberCanInvite(projectMember, &roles).
			Err(); err != nil {
			if errors.Is(err, repo.ErrorNotFound) {
				status = http.StatusNotFound
			}

			return &result, status, err
		}

	}

	return &result, status, err
}
