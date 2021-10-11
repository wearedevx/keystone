/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/config"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/spinner"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/client/auth"
	"github.com/wearedevx/keystone/cli/ui"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

var forceYes bool

// memberRmCmd represents the memberRm command
var memberRmCmd = &cobra.Command{
	Use:   "rm <member-id>...",
	Short: "Removes members from the current project",
	Long: `Removes members from the current project,
effectively preventing them from accessing future version
of the secrets and files.
`,
	Example: "ks member rm aster_23@github sam@gitlab",
	Args: func(_ *cobra.Command, args []string) error {
		r := regexp.MustCompile(`[\w-_.]+@(gitlab|github)`)

		if len(args) == 0 {
			return fmt.Errorf("missing member id")
		}

		for _, memberId := range args {
			if !r.Match([]byte(memberId)) {
				return fmt.Errorf("invalid member id: %s", memberId)
			}
		}

		return nil
	},
	Run: func(_ *cobra.Command, args []string) {
		var err error
		c, err := client.NewKeystoneClient()
		exitIfErr(err)

		projectID := ctx.GetProjectID()

		sp := spinner.Spinner(" Checking users exist...")
		sp.Start()
		r, err := c.Users().CheckUsersExist(args)
		sp.Stop()

		if err != nil {
			if errors.Is(err, auth.ErrorUnauthorized) {
				config.Logout()
				err = kserrors.InvalidConnectionToken(err)
			} else {
				err = kserrors.UnkownError(err)
			}
			exit(err)
		}

		if r.Error != "" {
			exit(kserrors.UsersDontExist(r.Error, nil))
		}

		membersToRevoke := make([]string, 0)

		for _, memberId := range args {
			revoke := true

			if !forceYes {
				revoke = prompts.Confirm("Revoke access to " + memberId)
			}

			if revoke {
				membersToRevoke = append(membersToRevoke, memberId)
			}
		}

		if len(membersToRevoke) == 0 {
			exit(nil)
		}

		sp = spinner.Spinner(" Removing members...")
		sp.Start()
		err = c.Project(projectID).RemoveMembers(membersToRevoke)
		sp.Stop()

		if err != nil {
			exit(kserrors.CannotRemoveMembers(err))
		}

		ui.Print(ui.RenderTemplate("removed members", `
{{ OK }} {{ "Revoked Access To Members" | green }}
`, nil))
	},
}

func init() {
	memberCmd.AddCommand(memberRmCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// memberRmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// memberRmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	memberRmCmd.Flags().BoolVarP(&forceYes, "yes", "y", false, "skip prompt say yes to all")
}
