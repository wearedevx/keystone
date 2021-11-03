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
	"github.com/spf13/cobra"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/spinner"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui/display"
)

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Renames an organization",
	Long: `Renames an organization.
`,
	Example: "ks orga rename my_orag my_orga",
	Args:    cobra.ExactArgs(2),
	Run: func(_ *cobra.Command, args []string) {
		var err error
		organizationName := args[0]
		newName := args[1]

		sp := spinner.Spinner(" ")
		sp.Start()

		c, err := client.NewKeystoneClient()
		exitIfErr(err)

		foundOrga, err := c.Organizations().GetByName(organizationName, client.OWNED_ONLY)
		if err != nil {
			handleClientError(err)
			exit(kserrors.YouDoNotOwnTheOrganization(organizationName, nil))
		}

		foundOrga.Name = newName

		organization, err := c.Organizations().UpdateOrganization(foundOrga)
		if err != nil {
			handleClientError(err)
			exit(err)
		}

		display.OrganizationRenamed(organizationName, organization.Name)
	},
}

func init() {
	orgaCmd.AddCommand(renameCmd)
}
