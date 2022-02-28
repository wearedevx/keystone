package errors

//go:generate go run ./generate_errors.go

import (
	"fmt"
	"os"
	"strings"

	"github.com/wearedevx/keystone/cli/ui"
)

type Error struct {
	name  string
	help  string
	cause error
	meta  map[string]interface{}
}

// NewError function creates a new keysone error
func NewError(
	name string,
	help string,
	meta map[string]interface{},
	cause error,
) *Error {
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

// SetCause method sets the original error that cause this one to be raised
func (e *Error) SetCause(err error) *Error {
	e.cause = err

	e.meta["Cause"] = e.cause.Error()

	return e
}

// Cause method return the originalerror that caused this one to be raised
func (e *Error) Cause() error {
	return e.cause
}

// Name method returns the name of the error
func (e *Error) Name() string {
	return e.name
}

// Help method returns the help part of the error
func (e *Error) Help() string {
	rendered := e.Error()
	s := strings.Split(rendered, "\n")
	helpLines := s[1:]

	return strings.Join(helpLines, "\n")
}

// Print method prints the error to stderr
func (e *Error) Print() {
	fmt.Fprintln(os.Stderr, e.Error())
	// os.Stderr.WriteString(e.Error())
}

// Error method returns the formatted error message
func (e *Error) Error() string {
	return ui.RenderTemplate(e.name, e.help, e.meta)
}

// IsKsError function checks if `err` is an instance of a keystone `Error`
func IsKsError(err error) bool {
	if err == nil {
		return false
	}

	kserrorPtrType := fmt.Sprintf("%T", &Error{})
	errType := fmt.Sprintf("%T", err)

	return kserrorPtrType == errType
}

// AsKsError function dynamically casts an `error` into a keystone `Error`.
// You should use `IsKsError(error) bool` beforehand to ensure this cast
// will wor
func AsKsError(err error) *Error {
	return (err).(*Error)
}
