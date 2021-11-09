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
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/config"
	"github.com/wearedevx/keystone/cli/internal/login"
	"github.com/wearedevx/keystone/cli/internal/spinner"
	"github.com/wearedevx/keystone/cli/ui/display"
)

var serviceName string

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login or sign up to your Keystone accounts",
	Long: `Login or sign up to your Keystone accounts.

When singing up, you will be asked to log into either your GitHub or Gitlab
account, to verify your identity.
We do not use any information other than your email address and your username.
	`,
	Example: `ks login
ks login --with=gitlab
ks login ––with=github`,
	Args: cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		currentAccount, accountIndex := config.GetCurrentAccount()

		// Already logged in
		if accountIndex >= 0 {
			display.AlreadyLoggedIn(currentAccount)
			exit(nil)
		}

		sp := spinner.Spinner("").Start()

		var url, name string
		ls := login.NewLoginService(serviceName).
			GetLoginLink(&url, &name)

		sp.Stop()
		exitIfErr(ls.Err())

		display.LoginLink(name, url)

		sp.Start()

		// Blocking call
		err := ls.WaitForExternalLogin().Err()

		sp.Stop()
		exitIfErr(err)

		existingUser := ls.
			PromptDeviceName(skipPrompts).
			FindAccount(&currentAccount, &accountIndex).
			PerformLogin(currentAccount, accountIndex)

		exitIfErr(ls.Err())

		if existingUser {
			display.WelcomeBack(currentAccount)
		} else {
			display.LoginSucces()
		}
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	loginCmd.Flags().
		StringVar(&serviceName, "with", "", "identity provider. Either github or gitlab")
}
