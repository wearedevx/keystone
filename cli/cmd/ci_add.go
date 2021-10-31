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
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/ci"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Configures a new CI service",
	Long: `Configures a new CI service.

Once you have configured a new CI service, you can send it secrets using:
` + "`" + `ks ci send --env prod` + "`\n",
	Example: `ks ci add`,
	Args:    cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		serviceName := prompts.ServiceIntegrationName()

		if _, nameExists := ci.FindCiServiceWithName(ctx, serviceName); nameExists {
			exit(kserrors.CiServiceAlreadyExists(serviceName, nil))
		}

		ciService, err := ci.PickCiService(serviceName, ctx, client.ApiURL)
		if err != nil {
			exit(kserrors.CouldNotAddService(serviceName, err))
		}

		if err = ciService.Setup().Error(); err != nil {
			exit(kserrors.CouldNotAddService(serviceName, err))
		}

		if err = ci.AddCiService(ctx, ciService); err != nil {
			exit(kserrors.CouldNotAddService(serviceName, err))
		}

		display.CiAdded()
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
