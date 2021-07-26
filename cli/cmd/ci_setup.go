package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/ui"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup CI service settings to push secrets.",
	Long: `Use this command to modify CI service specific settings
like API key and project name.`,

	Run: func(cmd *cobra.Command, args []string) {

		ciService, err := SelectCiService(*ctx)

		if err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		ciService = ciService.Setup()

		if ciService.Error() != nil {
			ui.PrintError(ciService.Error().Error())
			os.Exit(1)
		}
		ui.PrintSuccess("Ci service setup successfully")
	},
}

func init() {
	ciCmd.AddCommand(setupCmd)

}
