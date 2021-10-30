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
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/ui/display"
)

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env [environment]",
	Short: "Manages environments",
	Long: `Manages environments.

Displays a list of available environments:
` + "```" + `
$ ks env
 * dev
   staging
   prod
` + "```" + `

With an argument name, activates the environment:
` + "```" + `
$ ks env staging
` + "```" + `
`,
	Args: cobra.NoArgs,
	Run: func(_ *cobra.Command, args []string) {
		ctx.MustHaveEnvironment(currentEnvironment)

		if quietOutput {
			display.Environment(currentEnvironment)
		} else {
			environments := ctx.ListEnvironments()
			exitIfErr(ctx.Err())

			display.EnvironmentList(environments, currentEnvironment)
		}
	},
}

func init() {
	RootCmd.AddCommand(envCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// envCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// envCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
