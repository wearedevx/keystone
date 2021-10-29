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
	"path"

	"github.com/spf13/cobra"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set <path to a file>",
	Short: "Updates a file’s content for the current environment",
	Long: `Updates a file’s content for the current environment.

Changes the content of a file without altering other environments.
The local version of the file will be used.
`,
	Example: `ks file set ./config.php

# Change the content of ./config.php for the 'staging' environment:
ks --env staging file set ./config.php
`,
	Args: cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		var err error

		ctx.MustHaveEnvironment(currentEnvironment)
		ctx.MustHaveAccessToEnvironment(currentEnvironment)

		filePath := args[0]

		if !utils.FileExists(path.Join(ctx.Wd, filePath)) {
			exit(kserrors.
				CannotSetFile(filePath, errors.New("file not found")))
		}

		if !ctx.HasFile(filePath) {
			exit(kserrors.
				CannotSetFile(
					filePath,
					errors.New("file not added to project"),
				))
		}

		content, err := ctx.GetLocalFileContents(filePath)
		if err != nil {
			if err.Error() != "No contents" {
				exit(kserrors.CannotSetFile(filePath, err))
			}
		}

		changes, messageService := mustFetchMessages()

		err = ctx.
			CompareNewFileWhithChanges(filePath, changes).
			SetFile(filePath, content).
			// Local files should be kept during a file set
			FilesUseEnvironment(
				currentEnvironment,
				currentEnvironment,
				core.CTX_KEEP_LOCAL_FILES,
			).
			Err()
		exitIfErr(err)

		err = messageService.SendEnvironments(ctx.AccessibleEnvironments).Err()
		exitIfErr(err)

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
