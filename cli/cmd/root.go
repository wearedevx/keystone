/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/api/pkg/apierrors"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/config"
	"github.com/wearedevx/keystone/cli/internal/environments"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/internal/messages"
	"github.com/wearedevx/keystone/cli/ui"

	"github.com/wearedevx/keystone/cli/pkg/client/auth"
	"github.com/wearedevx/keystone/cli/pkg/core"
)

var ksauthURL string //= "http://localhost:9000"
var ksapiURL string  //= "http://localhost:9001"

var cfgFile string = ""
var currentEnvironment string
var quietOutput bool
var skipPrompts bool

var ctx *core.Context

var CWD string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ks <command> [sub-command] [inputs]...",
	Short: "A safe system for developers to store, share and use secrets.",
	Long:  `A safe system for developers to store, share and use secrets.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Usage(); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() int {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	return 0
}

var noEnvironmentCommands []string
var noProjectCommands []string
var noLoginCommands []string

func isIn(haystack []string, needle string) bool {
	for _, hay := range haystack {
		if hay == needle {
			return true
		}
	}

	return false
}

func findCurrentCommand(args []string) string {
	for _, candidate := range args[1:] {
		if !strings.HasPrefix(candidate, "-") &&
			!strings.HasPrefix(candidate, "/") &&
			candidate != "ks" {
			return candidate
		}
	}

	return ""
}

func Initialize() {
	if len(os.Args) <= 1 {
		return
	}

	checkEnvironment := true
	checkProject := true
	checkLogin := false

	askHelp := core.Contains(os.Args, "--help")

	if len(os.Args) > 1 {
		command := findCurrentCommand(os.Args)
		checkEnvironment = !isIn(noEnvironmentCommands, command) && !askHelp
		checkProject = !isIn(noProjectCommands, command) && !askHelp
		checkLogin = !isIn(noLoginCommands, command) && !askHelp

		if command == "init" {
			ctx = core.New(core.CTX_INIT)
		} else {
			ctx = core.New(core.CTX_RESOLVE)
		}
	}

	if checkProject {
		exitIfErr(ctx.Err())
	}

	isKeystoneFile := keystonefile.ExistsKeystoneFile(ctx.Wd)

	current := ctx.CurrentEnvironment()
	ctx.SetError(nil)

	if currentEnvironment == "" {
		currentEnvironment = current
	}

	if checkProject && !isKeystoneFile {
		exit(kserrors.NotAKeystoneProject(".", nil))
	}

	if checkProject && config.IsLoggedIn() {
		es := environments.NewEnvironmentService(ctx)
		exitIfErr(es.Err())

		ctx.AccessibleEnvironments = es.GetAccessibleEnvironments()
		exitIfErr(ctx.Err())

		// If no accessible environment, then user has no access to the project
		if len(ctx.AccessibleEnvironments) == 0 {
			exit(
				kserrors.ProjectDoesntExist(
					ctx.GetProjectName(),
					ctx.GetProjectID(),
					nil,
				),
			)
		}

		if isKeystoneFile {
			environmentsToSave := make([]models.Environment, 0)
			for _, env := range ctx.AccessibleEnvironments {
				env.VersionID = ""
				environmentsToSave = append(environmentsToSave, env)
			}

			ctx.Init(models.Project{
				Environments: environmentsToSave,
			})
		}
		ctx.RemoveForbiddenEnvironments(ctx.AccessibleEnvironments)

		if currentEnvironment == "" {
			currentEnvironment = ctx.CurrentEnvironment()
		}
	}

	if checkEnvironment && !ctx.HasEnvironment(currentEnvironment) {
		ctx.Init(models.Project{})
		if currentEnvironment == "" {
			ctx.SetCurrent("dev")
		}
	}

	if checkLogin && !config.IsLoggedIn() {
		exit(kserrors.MustBeLoggedIn(nil))
	}

}

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	CWD = cwd

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")
	RootCmd.PersistentFlags().BoolVarP(&quietOutput, "quiet", "q", false, "make the output machine readable")

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.config/keystone/keystone.yaml)")

	RootCmd.PersistentFlags().BoolVarP(&skipPrompts, "skip", "s", false, "skip prompts and use default")

	RootCmd.PersistentFlags().StringVarP(&currentEnvironment, "env", "", "", "environment to use instead of the current one")

	cobra.OnInitialize(func() {
		// Call directly initConfig. cobra doesn't call initConfig func.
		config.InitConfig(cfgFile)
		Initialize()
	})
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	noEnvironmentCommands = []string{
		"login", "logout", "documentation", "init", "whoami", "invite", "version", "device", "orga", "project",
	}

	noProjectCommands = noEnvironmentCommands

	noLoginCommands = []string{"login", "source", "documentation", "version"}
}

func fetchMessages(ms messages.MessageService) (core.ChangesByEnvironment, *kserrors.Error) {
	if ms == nil {
		ms = messages.NewMessageService(ctx)
	}
	changes := ms.GetMessages()

	err := ms.Err()
	if err != nil {
		config.CheckExpiredTokenError(err)
	}

	return changes, err
}

func mustFetchMessages(ms messages.MessageService) core.ChangesByEnvironment {
	changes, err := fetchMessages(ms)

	exitIfErr(err)

	return changes
}

func shouldFetchMessages(ms messages.MessageService) core.ChangesByEnvironment {
	changes, err := fetchMessages(ms)

	if err != nil {
		ui.PrintStdErr(
			"WARNING: Could not get messages (%s)",
			err.Error(),
		)
	}

	return changes
}

// handleClientError handles most of the errors returned by the KeystoneClent
// It prints the error, then exits the program when `err` is an error it can handle.
// If the the error is too generic to be handled here, it returns without
// printing or exiting the program, so that the caller can handle the error
// its own way.
func handleClientError(err error) {
	switch {
	case errors.Is(err, auth.ErrorUnauthorized):
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
		kserrors.PermissionDenied(currentEnvironment, err)

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
	case errors.Is(err, apierrors.ErrorNoDevice):
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

func exit(err error) {
	if err == nil {
		os.Exit(0)
		return
	}

	printError(err)

	os.Exit(1)
}

func printError(err error) bool {
	if err == nil {
		return false
	}

	if kserrors.IsKsError(err) {
		kserr := kserrors.AsKsError(err)
		if kserr == nil {
			return false
		}
		kserr.Print()
	} else {
		ui.PrintError(err.Error())
	}

	return true
}

func exitIfErr(err error) {
	if err == nil {
		return
	}
	if err != nil {
		if printError(err) {
			os.Exit(1)
		}
	}
}
