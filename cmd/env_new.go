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
	"path/filepath"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/manifoldco/promptui"
	"github.com/wearedevx/keystone/internal/errors"
	. "github.com/wearedevx/keystone/internal/utils"
	core "github.com/wearedevx/keystone/pkg/core"
	. "github.com/wearedevx/keystone/ui"

	"github.com/spf13/cobra"
)

// envNewCmd represents the new command
var envNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new environment",
	Long: `Creates a new environment with the given name.

Values for every known secret, and content for every tracked file will be asked.

Example:
  $ ks env new prod
`,

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err *errors.Error

		ctx := core.New(core.CTX_RESOLVE)

		ctx.MustHaveEnvironment(currentEnvironment)

		currentSecrets := ctx.GetAllSecrets(currentEnvironment)
		files := ctx.ListFiles()

		environmentName := args[0]
		ctx.CreateEnvironment(environmentName)

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}

		if !skipPrompts {
			secrets := map[string]string{}

			Print("Please enter the variables values for the environment:")
			for key, value := range currentSecrets {
				p := promptui.Prompt{
					Label:   key,
					Default: value,
				}

				result, _ := p.Run()
				secrets[key] = strings.Trim(result, " ")
			}

			ctx.SetAllSecrets(environmentName, secrets)
		} else {
			ctx.SetAllSecrets(environmentName, currentSecrets)
		}

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}

		for _, file := range files {
			var content []byte
			var erro error
			if !skipPrompts {
				Print("Enter the content of '%s' for the '%s' environment (any key to continue):", file, environmentName)
				// wait for any key
				keyboard.GetSingleKey()
				fmt.Println(file)

				extension := filepath.Ext(file)

				content, erro = CaptureInputFromEditor(
					GetPreferredEditorFromEnvironment,
					extension,
				)

			} else {
				content, erro = ioutil.ReadFile(file)
			}

			if erro != nil {
				panic(erro)
			}

			environmentFileMap := map[string][]byte{}
			environmentFileMap[environmentName] = content
			ctx.AddFile(file, environmentFileMap)

			if err = ctx.Err(); err != nil {
				err.Print()
				return
			}
		}

		Print(RenderTemplate("env created", `
{{ OK }} {{ .Message | bright_green }}

To start using it:
  $ ks env {{ .EnvName }}
`, map[string]string{
			"Message": fmt.Sprintf("The '%s' environment has been created", environmentName),
			"EnvName": environmentName,
		}))
	},
}

func init() {
	envCmd.AddCommand(envNewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RootCmd.Flags().BoolVarP(&skipPrompts, "skip", "s", false, "Skip questions and use defaults")

}
