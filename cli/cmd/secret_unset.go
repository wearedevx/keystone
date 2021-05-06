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
	"github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/pkg/core"
	. "github.com/wearedevx/keystone/cli/ui"
)

// secretsUnsetCmd represents the unset command
var secretsUnsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Clears a secret for the current environment",
	Long: `Clears a secret for the current environment.

Other environments will not be afftected.
The secret must not be required.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err *errors.Error

		ctx := core.New(core.CTX_RESOLVE)
		ctx.MustHaveEnvironment(currentEnvironment)

		secretName := args[0]

		if !ctx.HasSecret(secretName) {
			errors.SecretDoesNotExist(secretName, nil).Print()
			return
		}

		if ctx.SecretIsRequired(secretName) {
			errors.SecretRequired(secretName, nil).Print()
			return
		}

		ctx.UnsetSecret(currentEnvironment, secretName)

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}

		PrintSuccess(fmt.Sprintf("Secret '%s' updated for the '%s' environment", secretName, currentEnvironment))
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
