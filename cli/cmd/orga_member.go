package cmd

import (
	"errors"
	"fmt"
	"os"

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
	Run: func(cmd *cobra.Command, args []string) {

		argc := len(args)
		if argc == 0 || argc > 1 {
			ui.PrintError(fmt.Sprintf("invalid number of arguments. Expected 1, got %d", argc))
			os.Exit(1)
		}
		orgaName := args[0]
		c, kcErr := client.NewKeystoneClient()

		if kcErr != nil {
			fmt.Println(kcErr)
			os.Exit(1)
		}

		kf := keystonefile.KeystoneFile{}
		kf.Load(ctx.Wd)

		organizations, err := c.Organizations().GetAll()
		orga := models.Organization{}

		for _, organization := range organizations {
			if organization.Name == orgaName {
				orga = organization
			}
		}

		if orga.ID == 0 {
			ui.PrintError("Organization does not exist")
			os.Exit(1)
		}

		members, err := c.Organizations().GetMembers(orga)

		if err != nil {
			if errors.Is(err, auth.ErrorUnauthorized) {
				config.Logout()
				kserrors.InvalidConnectionToken(err).Print()
			} else {
				kserrors.UnkownError(err).Print()
			}
			os.Exit(1)
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
