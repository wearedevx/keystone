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
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/ui"
)

// optionalCmd represents the optional command
var optionalCmd = &cobra.Command{
	Use:   "optional <secret name>",
	Short: "Marks a secret as optional",
	Long: `Marks a secret as optional.

When a file is marked as optional, its absence or emptiness won’t cause
` + "`" + `ks source` + "`" + ` or ` + "`" + `ks ci send` + "`" + ` to fail.
`,
	Example: "ks secret optional PORT",
	Args:    cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		var err *kserrors.Error

		secretName := args[0]

		if !ctx.HasSecret(secretName) {
			kserrors.SecretDoesNotExist(secretName, nil).Print()
			return
		}

		ctx.MarkSecretRequired(secretName, false)

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}

		template := `Secret {{ .SecertName }} is now optional.`

		ui.Print(ui.RenderTemplate("set secret optional", template, struct{ SecretName string }{
			SecretName: secretName,
		}))
	},
}

func init() {
	secretsCmd.AddCommand(optionalCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// optionalCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// optionalCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
