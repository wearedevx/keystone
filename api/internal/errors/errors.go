package errors

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type Error struct {
	message string
	cause   error
}

func (e Error) Error() string {
	return fmt.Sprintf("%s", e.message)
}

func (e Error) Is(err error) bool {
	if strings.Contains(err.Error(), e.message) {
		return true
	}

	if strings.Contains(e.message, err.Error()) {
		return true
	}

	return errors.Is(e.cause, err)
}

func (e Error) Unwrap() error {
	return e.cause
}

func newError(message string, cause error) Error {
	log.Printf("[ERROR] %s: %v", message, cause)
	return Error{message: message, cause: cause}
}
