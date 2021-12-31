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
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/files"
	"github.com/wearedevx/keystone/cli/internal/gitignorehelper"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui/display"
)

// filesAddCmd represents the push command
var filesAddCmd = &cobra.Command{
	Use:   "add <path to a file>",
	Short: "Adds a file to secrets",
	Long: `Adds a file to secrets.

A secret file is a file which have content that can change
across environments, such as configuration files, credentials,
certificates and so on.

When adding a file, you will be asked for a version of its content
for all known environments – the current content will be used as default.
`,
	Example: `ks file add ./config/config.exs
ks file add ./wp-config.php
ks file add ./certs/my-website.cert

# Skip the prompts
ks file add -s ./credentials.json`,
	Args: cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		var err error

		ctx.MustHaveEnvironment(currentEnvironment)
		mustFetchMessages()

		filePath, err := cleanPathArgument(args[0], ctx.Wd)
		exitIfErr(err)

		fileservice := files.NewFileService(ctx)

		environments := ctx.AccessibleEnvironments
		environmentFileMap := map[string][]byte{}

		if !fileservice.AskToOverrideFilesInCache(filePath, skipPrompts) {
			absPath := filepath.Join(ctx.Wd, filePath)

			if !utils.FileExists(absPath) {
				exit(
					kserrors.
						CannotAddFile(filePath, errors.New("file not found")),
				)
			}

			/* #nosec
			 * Contents are read and copied, not ever run
			 */
			currentContent, err := ioutil.ReadFile(absPath)
			if err != nil {
				exit(kserrors.CannotAddFile(filePath, err))
			}

			environmentFileMap[currentEnvironment] = currentContent

			if !skipPrompts {
				exitIfErr(fileservice.AskContent(
					filePath,
					environments,
					environmentFileMap,
					currentContent,
					currentEnvironment,
				).Err())
			} else {
				for _, environment := range environments {
					environmentFileMap[environment.Name] = currentContent
				}
			}

			changes, messageService := mustFetchMessages()

			exitIfErr(
				ctx.CompareNewFileWhithChanges(filePath, changes).Err(),
			)

			file := keystonefile.FileKey{
				Path:   filePath,
				Strict: addOptional,
			}

			exitIfErr(ctx.AddFile(file, environmentFileMap).Err())

			err = gitignorehelper.GitIgnore(ctx.Wd, filePath)
			exitIfErr(err)

			exitIfErr(
				ctx.FilesUseEnvironment(
					currentEnvironment,
					currentEnvironment,
					core.CTX_KEEP_LOCAL_FILES,
				).Err(),
			)

			exitIfErr(
				messageService.SendEnvironments(ctx.AccessibleEnvironments).
					Err(),
			)
		} else {
			// just add file to keystone.yaml and keep old content

			file := keystonefile.FileKey{
				Path:   filePath,
				Strict: addOptional,
			}

			if err := new(keystonefile.KeystoneFile).
				Load(ctx.Wd).
				AddFile(file).
				Save().
				Err(); err != nil {
				exit(kserrors.FailedToUpdateKeystoneFile(err))
			}
		}

		display.FileAddSuccess(filePath, len(environments))
	},
}

func init() {
	filesCmd.AddCommand(filesAddCmd)
}
