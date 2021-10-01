package cmd

import (
	"errors"
	"fmt"
	"os"

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
		c, kcErr := client.NewKeystoneClient()

		if kcErr != nil {
			fmt.Println(kcErr)
			os.Exit(1)
		}

		kf := keystonefile.KeystoneFile{}
		kf.Load(ctx.Wd)

		projects, err := c.Project("").GetAll()

		if err != nil {
			if errors.Is(err, auth.ErrorUnauthorized) {
				config.Logout()
				kserrors.InvalidConnectionToken(err).Print()
			} else {
				kserrors.UnkownError(err).Print()
			}
			os.Exit(1)
		}
		ui.Print("You are part of %d project(s):", len(projects))

		fmt.Println()
		for _, project := range projects {
			fmt.Printf("  - %s, created at %s\n", project.Name, project.CreatedAt.Format("2006/01/02"))
		}

	},
}

func init() {
	RootCmd.AddCommand(projectCmd)
}
