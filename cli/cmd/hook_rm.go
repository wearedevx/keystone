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

	"github.com/wearedevx/keystone/cli/internal/config"
	"github.com/wearedevx/keystone/cli/ui/display"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

// rmCmd represents the rm command
var hookRmCmd = &cobra.Command{
	Use:     "rm",
	Short:   "Removes the hook",
	Long:    `Removes the hook.`,
	Args:    cobra.NoArgs,
	Example: "ks hook rm",
	Run: func(_ *cobra.Command, _ []string) {
		if hook, ok := ctx.GetHook(); ok {
			if prompts.ConfirmHookRemoval(hook) {
				config.RemoveHook()
				config.Write()
			}
		} else {
			display.ThereIsNoHookYet()
		}
	},
}

func init() {
	hookCmd.AddCommand(hookRmCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
