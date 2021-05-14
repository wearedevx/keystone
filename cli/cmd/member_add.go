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

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/pkg/client"
	core "github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

var membersFile string

// memberAddCmd represents the memberAdd command
var memberAddCmd = &cobra.Command{
	Use: "add <list of member ids>",
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
	Short: "Add members to the current project",
	Long: `Add members to the current project.

Passed arguments are list member ids, which users can 
obtain using ks whoami.

This will cause secrets to be encryted for all members, existing and new.`,
	Example: `ks member add john.doe@gitlab danny54@github helena@gitlab`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get keystone key.
		c, kcErr := client.NewKeystoneClient()

		if kcErr != nil {
			kcErr.Print()
			os.Exit(1)
		}

		// Read Roles from config
		ctx := core.New(core.CTX_RESOLVE)
		projectID := ctx.GetProjectID()

		r, err := c.Users().CheckUsersExist(args)

		if r.Error != "" {
			errors.UsersDontExist(r.Error, nil).Print()
			os.Exit(1)
		}

		roles, err := c.Roles().GetAll()

		if err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		memberRole := make(map[string]models.Role)

		for _, memberId := range args {
			role, err := prompts.PromptRole(memberId, roles)

			if err != nil {
				// TODO: Handle error
				fmt.Println(err)

				os.Exit(1)
			}

			memberRole[memberId] = role
		}

		err = c.Project(projectID).AddMembers(memberRole)

		if err != nil {
			errors.CannotAddMembers(err).Print()
			os.Exit(1)
		}

		ui.Print(ui.RenderTemplate("added members", `
{{ OK }} {{ "Members Added" | green }}
`, struct {
		}{}))
	},
}

func init() {
	memberCmd.AddCommand(memberAddCmd)

	memberAddCmd.Flags().StringVar(&membersFile, "from-file", "", "yml file to import a known list of members")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// memberAddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// memberAddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
