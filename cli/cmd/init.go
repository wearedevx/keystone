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
	"os"
	"path"
	"strings"

	"github.com/wearedevx/keystone/api/pkg/models"
<<<<<<< HEAD
	"github.com/wearedevx/keystone/cli/internal/config"
=======
>>>>>>> 5248307 (feat(init): allow keystone files without a project id)
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/internal/spinner"
	. "github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/client/auth"
	core "github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"

	"github.com/spf13/cobra"
)

var projectName string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init <project-name>",
	Short: "Creates Keystone config files and directories",
	Long: `Creates Keystone config files and directories.

Created files and directories:
 - keystone.yaml: the project’s config,
 - .keystone:    cache and various files for internal use. 
                 automatically added to .gitignore
`,
	Example: "ks init my-awesome-project",
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("a project name cannot be empty")
		}
		return nil
	},
	Run: func(_ *cobra.Command, args []string) {
		var err *kserrors.Error
		projectName = strings.Join(args, " ")

		// Retrieve working directry
		currentfolder, osError := os.Getwd()
		ctx := core.New(core.CTX_INIT)

		if osError != nil {
			err = kserrors.NewError(
				"OS Error",
				"Error when retrieving working directory",
				map[string]interface{}{},
				osError,
			)
			err.Print()
			os.Exit(1)
		}

		var ksfile *keystonefile.KeystoneFile

		if keystonefile.ExistsKeystoneFile(currentfolder) {
			ksfile = new(keystonefile.KeystoneFile).Load(currentfolder)

			// If there is already a keystone file around here,
			// inform the user they are in a keystone project
			if ksfile.ProjectId != "" && ksfile.ProjectName != projectName {
				// check if .keystone directory too
				if DirExists(path.Join(ctx.Wd, ".keystone")) {
					kserrors.AlreadyKeystoneProject(errors.New("")).Print()
					os.Exit(0)
				}
			}
		} else {
			ksfile = keystonefile.NewKeystoneFile(
				currentfolder,
				models.Project{},
			)
		}

		var project models.Project
		var initErr error
		c, err := client.NewKeystoneClient()
		if err != nil {
			err.Print()
			os.Exit(1)
		}

		// Ask for project name if keystone file doesn't exist.
		if ksfile.ProjectId == "" {
			initErr = createProject(c, &project, ksfile)
		} else {
			initErr = getProject(c, &project, ksfile)
		}

		if initErr != nil {
			if errors.Is(initErr, auth.ErrorUnauthorized) {
<<<<<<< HEAD
				config.Logout()
=======
>>>>>>> 5248307 (feat(init): allow keystone files without a project id)
				kserrors.InvalidConnectionToken(initErr).Print()
			} else {
				ui.PrintError(initErr.Error())
			}

			os.Exit(1)
		}

		// Setup the local files
		if err = ctx.Init(project).Err(); err != nil {
			err.Print()
			os.Exit(1)
		}

		ui.Print(ui.RenderTemplate("Init Success", `
{{ .Message | box | bright_green | indent 2 }}

{{ .Text | bright_black | indent 2 }}`, map[string]string{
			"Message": "All done!",
			"Text": `You can start adding environment variable with:
  $ ks secret add VARIABLE value

Load them with:
  $ eval $(ks source)

If you need help with anything:
  $ ks help [command]

`,
		}))
	},
}

func createProject(
	c client.KeystoneClient,
	project *models.Project,
	ksfile *keystonefile.KeystoneFile,
) (err error) {
	// Remote Project Creation
	sp := spinner.Spinner("Creating project...")
	sp.Start()
	*project, err = c.Project("").Init(projectName)
	sp.Stop()

	// Handle invalid token
	if err != nil {
		return err
	}

	// Update the ksfile
	// So that it keeps secrets and files
	// if the file exited without a project-id
	ksfile.ProjectId = project.UUID
	ksfile.ProjectName = project.Name

	if err := ksfile.Save().Err(); err != nil {
		return err
	}

	return nil
}

func getProject(
	c client.KeystoneClient,
	project *models.Project,
	ksfile *keystonefile.KeystoneFile,
) (err error) {
	// We have a keystone.yaml with a project-id, but no .keystone dir
	// So the project needs to be re-created.

	// Reconstruct the project
	// id and name are in the keystone file
	project.UUID = ksfile.ProjectId
	project.Name = ksfile.ProjectName

	// But environment data is still on the server
	sp := spinner.Spinner("Fetching project informations...")
	sp.Start()
	environments, err := c.Project(project.UUID).GetAccessibleEnvironments()
	sp.Stop()

	if err != nil {
		return err
	}

	project.Environments = environments

	return nil
}

func init() {
	RootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
