package repo

import (
	"github.com/wearedevx/keystone/api/pkg/models"
	"gorm.io/gorm/clause"
)

func (r *Repo) SaveActivityLog(al *models.ActivityLog) IRepo {
	r.err = r.GetDB().
		Omit(clause.Associations).
		Create(al).
		Error

	return r
}

// GetActivityLogs returns a list of all activity logs associated with
// the given project
func (r *Repo) GetActivityLogs(
	projectID string,
	options models.GetLogsOptions,
	logs *[]models.ActivityLog,
) IRepo {
	if r.Err() != nil {
		return r
	}

	req := r.GetDB().
		Model(&models.ActivityLog{}).
		Joins("inner join projects on activity_logs.project_id = projects.id")

	if len(options.Users) > 0 {
		req = req.
			Joins("inner join users on activity_logs.user_id = users.id").
			Where("users.user_id IN (?)", options.Users)
	}

	if len(options.Environments) > 0 {
		req = req.
			Joins("left join environments on activity_logs.environment_id = environments.id").
			Where("environments.name IN (?) OR activity_logs.environment_id IS NULL", options.Environments)
	}

	if len(options.Actions) > 0 {
		req = req.Where("activity_logs.action IN (?)", options.Actions)
	}

	var l = make([]models.ActivityLog, options.Limit)
	r.err = req.
		Where("projects.uuid = ?", projectID).
		Order("activity_logs.created_at DESC").
		Limit(int(options.Limit)).
		Preload(clause.Associations).
		Find(&l).
		Error

	if r.err == nil {
		*logs = reverseLogs(l)
	}

	return r
}

func reverseLogs(logs []models.ActivityLog) []models.ActivityLog {
	out := make([]models.ActivityLog, len(logs))

	for index, log := range logs {
		backIndex := len(out) - 1 - index
		out[backIndex] = log
	}

	return out
}
