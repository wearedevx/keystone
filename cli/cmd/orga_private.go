package cmd

import (
	"github.com/spf13/cobra"

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

		sp := spinner.Spinner("")
		sp.Start()

		c, err := client.NewKeystoneClient()
		exitIfErr(err)

		foundOrga, err := c.Organizations().GetByName(organizationName, client.OwnedOnly)
		if err != nil {
			handleClientError(err)
			exit(kserrors.YouDoNotOwnTheOrganization(organizationName, err))
		}

		foundOrga.Private = !foundOrga.Private

		organization, err := c.Organizations().UpdateOrganization(foundOrga)
		exitIfErr(err)

		display.OrganizationStatusUpdate(organization)
	},
}

func init() {
	orgaCmd.AddCommand(privateCmd)
}
