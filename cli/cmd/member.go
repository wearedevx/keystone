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

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/config"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/ui/display"

	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/client/auth"
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
		exitIfErr(kcErr)

		kf := keystonefile.KeystoneFile{}
		exitIfErr(kf.Load(ctx.Wd).Err())

		members, err := c.Project(kf.ProjectId).GetAllMembers()

		if err != nil {
			if errors.Is(err, auth.ErrorUnauthorized) {
				config.Logout()
				err = kserrors.InvalidConnectionToken(err)
			} else {
				err = kserrors.UnkownError(err)
			}

			exit(err)
		}

		display.MembersByRole(members)
	},
}

func init() {
	RootCmd.AddCommand(memberCmd)
}
