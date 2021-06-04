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
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/errors"
	core "github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"

	"github.com/spf13/cobra"
)

var addOptional bool = false

// secretAddCmd represents the set command
var secretAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a secret to all environments",
	Long: `Adds a secret to all environments.

Secrets are environment variables which value may vary
across environments, such as 'staging', 'prduction',
and 'development' environments.

The varible name will be added to all such environments,
you will be asked its value for each environment

Example:
  Add an environment variable PORT to all environments
  and set its value to 3000 for the current one.
  $ ks set PORT 3000

`,
	Args: cobra.ExactArgs(2),
	Run: func(_ *cobra.Command, args []string) {
		var err *errors.Error
		secretName, secretValue := args[0], args[1]
		environmentValueMap := make(map[string]string)

		ctx := core.New(core.CTX_RESOLVE)

		ctx.MustHaveEnvironment(currentEnvironment)

		environments := make([]string, 0)
		accessibleEnvironments := ctx.GetAccessibleEnvironments()

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}

		for _, environment := range accessibleEnvironments {
			environments = append(environments, environment.Name)
		}

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}

		// Ask value for each env
		if !skipPrompts {
			ui.Print(ui.RenderTemplate("ask new value for environment", `
Enter a values for {{ . }}:`, secretName))

			for _, environment := range environments {

				p := promptui.Prompt{
					Label:   environment,
					Default: secretValue,
				}

				result, err := p.Run()

				// Handle user cancelation
				// or prompt error
				if err != nil {
					if err.Error() != "^C" {
						ui.PrintError(err.Error())
						os.Exit(1)
					}
					os.Exit(0)
				}

				environmentValueMap[environment] = strings.Trim(result, " ")
			}

		} else {
			for _, environment := range environments {
				environmentValueMap[environment] = strings.Trim(secretValue, " ")
			}

		}

		environmentValueMap[currentEnvironment] = secretValue

		// Fetch new messages to see if added secret has changed
		messagesByEnvironment := &models.GetMessageByEnvironmentResponse{
			Environments: map[string]models.GetMessageResponse{},
		}

		fmt.Println("Syncing data...")
		fetchErr := ctx.FetchNewMessages(messagesByEnvironment)

		if fetchErr != nil {
			ui.PrintError(fetchErr.Error())
			return
		}

		changes, writeErr := ctx.WriteNewMessages(*messagesByEnvironment)

		if writeErr != nil {
			writeErr.Print()
			// err.SetCause(writeErr)
			// err.Print()
			return
		}

		if err = ctx.CompareNewSecretWithChanges(secretName, environmentValueMap, changes); err != nil {
			err.Print()
			return
		}

		flag := core.S_REQUIRED

		if addOptional {
			flag = core.S_OPTIONAL
		}

		if err = ctx.AddSecret(secretName, environmentValueMap, flag).Err(); err != nil {
			err.Print()
			return
		}

		// TODO
		// Format beautyiful error
		if pushErr := ctx.PushEnv(accessibleEnvironments); pushErr != nil {
			ui.PrintError(pushErr.Error())
			return
		}

		ui.PrintSuccess("Variable '%s' is set for %d environment(s)", secretName, len(environments))
	},
}

func init() {
	secretsCmd.AddCommand(secretAddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	secretAddCmd.Flags().BoolVarP(&addOptional, "optional", "o", false, "mark the secret as optional")
}
