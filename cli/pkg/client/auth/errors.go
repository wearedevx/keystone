package auth

import "errors"

var (
	ErrorUnauthorized        = errors.New("unauthorized")
	ErrorServiceNotAvailable = errors.New("connection refused")
	ErrorDeviceNotRegistered = errors.New("device not registered")
	ErrorRefreshNotFound     = errors.New("refresh not found")
	ErrorNoToken             = errors.New("no token in response")
	ErrorNoRefresh           = errors.New("no refresh token in response")
)
