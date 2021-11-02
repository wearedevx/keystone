package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/api/pkg/models"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui/display"
)

// projectCmd represents the project command
var orgaMemberCmd = &cobra.Command{
	Use:   "member <organization-name>",
	Short: "Lists members of an organization",
	Long: `Lists members of an organization.
`,
	Args:    cobra.ExactArgs(1),
	Example: "ks orga member my_organization",
	Run: func(_ *cobra.Command, args []string) {
		var err error
		orgaName := args[0]
		c, err := client.NewKeystoneClient()
		exitIfErr(err)

		kf := keystonefile.KeystoneFile{}
		exitIfErr(kf.Load(ctx.Wd).Err())

		// TODO: There should be a GetByName method, no ?
		organizations, err := c.Organizations().GetAll()
		if err != nil {
			handleClientError(err)
			exit(err)
		}
		orga := models.Organization{}

		for _, organization := range organizations {
			if organization.Name == orgaName {
				orga = organization
			}
		}

		if orga.ID == 0 {
			exit(kserrors.OrganizationDoesNotExist(nil))
		}

		members, err := c.Organizations().GetMembers(orga)
		if err != nil {
			handleClientError(err)
			exit(err)
		}

		display.OrganizationMembers(members)
	},
}

func init() {
	orgaCmd.AddCommand(orgaMemberCmd)
}
