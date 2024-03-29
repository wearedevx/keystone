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
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/messages"
	"github.com/wearedevx/keystone/cli/ui/display"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Sends all secrets and files to all members",
	Long: `Sends all secrets and files to all members.

Members will receive all secrets and files values for all the environments
they have access to.
`,
	Example: `ks env send`,
	Args:    cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		ctx.MustHaveEnvironment(currentEnvironment)

		environments := ctx.AccessibleEnvironments

		for i, env := range environments {
			localEnvironment := ctx.LoadEnvironmentsFile().GetByName(env.Name)

			environments[i] = models.Environment{
				Name:          localEnvironment.Name,
				VersionID:     localEnvironment.VersionID,
				EnvironmentID: localEnvironment.EnvironmentID,
			}
		}

		ms := messages.NewMessageService(ctx)

		exitIfErr(
			ms.SendEnvironments(environments).Err(),
		)

		display.EnvironmentSendSuccess()
	},
}

func init() {
	envCmd.AddCommand(sendCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sendCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sendCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
