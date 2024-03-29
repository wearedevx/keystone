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

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/ui/display"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

var (
	forcePrompts bool
	purgeFile    bool
)

// filesRmCmd represents the rm command
var filesRmCmd = &cobra.Command{
	Use:   "rm [path to a file]",
	Short: "Removes a file from secrets",
	Long: `Removes a file from secrets.

The file will no longer be gitignored and its content
will no longer be updated when changing environment.

The content of the file for all environments will be kept in the keystone project.
Files can be used again using "file add" command.
`,
	Example: `ks file rm config/old-test-config.php`,
	Args:    cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		filePath := args[0]

		if !utils.FileExists(filePath) {
			exit(errors.
				CannotRemoveFile(filePath, fmt.Errorf("file not found")))
		}

		if prompts.ConfirmFileRemove(
			filePath,
			ctx.CurrentEnvironment(),
			skipPrompts || !purgeFile,
		) {
			fmt.Println("icit")

			_, messageService := mustFetchMessages()

			exitIfErr(ctx.
				RemoveFile(
					filePath,
					forcePrompts,
					purgeFile,
					ctx.AccessibleEnvironments,
				).
				Err())

			if purgeFile {
				exitIfErr(messageService.
					SendEnvironments(ctx.AccessibleEnvironments).
					Err())
			} else {
				display.FileKept()
			}

			display.FileRemovedSuccess(filePath)
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
