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
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/config"
	"github.com/wearedevx/keystone/cli/internal/environments"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/internal/loggers"
	"github.com/wearedevx/keystone/cli/ui"

	"github.com/wearedevx/keystone/cli/pkg/constants"
	"github.com/wearedevx/keystone/cli/pkg/core"
)

var (
	ksauthURL string //= "http://localhost:9000"
	ksapiURL  string //= "http://localhost:9001"
)

var (
	cfgFile            string = ""
	currentEnvironment string
	quietOutput        bool
	skipPrompts        bool
	debug              bool
)

var ctx *core.Context

var CWD string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ks <command> [sub-command] [inputs]...",
	Short: "A safe system for developers to store, share and use secrets.",
	Long:  `A safe system for developers to store, share and use secrets.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, _ []string) {
		if err := cmd.Usage(); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() int {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	return 0
}

var (
	noEnvironmentCommands []string
	noProjectCommands     []string
	noLoginCommands       []string
)

func findCurrentCommand(args []string) string {
	for _, candidate := range args[1:] {
		if !strings.HasPrefix(candidate, "-") &&
			!strings.HasPrefix(candidate, "/") &&
			candidate != "ks" {
			return candidate
		}
	}

	return ""
}

// Initialize function performs initialization logic before executing the
// actual command
func Initialize() {
	if len(os.Args) <= 1 {
		return
	}

	loggers.SetOutput(ioutil.Discard)
	if debug {
		loggers.SetOutput(os.Stderr)
	}

	checkEnvironment := true
	checkProject := true
	checkLogin := false

	askHelp := core.Contains(os.Args, "--help")

	if len(os.Args) > 1 {
		command := findCurrentCommand(os.Args)
		checkEnvironment = !isIn(noEnvironmentCommands, command) && !askHelp
		checkProject = !isIn(noProjectCommands, command) && !askHelp
		checkLogin = !isIn(noLoginCommands, command) && !askHelp

		if command == "init" {
			ctx = core.New(core.CTX_INIT)
		} else {
			ctx = core.New(core.CTX_RESOLVE)
		}
	}

	if checkProject {
		exitIfErr(ctx.Err())
	}

	isKeystoneFile := keystonefile.ExistsKeystoneFile(ctx.Wd)

	current := ctx.CurrentEnvironment()
	ctx.SetError(nil)

	if currentEnvironment == "" {
		currentEnvironment = current
	}

	if checkProject && !isKeystoneFile {
		exit(kserrors.NotAKeystoneProject(".", nil))
	}

	if checkProject && config.IsLoggedIn() {
		es := environments.NewEnvironmentService(ctx)
		exitIfErr(es.Err())

		ctx.AccessibleEnvironments = es.GetAccessibleEnvironments()
		exitIfErr(ctx.Err())

		// If no accessible environment, then user has no access to the project
		if len(ctx.AccessibleEnvironments) == 0 {
			exit(
				kserrors.ProjectDoesntExist(
					ctx.GetProjectName(),
					ctx.GetProjectID(),
					nil,
				),
			)
		}

		if isKeystoneFile {
			environmentsToSave := make([]models.Environment, 0)
			for _, env := range ctx.AccessibleEnvironments {
				env.VersionID = ""
				environmentsToSave = append(environmentsToSave, env)
			}

			ctx.Init(models.Project{
				Environments: environmentsToSave,
			})
		}
		ctx.RemoveForbiddenEnvironments(ctx.AccessibleEnvironments)

		if currentEnvironment == "" {
			currentEnvironment = ctx.CurrentEnvironment()
		}
	}

	if checkEnvironment && !ctx.HasEnvironment(currentEnvironment) {
		ctx.Init(models.Project{})
		if currentEnvironment == "" {
			ctx.SetCurrent(string(constants.DEV))
		}
	}

	if checkLogin && !config.IsLoggedIn() {
		exit(kserrors.MustBeLoggedIn(nil))
	}
}

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	CWD = cwd

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")
	RootCmd.PersistentFlags().
		BoolVarP(&quietOutput, "quiet", "q", false, "make the output machine readable")

	RootCmd.PersistentFlags().
		BoolVar(&debug, "debug", false, "debug output")

	RootCmd.PersistentFlags().
		StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.config/keystone/keystone.yaml)")

	RootCmd.PersistentFlags().
		BoolVarP(&skipPrompts, "skip", "s", false, "skip prompts and use default")

	RootCmd.PersistentFlags().
		StringVarP(&currentEnvironment, "env", "", "", "environment to use instead of the current one")

	cobra.OnInitialize(func() {
		// Call directly initConfig. cobra doesn't call initConfig func.
		err := config.InitConfig(cfgFile)
		exitIfErr(err)

		Initialize()
	})
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	noEnvironmentCommands = []string{
		"login",
		"logout",
		"documentation",
		"completion",
		"__complete",
		"init",
		"whoami",
		"invite",
		"version",
		"device",
		"orga",
		"project",
		"hook",
	}

	noProjectCommands = noEnvironmentCommands

	noLoginCommands = []string{"login", "source", "documentation", "completion", "__complete", "version", "backup"}
}
