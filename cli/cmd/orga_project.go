package cmd

import (
	"github.com/spf13/cobra"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui/display"
)

// projectCmd represents the project command
var orgaProjectCmd = &cobra.Command{
	Use:   "project <organization-name>",
	Short: "Lists projects from organization",
	Long: `Lists projects from organization.
`,
	Example: "ks orga project my_organization",
	Args:    cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		var err error
		orgaName := args[0]
		c, err := client.NewKeystoneClient()
		exitIfErr(err)

		kf := keystonefile.KeystoneFile{}
		exitIfErr(kf.Load(ctx.Wd).Err())

		// FIXME: couldn't we merge thes two calls ?
		orga, err := c.Organizations().GetByName(orgaName, client.ALL_KNWON)
		if err != nil {
			handleClientError(err)
			exit(kserrors.OrganizationDoesNotExist(err))
		}

		projects, err := c.Organizations().GetProjects(orga)
		if err != nil {
			handleClientError(err)
			exit(err)
		}

		display.OrganizationAccessibleProjects(projects)
	},
}

func init() {
	orgaCmd.AddCommand(orgaProjectCmd)
}
