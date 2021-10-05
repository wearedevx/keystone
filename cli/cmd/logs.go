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
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Shows activity logs for the current project",
	Long: `Shows activity logs for the current project.
This functionnality requires a paid plan:

ks orga upgrade <your-organization>
`,
	Example: `ks logs`,
	Run: func(_ *cobra.Command, _ []string) {
		var kerr *kserrors.Error

		ks, kerr := client.NewKeystoneClient()

		if kerr != nil {
			kerr.Print()
			os.Exit(1)
		}

		logs := ks.Logs()

		allTheLogs, err := logs.GetAll()

		if err != nil {
			// TODO: have better error messages
			ui.PrintError(err.Error())
			os.Exit(1)

		}

		printAllTheLogs(allTheLogs)
	},
}

func printAllTheLogs(logs []models.ActivityLogLite) {
	if len(logs) == 0 {
		ui.PrintStdErr("No logs to display")
		os.Exit(0)
	}

	for _, log := range logs {
		printLog(log)
	}
}

func printLog(log models.ActivityLogLite) {
	fmt.Printf(
		"[%s] %s on %s (%s): %s | %s\n",
		log.CreatedAt.Format("2006-12-29 01:34:59"),
		log.UserID,
		log.ProjectName,
		log.EnvironmentName,
		log.Action,
		formatForLog(log.Success, log.ErrorMessage),
	)
}

func formatForLog(s bool, message string) string {
	if s {
		return "SUCCESS"
	} else {
		return fmt.Sprintf("FAILURE: %s", message)
	}
}

func init() {
	RootCmd.AddCommand(logsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
