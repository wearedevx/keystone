package repo

import "errors"

var (
	ErrorNotFound  = errors.New("not found")
	ErrorBadName   = errors.New("bad name")
	ErrorNameTaken = errors.New("name taken")
)
