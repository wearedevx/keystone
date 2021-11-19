package errors

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
	return newError("empty payload cannot be written", nil)
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

func ErrorNoDevice() error {
	return newError("no device", nil)
}

func ErrorBadDeviceName() error {
	return newError("bad device name", nil)
}

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

func ErrorFailedToCreateMailContent(cause error) error {
	return newError("failed to create mail content", cause)
}

func ErrorFailedToSendMail(cause error) error {
	return newError("failed to send mail", cause)
}

func ErrorFailedToSetRole(cause error) error {
	return newError("failed to set role", cause)
}

func ErrorFailedToGetPermission(cause error) error {
	return newError("failed to get permission", cause)
}

func ErrorFailedToWriteMessage(cause error) error {
	return newError("failed to write message", cause)
}

func ErrorFailedToSetEnvironmentVersion(cause error) error {
	return newError("failed to set environment version", cause)
}

func ErrorFailedToAddMembers(cause error) error {
	return newError("failed to add members", cause)
}

func ErrorMemberAlreadyInProject() error {
	return newError("member already in project", nil)
}

func ErrorNotAMember() error {
	return newError("not a member", nil)
}
