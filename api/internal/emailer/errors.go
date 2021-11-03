package emailer

import "errors"

var (
	ErrorEmailUnknown      = errors.New("unknown email error")
	ErrorEmailClientError  = errors.New("email client error")
	ErrorEmailServiceError = errors.New("email service error")
	ErrorEmailNoRecipient  = errors.New("no recipient for email")
)
