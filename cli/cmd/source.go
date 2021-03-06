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
	"os"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/messages"

	"github.com/wearedevx/keystone/cli/internal/utils"
	core "github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
)

// sourceCmd represents the source command
var sourceCmd = &cobra.Command{
	Use:   "source",
	Short: "Prints all environment variables",
	Long: `Prints all environment variables.

Environment variables values can then be loaded using eval, for example.

Example:
  $ ks source
  KEY=value
  OTHER_KEY=other_value

  $ eval "$(ks source)"
  $ echo $KEY
  value
  $ echo $OTHER_KEY
  other_value`,
	Run: func(_ *cobra.Command, _ []string) {
		var err *errors.Error

		ctx.MustHaveEnvironment(currentEnvironment)

		var printer = &ui.EchoPrinter{}

		ms := messages.NewMessageService(ctx, printer)
		ms.GetMessages()

		env := ctx.ListSecrets()
		// TODO: crash on missing required files
		ctx.FilesUseEnvironment(currentEnvironment)

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}

		for _, secretInfo := range env {
			value := secretInfo.Values[core.EnvironmentName(currentEnvironment)]

			checkErr := utils.CheckSecretContent(secretInfo.Name)
			if checkErr != nil {
				ui.Print(`echo %s`, checkErr.Error())
				os.Exit(1)
			}
			if secretInfo.Required && value == "" {
				ui.Print("echo \"Error: secret '%s' is required, but value is missing\"", secretInfo.Name)
				// make the eval crash in such situation
				os.Exit(1)
			}

			ui.Print(`export %s="%s"`, secretInfo.Name, value)

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
