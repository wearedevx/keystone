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

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/pkg/core"
)

// secretsCmd represents the secrets command
var secretsCmd = &cobra.Command{
	Use:   "secret",
	Short: "Manage secrets",
	Long: `Manage secrets.

Used without arguments, displays a table of secrets.`,
	Run: func(_ *cobra.Command, _ []string) {
		var err *errors.Error

		ctx := core.New(core.CTX_RESOLVE)
		ctx.MustHaveEnvironment(currentEnvironment)
		environments := ctx.ListEnvironments()

		secrets := ctx.ListSecrets()

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.SetStyle(table.StyleRounded)
		t.Style().Options.SeparateRows = true
		t.SetColumnConfigs([]table.ColumnConfig{
			{Number: 1, AutoMerge: true},
		})

		topHeader := table.Row{"Secret name"}
		envHeader := table.Row{""}

		for _, environment := range environments {
			topHeader = append(topHeader, "Environments")
			envHeader = append(envHeader, environment)
		}

		t.AppendHeader(topHeader, table.RowConfig{AutoMerge: true})
		t.AppendHeader(envHeader)

		for _, secret := range secrets {
			name := secret.Name

			if secret.Required {
				name = name + " *"
			}

			row := table.Row{name}

			for _, environment := range environments {
				value := secret.Values[core.EnvironmentName(environment)]

				row = append(row, value)
			}

			t.AppendRow(row)
		}

		t.Render()
		fmt.Println(" * Required secrets")

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
