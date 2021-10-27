package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/ci"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Removes all secrets and files from a CI service",
	Long: `Removes all secrets and files from a CI service.

This does not remove the CI service from the projcet. To do so, use:
` + "`" + `ks ci rm <service-name>` + "`",
	Example: `
# To remove everything regarding the current environment:
ks ci clean

# You can specify the target environment with the --env flag:
ks ci clean --env prod
`,
	Run: func(_ *cobra.Command, _ []string) {
		var err error
		ctx.MustHaveEnvironment(currentEnvironment)

		ciService, err := ci.SelectCiServiceConfiguration(
			"",
			ctx,
			client.ApiURL,
		)

		if err != nil {
			if errors.Is(err, ci.ErrorNoCIServices) {
				exit(kserrors.NoCIServices(nil))
			} else {
				exit(kserrors.CouldNotCleanService(err))
			}
		}

		if err = ciService.
			CleanSecret(currentEnvironment).
			Error(); err != nil {
			if strings.Contains(ciService.Error().Error(), "404") == true {
				ui.PrintSuccess(fmt.Sprintf("No secret found for environment %s in CI service", currentEnvironment))
				exit(nil)
			}

			exit(kserrors.CouldNotCleanService(err))
		}

		ui.PrintSuccess(fmt.Sprintf("Secrets successfully removed from CI service, environment %s.", currentEnvironment))
	},
}

func init() {
	ciCmd.AddCommand(cleanCmd)
}
