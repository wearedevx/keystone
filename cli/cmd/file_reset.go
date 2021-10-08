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
	"path"

	"github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/ui"
	"github.com/wearedevx/keystone/cli/ui/prompts"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/pkg/core"
)

var fileResetYes bool

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset [file path]...",
	Short: "Resets locally modified files with their cached contents",
	Long: `Resets locally modified files with their cached contents.

You can modify a file managed by Keystone without using ` + "`" + `ks file set` + "`." + `
However, you will no longer receive updates on that file, and won’t be able
to change environment.
To discard the changes you made, use ` + "`" + `ks file reset [file path]...` + "`" + `.`,
	Example: `# To reset one specific file
ks file reset ./config.js

# To reset all mananged files
ks file reset
`,
	Run: func(_ *cobra.Command, args []string) {
		ctx := core.New(core.CTX_RESOLVE)

		mustFetchMessages(nil)

		filesToReset := args
		if len(filesToReset) == 0 {
			for _, file := range ctx.ListFiles() {
				filesToReset = append(filesToReset, file.Path)
			}
		}

		ui.Print(ui.RenderTemplate(
			"careful reset",
			`{{ CAREFUL }} {{ "Local changes will be lost" | yellow }}
The content of the files you are resetting will be replaced by their cached content.`,
			nil,
		))

		doIt := fileResetYes
		if !fileResetYes {
			doIt = prompts.Confirm("Continue")
		}

		if doIt {
			for _, file := range filesToReset {
				if !ctx.HasFile(file) {
					ui.Print("File '" + file + "' is not managed by Keystone, ignoring")
					continue
				}

				cachedFilePath := path.Join(ctx.CachedEnvironmentFilesPath(currentEnvironment), file)
				filePath := path.Join(ctx.Wd, file)

				err := utils.CopyFile(cachedFilePath, filePath)
				if err != nil {
					ui.PrintError(err.Error())
				}
			}
		}
	},
}

func init() {
	filesCmd.AddCommand(resetCmd)

	resetCmd.Flags().BoolVarP(
		&fileResetYes,
		"yes",
		"y",
		false,
		"force yes to prompts",
	)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// resetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// resetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
