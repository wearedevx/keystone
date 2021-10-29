package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/ci"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

// ciSendCmd represents the pushCi command
var ciSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Sends environment to a CI service",
	Long: `Sends environment to a CI service.

This command will send all your secrets and files followed by keystone to your CI service.

The CI service must have been setup using: ` + "`" + `ks ci add` + "`" + `
`,
	Example: `# To send the current environment:
ks ci send

# To send a specific environment:
ks ci send --env prod
`,
	Run: func(_ *cobra.Command, _ []string) {
		var environment models.Environment
		ctx.MustHaveEnvironment(currentEnvironment)

		fetchMessages()

		for _, accessibleEnvironment := range ctx.AccessibleEnvironments {
			if accessibleEnvironment.Name == currentEnvironment {
				environment = accessibleEnvironment
			}
		}

		mustNotHaveMissingSecrets(environment)
		mustNotHaveMissingFiles(environment)

		if !prompts.ConfirmSendEnvironmentToCiService(environment.Name) {
			exit(nil)
		}

		message, err := ctx.PrepareMessagePayload(environment)
		if err != nil {
			exit(kserrors.PayloadErrors(err))
		}

		ciServices, err := ci.ListCiServices(ctx)
		exitIfErr(err)

		for _, serviceDef := range ciServices {
			ciService, err := ci.GetCiService(serviceDef.Name, ctx, client.ApiURL)
			exitIfErr(err)

			if err = ciService.CheckSetup().Error(); err != nil {
				if errors.Is(err, ci.ErrorMissinCiInformation) {
					err = kserrors.MissingCIInformation(serviceName, nil)
				}

				exit(err)
			}

			if err = ciService.
				PushSecret(message, currentEnvironment).
				Error(); err != nil {
				err = kserrors.CouldNotSendToCIService(err)
			}
			exitIfErr(err)

			// TODO: Printing should be done by the ui, or ui/prompts packages
			ciService.PrintSuccess(currentEnvironment)
		}
	},
}

func init() {
	ciCmd.AddCommand(ciSendCmd)

	ciSendCmd.Flags().StringVar(&serviceName, "with", "", "Ci service name.")
}
