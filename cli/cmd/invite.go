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
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/spinner"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui"
)

// inviteCmd represents the invite command
var inviteCmd = &cobra.Command{
	Use:   "invite <email address>",
	Short: "Sends an invitation to join Keystone",
	Long:  `Sends an invitation to join Keystone.`,
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

		email := args[0]

		if !emailRegex.Match([]byte(email)) {
			exit(fmt.Errorf("invalid email address: %s", email))
		}

		c, kcErr := client.NewKeystoneClient()
		exitIfErr(kcErr)

		sp := spinner.Spinner("Inviting user")
		sp.Start()

		projectName := ctx.GetProjectName()

		result, err := c.Users().InviteUser(email, projectName)
		exitIfErr(err)

		if len(result.UserUIDs) > 0 {

			ui.Print(ui.RenderTemplate("file add success", `
{{ OK }} {{ .Title | green }}

The email is associated with a Keystone account. They are registered as: {{ .Usernames | bright_green }}.

To add them to the project use "member add" command:
  $ ks member add <username>
`, map[string]string{
				"Title":     "User already on Keystone",
				"Usernames": fmt.Sprintf("%s", strings.Join(result.UserUIDs, ", ")),
			}))
		} else {
			ui.PrintSuccess("A email has been sent to %s, they will get back to you when their Keystone account will be created", email)
		}

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
