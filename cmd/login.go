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
	"context"
	"fmt"
	"os"

	"github.com/cossacklabs/themis/gothemis/keys"
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/internal/config"
	"github.com/wearedevx/keystone/pkg/client"
	. "github.com/wearedevx/keystone/ui"
)

func ShowAlreadyLoggedInAndExit(account map[string]string) {
	username := account["username"]
	if account["username"] == "" {
		username = account["email"]
	}

	Print(RenderTemplate("already logged in", `You are already logged in as {{ . }}`, username))
	os.Exit(0)
}

func LogIntoExisitingAccount(accountIndex int, currentAccount map[string]string, c client.AuthService) {
	config.SetCurrentAccount(accountIndex)

	publicKey := []byte(currentAccount["public_key"])
	_, jwtToken, err := c.Finish(publicKey)

	if err != nil {
		PrintError(err.Error())
		os.Exit(1)
	}

	config.SetAuthToken(jwtToken)
	config.Write()

	Print(RenderTemplate("login ok", `
{{ OK }} {{ . | bright_green }}
`, fmt.Sprintf("Welcome back, %s", currentAccount["username"])))
	os.Exit(0)
}

func CreateAccountAndLogin(c client.AuthService) {
	keyPair, err := keys.New(keys.TypeEC)

	if err != nil {
		PrintError(err.Error())
		os.Exit(1)
	}

	// Transfer credentials to the server
	// Create (or get) the user info
	user, jwtToken, err := c.Finish(keyPair.Public.Value)

	if err != nil {
		PrintError(err.Error())
		os.Exit(1)
	}

	// Save the user info in the local config
	accountIndex = config.AddAccount(
		map[string]string{
			"account_type": string(user.AccountType),
			"user_id":      user.UserID,
			"ext_id":       user.ExtID,
			"username":     user.Username,
			"fullname":     user.Fullname,
			"email":        user.Email,
			"public_key":   string(keyPair.Public.Value),
			"private_key":  string(keyPair.Private.Value),
		},
	)

	config.SetCurrentAccount(accountIndex)
	config.SetAuthToken(jwtToken)
	config.Write()

	Print(RenderTemplate("login success", `
{{ OK }} {{ . | bright_green }}

Thank you for using Keystone!

To start managing secrets for a project:
  $ cd <path-to-your-project>
  $ ks init

To invite collaborators:
  $ ks invite collaborator@email
`, "Thank you for using Keystone!"))
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login or sign up to your keystone accounts",
	Long:  `Login or sign up to your keystone accounts`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		currentAccount, accountIndex := config.GetCurrentAccount()

		// Already logged in
		if accountIndex >= 0 {
			ShowAlreadyLoggedInAndExit(currentAccount)
		}

		// Currently the only supported auth third party is github
		c := client.GitHubAuth(ctx)

		// Get OAuth connect url
		url, err := c.Start()

		if err != nil {
			PrintError(err.Error())
			os.Exit(1)
		}

		Print(RenderTemplate("login visit", `Visit the URL below to login with your GitHub account:

{{ . | indent 8 }}

Waiting for you to login with your GitHub Account...`, url))

		// Blocking call
		err = c.WaitForExternalLogin()

		if err != nil {
			PrintError(err.Error())
			os.Exit(1)
		}

		currentAccount, accountIndex = config.FindAccount(c)

		if accountIndex >= 0 {
			// Found an exiting matching account,
			// log into it
			LogIntoExisitingAccount(accountIndex, currentAccount, c)
		} else {
			CreateAccountAndLogin(c)
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
}
