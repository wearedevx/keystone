package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/ci"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui/display"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "edit [ci service name]",
	Short: "Configures an existing CI service integration",
	Long: `Configures an existing CI service integration.

Use this command to modify CI service specific settings,
like API key and project name.`,
	Example: `ks ci edit

# To avoid the prompt
ks ci edit my-gitub-ci-service`,
	Args: cobra.MaximumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
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
			ciService, err = ci.SelectCiServiceConfiguration(
				serviceName,
				ctx,
				client.ApiURL,
			)
			if err != nil {
				if errors.Is(err, ci.ErrorNoCIServices) {
					exit(kserrors.NoCIServices(nil))
				} else {
					exit(kserrors.CouldNotChangeService(serviceName, err))
				}
			}
		}

		if err = ciService.Setup().Error(); err != nil {
			exit(kserrors.CouldNotChangeService(serviceName, err))
		}

		display.CiServiceSetupSuccessfully()
	},
}

func init() {
	ciCmd.AddCommand(setupCmd)
}
