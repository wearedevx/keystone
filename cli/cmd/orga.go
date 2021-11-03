package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/config"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui/display"
)

// orgaCmd represents the orga command
var orgaCmd = &cobra.Command{
	Use:   "orga",
	Short: "Manages organizations",
	Long: `Manages organizations.

Used without arguments, displays a list of all members related to the organization
the current project belongs to, grouped by their role.`,
	Run: func(_ *cobra.Command, _ []string) {
		var err error

		c, err := client.NewKeystoneClient()
		exitIfErr(err)

		organizations, err := c.Organizations().GetAll()
		if err != nil {
			handleClientError(err)
			exit(err)
		}

		currentUser, _ := config.GetCurrentAccount()

		display.Organizations(organizations, currentUser)
	},
}

func init() {
	RootCmd.AddCommand(orgaCmd)
}
