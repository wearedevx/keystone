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
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui/display"
)

// secretsCmd represents the secrets command
var secretsCmd = &cobra.Command{
	Use:   "secret",
	Short: "Manages secrets",
	Long: `Manages secrets.

Used without arguments, displays a table of secrets.`,
	Run: func(_ *cobra.Command, _ []string) {
		ctx.MustHaveEnvironment(currentEnvironment)
		environments := ctx.ListEnvironments()

		shouldFetchMessages()

		secrets := ctx.ListSecrets()
		secretsFromCache := ctx.ListSecretsFromCache()
		secretsFromCache = core.FilterSecretsFromCache(
			secretsFromCache,
			secrets,
		)
		secrets = append(secrets, secretsFromCache...)

		exitIfErr(ctx.Err())

		display.SecretTable(secrets, environments)
	},
}

func init() {
	RootCmd.AddCommand(secretsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// secretsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// secretsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
