package auth

import "errors"

var (
	ErrorUnauthorized   = errors.New("Unauthorized")
	ErrorServiceNotAvailable = errors.New("connection refused")
	ErrorDeviceNotRegistered = errors.New("device not registered")
)
