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
var orgaProjectCmd = &cobra.Command{
	Use:   "project <organization-name>",
	Short: "Lists projects from organization",
	Long: `Lists projects from organization.
`,
	Example: "ks orga project my_organization",
	Args:    cobra.ExactArgs(1),
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

		projects, err := c.Organizations().GetProjects(orga)

		if err != nil {
			if errors.Is(err, auth.ErrorUnauthorized) {
				config.Logout()
				err = kserrors.InvalidConnectionToken(err)
			} else {
				err = kserrors.UnkownError(err)
			}
			exit(err)
		}

		ui.Print("You have access to %d project(s) in this organization :", len(projects))

		fmt.Println()
		for _, project := range projects {
			fmt.Printf("  - %s\n", project.Name)
		}
	},
}

func init() {
	orgaCmd.AddCommand(orgaProjectCmd)
}
