/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"github.com/wearedevx/keystone/api/pkg/models"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
)

// fileSetCmd represents the fileSet command
var fileSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Updates the content of a secret file for the current environment",
	Long: `Updates the content of a secret file for the current environment.

Changes the content of a file for the current environment without alering the others.

Examples:
  $ ks file set credentials.json

  Change the content of 'credentials.json' for the 'staging' environment:
  $ ks --env staging file set credentials.json
`,
	Run: func(cmd *cobra.Command, args []string) {
		var err *kserrors.Error

		ctx := core.New(core.CTX_RESOLVE)
		ctx.MustHaveEnvironment(currentEnvironment)
		environment := currentEnvironment

		accessibleEnvironments := ctx.GetAccessibleEnvironments()

		filePath := args[0]
		extension := filepath.Ext(filePath)

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

		content := currentContent

		if !skipPrompts {
			ui.Print(fmt.Sprintf("Enter content for file `%s` for the '%s' environment (Press any key to continue)", filePath, environment))
			_, _, err := keyboard.GetSingleKey()
			if err != nil {
				panic(err)
			}

			content, err = utils.CaptureInputFromEditor(
				utils.GetPreferredEditorFromEnvironment,
				extension,
			)

			if err != nil {
				panic(err)
			}

		}

		file := keystonefile.FileKey{
			Path:   filePath,
			Strict: false, // TODO
		}

		ctx.SetFile(file, environment, content)

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}

		ctx.FilesUseEnvironment(currentEnvironment)

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}

		// Fetch new messages to see if added secret has changed
		messagesByEnvironment := &models.GetMessageByEnvironmentResponse{
			Environments: map[string]models.GetMessageResponse{},
		}

		fmt.Println("Syncing data...")
		terr := ctx.FetchNewMessages(messagesByEnvironment)

		if terr != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		_, werr := ctx.WriteNewMessages(*messagesByEnvironment)

		if werr != nil {
			ui.PrintError(werr.Error())
			os.Exit(1)
		}

		// TODO
		// Format beautyiful error
		if pushErr := ctx.PushEnv(accessibleEnvironments); err != nil {
			ui.PrintError(pushErr.Error())
			return
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
	filesCmd.AddCommand(fileSetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fileSetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fileSetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
