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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/eiannone/keyboard"
	"github.com/spf13/cobra"
	envservice "github.com/wearedevx/keystone/cli/internal/environments"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/gitignorehelper"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/internal/messages"
	"github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
)

// filesAddCmd represents the push command
var filesAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a file to secrets",
	Long: `Adds a file to secrets

A secret file is a file which have content that can changge
across environments, such as configuration files, credentials,
certificates and so on.

When adding a file, you will be asked for a version of its content
for all known environments.

Examples:
  $ ks file add ./config/config.exs
  
  $ ks file add ./wp-config.php

  $ ks file add ./certs/my-website.cert
`,
	Args: cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		var err *kserrors.Error

		ctx := core.New(core.CTX_RESOLVE)
		ctx.MustHaveEnvironment(currentEnvironment)

		filePath := args[0]
		extension := filepath.Ext(filePath)

		environments := ctx.ListEnvironments()

		environmentFileMap := map[string][]byte{}

		if !utils.FileExists(filePath) {
			err = kserrors.CannotAddFile(filePath, errors.New("file not found"))
			err.Print()

			return
		}

		currentContent, erro := ioutil.ReadFile(filePath)
		if erro != nil {
			err = kserrors.CannotAddFile(filePath, erro)
			err.Print()

			return
		}

		environmentFileMap[currentEnvironment] = currentContent

		if !skipPrompts {
			for _, environment := range environments {
				if environment != currentEnvironment {
					ui.Print(fmt.Sprintf("Enter content for file `%s` for the '%s' environment (Press any key to continue)", filePath, environment))
					_, _, err := keyboard.GetSingleKey()
					if err != nil {
						panic(err)
					}

					content, err := utils.CaptureInputFromEditor(
						utils.GetPreferredEditorFromEnvironment,
						extension,
					)

					if err != nil {
						panic(err)
					}

					environmentFileMap[environment] = content
				}
			}
		} else {
			for _, environment := range environments {
				environmentFileMap[environment] = currentContent
			}
		}

		ms := messages.NewMessageService(ctx)
		ms.GetMessages()

		if err := ms.Err(); err != nil {
			err.Print()
			os.Exit(1)
		}

		es := envservice.NewEnvironmentService(ctx)
		accessibleEnvironments := es.GetAccessibleEnvironments()

		if err := es.Err(); err != nil {
			err.Print()
			os.Exit(1)
		}

		file := keystonefile.FileKey{
			Path:   filePath,
			Strict: false, // TODO
		}

		ctx.AddFile(file, environmentFileMap)

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}

		err_ := gitignorehelper.GitIgnore(ctx.Wd, filePath)
		if err_ != nil {
			ui.PrintError(err_.Error())
			return
		}

		ctx.FilesUseEnvironment(currentEnvironment)

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}

		if err := ms.SendEnvironments(accessibleEnvironments).Err(); err != nil {
			err.Print()
			os.Exit(1)
		}

		ui.Print(ui.RenderTemplate("file add success", `
{{ OK }} {{ .Title | green }}
The file has been added to all environments.
It has also been gitignored.`, map[string]string{
			"Title": fmt.Sprintf("Added '%s'", filePath),
		}))
	},
}

func init() {
	filesCmd.AddCommand(filesAddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pushCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pushCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// RootCmd.Flags().BoolVarP(&skipPrompts, "skip", "s", false, "Skip questions and use defaults")
}
