package apierrors

import (
	"errors"
	"strings"
)

var (
	ErrorUnknown                       error = errors.New("unknown")
	ErrorBadRequest                    error = errors.New("bad request")
	ErrorPermissionDenied              error = errors.New("permission denied")
	ErrorEmptyPayload                  error = errors.New("empty payload cannot be wrtten")
	ErrorFailedToCreateResource        error = errors.New("failed to create resource")
	ErrorFailedToGetResource           error = errors.New("failed to get")
	ErrorFailedToUpdateResource        error = errors.New("failed to update")
	ErrorFailedToDeleteResource        error = errors.New("failed to delete")
	ErrorNeedsUpgrade                  error = errors.New("needs upgrade")
	ErrorAlreadySubscribed             error = errors.New("already subscribed")
	ErrorFailedToStartCheckout         error = errors.New("failed to start checkout")
	ErrorFailedToGetManagementLink     error = errors.New("failed to get management link")
	ErrorFailedToUpdateSubscription    error = errors.New("failed to update subscription")
	ErrorCheckoutCompleteFailed        error = errors.New("checkout complete failed")
	ErrorSubscriptionPaidFailed        error = errors.New("subscription paid failed")
	ErrorSubscriptionUnpaidFailed      error = errors.New("subscription unpaid failed")
	ErrorSubscriptionCanceledFailed    error = errors.New("subscription canceled failed")
	ErrorNoDevice                      error = errors.New("no device")
	ErrorBadDeviceName                 error = errors.New("bad device name")
	ErrorBadOrganizationName           error = errors.New("bad organization name")
	ErrorOrganizationNameAlreadyTaken  error = errors.New("organization name already taken")
	ErrorNotOrganizationOwner          error = errors.New("not organization owner")
	ErrorOrganizationWithoutAnAdmin    error = errors.New("organization without an admin")
	ErrorFailedToCreateMailContent     error = errors.New("failed to create mail content")
	ErrorFailedToSendMail              error = errors.New("failed to send mail")
	ErrorFailedToSetRole               error = errors.New("failed to set role")
	ErrorFailedToGetPermission         error = errors.New("failed to get permission")
	ErrorFailedToWriteMessage          error = errors.New("failed to write message")
	ErrorFailedToSetEnvironmentVersion error = errors.New("failed to set environment version")
	ErrorFailedToAddMembers            error = errors.New("failed to add members")
	ErrorMemberAlreadyInProject        error = errors.New("member already in project")
	ErrorNotAMember                    error = errors.New("not a member")
)

func FromString(s string) error {
	s = strings.TrimSpace(s)

	switch s {
	case ErrorUnknown.Error():
		return ErrorUnknown
	case ErrorBadRequest.Error():
		return ErrorBadRequest
	case ErrorPermissionDenied.Error():
		return ErrorPermissionDenied
	case ErrorEmptyPayload.Error():
		return ErrorEmptyPayload
	case ErrorFailedToCreateResource.Error():
		return ErrorFailedToCreateResource
	case ErrorFailedToGetResource.Error():
		return ErrorFailedToGetResource
	case ErrorFailedToUpdateResource.Error():
		return ErrorFailedToUpdateResource
	case ErrorFailedToDeleteResource.Error():
		return ErrorFailedToDeleteResource
	case ErrorNeedsUpgrade.Error():
		return ErrorNeedsUpgrade
	case ErrorAlreadySubscribed.Error():
		return ErrorAlreadySubscribed
	case ErrorFailedToStartCheckout.Error():
		return ErrorFailedToStartCheckout
	case ErrorFailedToGetManagementLink.Error():
		return ErrorFailedToGetManagementLink
	case ErrorFailedToUpdateSubscription.Error():
		return ErrorFailedToUpdateSubscription
	case ErrorCheckoutCompleteFailed.Error():
		return ErrorCheckoutCompleteFailed
	case ErrorSubscriptionPaidFailed.Error():
		return ErrorSubscriptionPaidFailed
	case ErrorSubscriptionUnpaidFailed.Error():
		return ErrorSubscriptionUnpaidFailed
	case ErrorSubscriptionCanceledFailed.Error():
		return ErrorSubscriptionCanceledFailed
	case ErrorNoDevice.Error():
		return ErrorNoDevice
	case ErrorBadDeviceName.Error():
		return ErrorBadDeviceName
	case ErrorBadOrganizationName.Error():
		return ErrorBadOrganizationName
	case ErrorOrganizationNameAlreadyTaken.Error():
		return ErrorOrganizationNameAlreadyTaken
	case ErrorNotOrganizationOwner.Error():
		return ErrorNotOrganizationOwner
	case ErrorOrganizationWithoutAnAdmin.Error():
		return ErrorOrganizationWithoutAnAdmin
	case ErrorFailedToCreateMailContent.Error():
		return ErrorFailedToCreateMailContent
	case ErrorFailedToSendMail.Error():
		return ErrorFailedToSendMail
	case ErrorFailedToSetRole.Error():
		return ErrorFailedToSetRole
	case ErrorFailedToGetPermission.Error():
		return ErrorFailedToGetPermission
	case ErrorFailedToWriteMessage.Error():
		return ErrorFailedToWriteMessage
	case ErrorFailedToSetEnvironmentVersion.Error():
		return ErrorFailedToSetEnvironmentVersion
	case ErrorFailedToAddMembers.Error():
		return ErrorFailedToAddMembers
	case ErrorMemberAlreadyInProject.Error():
		return ErrorMemberAlreadyInProject
	case ErrorNotAMember.Error():
		return ErrorNotAMember
	default:
		return errors.New(s)
	}
}
