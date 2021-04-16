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
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/wearedevx/keystone/internal/errors"
	core "github.com/wearedevx/keystone/pkg/core"
	. "github.com/wearedevx/keystone/ui"

	"github.com/spf13/cobra"
)

// envRmCmd represents the remove command
var envRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Removes and environment",
	Long: `Permanently removes an environment.

Every secret and tracked file content will be lost.
This is permanent and cannot be undone.

Example:
  $ ks env rm temp
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err *errors.Error

		ctx := core.New(core.CTX_RESOLVE)

		envName := args[0]

		currentEnvironment := ctx.CurrentEnvironment()

		if currentEnvironment == envName {
			err = errors.CannotRemoveCurrentEnvironment(currentEnvironment, nil)
			err.Print()
			return
		}

		promptResult := "n"

		if !skipPrompts {
			title := fmt.Sprintf("You are about to remove the '%s' environment", envName)
			Print(RenderTemplate("remove confirm", `
{{ CAREFUL }} {{ . | yellow }}
The data for the environment will be lost.
This can not be undone.
`, title))

			p := promptui.Prompt{
				Label:     "Continue",
				IsConfirm: true,
				Default:   "n",
			}

			promptResult, _ = p.Run()
		}

		if skipPrompts || promptResult == "y" {
			ctx.RemoveEnvironment(envName)

			if err = ctx.Err(); err != nil {
				err.Print()
				return
			}

			Print(RenderTemplate("env removed", `
{{ OK }} {{ .Message | bright_green }}
`, map[string]string{
				"Message": fmt.Sprintf("The '%s' environment has been removed", envName),
				"EnvName": envName,
			}))
		}
	},
}

func init() {
	envCmd.AddCommand(envRmCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RootCmd.PersistentFlags().BoolVarP(&skipPrompts, "y", "y", false, "Skip confirm and say yes")
}
