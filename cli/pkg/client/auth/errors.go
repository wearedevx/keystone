package auth

import "errors"

var (
	ErrorUnauthorized   = errors.New("Unauthorized")
	ServiceNotAvailable = errors.New("connection refused")
)
