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
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/ui/display"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

// addCmd represents the add command
var hookAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a hook",
	Long: `Adds a hook.

It must be executable.
It will receive the project UUID and the path to the .keystone folder as
parameters.

The it will be run when environments change.`,
	Args:    cobra.ExactArgs(1),
	Example: "ks hook add backup-secrets.sh",
	Run: func(_ *cobra.Command, args []string) {
		if hook, ok := ctx.GetHook(); ok {
			if !prompts.ConfirmHookOverwrite(hook) {
				exit(nil)
			}
		}

		command := args[0]
		asAbsolutePath, err := filepath.Abs(command)
		exitIfErr(err)

		if !utils.FileExists(asAbsolutePath) {
			display.HookPathDoesNotExist(asAbsolutePath)
			exit(nil)
		}

		ctx.AddHook(asAbsolutePath)

		display.HookAddedSuccessfully()
	},
}

func init() {
	hookCmd.AddCommand(hookAddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
