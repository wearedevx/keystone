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

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/config"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/spinner"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/client/auth"
	"github.com/wearedevx/keystone/cli/ui"
	"github.com/wearedevx/keystone/cli/ui/prompts"
	"gopkg.in/yaml.v2"
)

var (
	membersFile string
	oneRole     string
	manyMembers []string
)

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

		if len(args) == 0 && membersFile == "" && oneRole == "" &&
			len(manyMembers) == 0 {
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

Passed arguments are a list of member ids, which users can 
obtain using ` + "`" + `ks whoami` + "`" + `.

Added members will be able to access secrets after a ` + "`" + `ks source` + "`" + `.
`,
	Example: `# Add a list of members:
ks member add john.doe@gitlab danny54@github helena@gitlab

# Add members and defining their roles from the command line:
ks member add -r developer -u john.doe@gitlab -u danny54@gitlab

# Add members with their roles from a file:
ks member add --from-file team.yaml
`,
	Run: func(_ *cobra.Command, _ []string) {
		var err error
		// Auth check
		projectID := ctx.GetProjectID()

		c, err := client.NewKeystoneClient()
		exitIfErr(err)

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

		err = c.Project(projectID).AddMembers(memberRoles)
		sp.Stop()

		if err != nil {
			if errors.Is(err, auth.ErrorUnauthorized) {
				config.Logout()
				err = kserrors.InvalidConnectionToken(err)
			} else {
				err = kserrors.CannotAddMembers(err)
			}

			exit(err)
		}

		ui.Print(ui.RenderTemplate("added members", `
{{ OK }} {{ "Members Added" | green }}

To send secrets and files to new member, use "member add" command.
  $ ks member send-env --all-env <member-id>
`, struct{}{}))
	},
}

func getMemberRolesFromFile(
	c client.KeystoneClient,
	filepath string,
) map[string]models.Role {
	var err error
	memberRoleNames := make(map[string]string)

	/* #nosec
	 * the file is going to be parsed, not executed in anyway
	 */
	dat, err := ioutil.ReadFile(filepath)
	exitIfErr(err)

	err = yaml.Unmarshal(dat, &memberRoleNames)
	exitIfErr(err)

	memberIDs := make([]string, 0)
	for m := range memberRoleNames {
		memberIDs = append(memberIDs, m)
	}

	mustMembersExist(c, memberIDs)
	roles := mustGetRoles(c)

	warningFreeOrga(roles)

	memberRoles := mapRoleNamesToRoles(memberRoleNames, roles)

	return memberRoles
}

func getMemberRolesFromArgs(
	c client.KeystoneClient,
	roleName string,
	memberIDs []string,
) map[string]models.Role {
	mustMembersExist(c, memberIDs)
	roles := mustGetRoles(c)
	foundRole := &models.Role{}

	warningFreeOrga(roles)

	for _, role := range roles {
		if role.Name == roleName {
			*foundRole = role
		}
	}

	memberRoles := make(map[string]models.Role)

	for _, member := range memberIDs {
		memberRoles[member] = *foundRole
	}

	return memberRoles
}

// TODO: to ui package
func getMemberRolesFromPrompt(
	c client.KeystoneClient,
	memberIDs []string,
) map[string]models.Role {
	mustMembersExist(c, memberIDs)
	roles := mustGetRoles(c)

	warningFreeOrga(roles)

	memberRole := make(map[string]models.Role)

	for _, memberId := range memberIDs {
		role, err := prompts.PromptRole(memberId, roles)
		exitIfErr(err)

		memberRole[memberId] = role
	}

	return memberRole
}

// TODO: to ui package
func mapRoleNamesToRoles(
	memberRoleNames map[string]string,
	roles []models.Role,
) map[string]models.Role {
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
			exit(
				kserrors.UnkownError(
					fmt.Errorf("role %s does not exist", roleName),
				),
			)
		}

		memberRoles[member] = *foundRole
	}

	return memberRoles
}

// TODO: to ui package
func warningFreeOrga(roles []models.Role) {
	if len(roles) == 1 {
		fmt.Fprintln(
			os.Stderr,
			"WARNING: You are not allowed to set role other than admin for free organization",
		)
		ui.Print("To learn more: https://keystone.sh")
		ui.Print("")
	}
}

func init() {
	memberCmd.AddCommand(memberAddCmd)

	memberAddCmd.Flags().
		StringVar(&membersFile, "from-file", "", "yaml file to import a known list of members")

	memberAddCmd.Flags().
		StringVarP(&oneRole, "role", "r", "", "role to set users, when not using the prompt")
	memberAddCmd.Flags().
		StringSliceVarP(&manyMembers, "user", "u", []string{}, "user to add, when not using the prompt")
}
