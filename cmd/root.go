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
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/internal/config"
	"github.com/wearedevx/keystone/internal/errors"
	"github.com/wearedevx/keystone/pkg/core"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var ksauthURL string //= "http://localhost:9000"
var ksapiURL string  //= "http://localhost:9001"

var cfgFile string = ""
var currentEnvironment string
var quietOutput bool

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
		fmt.Println(err)
		os.Exit(1)
	}
	return 0
}

func Initialize() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	ctx := core.New(core.CTX_RESOLVE)
	environments := ctx.ListEnvironments()
	current := ctx.CurrentEnvironment()

	RootCmd.PersistentFlags().StringVar(&currentEnvironment, "env", current, "environment to use instead of the current one")

	checkEnvironment := true
	checkLogin := false

	if len(os.Args) > 1 {
		if os.Args[1] == "login" || os.Args[1] == "documentation" || os.Args[1] == "init" {
			checkEnvironment = false
		}

		if os.Args[1] != "login" {
			checkLogin = true
		}

	}

	if checkEnvironment && !ctx.HasEnvironment(currentEnvironment) {
		errors.EnvironmentDoesntExist(currentEnvironment, strings.Join(environments, ", "), nil).Print()
		os.Exit(1)
	}

	if checkLogin && !config.IsLoggedIn() {
		errors.MustBeLoggedIn(nil).Print()
		os.Exit(1)
	}
}

func init() {
	// Call directly initConfig. cobra doesn't call initConfig func.
	initConfig()
	// cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")
	RootCmd.PersistentFlags().BoolVarP(&quietOutput, "quiet", "q", false, "make the output machine readable")

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/keystone.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// Create conf file if not exist
func createFileIfNotExist(filePath string) {

	// Check if need to create file
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// path/to/whatever does not exist

		err := ioutil.WriteFile(filePath, []byte(""), 0755)

		if err != nil {
			fmt.Printf("Unable to write file: %v", err)
		}
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".keystone" (without extension).
		viper.AddConfigPath(path.Join(home, ".config"))
		viper.SetConfigName("keystone")
		viper.SetConfigType("yaml")

		createFileIfNotExist(path.Join(home, ".config", "keystone.yaml"))
	}

	defaultAccounts := make([]map[string]string, 0)

	viper.SetDefault("current", -1)
	viper.SetDefault("accounts", defaultAccounts)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		viper.WriteConfig()
	}
}

func WriteConfig() error {
	var err error

	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	configPath := path.Join(home, ".config", "keystone.yml")
	if err = viper.WriteConfigAs(configPath); err != nil {
		if os.IsNotExist(err) {
			err = viper.WriteConfigAs(configPath)
		}
	}

	return err
}
