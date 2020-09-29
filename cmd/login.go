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

	"github.com/manifoldco/promptui"
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

func completeLogin(tok *oauth2.Token, pk string) User {
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

		keys, noKeysMessage, err := c.GetPublicKeys()

		if err != nil {
			PrintError(err.Error())
			os.Exit(1)
		}

		// TODO: Allow key generation
		if err == nil && len(keys) == 0 {
			Print(RenderTemplate("no ssh keys", `
{{ ERROR }} {{ "We could not find any SSH keys associated with your GitHub account" | red }}

You can create a new key pair with:
  $ ssh-keygen -t rsa -b 4096
{{ . }}
After it is done, try ks login again.
`, noKeysMessage))
			os.Exit(1)
		}

		// Let the user select the SSH key they want to use
		var selected_key client.PublicKey
		if len(keys) == 1 {
			selected_key = keys[0]
		} else {
			p := promptui.Select{
				Label: "Select an encryption key to identify yourself",
				Items: keys,
				Templates: &promptui.SelectTemplates{
					Active:   `{{ " * " | blue }} {{- .Typ | yellow }} {{ .KeyID | yellow }}`,
					Inactive: `   {{ .Typ }} {{ .KeyID }}`,
				},
				IsVimMode:    true,
				HideSelected: true,
			}

			i, _, err := p.Run()

			if err == nil {
				selected_key = keys[i]
			}

		}

		Print(RenderTemplate("using gpg", `Using the {{ .Typ }} key {{ .KeyID }}, to identify you across projects.`, selected_key))

		// Transfer credentials to the server
		// Create (or get) the user info
		user, err := c.Finish(selected_key)

		if err != nil {
			PrintError(err.Error())
			os.Exit(1)
		}

		// Save the user info in the local config
		viper.Set("account_type", user.AccountType)
		viper.Set("user_id", user.UserID)
		viper.Set("ext_id", user.ExtID)
		viper.Set("username", user.Username)
		viper.Set("fullname", user.Fullname)
		viper.Set("email", user.Email)
		viper.Set("sign_key", user.Keys.Sign)
		viper.Set("cipher_key", user.Keys.Cipher)

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
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
