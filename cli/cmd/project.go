package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/config"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/client/auth"
	"github.com/wearedevx/keystone/cli/ui"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "List project you are in",
	Long:  `List project you are in`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		c, err := client.NewKeystoneClient()
		exitIfErr(err)

		kf := keystonefile.KeystoneFile{}
		exitIfErr(kf.Load(ctx.Wd).Err())

		projects, err := c.Project("").GetAll()

		if err != nil {
			if errors.Is(err, auth.ErrorUnauthorized) {
				config.Logout()
				err = kserrors.InvalidConnectionToken(err)
			} else {
				err = kserrors.UnkownError(err)
			}
			exit(err)
		}

		ui.Print("You are part of %d project(s):", len(projects))

		fmt.Println()
		for _, project := range projects {
			fmt.Printf(" - %s, created at %s\n", project.Name, project.CreatedAt.Format("2006/01/02"))
		}

	},
}

func init() {
	RootCmd.AddCommand(projectCmd)
}
