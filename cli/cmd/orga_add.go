package cmd

import (
	"os"

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
		organizationName := args[0]

		sp := spinner.Spinner(" ")
		sp.Start()

		c, err := client.NewKeystoneClient()
		if err != nil {
			err.Print()
			os.Exit(1)
		}

		organization, createErr := c.Organizations().CreateOrganization(organizationName, private)

		if createErr != nil {
			ui.PrintError(createErr.Error())
			os.Exit(1)
		}

		ui.PrintSuccess("Organization %s has been created", organization.Name)
	},
}

func init() {
	orgaCmd.AddCommand(orgaAddCmd)
	orgaAddCmd.Flags().BoolVarP(&private, "private", "p", false, "Make the organization private.")
}
