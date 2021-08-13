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

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/messages"
	"github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/ui"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

var forcePrompts bool
var purgeFile bool

// filesRmCmd represents the rm command
var filesRmCmd = &cobra.Command{
	Use:   "rm [path to a file]",
	Short: "Removes a file from secrets",
	Long: `Removes a file from secrets.

The file will no longer be gitignored and its content
will no longer be updated when changing environment.

The content of the file for other environments *will be lost*.
This is permanent, and cannot be undone.
`,
	Example: `ks file rm config/old-test-config.php`,
	Args:    cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		var err *errors.Error

		filePath := args[0]

		if !utils.FileExists(filePath) {
			err = errors.CannotRemoveFile(filePath, fmt.Errorf("file not found"))
			err.Print()
			return
		}

		if promptYesNo(filePath) {
			var printer = &ui.UiPrinter{}
			ms := messages.NewMessageService(ctx, printer)
			ms.GetMessages()

			if err := ms.Err(); err != nil {
				err.Print()
				os.Exit(1)
			}

			ctx.RemoveFile(filePath, forcePrompts, purgeFile, ctx.AccessibleEnvironments)
			if err = ctx.Err(); err != nil {
				err.Print()
				return
			}

			if purgeFile {
				if err := ms.SendEnvironments(ctx.AccessibleEnvironments).Err(); err != nil {
					err.Print()
					os.Exit(1)
				}
			} else {
				ui.Print("The file is kept in your keystone project for all the environments, in case you need it again.")
				ui.Print("If you want to remove it from your device, use --purge")

			}

			ui.PrintSuccess("%s has been removed from the secret files.", filePath)
		}

	},
}

func init() {
	filesCmd.AddCommand(filesRmCmd)

	filesRmCmd.Flags().BoolVarP(
		&forcePrompts,
		"force",
		"f",
		false,
		"force remove file on system.",
	)

	filesRmCmd.Flags().BoolVarP(
		&purgeFile,
		"purge",
		"p",
		false,
		"purge file content from all environments",
	)
}

func promptYesNo(filePath string) bool {
	if skipPrompts {
		return true
	}
	if !purgeFile {
		return true
	}

	ui.Print(ui.RenderTemplate("confirm files rm",
		`{{ CAREFUL }} You are about to remove {{ .Path }} from the secret files.
Content for the current environment ({{ .Environment }}) will be kept.
Its content for other environments will be lost, it will no longer be gitignored.
This is permanent, and cannot be undone.`, map[string]string{
			"Path":        filePath,
			"Environment": ctx.CurrentEnvironment(),
		}))

	return prompts.Confirm("Continue")
}
