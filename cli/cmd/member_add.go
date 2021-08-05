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
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/api/pkg/models"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/spinner"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/client/auth"
	"github.com/wearedevx/keystone/cli/ui"
	"github.com/wearedevx/keystone/cli/ui/prompts"
	"gopkg.in/yaml.v2"
)

var membersFile string
var oneRole string
var manyMembers []string

type Flow int

const (
	PromptFlow = iota
	FileFlow
	ArgsFlow
)

var flow Flow

// memberAddCmd represents the memberAdd command
var memberAddCmd = &cobra.Command{
	Use: "add <member-id>...",
	Args: func(_ *cobra.Command, args []string) error {
		// r := regexp.MustCompile(`[\w-_.]+@(gitlab|github)`)
		flow = PromptFlow

		if len(args) == 0 && membersFile == "" && oneRole == "" && len(manyMembers) == 0 {
			return fmt.Errorf("missing arguments")
		}

		if membersFile != "" {
			flow = FileFlow
		} else if oneRole != "" && len(manyMembers) > 0 {
			flow = ArgsFlow
		} else {
			manyMembers = args
		}

		return nil
	},
	Short: "Adds members to the current project",
	Long: `Adds members to the current project.

Passed arguments are list member ids, which users can 
obtain using ks whoami.

This will cause secrets to be encryted for all members, existing and new.`,
	Example: `ks member add john.doe@gitlab danny54@github helena@gitlab
ks member add --from-file team.yml
ks member add -r developer -u john.doe@gitlab -u danny54@gitlab
`,
	Run: func(_ *cobra.Command, _ []string) {
		// Auth check
		projectID := ctx.GetProjectID()

		c, kcErr := client.NewKeystoneClient()

		if kcErr != nil {
			kcErr.Print()
			os.Exit(1)
		}

		var memberRoles map[string]models.Role

		switch flow {
		case FileFlow:
			memberRoles = getMemberRolesFromFile(c, membersFile)
		case ArgsFlow:
			memberRoles = getMemberRolesFromArgs(c, oneRole, manyMembers)
		default:
			memberRoles = getMemberRolesFromPrompt(c, manyMembers)
		}

		sp := spinner.Spinner(" ")
		sp.Start()

		err := c.Project(projectID).AddMembers(memberRoles)
		sp.Stop()

		if err != nil {
			if errors.Is(err, auth.ErrorUnauthorized) {
				kserrors.InvalidConnectionToken(err)
			} else {
				kserrors.CannotAddMembers(err).Print()
			}

			os.Exit(1)
		}

		ui.Print(ui.RenderTemplate("added members", `
{{ OK }} {{ "Members Added" | green }}
`, struct {
		}{}))
	},
}

func getMemberRolesFromFile(c client.KeystoneClient, filepath string) map[string]models.Role {
	memberRoleNames := make(map[string]string)

	dat, err := ioutil.ReadFile(filepath)
	if err != nil {
		ui.PrintError(err.Error())
		os.Exit(1)
	}

	if err = yaml.Unmarshal(dat, &memberRoleNames); err != nil {
		ui.PrintError(err.Error())
		os.Exit(1)
	}

	memberIDs := make([]string, 0)
	for m := range memberRoleNames {
		memberIDs = append(memberIDs, m)
	}

	mustMembersExist(c, memberIDs)
	roles := mustGetRoles(c)

	memberRoles := mapRoleNamesToRoles(memberRoleNames, roles)

	return memberRoles
}

func getMemberRolesFromArgs(c client.KeystoneClient, roleName string, memberIDs []string) map[string]models.Role {
	mustMembersExist(c, memberIDs)
	roles := mustGetRoles(c)
	foundRole := &models.Role{}

	for _, role := range roles {
		if role.Name == roleName {
			*foundRole = role
		}
	}

	if foundRole == nil {
		roleNames := []string{}

		for _, role := range roles {
			roleNames = append(roleNames, role.Name)
		}

		kserrors.RoleDoesNotExist(roleName, strings.Join(roleNames, ", "), nil).Print()
		os.Exit(1)
	}

	memberRoles := make(map[string]models.Role)

	for _, member := range memberIDs {
		memberRoles[member] = *foundRole
	}

	return memberRoles
}

func getMemberRolesFromPrompt(c client.KeystoneClient, memberIDs []string) map[string]models.Role {
	mustMembersExist(c, memberIDs)
	roles := mustGetRoles(c)

	memberRole := make(map[string]models.Role)

	for _, memberId := range memberIDs {
		role, err := prompts.PromptRole(memberId, roles)

		if err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		memberRole[memberId] = role
	}

	return memberRole
}

func mustMembersExist(c client.KeystoneClient, memberIDs []string) {
	r, err := c.Users().CheckUsersExist(memberIDs)
	if err != nil {
		// The HTTP request must have failed
		kserrors.UnkownError(err).Print()
		os.Exit(1)
	}

	if r.Error != "" {
		kserrors.UsersDontExist(r.Error, nil).Print()
		os.Exit(1)
	}
}

func mustGetRoles(c client.KeystoneClient) []models.Role {
	roles, err := c.Roles().GetAll()
	if err != nil {
		ui.PrintError(err.Error())
		os.Exit(1)
	}

	return roles
}

func mapRoleNamesToRoles(memberRoleNames map[string]string, roles []models.Role) map[string]models.Role {
	memberRoles := make(map[string]models.Role)

	for member, roleName := range memberRoleNames {
		var foundRole *models.Role
		for _, role := range roles {
			if role.Name == roleName {
				*foundRole = role

				break
			}
		}

		if foundRole == nil {
			roleNames := []string{}

			for _, role := range roles {
				roleNames = append(roleNames, role.Name)
			}

			kserrors.RoleDoesNotExist(roleName, strings.Join(roleNames, ", "), nil).Print()
			os.Exit(1)
		}

		memberRoles[member] = *foundRole
	}

	return memberRoles
}

func init() {
	memberCmd.AddCommand(memberAddCmd)

	memberAddCmd.Flags().StringVar(&membersFile, "from-file", "", "yml file to import a known list of members")

	memberAddCmd.Flags().StringVarP(&oneRole, "role", "r", "", "role to set users, when not using the prompt")
	memberAddCmd.Flags().StringSliceVarP(&manyMembers, "user", "u", []string{}, "user to add, when not using the prompt")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// memberAddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// memberAddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
