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

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/config"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/members"
	"github.com/wearedevx/keystone/cli/internal/spinner"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/client/auth"
	"github.com/wearedevx/keystone/cli/ui/display"
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

		roles := mustGetRoles(c)

		var memberRoles map[string]models.Role

		switch flow {
		case FileFlow:
			memberRoles, err = members.GetMemberRolesFromFile(
				c,
				membersFile,
				roles,
			)
		case ArgsFlow:
			memberRoles, err = members.GetMemberRolesFromArgs(
				c,
				oneRole,
				manyMembers,
				roles,
			)
		default:
			memberRoles, err = members.GetMemberRolesFromPrompt(
				c,
				manyMembers,
				roles,
			)
		}
		exitIfErr(err)

		sp := spinner.Spinner("").Start()

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

		display.MembersAdded()
	},
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
