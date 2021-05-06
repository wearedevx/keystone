package utils

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type runner struct {
	err     error
	status  int
	actions []RunnerAction
}

type Runner interface {
	Run() Runner
	Status() int
	Error() error
}

type runnerAction struct {
	statusError   int
	statusSuccess int
	fn            func() error
}

type RunnerAction interface {
	SetStatusError(int) RunnerAction
	SetStatusSuccess(int) RunnerAction
	run() error
	success() int
	error() int
}

func NewAction(fn func() error) RunnerAction {
	return runnerAction{
		fn:            fn,
		statusSuccess: http.StatusOK,
		statusError:   http.StatusInternalServerError,
	}
}

func (a runnerAction) SetStatusSuccess(s int) RunnerAction {
	a.statusSuccess = s

	return a
}

func (a runnerAction) SetStatusError(s int) RunnerAction {
	a.statusError = s

	return a
}

func (a runnerAction) run() error {
	return a.fn()
}

func (a runnerAction) success() int {
	return a.statusSuccess
}

func (a runnerAction) error() int {
	return a.statusError
}

func NewRunner(actions []RunnerAction) Runner {
	return runner{
		status:  http.StatusOK,
		actions: actions,
	}
}

func (r runner) Run() Runner {
	for _, a := range r.actions {
		// Run the action,
		// if err: set status code to action.StatusError, or InternalServerError
		// else: set status code to action.StatusSuccess, or leave it as it were
		r.err = a.run()
		if r.err != nil {
			if errors.Is(r.err, gorm.ErrRecordNotFound) {
				r.status = http.StatusNotFound
			} else {
				r.status = a.error()
			}
			break
		} else {
			r.status = a.success()
		}
	}

	return r
}

func (r runner) Status() int {
	return r.status
}

func (r runner) Error() error {
	return r.err
}
