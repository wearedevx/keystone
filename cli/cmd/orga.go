package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/config"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui"
)

// orgaCmd represents the orga command
var orgaCmd = &cobra.Command{
	Use:   "orga",
	Short: "Manage organizations",
	Long: `Manages organizations.

Used without arguments, displays a list of all members related to the organization
the current project belongs to, grouped by their role.`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.NewKeystoneClient()
		exitIfErr(err)

		organizations, _ := c.Organizations().GetAll()

		ui.Print("Organizations your are in:")
		ui.Print("---")
		currentUser, _ := config.GetCurrentAccount()

		for _, orga := range organizations {
			orgaString := orga.Name
			if orga.User.UserID == currentUser.UserID {
				orgaString += " 👑"
			}
			if orga.Private {
				orgaString += " P"
			}
			ui.Print(orgaString)
		}

		ui.Print("")
		ui.Print(" 👑 : You own; P : private")
	},
}

func init() {
	RootCmd.AddCommand(orgaCmd)
}
