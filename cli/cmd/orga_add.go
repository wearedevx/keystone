package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/spinner"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui/display"
)

var private bool

// addCmd represents the add command
var orgaAddCmd = &cobra.Command{
	Use:   "add <organization-name>",
	Short: "Creates a new organization",
	Long: `Creates a new oranization for your projects.
Organzation names must be unique and composed of only alphanumeric characters,
. (period), - (dash), and _ (underscore). No space are allowed
`,
	Example: `ks orga add my_new_organization`,
	Args:    cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		var err error
		organizationName := args[0]

		sp := spinner.Spinner(" ")
		sp.Start()

		c, err := client.NewKeystoneClient()
		exitIfErr(err)

		organization, err := c.Organizations().CreateOrganization(organizationName, private)
		if err != nil {
			handleClientError(err)
			exit(err)
		}

		display.OrganizationCreated(organization)
	},
}

func init() {
	orgaCmd.AddCommand(orgaAddCmd)
	orgaAddCmd.Flags().BoolVarP(&private, "private", "p", false, "Make the organization private.")
}
