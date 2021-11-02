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
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui/display"
)

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade [organization-name]",
	Short: "Ugrade an organization plan to a paid one",
	Long: `Ugrade an organization plan to a paid one.
To benefit from all of Keystone's fonctionalities, you should upgrade your
organization to a paid plan using this command.
`,
	Example: "ks orga upgrade my-organization",
	Args:    cobra.MaximumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		var organizationName string
		var err error

		kc, err := client.NewKeystoneClient()
		exitIfErr(err)

		o := kc.Organizations()

		organizationName = mustGetOrganizationName(o, args)

		url, err := o.GetUpgradeUrl(organizationName)
		exitIfErr(err)

		display.UpgradeUrl(url)
	},
}

func init() {
	orgaCmd.AddCommand(upgradeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upgradeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upgradeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
