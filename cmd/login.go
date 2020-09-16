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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-github/v32/github"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	. "github.com/wearedevx/keystone/internal/models"
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

func pollLoginRequest(code string, c chan pollResult) {
	var done bool = false

	for !done {
		time.Sleep(5 * time.Second)
		timeout := time.Duration(20 * time.Second)
		client := http.Client{
			Timeout: timeout,
		}

		u, _ := url.Parse(ksauthURL + "/login-request")
		q := u.Query()
		q.Add("code", code)

		fmt.Println(u.String() + "?" + q.Encode())

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

	}

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

		lr, _ := getLoginRequest()

		conf := &oauth2.Config{
			ClientID:     "60165e42468cf5e34aa8",
			ClientSecret: "d68dba1a9fa7e21d3bdad5e641065b641543587e",
			Scopes:       []string{"user", "read:public_key", "read:emails"},
			RedirectURL:  ksauthURL + "/auth-redirect/" + lr.TemporaryCode + "/",
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://github.com/login/oauth/authorize",
				TokenURL: "https://github.com/login/oauth/access_token",
			},
		}
		url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)

		Print(RenderTemplate("login visit", `Visit the URL below to login with your GitHub account:

{{ . | indent 8 }}
`, url))

		c := make(chan pollResult)
		var result pollResult

		go pollLoginRequest(lr.TemporaryCode, c)

		result = <-c

		if result.err != nil {
			panic(result.err)
		}

		token, err := conf.Exchange(ctx, result.authCode)

		if err != nil {
			panic(err)
		}

		// TODO: Send token to server
		// Server should be fetching the user data to
		// either login or create the user

		ts := oauth2.StaticTokenSource(token)
		tc := oauth2.NewClient(ctx, ts)

		client := github.NewClient(tc)

		user, _, _ := client.Users.Get(ctx, "")
		Print("\nHi, %s!\n\n", *user.Name)

		viper.Set("username", *user.Login)
		viper.Set("fullname", *user.Name)

		userPublicKeys, _, _ := client.Users.ListKeys(ctx, "", nil)

		fmt.Print(userPublicKeys)

		keys := make([]PublicKey, 0)

		for _, githubKey := range userPublicKeys {
			keys = append(keys, PublicKey{
				Typ:       "SSH",
				KeyID:     fmt.Sprintf("%s (ssh)", githubKey.GetTitle()),
				PublicKey: githubKey.GetKey(),
				Email:     githubKey.GetTitle(),
			})
		}

		var selected_key PublicKey
		if len(keys) == 1 {
			selected_key = keys[0]
		} else {
			p := promptui.Select{
				Label: "Select an encryption key to identify yourself",
				Items: keys,
				Templates: &promptui.SelectTemplates{
					Active:   `{{ " * " | blue }} {{- .Typ | yellow }} {{ .KeyID | yellow }} <{{ .Email }}>`,
					Inactive: `   {{ .Typ }} {{ .KeyID }} <{{ .Email }}>`,
				},
				IsVimMode:    true,
				HideSelected: true,
			}

			i, _, err := p.Run()

			if err == nil {
				selected_key = keys[i]
			}

		}

		Print(RenderTemplate("using gpg", `Using the {{ .Typ }} key {{ .KeyID }} <{{ .Email }}>, to identify you across projects.`, selected_key))

		viper.Set("email", selected_key.Email)
		viper.Set("public_key", selected_key)

		if err := WriteConfig(); err != nil {
			panic(err)
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

type PublicKey struct {
	Typ       string
	KeyID     string
	PublicKey string
	Email     string
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
