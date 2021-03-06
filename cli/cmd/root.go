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
	"path"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/config"
	"github.com/wearedevx/keystone/cli/internal/environments"
	"github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"

	"github.com/wearedevx/keystone/cli/pkg/core"

	"github.com/spf13/viper"
)

var ksauthURL string //= "http://localhost:9000"
var ksapiURL string  //= "http://localhost:9001"

var cfgFile string = ""
var currentEnvironment string
var quietOutput bool
var skipPrompts bool

var ctx *core.Context

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ks",
	Short: "A safe system for developers to store, share and use secrets.",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() int {
	Initialize()
	if err := RootCmd.Execute(); err != nil {
		// fmt.Println(err)
		os.Exit(1)
	}
	return 0
}

var noEnvironmentCommands []string
var noProjectCommands []string
var noLoginCommands []string

func isIn(haystack []string, needle string) bool {
	for _, hay := range haystack {
		if hay == needle {
			return true
		}
	}

	return false
}

func Initialize() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	if len(os.Args) > 1 {
		command := os.Args[1]
		if command == "init" {
			ctx = core.New(core.CTX_INIT)
		} else {
			ctx = core.New(core.CTX_RESOLVE)
		}
	}

	currentfolder, _ := os.Getwd()

	isKeystoneFile := keystonefile.ExistsKeystoneFile(currentfolder)

	current := ctx.CurrentEnvironment()
	RootCmd.PersistentFlags().StringVar(&currentEnvironment, "env", current, "environment to use instead of the current one")

	checkEnvironment := true
	checkProject := true
	checkLogin := false

	if len(os.Args) > 1 {
		command := os.Args[1]
		checkEnvironment = !isIn(noEnvironmentCommands, command)
		checkProject = !isIn(noProjectCommands, command)
		checkLogin = !isIn(noLoginCommands, command)
	}

	if checkProject && !isKeystoneFile {
		errors.NotAKeystoneProject(".", nil).Print()
		os.Exit(1)
	}

	if checkProject && config.IsLoggedIn() {
		es := environments.NewEnvironmentService(ctx)
		if err := es.Err(); err != nil {
			err.Print()
			os.Exit(1)
		}

		ctx.AccessibleEnvironments = es.GetAccessibleEnvironments()

		if err := es.Err(); err != nil {
			err.Print()
			return
		}

		// If no accessible environment, then user has no access to the project
		if len(ctx.AccessibleEnvironments) == 0 {
			errors.ProjectDoesntExist(ctx.GetProjectName(), ctx.GetProjectID(), nil).Print()
			os.Exit(1)
		}

		if isKeystoneFile {
			ctx.Init(models.Project{
				Environments: ctx.AccessibleEnvironments,
			})
		}
		ctx.RemoveForbiddenEnvironments(ctx.AccessibleEnvironments)
	}

	if checkEnvironment && !ctx.HasEnvironment(currentEnvironment) {
		ctx.Init(models.Project{})
		// errors.EnvironmentDoesntExist(currentEnvironment, strings.Join(environments, ", "), nil).Print()
		// os.Exit(1)
	}

	if checkLogin && !config.IsLoggedIn() {
		errors.MustBeLoggedIn(nil).Print()
		os.Exit(1)
	}

}

func init() {
	// Call directly initConfig. cobra doesn't call initConfig func.
	config.InitConfig(cfgFile)
	// cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")
	RootCmd.PersistentFlags().BoolVarP(&quietOutput, "quiet", "q", false, "make the output machine readable")

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/keystone.yaml)")

	RootCmd.PersistentFlags().BoolVarP(&skipPrompts, "skip", "s", false, "skip prompts and use default")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	noEnvironmentCommands = []string{
		"login", "logout", "documentation", "init", "whoami",
	}
	noProjectCommands = noEnvironmentCommands

	noLoginCommands = []string{"login"}

}

func WriteConfig() error {
	var err error

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}

	configPath := path.Join(home, ".config", "keystone.yml")
	if err = viper.WriteConfigAs(configPath); err != nil {
		if os.IsNotExist(err) {
			err = viper.WriteConfigAs(configPath)
		}
	}

	return err
}
