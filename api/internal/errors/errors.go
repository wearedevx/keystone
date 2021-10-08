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

// General Errors
func ErrorUnknown(cause error) error {
	return newError("unknown", cause)
}

func ErrorBadRequest(cause error) error {
	return newError("bad request", cause)
}

func ErrorPermissionDenied() error {
	return newError("permission denied", nil)
}

func ErrorEmptyPayload() error {
	return newError("empty payload cannot be wrtten", nil)
}

func ErrorFailedToCreateResource(cause error) error {
	return newError("failed to create resource", cause)
}

func ErrorFailedToGetResource(cause error) error {
	return newError("failed to get", cause)
}

func ErrorFailedToUpdateResource(cause error) error {
	return newError("failed to update", cause)
}

func ErrorFailedToDeleteResource(cause error) error {
	return newError("failed to delete", cause)
}

// Subscription Errors
func ErrorNeedsUpgrade() error {
	return newError("needs upgrade", nil)
}

func ErrorAlreadySubscribed() error {
	return newError("already subscribed", nil)
}

func ErrorFailedToStartCheckout(cause error) error {
	return newError("failed to start checkout", cause)
}

func ErrorFailedToGetManagementLink(cause error) error {
	return newError("failed to get management link", cause)
}

func ErrorFailedToUpdateSubscription(cause error) error {
	return newError("failed to update subscription", cause)
}

func ErrorCheckoutCompleteFailed(cause error) error {
	return newError("checkout complete failed", cause)
}

func ErrorSubscriptionPaidFailed(cause error) error {
	return newError("subscription paid failed", cause)
}

func ErrorSubscriptionUnpaidFailed(cause error) error {
	return newError("subscription unpaid failed", cause)
}

func ErrorSubscriptionCanceledFailed(cause error) error {
	return newError("subscription canceled failed", cause)
}

// Device Errors
func ErrorNoDevice() error {
	return newError("no device", nil)
}

func ErrorBadDeviceName() error {
	return newError("bad device name", nil)
}

// Organization Errors
func ErrorBadOrganizationName() error {
	return newError("bad organization name", nil)
}

func ErrorOrganizationNameAlreadyTaken() error {
	return newError("organization name already taken", nil)
}

func ErrorNotOrganizationOwner() error {
	return newError("not organization owner", nil)
}

func ErrorOrganizationWithoutAnAdmin(cause error) error {
	return newError("organization without an admin", cause)
}

// Invite Errors
func ErrorFailedToCreateMailContent(cause error) error {
	return newError("failed to create mail content", cause)
}

func ErrorFailedToSendMail(cause error) error {
	return newError("failed to send mail", cause)
}

// Role Errors
func ErrorFailedToSetRole(cause error) error {
	return newError("failed to set role", cause)
}

func ErrorFailedToGetPermission(cause error) error {
	return newError("failed to get permission", cause)
}

// Messages Errors
func ErrorFailedToWriteMessage(cause error) error {
	return newError("failed to write message", cause)
}

// Environment Errors
func ErrorFailedToSetEnvironmentVersion(cause error) error {
	return newError("failed to set environment version", cause)
}

// Members Errors
func ErrorFailedToAddMembers(cause error) error {
	return newError("failed to add members", cause)
}

func ErrorMemberAlreadyInProject() error {
	return newError("member already in project", nil)
}

func ErrorNotAMember() error {
	return newError("not a member", nil)
}
