/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
	"github.com/wearedevx/keystone/cli/internal/backup"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/ui/display"
)

var backupSetup bool

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Creates a local backup of your secrets",
	Long: `Creates a local backup of your secrets
if a backup strategy has been setup. Sets one up if none exist.

The ` + "`" + `--setup` + "`" + ` allows to change the current settings.`,
	Args: cobra.NoArgs,
	Example: `# Manually perform a backup for the current project
ks backup

# Change the backup settings
ks backup --setup`,
	Run: func(cmd *cobra.Command, args []string) {
		bs := backup.NewBackupService(ctx)

		if backupSetup || !bs.IsSetup() {
			// setup a backup strategy
			exitIfErr(bs.Setup())
		}

		// Perform the backup
		backupFileFull, _, err := bs.Backup()
		if err != nil {
			exitIfErr(kserrors.BackupFailed(err))
		}

		display.BackupCreated(backupFileFull, false)
	},
}

func init() {
	RootCmd.AddCommand(backupCmd)

	// ks backup --setup
	backupCmd.
		Flags().
		BoolVar(
			&backupSetup,
			"setup",
			false,
			"Setup or change a backup strategy",
		)
}
