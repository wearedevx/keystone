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
	"sort"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/config"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/client/auth"
	"github.com/wearedevx/keystone/cli/ui"
)

// memberCmd represents the member command
var memberCmd = &cobra.Command{
	Use:   "member",
	Args:  cobra.NoArgs,
	Short: "Manages members",
	Long: `Manages members.

Used without arguments, displays a list of all members,
grouped by their role.`,
	Run: func(_ *cobra.Command, _ []string) {

		c, kcErr := client.NewKeystoneClient()

		if kcErr != nil {
			fmt.Println(kcErr)
			os.Exit(1)
		}

		kf := keystonefile.KeystoneFile{}
		kf.Load(ctx.Wd)

		members, err := c.Project(kf.ProjectId).GetAllMembers()

		if err != nil {
			if errors.Is(err, auth.ErrorUnauthorized) {
				config.Logout()
				kserrors.InvalidConnectionToken(err).Print()
			} else {
				kserrors.UnkownError(err).Print()
			}
			os.Exit(1)
		}

		grouped := groupMembersByRole(members)

		for _, role := range getSortedRoles(grouped) {
			members := grouped[role]
			printRole(role, members)
		}
	},
}

func getSortedRoles(m map[models.Role][]models.ProjectMember) []models.Role {
	roles := make([]models.Role, 0)
	for r := range m {
		roles = append(roles, r)
	}

	s := models.NewRoleSorter(roles)
	return s.Sort()
}

func groupMembersByRole(pmembers []models.ProjectMember) map[models.Role][]models.ProjectMember {
	result := make(map[models.Role][]models.ProjectMember)

	for _, member := range pmembers {
		membersWithSameRole := result[member.Role]

		result[member.Role] = append(membersWithSameRole, member)
	}

	return result
}

func printRole(role models.Role, members []models.ProjectMember) {
	ui.Print("%s: %s", role.Name, role.Description)
	ui.Print("---")

	memberIDs := make([]string, len(members))
	for idx, member := range members {
		// FIXME: there should not be members with an empty UserID here
		// but still it happens in tests
		if member.User.UserID != "" {
			memberIDs[idx] = member.User.UserID
		}
	}

	sort.Strings(memberIDs)

	for _, member := range memberIDs {
		ui.Print("%s", member)
	}

	ui.Print("")
}

var envs []string

func init() {
	RootCmd.AddCommand(memberCmd)

	envs = []string{
		"dev",
		"staging",
		"prod",
	}

}

func isProjectOrganizationPaid(c client.KeystoneClient) bool {

	projectID := ctx.GetProjectID()
	organization, err := c.Project(projectID).GetProjectsOrganization()

	if err != nil {
		ui.PrintError(err.Error())
		os.Exit(1)
	}
	return organization.Paid
}
