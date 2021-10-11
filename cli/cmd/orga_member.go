package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/config"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/client/auth"
	"github.com/wearedevx/keystone/cli/ui"
)

// projectCmd represents the project command
var orgaMemberCmd = &cobra.Command{
	Use:   "member",
	Short: "List members of an organization",
	Long:  `List members of an organization`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		orgaName := args[0]
		c, err := client.NewKeystoneClient()
		exitIfErr(err)

		kf := keystonefile.KeystoneFile{}
		exitIfErr(kf.Load(ctx.Wd).Err())

		organizations, err := c.Organizations().GetAll()
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
			if errors.Is(err, auth.ErrorUnauthorized) {
				config.Logout()
				exit(kserrors.InvalidConnectionToken(err))
			} else {
				exit(kserrors.UnkownError(err))
			}
		}
		ui.Print("%d members are in projects that belong to this organization:", len(members))

		fmt.Println()
		for _, member := range members {
			fmt.Printf("  - %s\n", member.User.UserID)
		}
	},
}

func init() {
	orgaCmd.AddCommand(orgaMemberCmd)
}
