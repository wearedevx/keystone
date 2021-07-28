package cmd

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/ci"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
)

// ciSendCmd represents the pushCi command
var ciSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Sends environment to a CI service",
	Long: `Sends environment to a CI service.

The CI service must have been setup using:
  $ ks ci setup
`,
	Run: func(_ *cobra.Command, _ []string) {
		var environment models.Environment
		ctx.MustHaveEnvironment(currentEnvironment)

		for _, accessibleEnvironment := range ctx.AccessibleEnvironments {
			if accessibleEnvironment.Name == currentEnvironment {
				environment = accessibleEnvironment
			}
		}

		message, err := ctx.PrepareMessagePayload(environment)

		if err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		ciService, err := SelectCiService(*ctx)

		if err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		// ciService = askForKeys(ciService)
		ciService.CheckSetup()
		if ciService.Error() != nil {
			ui.PrintError(ciService.Error().Error())
			os.Exit(1)
		}

		ciService.PushSecret(message, currentEnvironment)

		if ciService.Error() != nil {
			ui.PrintError(ciService.Error().Error())
			os.Exit(1)
		}
		ui.PrintSuccess(fmt.Sprintf("Secrets successfully sent to CI service, environment %s. See https://github.com/wearedevx/keystone-action to use them.", currentEnvironment))
	},
}

func init() {
	ciCmd.AddCommand(ciSendCmd)

	ciSendCmd.Flags().StringVar(&serviceName, "with", "", "identity provider. Either github or gitlab")
}

func SelectCiService(ctx core.Context) (ci.CiService, error) {
	var err error

	if serviceName == "" {
		prompt := promptui.Select{
			Label: "Select a ci service",
			Items: []string{
				"github",
			},
		}

		_, serviceName, err = prompt.Run()

		if err != nil {
			return nil, err
		}
	}

	return ci.GetCiService(serviceName, ctx, client.ApiURL)
}
