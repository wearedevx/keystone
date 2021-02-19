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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/cossacklabs/themis/gothemis/keys"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wearedevx/keystone/internal/crypto"
	. "github.com/wearedevx/keystone/internal/models"
	"github.com/wearedevx/keystone/pkg/client"
	. "github.com/wearedevx/keystone/ui"
	"golang.org/x/oauth2"
)

func getLoginRequest() (LoginRequest, error) {
	timeout := time.Duration(20 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("POST", ksauthURL+"/login-request", nil)
	request.Header.Set("Accept", "application/json; charset=utf-8")

	if err != nil {
		panic(err)
	}

	resp, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	var loginRequest LoginRequest

	if resp.StatusCode == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(&loginRequest)

		return loginRequest, err
	} else {
		message := fmt.Sprintf("Request Error: %s", resp.Status)
		fmt.Println(message)

		return loginRequest, errors.New(message)
	}
}

type pollResult struct {
	authCode string
	err      error
}

const MAX_ATTEMPTS int = 12

func pollLoginRequest(code string, c chan pollResult) {
	var done bool = false
	attemps := 0

	for !done {
		attemps = attemps + 1

		time.Sleep(5 * time.Second)

		timeout := time.Duration(20 * time.Second)
		client := http.Client{
			Timeout: timeout,
		}

		u, _ := url.Parse(ksauthURL + "/login-request")
		q := u.Query()
		q.Add("code", code)

		request, err := http.NewRequest("GET", u.String()+"?"+q.Encode(), nil)
		request.Header.Set("Accept", "application/json; charset=utf-8")

		if err != nil {
			panic(err)
		}

		resp, err := client.Do(request)

		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		var loginRequest LoginRequest

		if resp.StatusCode == http.StatusOK {
			err = json.NewDecoder(resp.Body).Decode(&loginRequest)

			if loginRequest.AuthCode != "" {
				r := pollResult{
					authCode: loginRequest.AuthCode,
				}

				c <- r

				done = true
			}

		} else {
			message := fmt.Sprintf("Request Error: %s", resp.Status)
			fmt.Println(message)

			r := pollResult{
				err: errors.New(message),
			}

			c <- r

			done = true
		}

		if attemps == MAX_ATTEMPTS {
			done = true
		}
	}

}

func completeLogin(tok *oauth2.Token, pk []byte) User {
	payload := LoginPayload{
		AccountType: GitHubAccountType,
		Token:       tok,
		PublicKey:   pk,
	}

	requestPayload := make([]byte, 0)
	buf := bytes.NewBuffer(requestPayload)
	json.NewEncoder(buf).Encode(&payload)

	req, err := http.NewRequest("POST", ksauthURL+"/complete", buf)
	req.Header.Add("Accept", "application/octet-stream")

	if err != nil {
		panic(err)
	}

	timeout := time.Duration(20 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		panic(fmt.Errorf("Failed to complete login: %s", resp.Status))
	}

	p, err := ioutil.ReadAll(resp.Body)
	var user User

	if err != nil {
		panic(err)
	}

	crypto.DecryptWithPublicKey(pk, p, &user)

	return user
}

// finds an account matching `user` in the `account` slice
func findAccountIndex(accounts []map[string]string, c client.AuthService) int {
	current := -1

	for i, account := range accounts {
		isAccount, _ := c.CheckAccount(account)

		if isAccount {
			current = i
			break
		}
	}

	return current
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		accountIndex := viper.Get("current").(int)
		accounts := viper.Get("accounts").([]map[string]string)

		if accountIndex >= 0 && len(accounts) > accountIndex {
			account := accounts[accountIndex]
			username := account["username"]

			if username != "" {
				Print(RenderTemplate("already logged in", `You are already logged in as {{ . }}`, username))
			}
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

		accountIndex = findAccountIndex(accounts, c)

		if accountIndex >= 0 {
			viper.Set("current", accountIndex)

			if viper.WriteConfig(); err != nil {
				Print(RenderTemplate("config write error", `
{{ ERROR }} {{ . | red }}

You have been successfully logged in, but the configuration file could not be written
`, err.Error()))
				os.Exit(1)
			}

			Print(RenderTemplate("login ok", `
{{ OK }} {{ . | bright_green }}
`, fmt.Sprintf("Welcome back, %s", accounts[accountIndex]["username"])))

			os.Exit(0)
		}

		keyPair, err := keys.New(keys.TypeEC)

		if err != nil {
			PrintError(err.Error())
			os.Exit(1)
		}

		// Transfer credentials to the server
		// Create (or get) the user info
		user, err := c.Finish(keyPair.Public.Value)

		if err != nil {
			PrintError(err.Error())
			os.Exit(1)
		}

		// Save the user info in the local config
		account := map[string]string{
			"account_type": string(user.AccountType),
			"user_id":      user.UserID,
			"ext_id":       user.ExtID,
			"username":     user.Username,
			"fullname":     user.Fullname,
			"email":        user.Email,
			"public_key":   string(keyPair.Public.Value),
			"private_key":  string(keyPair.Private.Value),
		}

		accounts = append(accounts, account)

		viper.Set("current", len(accounts)-1)
		viper.Set("accounts", accounts)

		if err := WriteConfig(); err != nil {
			Print(RenderTemplate("config write error", `
{{ ERROR }} {{ . | red }}

You have been successfully logged in but the configuration file could not be written
`, err.Error()))
			os.Exit(1)
		}

		Print(RenderTemplate("login success", `
{{ OK }} {{ . | bright_green }}

Thank you for using Keystone!

To start managing secrets for a project:
  $ cd <path-to-your-project>
  $ ks init

To invite collaborators:
  $ ks invite collaborator@email
`, "Thank you for using Keystone!"))
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
