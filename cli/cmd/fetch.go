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

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/config"
	"github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Get remote modifications.",
	Long: `Get remote modifications.
Get info from your team:
  $ ks fetch
`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var err *errors.Error

		account, index := config.GetCurrentAccount()

		if index < 0 {
			ui.Print(errors.MustBeLoggedIn(nil).Error())
		}
		token := config.GetAuthToken()

		ctx := core.New(core.CTX_RESOLVE)

		ctx.MustHaveEnvironment(currentEnvironment)

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}

		ctx.MustHaveProject()

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}

		projectID := ctx.GetProjectID()

		c := client.NewKeystoneClient(account["user_id"], token)

		fmt.Println("Fetching new data...")
		result, _ := c.Messages().GetMessages(projectID)

		ctx.SaveMessages(result)

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}

		for environmentName, environment := range result.Environments {
			messageID := environment.Message.ID
			if messageID != 0 {
				fmt.Println("Environment", environmentName, "updated")
				response, _ := c.Messages().DeleteMessage(environment.Message.ID)
				if !response.Success {
					fmt.Println("Can't delete message", response.Error)
				}
			} else {
				environmentChanged := ctx.EnvironmentVersionHasChanged(environmentName, environment.VersionID)
				if environmentChanged {
					fmt.Println("Environment", environmentName, "has changed but no message available. Ask someone to push their secret.")
				} else {
					fmt.Println("Environment", environmentName, "up to date")
				}
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(fetchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fetchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
