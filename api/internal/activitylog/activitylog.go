package activitylog

import (
	"fmt"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

type activityLogger struct {
	err  error
	repo repo.IRepo
}

type ActivityLogger interface {
	Save(err error) ActivityLogger
	Err() error
}

func NewActivityLogger(repo repo.IRepo) ActivityLogger {
	return &activityLogger{repo: repo}
}

func (logger *activityLogger) Err() error {
	return logger.err
}

func (logger *activityLogger) Save(err error) ActivityLogger {
	if logger.err != nil {
		return logger
	}

	if models.ErrorIsActivityLog(err) {
		log := err.(*models.ActivityLog)
		logger.err = logger.repo.SaveActivityLog(log).Err()
		fmt.Printf("log: %+v %s\n", log.Success, log.Message)
		fmt.Printf("logger.err: %+v\n", logger.err)
	}

	return logger
}
