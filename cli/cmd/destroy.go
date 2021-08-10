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

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/config"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy the whole Keystone project",
	Long: `Destroy the whole Keystone project.

The project will be deleted, members won’t be able to send or receive
updates about it. 
This is irreversible.
`,
	Example: "ks destroy",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := core.New(core.CTX_RESOLVE)
		c, kcErr := client.NewKeystoneClient()
		if kcErr != nil {
			fmt.Println(kcErr)
			os.Exit(1)
		}
		projectId := ctx.GetProjectID()
		projectService := c.Project(projectId)

		mustBeAdmin(projectService)

		projectName := ctx.GetProjectName()

		ui.Print(ui.RenderTemplate("confirm project destroy",
			`{{ CAREFUL }} You are about to destroy the {{ .Project }} project.
Secrets and files managed by Keystone WILL BE LOST. Make sure you have backups.

Members of the project will no longer be able to get the latest updates,
or share secrets between them.

This is permanent, and cannot be undone.
`, map[string]string{
				"Project": projectName,
			}))

		p := promptui.Prompt{
			Label: "Type the project name to confirm its destruction",
		}

		result, err := p.Run()
		if err != nil {
			kserrors.UnkownError(err).Print()
			os.Exit(1)
			return
		}

		// expect result to be the project name
		if projectName != result {
			kserrors.UnkownError(errors.New("Invalid Project Name")).Print()
			os.Exit(1)
			return
		}

		err = projectService.Destroy()

		if err != nil {
			kserrors.UnkownError(err).Print()
			os.Exit(1)
		}

		ctx.Destroy()

		ui.Print(ui.RenderTemplate("deletion ok",
			`{{ OK }} The project {{ .Project }} has successfully been destroyed.
Secrets and files are no longer accessible.
You may need to remove entries from your .gitignore file`,
			map[string]string{
				"Project": projectName,
			},
		))
	},
}

func mustBeAdmin(projectService *client.Project) {
	members, err := projectService.GetAllMembers()
	if err != nil {
		kserrors.UnkownError(err).Print()
		os.Exit(1)
	}

	account, _ := config.GetCurrentAccount()

	for _, member := range members {
		if member.User.UserID == account.UserID {
			if member.Role.Name == "admin" {
				return
			}
		}
	}

	kserrors.UnkownError(errors.New("Not allowed")).Print()
	os.Exit(1)
}

func init() {
	RootCmd.AddCommand(destroyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// destroyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// destroyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
