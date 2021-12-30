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
	"github.com/wearedevx/keystone/cli/ui/display"
)

// hookCmd represents the hook command
var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "Manages hook",
	Long: `Manages hook.
Used without arguments nor parameters, shows the currently registered hook.

A hook is a command or shell script that gets executed every time a change
is detected among an environment secrets or files, and every time changes
you've made are sent to other project members.

It receives the project UUID and the path to the .keystone folder as parameters.

The hook is global, meaning that once you've set it up, it will run for every
project.
It is also unique (there can only be one hook).`,
	Example: "ks hook",
	Run: func(cmd *cobra.Command, args []string) {
		if hook, ok := ctx.GetHook(); ok {
			display.HookCommand(hook.Command)
		} else {
			display.ThereIsNoHookYet()
		}
	},
}

func init() {
	RootCmd.AddCommand(hookCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hookCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hookCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
