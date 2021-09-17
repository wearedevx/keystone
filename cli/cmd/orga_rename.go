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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/spinner"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui"
)

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		argc := len(args)
		if argc != 2 {
			ui.PrintError(fmt.Sprintf("invalid number of arguments. Expected 2, got %d", argc))
			os.Exit(1)
		}

		organizationName := args[0]
		newName := args[1]

		sp := spinner.Spinner(" ")
		sp.Start()

		c, err := client.NewKeystoneClient()
		if err != nil {
			err.Print()
			os.Exit(1)
		}

		organizations, _ := c.Organizations().GetAll()

		foundOrga := models.Organization{}

		for _, orga := range organizations {
			if orga.Name == organizationName {
				foundOrga = orga
			}
		}

		if foundOrga.ID == 0 {
			ui.PrintError("You don't own an organization named %s", organizationName)
			ui.Print("To see organizations you own, use : $ ks orga")
			os.Exit(1)
		}

		foundOrga.Name = newName

		organization, updateErr := c.Organizations().UpdateOrganization(foundOrga)

		if updateErr != nil {
			ui.PrintError(updateErr.Error())
			os.Exit(1)
		}
		ui.PrintSuccess("Organization %s has been updated to %s", organizationName, organization.Name)
	},
}

func init() {
	orgaCmd.AddCommand(renameCmd)
}
