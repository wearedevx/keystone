package errors

//go:generate go run ./generate_errors.go

import (
	"fmt"
	"os"

	"github.com/wearedevx/keystone/cli/ui"
)

type Error struct {
	name  string
	help  string
	cause error
	meta  map[string]interface{}
}

func NewError(name string, help string, meta map[string]interface{}, cause error) *Error {
	err := new(Error)

	err.name = name
	err.help = help
	err.cause = cause

	meta["Name"] = name

	if err.cause != nil {
		meta["Cause"] = err.cause.Error()
	} else {
		meta["Cause"] = "<unkown>"
	}

	err.meta = meta

	return err
}

func (e *Error) SetCause(err error) *Error {
	e.cause = err

	e.meta["Cause"] = e.cause.Error()

	return e
}

func (e *Error) Cause() error {
	return e.cause
}

func (e *Error) Print() {
	fmt.Fprintln(os.Stderr, e.Error())
	// os.Stderr.WriteString(e.Error())
}

func (e *Error) Error() string {
	return ui.RenderTemplate(e.name, e.help, e.meta)
}
