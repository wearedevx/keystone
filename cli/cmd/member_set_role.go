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
	"strings"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

var memberId string
var roleName string

// memberSetRoleCmd represents the memberSetRole command
var memberSetRoleCmd = &cobra.Command{
	Use:   "set-role <member id> [role]",
	Short: "Sets the role for a member",
	Long: `Sets the role for a member.
If no role argument is provided, it will be prompted.

Roles determine access rights to environments.`,
	Example: `ks member set-role john@gitlab devops

ks member set-role sandra@github`,
	Args: func(_ *cobra.Command, args []string) error {
		r := regexp.MustCompile(`[\w-_.]+@(gitlab|github)`)
		argc := len(args)

		if argc == 0 || argc > 2 {
			return fmt.Errorf("invalid number of arguments. Expected 1 or 2, got %d", argc)
		}

		if argc >= 1 {
			memberId = args[0]
		}

		if argc == 2 {
			roleName = args[1]
		}

		if !r.Match([]byte(memberId)) {
			return fmt.Errorf("invalid member id: %s", memberId)
		}

		return nil
	},
	Run: func(_ *cobra.Command, _ []string) {
		// Auth check
		c, kcErr := client.NewKeystoneClient()

		if kcErr != nil {
			kcErr.Print()
			os.Exit(1)
		}

		projectID := ctx.GetProjectID()
		// Ensure member exists
		r, err := c.Users().CheckUsersExist([]string{memberId})
		if r.Error != "" || err != nil {
			errors.UsersDontExist(r.Error, err).Print()
			os.Exit(1)
		}

		// Get all roles, te make sure the role exists
		// And to be able to list them in the prompt
		roles, err := c.Roles().GetAll()
		if err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		// If user didnot provide a role,
		// prompt it
		if roleName == "" {
			r, err := prompts.PromptRole(memberId, roles)
			if err != nil {
				ui.PrintError(err.Error())
				os.Exit(1)
			}

			roleName = r.Name
		}

		// If the role exsists, do the work
		if _, ok := getRoleWithName(roleName, roles); ok {
			err = c.Project(projectID).SetMemberRole(memberId, roleName)

			if err != nil {
				ui.PrintError(err.Error())
				os.Exit(1)
			}
		} else {
			roleNames := make([]string, len(roles))
			for i, r := range roles {
				roleNames[i] = r.Name
			}

			errors.RoleDoesNotExist(roleName, strings.Join(roleNames, ", "), nil).Print()
			os.Exit(1)
		}

		ui.Print(ui.RenderTemplate("set role ok", `
{{ OK }} {{ "Role set" | green }}
`, struct {
		}{}))
	},
}

func getRoleWithName(roleName string, roles []models.Role) (models.Role, bool) {
	found := false
	var role models.Role

	for _, existingRole := range roles {
		if existingRole.Name == roleName {
			found = true
			role = existingRole
			break
		}
	}

	return role, found
}

func init() {
	memberCmd.AddCommand(memberSetRoleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// memberSetRoleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// memberSetRoleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
