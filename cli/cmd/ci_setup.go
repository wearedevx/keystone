package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/ci"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup [ci service name]",
	Short: "Sets up CI service integration",
	Long: `Sets up CI service integration.

Use this command to modify CI service specific settings
like API key and project name.`,
	Example: "ks ci setup",
	Args:    cobra.MaximumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		ctx := core.New(core.CTX_RESOLVE)

		var serviceName string
		var ciService ci.CiService
		var found bool
		var err error

		if len(args) == 1 {
			serviceName = args[0]
			ciService, err = ci.GetCiService(serviceName, ctx, client.ApiURL)
			found = err == nil
		}

		if !found {
			ciService, err = SelectCiService(ctx)

			if err != nil {
				ui.PrintError(err.Error())
				os.Exit(1)
			}
		}

		if err = ciService.Setup().Error(); err != nil {
			ui.PrintError(ciService.Error().Error())
			os.Exit(1)
		}
		ui.PrintSuccess("Ci service setup successfully")
	},
}

func init() {
	ciCmd.AddCommand(setupCmd)

}
