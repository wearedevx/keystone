package controllers

import (
	"errors"
	"io"
	"net/http"

	apierrors "github.com/wearedevx/keystone/api/internal/errors"
	"github.com/wearedevx/keystone/api/internal/router"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

func GetActivityLogs(
	params router.Params,
	body io.ReadCloser,
	Repo repo.IRepo,
	_ models.User,
) (_ router.Serde, status int, err error) {
	var response models.GetActivityLogResponse
	var logs []models.ActivityLog
	var organization models.Organization
	status = http.StatusOK

	projectID := params.Get("projectID").(string)
	options := models.GetLogsOptions{}
	options.Deserialize(body)

	if err = Repo.
		GetProjectsOrganization(projectID, &organization).
		Err(); err != nil {
		if errors.Is(err, repo.ErrorNotFound) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
			err = apierrors.ErrorFailedToGetResource(err)
		}

		goto done
	}

	// Only paid organiations can access this resource
	if !organization.Paid {
		status = http.StatusForbidden
		err = apierrors.ErrorNeedsUpgrade()
		goto done
	}

	if err = Repo.
		GetActivityLogs(projectID, options, &logs).
		Err(); err != nil {
		status = http.StatusInternalServerError
		err = apierrors.ErrorFailedToGetResource(err)
		goto done
	}

	response.Logs = make([]models.ActivityLogLite, len(logs))

	for index, log := range logs {
		response.Logs[index] = log.Lite()
	}

done:
	return &response, status, err
}
