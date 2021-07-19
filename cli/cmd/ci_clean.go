package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/ui"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove secrets from a CI service.",
	Long:  `Remove secrets from a CI service.`,
	Run: func(cmd *cobra.Command, args []string) {

		ctx.MustHaveEnvironment(currentEnvironment)

		ciService, err := SelectCiService(*ctx)

		if err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		ciService.CleanSecret(currentEnvironment)

		if ciService.Error() != nil {
			if strings.Contains(ciService.Error().Error(), "404") == true {
				ui.PrintSuccess(fmt.Sprintf("No secret found for environment %s in CI service", currentEnvironment))
				os.Exit(0)
			}

			ui.PrintError(ciService.Error().Error())
			os.Exit(1)
		}
		ui.PrintSuccess(fmt.Sprintf("Secrets successfully removed from CI service, environment %s.", currentEnvironment))
	},
}

func init() {
	ciCmd.AddCommand(cleanCmd)
}
