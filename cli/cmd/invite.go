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
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/spinner"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui"
)

// inviteCmd represents the invite command
var inviteCmd = &cobra.Command{
	Use:   "invite",
	Short: "Sends an invitation to join Keystone",
	Long:  `Sends an invitation to join Keystone.`,
	Run: func(_ *cobra.Command, args []string) {
		emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

		argc := len(args)

		if argc == 0 || argc > 1 {
			ui.PrintError("invalid number of arguments. Expected 1 got %d", argc)
			os.Exit(1)
		}
		email := args[0]

		if !emailRegex.Match([]byte(email)) {
			ui.PrintError("invalid email address: %s", email)
			os.Exit(1)
		}
		c, kcErr := client.NewKeystoneClient()

		sp := spinner.Spinner("Inviting user")
		sp.Start()

		if kcErr != nil {
			sp.Stop()
			kcErr.Print()
			os.Exit(1)
		}

		err := c.Users().InviteUser(email)

		if err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}
		ui.PrintSuccess("A email has been sent to %s, they will get back to you when their Keystone account will be created", email)

	},
}

func init() {
	RootCmd.AddCommand(inviteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// inviteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// inviteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
