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
	"github.com/wearedevx/keystone/cli/ui/display"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm [service name]",
	Short: "Removes a CI service configuration",
	Long: `Removes a CI service configuration.

` + "`" + `ks ci send` + "`" + ` will no longer send secrets and files to the service.
However, secrets and files sent before calling ` + "`" + `ks ci send` + "`" + ` will
not be cleaned from the service.`,
	Example: `ks ci rm

# To avoid the prompt
ks ci rm my-github-ci-service
`,
	Args: cobra.MaximumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		serviceName := getServiceNameToRemove(args)

		s, ok := ci.FindCiServiceWithName(ctx, serviceName)
		if !ok {
			exit(kserrors.NoSuchService(serviceName, nil))
		}

		if prompts.ConfirmCiConfigurationRemoval(s.Name) {
			if err := ci.RemoveCiService(ctx, s.Name); err != nil {
				exit(kserrors.CouldNotRemoveService(err))
			}
		}

		display.CiServiceRemoved(s.Name)
	},
}

func getServiceNameToRemove(args []string) (serviceName string) {
	if len(args) == 1 {
		serviceName = args[0]
	} else {
		serviceName = prompts.ServiceConfigurationToRemove()
	}

	return serviceName
}

func init() {
	ciCmd.AddCommand(rmCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
