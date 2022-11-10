package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	d "runtime/debug"
	"strings"

	"github.com/wearedevx/keystone/api/pkg/apierrors"
	"github.com/wearedevx/keystone/api/pkg/models"

	"github.com/wearedevx/keystone/cli/internal/config"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/messages"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/client/auth"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
	"github.com/wearedevx/keystone/cli/ui/display"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

/// Exists the program.
/// If err is nil, exits with a 0 status code.
/// If err is an error, prints it and exits with a 1 status code
func exit(err error) {
	exitIfErr(err)

	os.Exit(0)
}

/// If err is not nil, prints the error and exits with a 1 status code
/// Does nothing other wise
func exitIfErr(err error) {
	if err == nil {
		return
	} else if display.Error(err) {
		if debug {
			fmt.Fprintln(os.Stderr, "\nStacktrace for the previous error:")
			d.PrintStack()
		}
		os.Exit(1)
	}
}

/// Get messages and print the changes
func fetchMessages() (
	core.ChangesByEnvironment,
	messages.MessageService,
	*kserrors.Error,
) {
	ms := messages.NewMessageService(ctx)
	changes := ms.GetMessages()

	err := ms.Err()
	if err != nil {
		config.CheckExpiredTokenError(err)
	}

	if err == nil {
		display.Changes(changes)
	}

	return changes, ms, err
}

/// Gets messages and print the changes
/// Exits the program on error
func mustFetchMessages() (core.ChangesByEnvironment, messages.MessageService) {
	changes, ms, err := fetchMessages()

	exitIfErr(err)

	return changes, ms
}

/// Gets messages and print the changes
/// Prints a warnning on error
func shouldFetchMessages() (core.ChangesByEnvironment, messages.MessageService) {
	changes, ms, err := fetchMessages()
	if err != nil {
		ui.PrintStdErr(
			"WARNING: Could not get messages (%s)",
			err.Name(),
		)
		ui.PrintStdErr(err.Help())
	}

	return changes, ms
}

// handleClientError handles most of the errors returned by the KeystoneClent
// It prints the error, then exits the program when `err` is an error it can handle.
// If the the error is too generic to be handled here, it returns without
// printing or exiting the program, so that the caller can handle the error
// its own way.
func handleClientError(err error) {
	switch {
	case errors.Is(err, auth.ErrorUnauthorized),
		errors.Is(err, auth.ErrorRefreshNotFound),
		errors.Is(err, auth.ErrorNoToken),
		errors.Is(err, auth.ErrorNoRefresh):
		config.Logout()
		kserrors.InvalidConnectionToken(err).Print()

		// Errors That should never happen
	case errors.Is(err, apierrors.ErrorUnknown),
		errors.Is(err, apierrors.ErrorFailedToGetPermission),
		errors.Is(err, apierrors.ErrorFailedToWriteMessage),
		errors.Is(err, apierrors.ErrorFailedToSetEnvironmentVersion),
		errors.Is(err, apierrors.ErrorOrganizationWithoutAnAdmin):
		kserrors.UnkownError(err).Print()

		// General Errors
	case errors.Is(err, apierrors.ErrorPermissionDenied):
		kserrors.PermissionDenied(currentEnvironment, err).Print()

		// These should be handled by the controller/service
	case errors.Is(err, apierrors.ErrorBadRequest),
		errors.Is(err, apierrors.ErrorEmptyPayload),
		errors.Is(err, apierrors.ErrorFailedToCreateResource),
		errors.Is(err, apierrors.ErrorFailedToGetResource),
		errors.Is(err, apierrors.ErrorFailedToUpdateResource),
		errors.Is(err, apierrors.ErrorFailedToDeleteResource),
		errors.Is(err, apierrors.ErrorMemberAlreadyInProject),
		errors.Is(err, apierrors.ErrorNotAMember):
		return

		// Subscription Errors
	case errors.Is(err, apierrors.ErrorNeedsUpgrade):
		kserrors.FeatureRequiresToUpgrade(err).Print()

	case errors.Is(err, apierrors.ErrorAlreadySubscribed):
		kserrors.AlreadySubscribed(err).Print()

	case errors.Is(err, apierrors.ErrorFailedToStartCheckout):
		kserrors.CannotUpgrade(err).Print()

	case errors.Is(err, apierrors.ErrorFailedToGetManagementLink):
		kserrors.ManagementInaccessible(err).Print()

		// Device Errors
	case errors.Is(err, apierrors.ErrorNoDevice),
		errors.Is(err, auth.ErrorDeviceNotRegistered),
		// There is a an undetermined path where only this seem to work…
		strings.Contains(err.Error(), auth.ErrorDeviceNotRegistered.Error()):
		kserrors.DeviceNotRegistered(err).Print()

	case errors.Is(err, apierrors.ErrorBadDeviceName):
		kserrors.BadDeviceName(err).Print()

		// Organization Errors
	case errors.Is(err, apierrors.ErrorBadOrganizationName):
		kserrors.BadOrganizationName(err).Print()

	case errors.Is(err, apierrors.ErrorOrganizationNameAlreadyTaken):
		kserrors.OrganizationNameAlreadyTaken(err).Print()

	case errors.Is(err, apierrors.ErrorNotOrganizationOwner):
		kserrors.MustOwnTheOrganization(err).Print()

		// Invite Errors
	case errors.Is(err, apierrors.ErrorFailedToCreateMailContent),
		errors.Is(err, apierrors.ErrorFailedToSendMail):
		kserrors.CouldntSendInvite(err).Print()

		// Role Errors
	case errors.Is(err, apierrors.ErrorFailedToSetRole):
		kserrors.CouldntSetRole(err).Print()

		// Members Errors
	case errors.Is(err, apierrors.ErrorFailedToAddMembers):
		kserrors.CannotAddMembers(err).Print()

	default:
		ui.PrintError(err.Error())
	}

	os.Exit(1)
}

// Exits the program if required secrets are missing
func mustNotHaveMissingSecrets(environment models.Environment) {
	missing, hasMissing := ctx.MissingSecretsForEnvironment(
		environment.Name,
	)

	if hasMissing {
		exit(
			kserrors.RequiredSecretsAreMissing(missing, environment.Name, nil),
		)
	}
}

// Exits the program if required files ar missing
func mustNotHaveMissingFiles(environment models.Environment) {
	missing, hasMissing := ctx.MissingFilesForEnvironment(
		environment.Name,
	)

	if hasMissing {
		exit(
			kserrors.RequiredFilesAreMissing(missing, environment.Name, nil),
		)
	}
}

// Exists the program if required secret or files are missing
// after printing a list of all the missing things
func mustNotHaveAnyRequiredThingMissing(ctx *core.Context) {
	missingSecrets, hasMisssingSecrets := ctx.
		MissingSecretsForEnvironment(currentEnvironment)

	for _, ms := range missingSecrets {
		fmt.Fprintf(os.Stderr, "Required Secret is missing: %s\n", ms)
	}

	missingFiles, hasMissingFiles := ctx.
		MissingFilesForEnvironment(currentEnvironment)

	for _, mf := range missingFiles {
		fmt.Fprintf(os.Stderr, "Required file is missing or empty: %s\n", mf)
	}

	if hasMissingFiles || hasMisssingSecrets {
		os.Exit(1)
	}
}

// Exits the program if the user is not admin on the proec
func mustBeAdmin(projectService *client.Project) {
	members, err := projectService.GetAllMembers()
	if err != nil {
		exit(kserrors.UnkownError(err))
	}

	account, _ := config.GetCurrentAccount()

	for _, member := range members {
		if member.User.UserID == account.UserID {
			if member.Role.Name == "admin" {
				return
			}
		}
	}

	exit(kserrors.PermissionDenied(currentEnvironment, nil))
}

// Removes the project directory path form `filePathArg`,
// and also any `./` or `/`
func cleanPathArgument(
	filePathArg string,
	wd string,
) (filePath string, err error) {
	filePathInCwd := filepath.Join(CWD, filePathArg)
	filePathClean := filepath.Clean(filePathInCwd)

	if !strings.HasPrefix(filePathClean, wd) {
		return "", fmt.Errorf("file %s not in project", filePathArg)
	}

	filePath = strings.TrimPrefix(filePathClean, ctx.Wd)
	filePath = strings.TrimPrefix(filePath, "/")

	return filePath, nil
}

// Exists if the at least one of the memberIDs in `memberIDs`
// does not exist.
func mustMembersExist(c client.KeystoneClient, memberIDs []string) {
	r, err := c.Users().CheckUsersExist(memberIDs)
	if err != nil {
		// The HTTP request must have failed
		exit(kserrors.UnkownError(err))
	}

	if r.Error != "" {
		exit(kserrors.UsersDontExist(r.Error, nil))
	}
}

// Get all the roles available to the user/organization,
// and exits the program on error
func mustGetRoles(c client.KeystoneClient) []models.Role {
	projectID := ctx.GetProjectID()
	roles, err := c.Roles().GetAll(projectID)
	exitIfErr(err)

	return roles
}

// Asks the user to select an organization from the ones they belong to.
// Exits the program on error
func mustGetOrganizationName(
	o *client.Organizations,
	args []string,
) (organizationName string) {
	if len(args) == 1 {
		organizationName = args[0]
	} else {
		organizations, err := o.GetAll()
		exitIfErr(err)

		organization := prompts.OrganizationsSelect(organizations)
		organizationName = organization.Name
	}

	return organizationName
}

// Returns true if `needle` is in `haytack`
func isIn(haystack []string, needle string) bool {
	for _, hay := range haystack {
		if hay == needle {
			return true
		}
	}

	return false
}
