package display

import "github.com/wearedevx/keystone/cli/ui"

func BackupCreated(backupName string) {
	ui.PrintSuccess("Backup created: %s", backupName)
}
