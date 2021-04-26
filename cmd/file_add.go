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
	"path/filepath"

	"github.com/eiannone/keyboard"

	"github.com/spf13/cobra"
	. "github.com/wearedevx/keystone/internal/errors"
	. "github.com/wearedevx/keystone/internal/gitignorehelper"
	. "github.com/wearedevx/keystone/internal/utils"
	"github.com/wearedevx/keystone/pkg/core"
	. "github.com/wearedevx/keystone/ui"
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
	Run: func(cmd *cobra.Command, args []string) {
		var err *Error

		ctx := core.New(core.CTX_RESOLVE)
		ctx.MustHaveEnvironment(currentEnvironment)

		filePath := args[0]
		extension := filepath.Ext(filePath)

		environments := ctx.ListEnvironments()

		environmentFileMap := map[string][]byte{}

		if !FileExists(filePath) {
			err = CannotAddFile(filePath, errors.New("File not found"))
			err.Print()

			return
		}

		currentContent, erro := ioutil.ReadFile(filePath)
		if erro != nil {
			err = CannotAddFile(filePath, erro)
			err.Print()

			return
		}

		environmentFileMap[currentEnvironment] = currentContent

		if !skipPrompts {
			for _, environment := range environments {
				if environment != currentEnvironment {
					Print(fmt.Sprintf("Enter content for file `%s` for the '%s' environment (Press any key to continue)", filePath, environment))
					keyboard.GetSingleKey()

					content, err := CaptureInputFromEditor(
						GetPreferredEditorFromEnvironment,
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

		ctx.AddFile(filePath, environmentFileMap)

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}

		GitIgnore(ctx.Wd, filePath)

		ctx.FilesUseEnvironment(currentEnvironment)

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}

		Print(RenderTemplate("file add success", `
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