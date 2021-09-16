package activitylog

import (
	"unsafe"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

type activityLogger struct {
	err  error
	repo repo.IRepo
}

type ActivityLogger interface {
	Save(err unsafe.Pointer) ActivityLogger
	Err() error
}

func NewActivityLogger(repo repo.IRepo) ActivityLogger {
	return &activityLogger{repo: repo}
}

func (logger *activityLogger) Err() error {
	return logger.err
}

func (logger *activityLogger) Save(err unsafe.Pointer) ActivityLogger {
	if logger.err != nil {
		return logger
	}

	if models.ErrorIsActivityLog(err) {
		log := (*models.ActivityLog)(err)
		logger.err = logger.repo.SaveActivityLog(log).Err()
	}

	return logger
}
