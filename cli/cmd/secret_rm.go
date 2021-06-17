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
	"github.com/wearedevx/keystone/cli/internal/environments"
	"github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/messages"
	"github.com/wearedevx/keystone/cli/internal/utils"
	core "github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
)

// secretsRmCmd represents the unset command
var secretsRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Removes a secret from all environments",
	Long: `Removes a secret from all environments.

Removes the given secret from all environments.

Exemple:
  $ ks rmove PORT`,
	Args: cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		var err *errors.Error
		secretName := args[0]
		environmentValueMap := make(map[string]string)

		checkSecretErr := utils.CheckSecretContent(secretName)

		if checkSecretErr != nil {
			ui.PrintError(checkSecretErr.Error())
			os.Exit(1)
		}

		ctx := core.New(core.CTX_RESOLVE)

		ctx.MustHaveEnvironment(currentEnvironment)

		es := environments.NewEnvironmentService(ctx)
		if err = es.Err(); err != nil {
			err.Print()
			os.Exit(1)
		}

		environmentNames := make([]string, 0)
		accessibleEnvironments := es.GetAccessibleEnvironments()

		if err = es.Err(); err != nil {
			err.Print()
			return
		}

		for _, environment := range accessibleEnvironments {
			environmentNames = append(environmentNames, environment.Name)
		}

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}

		var printer = &ui.UiPrinter{}
		ms := messages.NewMessageService(ctx, printer)
		changes := ms.GetMessages()

		if err = ms.Err(); err != nil {
			err.Print()
			os.Exit(1)
		}

		if err = ctx.CompareNewSecretWithChanges(secretName, environmentValueMap, changes); err != nil {
			err.Print()
			os.Exit(1)
			return
		}

		ctx.RemoveSecret(secretName)

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}
		// TODO
		// Format beautyiful error
		if err := ms.SendEnvironments(accessibleEnvironments).Err(); err != nil {
			err.Print()
			os.Exit(1)
			return
		}

		ui.PrintSuccess("Variable '%s' unset for all environments", secretName)
	},
}

func init() {
	secretsCmd.AddCommand(secretsRmCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// unsetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// unsetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
