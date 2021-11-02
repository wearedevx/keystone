package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/api/pkg/models"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/spinner"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui/display"
)

// privateCmd represents the private command
var privateCmd = &cobra.Command{
	Use:   "private <organization-name>",
	Short: "Toggles an organization privacy",
	Long:  `Toggles an organization privacy.`,
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		var err error
		organizationName := args[0]

		sp := spinner.Spinner(" ")
		sp.Start()

		c, err := client.NewKeystoneClient()
		exitIfErr(err)

		foundOrga := getUserOwnedOrganization(c, organizationName)

		foundOrga.Private = !foundOrga.Private

		organization, err := c.Organizations().UpdateOrganization(foundOrga)
		exitIfErr(err)

		display.OrganizationStatusUpdate(organization)
	},
}

func init() {
	orgaCmd.AddCommand(privateCmd)
}

// TODO: This should be an API thing
func getUserOwnedOrganization(
	c client.KeystoneClient,
	organizationName string,
) models.Organization {
	organizations, _ := c.Organizations().GetAll()

	foundOrga := models.Organization{}

	for _, orga := range organizations {
		if orga.Name == organizationName {
			foundOrga = orga
		}
	}

	if foundOrga.ID == 0 {
		exit(kserrors.YouDoNotOwnTheOrganization(organizationName, nil))
	}
	return foundOrga
}
