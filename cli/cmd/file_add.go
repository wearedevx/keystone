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
	"path"
	"path/filepath"
	"reflect"

	"github.com/eiannone/keyboard"
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/api/pkg/models"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/gitignorehelper"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/internal/messages"
	"github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
	"github.com/wearedevx/keystone/cli/ui/prompts"
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
for all known environments – the current contend will be used as default.
`,
	Example: `ks file add ./config/config.exs
ks file add ./wp-config.php
ks file add ./certs/my-website.cert`,
	Args: cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		var err *kserrors.Error

		ctx.MustHaveEnvironment(currentEnvironment)

		filePath := args[0]

		environments := ctx.AccessibleEnvironments

		environmentFileMap := map[string][]byte{}

		useOldFile := checkFileAlreadyInCache(filePath)
		if !useOldFile {

			if !utils.FileExists(path.Join(ctx.Wd, filePath)) {
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
				askContentOfFile(
					environments,
					filePath,
					environmentFileMap,
					currentContent,
				)
			} else {
				for _, environment := range environments {
					environmentFileMap[environment.Name] = currentContent
				}
			}

			var printer = &ui.UiPrinter{}
			ms := messages.NewMessageService(ctx, printer)
			changes := ms.GetMessages()

			if err := ms.Err(); err != nil {
				err.Print()
				os.Exit(1)
			}

			if err = ctx.CompareNewFileWhithChanges(filePath, changes).Err(); err != nil {
				err.Print()
				return
			}

			file := keystonefile.FileKey{
				Path:   filePath,
				Strict: addOptional,
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

			ctx.FilesUseEnvironment(currentEnvironment, currentEnvironment, core.CTX_KEEP_LOCAL_FILES)

			if err = ctx.Err(); err != nil {
				err.Print()
				return
			}

			if err := ms.SendEnvironments(ctx.AccessibleEnvironments).Err(); err != nil {
				err.Print()
				os.Exit(1)
			}
		} else {
			// just add file to keystone.yaml and keep old content

			file := keystonefile.FileKey{
				Path:   filePath,
				Strict: addOptional,
			}
			if err := new(keystonefile.KeystoneFile).Load(ctx.Wd).AddFile(file).Save().Err(); err != nil {
				kserrors.FailedToUpdateKeystoneFile(err).Print()
				os.Exit(1)
			}
		}

		ui.Print(ui.RenderTemplate("file add success", `
{{ OK }} {{ .Title | green }}
The file has been added to {{ .NumberEnvironments }} environment(s).
It has also been gitignored.`, map[string]string{
			"Title":              fmt.Sprintf("Added '%s'", filePath),
			"NumberEnvironments": fmt.Sprintf("%d", len(environments)),
		}))
	},
}

func init() {
	filesCmd.AddCommand(filesAddCmd)
}

func askContentOfFile(
	environments []models.Environment,
	filePath string,
	environmentFileMap map[string][]byte,
	currentContent []byte,
) {
	extension := filepath.Ext(filePath)

	for _, environment := range environments {
		if environment.Name != currentEnvironment {
			ui.Print(fmt.Sprintf("Enter content for file `%s` for the '%s' environment (Press any key to continue)", filePath, environment.Name))
			_, _, err := keyboard.GetSingleKey()
			if err != nil {
				errmsg := fmt.Sprintf("Failed to read user input (%s)", err.Error())
				println(errmsg)
				os.Exit(1)
				return
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
				return
			}

			environmentFileMap[environment.Name] = content
		}
	}
}

func checkFileAlreadyInCache(fileName string) bool {
	files := ctx.ListFilesFromCache()
	var found keystonefile.FileKey
	for _, file := range files {
		if file.Path == fileName {
			found = file
		}
	}
	if !reflect.ValueOf(found).IsZero() {
		ui.Print(`The file already exist but is not used.`)
		for _, env := range ctx.AccessibleEnvironments {
			content, err := ctx.GetFileContents(fileName, env.Name)

			ui.Print("\n")
			ui.Print("------------------" + env.Name + "------------------")
			ui.Print("\n")
			if err != nil {
				ui.Print("File not found for this environment")
			}

			ui.Print(string(content))
		}

		override := false

		if !skipPrompts {
			override = prompts.Confirm("Do you want to override the contents")
		}

		return !override
	}
	return false
}
