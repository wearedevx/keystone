package repo

import (
	"github.com/wearedevx/keystone/api/pkg/models"
	"gorm.io/gorm/clause"
)

func (r *Repo) SaveActivityLog(al *models.ActivityLog) IRepo {
	r.err = r.GetDb().
		Omit(clause.Associations).
		Create(al).
		Error

	return r
}

// GetActivityLogs returns a list of all activity logs associated with
// the given project
func (r *Repo) GetActivityLogs(projectID string, logs *[]models.ActivityLog) IRepo {
	if r.Err() != nil {
		return r
	}

	r.err = r.GetDb().
		Model(&models.ActivityLog{}).
		Joins("inner join projects on activity_logs.project_id = project.id").
		Where("projects.uuid = ?", projectID).
		Preload("Project").
		Preload("User").
		Preload("Envrionment").
		Find(logs).
		Error

	return r
}
