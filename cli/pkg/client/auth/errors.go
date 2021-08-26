package auth

import "errors"

var (
	ErrorUnauthorized   = errors.New("Unauthorized")
	ServiceNotAvailable = errors.New("connection refused")
	DeviceNotRegistered = errors.New("Device not registered")
)
