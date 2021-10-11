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
	"github.com/wearedevx/keystone/api/pkg/models"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/messages"
	"github.com/wearedevx/keystone/cli/ui"
)

// secretsSetCmd represents the set command
var secretsSetCmd = &cobra.Command{
	Use:   "set <secret name> [secret value]",
	Short: "Updates a secret’s value for the current environment",
	Long: `Updates a secret’s value for the current environment.

Changes the value of a secret without altering other environments.
`,
	Example: `ks secret set PORT 3000

# Change the value of PORT for the 'staging' environment:
ks --env staging PORT 4545
`,
	Args: cobra.ExactArgs(2),
	Run: func(_ *cobra.Command, args []string) {
		ctx.MustHaveEnvironment(currentEnvironment)

		secretName := args[0]
		secretValue := args[1]

		if !ctx.HasSecret(secretName) {
			exit(kserrors.SecretDoesNotExist(secretName, nil))
		}

		ms := messages.NewMessageService(ctx)
		mustFetchMessages(ms)

		ctx.SetSecret(currentEnvironment, secretName, secretValue)
		exitIfErr(ctx.Err())

		localEnvironment := ctx.LoadEnvironmentsFile().GetByName(currentEnvironment)
		environment := []models.Environment{{
			Name:          localEnvironment.Name,
			VersionID:     localEnvironment.VersionID,
			EnvironmentID: localEnvironment.EnvironmentID,
		}}

		exitIfErr(
			ms.SendEnvironments(environment).Err(),
		)

		ui.PrintSuccess(fmt.Sprintf("Secret '%s' updated for the '%s' environment", secretName, currentEnvironment))
	},
}

func init() {
	secretsCmd.AddCommand(secretsSetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
