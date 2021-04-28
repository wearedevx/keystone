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
	"fmt"
	"os"
	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/internal/config"
	"github.com/wearedevx/keystone/internal/errors"
	"github.com/wearedevx/keystone/pkg/client"
	core "github.com/wearedevx/keystone/pkg/core"
	"github.com/wearedevx/keystone/ui"
)

var forceYes bool

// memberRmCmd represents the memberRm command
var memberRmCmd = &cobra.Command{
	Use:   "rm <list of member ids>",
	Short: "Removes members from the current project",
	Long: `Removes members from the current project,
effecively preventing them from accessing future version
of the secrets and files.

This causes secrets to be re-crypted for the remainig members.`,
	Example: "ks member rm aster_23@github sam@gitlab",
	Args: func(cmd *cobra.Command, args []string) error {
		r := regexp.MustCompile("[\\w-_.]+@(gitlab|github)")

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
	Run: func(cmd *cobra.Command, args []string) {
		// Auth check
		account, index := config.GetCurrentAccount()
		token := config.GetAuthToken()

		if index < 0 {
			ui.Print(errors.MustBeLoggedIn(nil).Error())
		}

		ctx := core.New(core.CTX_RESOLVE)
		projectID := ctx.GetProjectID()

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

		c := client.NewKeystoneClient(account["user_id"], token)
		err := c.ProjectRemoveMembers(projectID, membersToRevoke)

		if err != nil {
			errors.CannotRemoveMembers(err).Print()
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
