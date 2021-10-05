package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/spinner"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui"
)

// privateCmd represents the private command
var privateCmd = &cobra.Command{
	Use:   "private",
	Short: "Toggle an organization privacy",
	Long:  `Toggle an organization privacy.`,
	Run: func(cmd *cobra.Command, args []string) {
		argc := len(args)
		if argc != 1 {
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

		foundOrga := getUserOwnedOrganization(c, organizationName)

		foundOrga.Private = !foundOrga.Private

		organization, updateErr := c.Organizations().UpdateOrganization(foundOrga)

		if updateErr != nil {
			ui.PrintError(updateErr.Error())
			os.Exit(1)
		}
		if organization.Private {
			ui.PrintSuccess("Organization %s is now private", organization.Name)
		} else {
			ui.PrintSuccess("Organization %s is now not private", organization.Name)
		}
	},
}

func init() {
	orgaCmd.AddCommand(privateCmd)
}

func getUserOwnedOrganization(c client.KeystoneClient, organizationName string) models.Organization {
	organizations, _ := c.Organizations().GetAll()

	foundOrga := models.Organization{}

	for _, orga := range organizations {
		if orga.Name == organizationName {
			foundOrga = orga
		}
	}

	if foundOrga.ID == 0 {
		ui.PrintError("You don't own an organization named %s", organizationName)
		ui.Print("To see organizations you own, use : $ ks orga")
		os.Exit(1)
	}
	return foundOrga
}
