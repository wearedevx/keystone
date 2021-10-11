package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/spinner"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui"
)

var private bool

// addCmd represents the add command
var orgaAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Create a new organization",
	Long:  `Create a new oranization for your projects`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		organizationName := args[0]

		sp := spinner.Spinner(" ")
		sp.Start()

		c, err := client.NewKeystoneClient()
		exitIfErr(err)

		organization, err := c.Organizations().CreateOrganization(organizationName, private)
		exitIfErr(err)

		ui.PrintSuccess("Organization %s has been created", organization.Name)
	},
}

func init() {
	orgaCmd.AddCommand(orgaAddCmd)
	orgaAddCmd.Flags().BoolVarP(&private, "private", "p", false, "Make the organization private.")
}
