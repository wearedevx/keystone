package emailer

import "errors"

var (
	EmailErrorUnknown      = errors.New("Unknown email error")
	EmailErrorClientError  = errors.New("Email client error")
	EmailErrorServiceError = errors.New("Email service error")
	EmailErrorNoRecipient  = errors.New("No recipient for email")
)
