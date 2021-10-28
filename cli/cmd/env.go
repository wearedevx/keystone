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
	"github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"

	"github.com/spf13/cobra"
)

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env [environment]",
	Short: "Manages environments",
	Long: `Manages environments.

Displays a list of available environments:
` + "```" + `
$ ks env
 * dev
   staging
   prod
` + "```" + `

With an argument name, activates the environment:
` + "```" + `
$ ks env staging
` + "```" + `
`,
	Args: cobra.NoArgs,
	Run: func(_ *cobra.Command, args []string) {
		var err *errors.Error

		ctx.MustHaveEnvironment(currentEnvironment)

		// If no argument is given show a list of environments
		if len(args) == 0 {
			listEnv(ctx, err)
			return
		}

	},
}

type EnvListViewModel struct {
	Current      string
	Environments []string
}

// Prints an environment list
// The current environment is marked with an asterisk
// TODO: move this to the ui package
func listEnv(ctx *core.Context, _ *errors.Error) {
	if quietOutput {
		ui.Print(currentEnvironment)
		return
	}

	environments := ctx.ListEnvironments()
	exitIfErr(ctx.Err())

	template := `{{ range  .Environments }}
{{ if eq . $.Current }} {{ "*" | blue }} {{ . | yellow }} {{ else }}   {{ . }} {{ end }} {{ end }}`

	ui.Print(ui.RenderTemplate("list env", template, EnvListViewModel{
		Current:      currentEnvironment,
		Environments: environments,
	}))
}

func init() {
	RootCmd.AddCommand(envCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// envCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// envCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
