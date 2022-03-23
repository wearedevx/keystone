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
	"errors"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/backup"
	"github.com/wearedevx/keystone/cli/internal/config"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/ui/display"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restores a local backup",
	Long: `Restores a local backup.

A backup strategy must have been setup (see ` + "`" + `ks backup` + "`" + `).`,
	Args:    cobra.NoArgs,
	Example: `ks restore`,
	Run: func(cmd *cobra.Command, args []string) {
		projectName := ctx.GetProjectName()
		bs := backup.NewBackupService(ctx)

		if !bs.IsSetup() {
			exit(kserrors.BackupNotSetUp(nil))
		}

		err := bs.Restore()
		if errors.Is(err, backup.ErrorBackupNotFound) {
			_, bp := config.GetBackupStrategy()
			exit(kserrors.NoBackup(projectName, bp, nil))
		} else {
			exit(kserrors.RestoreFailed(err))
		}

		display.BackupRestored()
	},
}

func init() {
	RootCmd.AddCommand(restoreCmd)
}
