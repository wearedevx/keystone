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
	"os"
	"path"
	"path/filepath"

	"github.com/eiannone/keyboard"
	"github.com/spf13/cobra"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/messages"
	"github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/ui"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set <path to a file>",
	Short: "Updates a file's content for the current environment",
	Long: `Updates a file's content for the current environment.

Changes the content of a file without altering other environments.
`,
	Example: `ks file set ./config.php

# Change the content of ./config.php for the 'staging' environment:
ks --env staging file set ./config.php
`,
	Run: func(cmd *cobra.Command, args []string) {
		var err *kserrors.Error

		ctx.MustHaveEnvironment(currentEnvironment)
		ctx.MustHaveAccessToEnvironment(currentEnvironment)

		filePath := args[0]

		if !utils.FileExists(path.Join(ctx.Wd, filePath)) {
			err = kserrors.CannotSetFile(filePath, errors.New("file not found"))
			err.Print()

			return
		}

		if !ctx.HasFile(filePath) {
			kserrors.
				CannotSetFile(filePath, errors.New("file not added to project"))

			os.Exit(1)
		}

		currentContent, erro := ctx.GetFileContents(filePath, currentEnvironment)
		if erro != nil {
			if erro.Error() != "No contents" {
				err = kserrors.CannotSetFile(filePath, erro)
				err.Print()

				os.Exit(1)
			}
		}

		content := currentContent
		if !skipPrompts {
			content = askContent(filePath, currentContent)
		}

		var printer = &ui.UiPrinter{}
		ms := messages.NewMessageService(ctx, printer)
		changes := ms.GetMessages()

		if err := ms.Err(); err != nil {
			err.Print()
			os.Exit(1)
		}

		if err = ctx.
			CompareNewFileWhithChanges(filePath, changes).
			SetFile(filePath, content).
			FilesUseEnvironment(currentEnvironment).
			Err(); err != nil {
			err.Print()
			os.Exit(1)
		}

		if err := ms.SendEnvironments(ctx.AccessibleEnvironments).Err(); err != nil {
			err.Print()
			os.Exit(1)
		}

		ui.Print(ui.RenderTemplate("file set success", `
{{ OK }} {{ .Title | green }}
`, map[string]string{
			"Title": fmt.Sprintf("Modified '%s'", filePath),
		}))
	},
}

func init() {
	filesCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func askContent(filePath string, currentContent []byte) []byte {
	extension := filepath.Ext(filePath)

	ui.Print(fmt.Sprintf("Enter content for file '%s' for the '%s' environment (Press any key to continue)", filePath, currentEnvironment))
	_, _, err := keyboard.GetSingleKey()
	if err != nil {
		errmsg := fmt.Sprintf("Failed to read user input (%s)", err.Error())
		println(errmsg)
		os.Exit(1)
		return []byte{}
	}

	content, err := utils.CaptureInputFromEditor(
		utils.GetPreferredEditorFromEnvironment,
		extension,
		string(currentContent),
	)

	if err != nil {
		errmsg := fmt.Sprintf("Failed to get content from editor (%s)", err.Error())
		println(errmsg)
		os.Exit(1)
		return []byte{}
	}

	return content
}
