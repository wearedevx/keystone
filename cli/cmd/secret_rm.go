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
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/messages"
	"github.com/wearedevx/keystone/cli/ui"
)

var purgeSecret bool

// secretsRmCmd represents the unset command
var secretsRmCmd = &cobra.Command{
	Use:   "rm <secret name>",
	Short: "Removes a secret from all environments",
	Long: `Removes a secret from all environments.

Removes the given secret from all environments.
`,
	Example: "ks rm PORT",
	Args:    cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		var err *kserrors.Error
		secretName := args[0]

		ctx.MustHaveEnvironment(currentEnvironment)

		if !ctx.HasSecret(secretName) && !purgeSecret {
			kserrors.SecretDoesNotExist(secretName, nil).Print()
			return
		}

		ms := messages.NewMessageService(ctx)
		changes := mustFetchMessages(ms)

		if err = ctx.
			CompareRemovedSecretWithChanges(secretName, changes).
			RemoveSecret(secretName, purgeSecret).
			Err(); err != nil {
			err.Print()
			os.Exit(1)
			return
		}

		if purgeSecret {
			if err := ms.SendEnvironments(ctx.AccessibleEnvironments).Err(); err != nil {
				err.Print()
				os.Exit(1)
			}
		}

		ui.PrintSuccess("Variable '%s' removed", secretName)
	},
}

func init() {
	secretsCmd.AddCommand(secretsRmCmd)

	secretsRmCmd.Flags().BoolVarP(&purgeSecret, "purge", "p", false, "purge all values from all environments aswell")
}
