/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/ci"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "Configures a new CI service",
	Long:    `Configures a new CI service.`,
	Example: `ks ci add`,
	Run: func(_ *cobra.Command, _ []string) {
		ctx := core.New(core.CTX_RESOLVE)

		serviceName := prompts.StringInput(
			"Enter a name for your integration",
			"",
		)

		if _, nameExists := ci.FindCiServiceWithName(ctx, serviceName); nameExists {
			// TODO: add a Ci service named {{.Name}} already exists
			ui.PrintError(fmt.Sprintf(
				"A CI service named %s already exists",
				serviceName,
			))
		}

		ciService, err := ci.PickCiService(serviceName, ctx, client.ApiURL)

		if err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		if err = ciService.Setup().Error(); err != nil {
			ui.PrintError(ciService.Error().Error())
			os.Exit(1)
		}

		if err = ci.AddCiService(ctx, ciService); err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		ui.PrintSuccess("CI service added successfully")
	},
}

func init() {
	ciCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
