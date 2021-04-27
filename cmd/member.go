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

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/internal/config"
	"github.com/wearedevx/keystone/internal/errors"
	"github.com/wearedevx/keystone/internal/keystonefile"

	. "github.com/wearedevx/keystone/internal/models"
	"github.com/wearedevx/keystone/pkg/client"
	"github.com/wearedevx/keystone/pkg/core"
	"github.com/wearedevx/keystone/ui"
)

// memberCmd represents the member command
var memberCmd = &cobra.Command{
	Use:   "member",
	Args:  cobra.NoArgs,
	Short: "Manage members",
	Long: `Manage members.

Used without arguments, displays a list of all members,
grouped by their role, with indication of their ownership.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := core.New(core.CTX_RESOLVE)

		currentUser, index := config.GetCurrentAccount()

		if index < 0 {
			ui.Print(errors.MustBeLoggedIn(nil).Error())
			os.Exit(1)
		}

		token := config.GetAuthToken()

		c := client.NewKeystoneClient(currentUser["user_id"], token)
		kf := keystonefile.KeystoneFile{}
		kf.Load(ctx.Wd)

		members, err := c.ProjectMembers(kf.ProjectId)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		grouped := groupMembersByEnv(members)

		for env, members := range grouped {
			printEnv(env, members)
		}
	},
}

func groupMembersByEnv(pmembers []ProjectMember) map[string][]ProjectMember {
	result := make(map[string][]ProjectMember)

	for _, env := range envs {
		members := make([]ProjectMember, 0)

		for _, member := range pmembers {
			if member.Environment.Name == env {
				members = append(members, member)
			}
		}

		result[env] = members
	}

	return result
}

func printEnv(env string, members []ProjectMember) {
	ui.Print(env)
	ui.Print("---")

	for _, member := range members {
		printMember(member)
	}

	ui.Print("")
}

func printMember(member ProjectMember) {
	ui.Print("%s (%s)", member.User.UserID, member.Role)
}

var envs []string

func init() {
	RootCmd.AddCommand(memberCmd)

	envs = []string{
		"dev",
		"staging",
		"prod",
		"ci",
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// memberCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// memberCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
