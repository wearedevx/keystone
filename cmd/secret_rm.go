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
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/internal/errors"
	core "github.com/wearedevx/keystone/pkg/core"
	. "github.com/wearedevx/keystone/ui"
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
	Run: func(cmd *cobra.Command, args []string) {
		var err *errors.Error
		secretName := args[0]

		ctx := core.New(core.CTX_RESOLVE).RemoveSecret(secretName)

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}

		PrintSuccess("Variable '%s' unset for all environments", secretName)
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
