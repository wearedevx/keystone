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
	"github.com/spf13/cobra"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui/display"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroys the whole Keystone project",
	Long: `Destroys the whole Keystone project.

The project will be deleted, members won’t be able to send nor receive
updates about it. 

All secrets and files managed by Keystone *WILL BE LOST*.
It is highly recommended that you backup everything up beforehand.

This is irreversible.
`,
	Example: "ks destroy",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		c, kcErr := client.NewKeystoneClient()
		exitIfErr(kcErr)

		projectId := ctx.GetProjectID()
		projectService := c.Project(projectId)

		mustBeAdmin(projectService)

		projectName := ctx.GetProjectName()

		if !prompts.ConfirmProjectDestruction(projectName) {
			exit(kserrors.NameDoesNotMatch(nil))
		}

		if err = projectService.Destroy(); err != nil {
			handleClientError(err)
			exit(err) // if handle hasn't
		}

		if err = ctx.Destroy(); err != nil {
			exit(kserrors.CouldNotRemoveLocalFiles(err))
		}

		display.DeletionSuccess(projectName)
	},
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
