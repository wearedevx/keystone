package cmd

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/config"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/internal/spinner"
	. "github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/client/auth"
	core "github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
	"github.com/wearedevx/keystone/cli/ui/prompts"

	"github.com/spf13/cobra"
)

var projectName string
var organizationName string

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
		ctx := core.New(core.CTX_INIT)

		if err != nil {
			exit(kserrors.NoWorkingDirectory(err))
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
					exit(nil)
				}
			}
		} else {
			ksfile = keystonefile.NewKeystoneFile(
				currentfolder,
				models.Project{},
			)
		}

		var project models.Project
		c, err := client.NewKeystoneClient()
		// fmt.Printf("err: %+v %+v\n", err == nil, err)
		exitIfErr(err)

		// Ask for project name if keystone file doesn't exist.
		if ksfile.ProjectId == "" {
			err = createProject(c, &project, ksfile)
		} else {
			err = getProject(c, &project, ksfile)
		}

		if err != nil {
			if errors.Is(err, auth.ErrorUnauthorized) {
				config.Logout()
				err = kserrors.InvalidConnectionToken(err)
			} else {
				ui.PrintError(err.Error())
			}

			exit(err)
		}

		// Setup the local files
		exitIfErr(
			ctx.Init(project).Err(),
		)

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

	organizationID := addOrganizationToProject(c, project)
	// Remote Project Creation
	sp := spinner.Spinner("Creating project...")
	sp.Start()
	*project, err = c.Project("").Init(projectName, organizationID)
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

func addOrganizationToProject(c client.KeystoneClient, project *models.Project) uint {
	organizations, err := c.Organizations().GetAll()
	if err != nil {
		ui.PrintError("Error getting organizations: ", err.Error())
		os.Exit(1)
	}
	orga := models.Organization{}
	if organizationName == "" {
		orga = prompts.OrganizationsSelect(organizations)
		project.OrganizationID = orga.ID
		return orga.ID
	} else {
		for _, o := range organizations {
			if organizationName == o.Name {
				orga = o
			}
		}
		if orga.ID == 0 {
			ui.PrintError("Organization does not exist")
			os.Exit(1)
		}
		project.OrganizationID = orga.ID
		return orga.ID
	}
}

func init() {
	RootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&organizationName, "orga", "o", "", "identity provider. Either github or gitlab")
}
