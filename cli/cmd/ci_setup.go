package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/ui"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		ciService, err := SelectCiService(*ctx)

		if err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		ciService = ciService.Setup()
		// ciService = askForApiKey(ciService)
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
