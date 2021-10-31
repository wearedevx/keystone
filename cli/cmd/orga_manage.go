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
	"github.com/wearedevx/keystone/cli/ui"
)

// manageCmd represents the manage command
var manageCmd = &cobra.Command{
	Use:   "manage [organization_name]",
	Short: "Manage your subscription",
	Long: `Manage your subscription.
Gives you a link to update your payment method or cancel your plan.
`,
	Args: cobra.MaximumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		var organizationName string
		var err error

		kc, err := client.NewKeystoneClient()
		exitIfErr(err)

		o := kc.Organizations()

		organizationName = mustGetOrganizationName(o, args)

		url, err := o.GetManagementUrl(organizationName)
		exitIfErr(err)

		ui.Print(
			ui.RenderTemplate(
				"upgrade-url",
				`To manage your organization plan, visit the following link:

        {{.UpgradeURL }}`,
				map[string]string{
					"UpgradeURL": url,
				},
			),
		)
	},
}

func init() {
	orgaCmd.AddCommand(manageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// manageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// manageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
