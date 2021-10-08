package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/ci"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "edit [ci service name]",
	Short: "Configures an existing CI service integration",
	Long: `Configures an existing CI service integration.

Use this command to modify CI service specific settings
like API key and project name.`,
	Example: `ks ci edit

# To avoid the prompt
ks ci edit my-gitub-ci-service`,
	Args: cobra.MaximumNArgs(1),
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
				kserrors.CouldNotChangeService(serviceName, err).Print()
				os.Exit(1)
			}
		}

		if err = ciService.Setup().Error(); err != nil {
			kserrors.CouldNotChangeService(serviceName, err).Print()
			os.Exit(1)
		}
		ui.PrintSuccess("CI service setup successfully")
	},
}

func init() {
	ciCmd.AddCommand(setupCmd)

}
