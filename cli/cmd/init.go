package cmd

import (
	"errors"
	"os"
	"strings"

	"github.com/wearedevx/keystone/api/pkg/models"

	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/projects"
	"github.com/wearedevx/keystone/cli/internal/spinner"

	core "github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui/display"

	"github.com/spf13/cobra"
)

var (
	projectName      string
	organizationName string
)

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
		var err error
		projectName = strings.Join(args, " ")

		// Retrieve working directry
		currentfolder, err := os.Getwd()
		ctx = core.New(core.CTX_INIT)

		if err != nil {
			exit(kserrors.NoWorkingDirectory(err))
		}

		pservice := projects.NewProjectService(ctx, currentfolder, projectName)
		if errors.Is(pservice.Err(), projects.ErrorAlreadyInKeystoneProject) {
			exit(
				kserrors.AlreadyKeystoneProject(nil),
			)
		}
		exitIfErr(pservice.Err())

		project := models.Project{
			Name: projectName,
		}

		sp := spinner.Spinner("Creating project...")
		sp.Start()

		err = pservice.GetOrCreate(&project, organizationName).Err()

		sp.Stop()

		if err != nil {
			handleClientError(err)
			exit(err)
		}

		// Setup the local files
		exitIfErr(
			ctx.Init(project).Err(),
		)

		display.ProjectInitSuccess()
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(
		&organizationName,
		"orga",
		"o",
		"",
		"identity provider. Either github or gitlab",
	)
}
