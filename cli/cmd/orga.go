package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui"
)

// orgaCmd represents the orga command
var orgaCmd = &cobra.Command{
	Use:   "orga",
	Short: "Manage organizations",
	Long: `Manages organizations.

Used without arguments, displays a list of all members,
grouped by their role.`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.NewKeystoneClient()
		if err != nil {
			err.Print()
			os.Exit(1)
		}

		organizations, _ := c.Organizations().GetAll()

		ui.Print("Organizations your are in:")
		ui.Print("---")
		for _, orga := range organizations {
			ui.Print(orga.Name)
		}

		ui.Print("")
	},
}

func init() {
	RootCmd.AddCommand(orgaCmd)
}
