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

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/config"

	"github.com/wearedevx/keystone/cli/internal/utils"
	core "github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
)

// sourceCmd represents the source command
var sourceCmd = &cobra.Command{
	Use:   "source",
	Short: "Echo a script to load secrets, and writes files",
	Long: `Echo a script to load secrets, and writes files.

Environment variables values can then be loaded using eval, for example.

Example:
` + "```" + `
$ ks source
export KEY="value"
export OTHER_KEY="other_value"

$ eval "$(ks source)"
$ echo $KEY
value
$ echo $OTHER_KEY
other_value
` + "```" + `
`,
	Run: func(_ *cobra.Command, _ []string) {
		ctx.MustHaveEnvironment(currentEnvironment)

		if config.IsLoggedIn() {
			shouldFetchMessages()
		}

		env := ctx.ListSecrets()

		exitIfErr(ctx.
			FilesUseEnvironment(
				currentEnvironment,
				currentEnvironment,
				core.CTX_KEEP_LOCAL_FILES,
			).
			Err())

		mustNotHaveAnyRequiredThingMissing(ctx)

		for _, secretInfo := range env {
			value := secretInfo.Values[core.EnvironmentName(currentEnvironment)]

			exitIfErr(
				utils.CheckSecretContent(secretInfo.Name),
			)

			if secretInfo.Required && value == "" {
				// make the eval crash in such situation
				exit(fmt.Errorf(
					"secret '%s' is required, but value is missing",
					secretInfo.Name,
				))
			}

			escapedValue := utils.DoubleQuoteEscape(string(value))

			ui.Print(`export %s="%s"`, secretInfo.Name, escapedValue)

		}
	},
}

func init() {
	RootCmd.AddCommand(sourceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sourceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sourceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
