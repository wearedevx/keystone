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

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/pkg/core"
	. "github.com/wearedevx/keystone/cli/ui"
)

// requireCmd represents the require command
var requireCmd = &cobra.Command{
	Use:   "require",
	Short: "Marks a secret as required",
	Long: `Marks a secret as required.

Secrets marked as required cannot be unset or set to blank value.
If they are, 'ks source' will exit with a non-zero exit code.
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err *errors.Error

		ctx := core.New(core.CTX_RESOLVE)
		secretName := args[0]

		if !ctx.HasSecret(secretName) {
			errors.SecretDoesNotExist(secretName, nil).Print()
			return
		}

		// Check for blank values in environments
		environments := ctx.ListEnvironments()

		secret := ctx.GetSecret(secretName)

		for _, environment := range environments {
			value := string(secret.Values[core.EnvironmentName(environment)])

			for len(value) == 0 {
				Print("Enter the value of '%s' for the '%s' environment", secretName, environment)

				p := promptui.Prompt{
					Label: secretName,
				}

				result, _ := p.Run()
				value = result
			}

			ctx.SetSecret(environment, secretName, value)
		}

		// All is OK, set is as optional
		ctx.MarkSecretRequired(secretName, true)

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}

		PrintSuccess(fmt.Sprintf("Secret '%s' is now required.", secretName))
	},
}

func init() {
	secretsCmd.AddCommand(requireCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// requireCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// requireCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
