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
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/ui/display"
)

func pathList(files []keystonefile.FileKey) []string {
	r := make([]string, len(files))

	for index, file := range files {
		r[index] = file.Path
	}

	return r
}

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch <environment>",
	Short: "Changes the current environment",
	Long: `Changes the current envrionment.

Next time ` + "`" + `ks source` + "+" + ` is executed, it will use values
from <environment>.

Valid values for environment are: "dev", "staging", and "prod"`,
	Example: `ks env switch prod`,
	Args:    cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		fetchMessages()

		locallyModified := ctx.LocallyModifiedFiles(currentEnvironment)
		if len(locallyModified) != 0 {
			exit(
				kserrors.YouHaveLocallyModifiedFiles(
					currentEnvironment,
					pathList(locallyModified),
					nil,
				),
			)
		}

		// Set the current environment
		environmentName := args[0]
		exitIfErr(ctx.
			MustHaveAccessToEnvironment(environmentName).
			SetCurrent(environmentName).
			Err())

		display.EnvironmentUsing(environmentName)
	},
}

func init() {
	envCmd.AddCommand(switchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// switchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// switchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
