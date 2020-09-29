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
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/wearedevx/keystone/internal/errors"
	core "github.com/wearedevx/keystone/pkg/core"
	. "github.com/wearedevx/keystone/ui"

	"github.com/spf13/cobra"
)

var addOptional bool = false

// secretsAddCmd represents the set command
var secretsAddCmd = &cobra.Command{
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
	Run: func(cmd *cobra.Command, args []string) {
		var err *errors.Error
		secretName, secretValue := args[0], args[1]
		environmentValueMap := make(map[string]string)

		ctx := core.New(core.CTX_RESOLVE)
		ctx.MustHaveEnvironment(currentEnvironment)

		environments := ctx.ListEnvironments()

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}

		environmentValueMap[currentEnvironment] = secretValue

		Print(RenderTemplate("ask new value for environment", `
Enter a values for {{ . }}:`, secretName))

		for _, environment := range environments {

			p := promptui.Prompt{
				Label:   environment,
				Default: secretValue,
			}

			result, _ := p.Run()

			environmentValueMap[environment] = strings.Trim(result, " ")
		}

		flag := core.S_REQUIRED

		if addOptional {
			flag = core.S_OPTIONAL
		}

		if err = ctx.AddSecret(secretName, environmentValueMap, flag).Err(); err != nil {
			err.Print()
			return
		}

		PrintSuccess("Variable '%s' is set for %d environment(s)", secretName, len(environments))
	},
}

func init() {
	secretsCmd.AddCommand(secretsAddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	secretsAddCmd.Flags().BoolVarP(&addOptional, "optional", "o", false, "mark the secret as optional")
}