package activitylog

import (
	"fmt"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/api/pkg/repo"
)

type Context struct {
	UserID        uint
	ProjectID     uint
	EnvironmentID uint
	Action        string
	Success       bool
}

func (cerr Context) IntoError(err error) error {
	msg := ""
	success := true

	if err != nil {
		msg = err.Error()
		success = false
	}

	return fmt.Errorf(
		"user=%d project=%d environment=%d action=\"%s\" success=%t error=\"%s\" %w",
		cerr.UserID,
		cerr.ProjectID,
		cerr.EnvironmentID,
		cerr.Action,
		success,
		msg,
		err,
	)
}

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

	log := new(models.ActivityLog)
	if log.FromString(err.Error()) {
		logger.err = logger.repo.SaveActivityLog(log).Err()
	}

	return logger
}
