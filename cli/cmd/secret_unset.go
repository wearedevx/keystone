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
	"os"

	"github.com/spf13/cobra"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/messages"
	"github.com/wearedevx/keystone/cli/ui"
)

// secretsUnsetCmd represents the unset command
var secretsUnsetCmd = &cobra.Command{
	Use:   "unset <secret name>",
	Short: "Clears a secret for the current environment",
	Long: `Clears a secret for the current environment.

Other environments will not be affected.
The secret must not be required.`,
	Example: "ks unset PORT",
	Args:    cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		var err *kserrors.Error

		ctx.MustHaveEnvironment(currentEnvironment)

		secretName := args[0]

		if !ctx.HasSecret(secretName) {
			kserrors.SecretDoesNotExist(secretName, nil).Print()
			return
		}

		if ctx.SecretIsRequired(secretName) {
			kserrors.SecretRequired(secretName, nil).Print()
			return
		}

		var printer = &ui.UiPrinter{}
		ms := messages.NewMessageService(ctx, printer)

		changes := ms.GetMessages()

		if err = ms.Err(); err != nil {
			err.Print()
			os.Exit(1)
		}

		if err = ctx.
			CompareRemovedSecretWithChanges(secretName, changes).
			UnsetSecret(currentEnvironment, secretName).
			Err(); err != nil {
			err.Print()
			os.Exit(1)
			return
		}

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}

		if err := ms.SendEnvironments(ctx.AccessibleEnvironments).Err(); err != nil {
			err.Print()
			os.Exit(1)
			return
		}

		ui.PrintSuccess(fmt.Sprintf("Secret '%s' updated for the '%s' environment", secretName, currentEnvironment))
	},
}

func init() {
	secretsCmd.AddCommand(secretsUnsetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// unsetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// unsetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
