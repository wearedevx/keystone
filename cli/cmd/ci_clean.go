package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/ui"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Removes all secrets and files from a CI service",
	Long:  `Removes all secrets and files from a CI service.`,
	Example: `
# To remove everything regarding the current environment:
ks ci clean

# You can specify the target environment with the --env flag:
ks ci clean --env prod
`,
	Run: func(_ *cobra.Command, _ []string) {
		var err error
		ctx.MustHaveEnvironment(currentEnvironment)

		ciService, err := SelectCiService(ctx)
		if err != nil {
			exit(kserrors.CouldNotCleanService(err))
		}

		if err = ciService.
			CleanSecret(currentEnvironment).
			Error(); err != nil {
			if strings.Contains(ciService.Error().Error(), "404") == true {
				ui.PrintSuccess(fmt.Sprintf("No secret found for environment %s in CI service", currentEnvironment))
				os.Exit(0)
			}

			exit(kserrors.CouldNotCleanService(err))
		}

		ui.PrintSuccess(fmt.Sprintf("Secrets successfully removed from CI service, environment %s.", currentEnvironment))
	},
}

func init() {
	ciCmd.AddCommand(cleanCmd)
}
