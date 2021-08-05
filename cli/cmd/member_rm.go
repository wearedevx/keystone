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
	"os"
	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/spinner"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/client/auth"
	"github.com/wearedevx/keystone/cli/ui"
)

var forceYes bool

// memberRmCmd represents the memberRm command
var memberRmCmd = &cobra.Command{
	Use:   "rm <member-id>...",
	Short: "Removes members from the current project",
	Long: `Removes members from the current project,
effecively preventing them from accessing future version
of the secrets and files.

This causes secrets to be re-crypted for the remainig members.`,
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
		// Auth check

		c, kcErr := client.NewKeystoneClient()

		if kcErr != nil {
			kcErr.Print()
			os.Exit(1)
		}
		projectID := ctx.GetProjectID()

		sp := spinner.Spinner(" Checking users exist...")
		sp.Start()
		r, err := c.Users().CheckUsersExist(args)
		sp.Stop()

		if err != nil {
			if errors.Is(err, auth.ErrorUnauthorized) {
				kserrors.InvalidConnectionToken(err)
			} else {
				kserrors.UnkownError(err)
			}
			os.Exit(1)
		}

		if r.Error != "" {
			kserrors.UsersDontExist(r.Error, nil).Print()

			os.Exit(1)
		}

		membersToRevoke := make([]string, 0)

		for _, memberId := range args {
			result := "y"

			if !forceYes {
				prompt := promptui.Prompt{
					Label: "Revoke access to " + memberId + "? [y/n]",
				}

				result, _ = prompt.Run()
			}

			if result == "y" {
				membersToRevoke = append(membersToRevoke, memberId)
			}
		}

		if len(membersToRevoke) == 0 {
			os.Exit(0)
		}

		sp = spinner.Spinner(" Removing members...")
		sp.Start()
		err = c.Project(projectID).RemoveMembers(membersToRevoke)
		sp.Stop()

		if err != nil {
			kserrors.CannotRemoveMembers(err).Print()
			os.Exit(1)
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
