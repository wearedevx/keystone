package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/spinner"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui"
)

// addCmd represents the add command
var orgaAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Create a new organization",
	Long:  `Create a new oranization for your projects`,
	Run: func(cmd *cobra.Command, args []string) {
		argc := len(args)
		if argc == 0 || argc > 1 {
			ui.PrintError(fmt.Sprintf("invalid number of arguments. Expected 1, got %d", argc))
			os.Exit(1)
		}

		organizationName := args[0]

		sp := spinner.Spinner(" ")
		sp.Start()

		c, err := client.NewKeystoneClient()
		if err != nil {
			err.Print()
			os.Exit(1)
		}

		organization, createErr := c.Organizations().CreateOrganization(organizationName)

		if createErr != nil {
			ui.PrintError(createErr.Error())
			os.Exit(1)
		}
		ui.PrintSuccess("Organization %s has been created", organization.Name)
	},
}

func init() {
	orgaCmd.AddCommand(orgaAddCmd)
}
