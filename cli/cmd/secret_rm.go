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
	"github.com/wearedevx/keystone/cli/ui"
)

var purge bool

// secretsRmCmd represents the unset command
var secretsRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Removes a secret from all environments",
	Long: `Removes a secret from all environments.

Removes the given secret from all environments.

Exemple:
  $ ks rm PORT`,
	Args: cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		var err *errors.Error
		secretName := args[0]

		ctx.MustHaveEnvironment(currentEnvironment)

		if !ctx.HasSecret(secretName) {
			errors.SecretDoesNotExist(secretName, nil).Print()
			return
		}

		var printer = &ui.UiPrinter{}
		ms := messages.NewMessageService(ctx, printer)

		changes := ms.GetMessages()

		if err = ms.Err(); err != nil {
			err.Print()
			os.Exit(1)
		}

		if err = ctx.CompareRemovedSecretWithChanges(secretName, changes); err != nil {
			err.Print()
			os.Exit(1)
			return
		}

		ctx.RemoveSecret(secretName)

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}

		// Unlike most of the other commands, we do not need
		// to send the environment to other users
		// because the only thing that needs to change is
		// the keystone.yml file

		ui.PrintSuccess("Variable '%s' removed", secretName)
	},
}

func init() {
	secretsCmd.AddCommand(secretsRmCmd)

	secretsRmCmd.Flags().BoolVarP(&purge, "purge", "p", false, "purge all values from all environments aswell")
}
