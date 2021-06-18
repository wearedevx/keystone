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

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/messages"
	"github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/ui"
)

var forcePrompts bool

// filesRmCmd represents the rm command
var filesRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Removes a file from secrets",
	Long: `Removes a file from secrets.

The file will no longer be gitignored and its content
will no longer be updated when changing environment.

The content of the file for other environments will be lost.
This is permanent, and cannot be undone.

Example:
  $ ks file rm config/old-test-config.php`,
	Args: cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		var err *errors.Error

		filePath := args[0]

		if !utils.FileExists(filePath) {
			err = errors.CannotRemoveFile(filePath, fmt.Errorf("file not found"))
			err.Print()
			return
		}

		var printer = &ui.UiPrinter{}
		ms := messages.NewMessageService(ctx, printer)
		ms.GetMessages()

		if err := ms.Err(); err != nil {
			err.Print()
			os.Exit(1)
		}

		ui.Print(ui.RenderTemplate("confirm files rm", `{{ CAREFUL }} You are about to remove {{ .Path }} from the secret files.
Content for the current environment ({{ .Environment }}) will be kept.
Its content for other environments will be lost, it will no longer be gitignored.
This is permanent, and cannot be undone.`, map[string]string{
			"Path":        filePath,
			"Environment": ctx.CurrentEnvironment(),
		}))

		result := "y"

		if !skipPrompts {
			p := promptui.Prompt{
				Label:     "Continue",
				IsConfirm: true,
			}

			result, _ = p.Run()
		}

		if result == "y" {
			ctx.RemoveFile(filePath, forcePrompts, ctx.AccessibleEnvironments)
		}

		if err = ctx.Err(); err != nil {
			err.Print()
			return
		}

		if result == "y" {
			ui.PrintSuccess("%s has been removed from the secret files.", filePath)
		}
	},
}

func init() {
	filesCmd.AddCommand(filesRmCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RootCmd.PersistentFlags().BoolVarP(&forcePrompts, "force", "f", false, "force remove file on system.")
}
