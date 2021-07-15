package cmd

import (
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
	Use:   "ci send",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(_ *cobra.Command, _ []string) {

		var environment models.Environment

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
		// ciService = askForApiKey(ciService)

		ciService.PushSecret(message)
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
