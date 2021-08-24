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
	"github.com/spf13/viper"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/config"
	"github.com/wearedevx/keystone/cli/internal/spinner"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/client/auth"
	"github.com/wearedevx/keystone/cli/ui"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

var serviceName string

func ShowAlreadyLoggedInAndExit(account models.User) {
	username := account.Username
	if username == "" {
		username = account.Email
	}

	ui.Print(ui.RenderTemplate("already logged in", `You are already logged in as {{ . }}`, username))
	os.Exit(0)
}

func LogIntoExisitingAccount(accountIndex int, currentAccount models.User, c auth.AuthService) {
	config.SetCurrentAccount(accountIndex)

	publicKey, _ := config.GetCurrentUserPublicKey()
	// publicKey := []byte(currentAccount["public_key"])
	_, jwtToken, err := c.Finish(publicKey, config.GetDeviceName(), config.GetDeviceUID())

	if err != nil {
		ui.PrintError(err.Error())
		os.Exit(1)
	}

	config.SetAuthToken(jwtToken)
	config.Write()

	ui.Print(ui.RenderTemplate("login ok", `
{{ OK }} {{ . | bright_green }}
`, fmt.Sprintf("Welcome back, %s", currentAccount.Username)))
	os.Exit(0)
}

func CreateAccountAndLogin(c auth.AuthService) {
	keyPair, err := keys.New(keys.TypeEC)

	if err != nil {
		ui.PrintError(err.Error())
		os.Exit(1)
	}

	// Transfer credentials to the server
	// Create (or get) the user info
	user, jwtToken, err := c.Finish(keyPair.Public.Value, config.GetDeviceName(), config.GetDeviceUID())

	if err != nil {
		ui.PrintError(err.Error())
		os.Exit(1)
	}

	// Save the user info in the local config
	accountIndex := config.AddAccount(
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

	ui.Print(ui.RenderTemplate("login success", `
{{ OK }} {{ . | bright_green }}

Thank you for using Keystone!

To start managing secrets for a project:
  $ cd <path-to-your-project>
  $ ks init <your-project-name>

To invite collaborators:
  $ ks invite collaborator@email
`, "Thank you for using Keystone!"))
}

func selectDeviceName() string {
	existingName := config.GetDeviceName()

	if existingName == "" {
		defaultName := ""
		if hostname, err := os.Hostname(); err == nil {
			defaultName = hostname
		}
		deviceName := prompts.StringInput(
			"Enter the name you want this device to have",
			defaultName,
		)
		return deviceName
	}

	return existingName
}

func SelectAuthService(ctx context.Context) (auth.AuthService, error) {
	if serviceName == "" {
		_, serviceName = prompts.Select(
			"Select an identity provider",
			[]string{
				"github",
				"gitlab",
			},
		)
	}

	return auth.GetAuthService(serviceName, ctx, client.ApiURL)
}

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
		ctx := context.Background()

		currentAccount, accountIndex := config.GetCurrentAccount()

		// Already logged in
		if accountIndex >= 0 {
			ShowAlreadyLoggedInAndExit(currentAccount)
		}

		c, err := SelectAuthService(ctx)

		if err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		// Get OAuth connect url
		url, err := c.Start()

		if err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		ui.Print(ui.RenderTemplate("login visit", `Visit the URL below to login with your {{ .Service }} account:

{{ .Url | indent 8 }}

Waiting for you to login with your {{ .Service }} Account...`, struct {
			Service string
			Url     string
		}{
			Service: c.Name(),
			Url:     url,
		}))

		sp := spinner.Spinner(" ")
		sp.Start()

		// Blocking call
		err = c.WaitForExternalLogin()
		sp.Stop()

		if err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		deviceName := selectDeviceName()
		viper.Set("device", deviceName)

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

	loginCmd.Flags().StringVar(&serviceName, "with", "", "identity provider. Either github or gitlab")
}
