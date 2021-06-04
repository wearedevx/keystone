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

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/api/pkg/models"
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

		c, kcErr := client.NewKeystoneClient()

		if kcErr != nil {
			kcErr.Print()
			os.Exit(1)
		}

		fmt.Println("Fetching new data...")

		messagesByEnvironment := &models.GetMessageByEnvironmentResponse{
			Environments: map[string]models.GetMessageResponse{},
		}

		fetchErr := ctx.FetchNewMessages(messagesByEnvironment)
		if fetchErr != nil {
			err.SetCause(fetchErr)
			err.Print()
		}

		_, writeErr := ctx.WriteNewMessages(*messagesByEnvironment)

		if writeErr != nil {
			writeErr.Print()
		}

		for _, msgEnv := range messagesByEnvironment.Environments {
			response, _ := c.Messages().DeleteMessage(msgEnv.Message.ID)

			if !response.Success {
				ui.Print("Can't delete message " + response.Error)
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
