/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"os"
	"path"
	"strings"

	kerrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	. "github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/pkg/client"
	core "github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"

	"github.com/spf13/cobra"
)

var projectName string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [project name]",
	Short: "Creates Keystone config files and directories",
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("a project name cannot be empty")
		}
		return nil
	},
	Long: `Creates Keystone config files and directories.

Created files and directories:
 - keystone.yml: the project's config,
 - .keystone:    cache and various files for internal use. 
                 automatically added to .gitignore
`,
	Run: func(_ *cobra.Command, args []string) {
		var err *kerrors.Error
		projectName = strings.Join(args, " ")

		// Retrieve working directry
		currentfolder, osError := os.Getwd()
		ctx := core.New(core.CTX_INIT)

		if osError != nil {
			err = kerrors.NewError("OS Error", "Error when retrieving working directory", map[string]string{}, osError)
			err.Print()
			os.Exit(1)
		}

		// Ask for project name if keystone file doesn't exist.
		if !keystonefile.ExistsKeystoneFile(currentfolder) {
			c, err := client.NewKeystoneClient()
			if err != nil {
				err.Print()
				os.Exit(1)
			}

			project, initErr := c.Project("").Init(projectName)

			if initErr != nil {
				ui.PrintError(initErr.Error())
				os.Exit(1)
			}

			if err = ctx.Init(project).Err(); err != nil {
				err.Print()
				os.Exit(1)
			}

			ui.Print(ui.RenderTemplate("Init Success", `
{{ .Message | box | bright_green | indent 2 }}

{{ .Text | bright_black | indent 2 }}`, map[string]string{
				"Message": "All done!",
				"Text": `You can start adding environment variable with:
  $ ks secret add VARIABLE value

Load them with:
  $ eval $(ks source)

If you need help with anything:
  $ ks help [command]

`,
			}))
		} else {
			if ctx.GetProjectName() != projectName {
				// check if .keystone directory too
				if DirExists(path.Join(ctx.Wd, ".keystone")) {
					kerrors.AlreadyKeystoneProject(errors.New("")).Print()
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initCmd.Flags().StringVar(&projectName, "project", "", "Define the project name")
}
