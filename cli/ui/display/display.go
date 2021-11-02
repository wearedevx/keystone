package display

import (
	"fmt"

	"github.com/wearedevx/keystone/api/pkg/models"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/pkg/constants"
	"github.com/wearedevx/keystone/cli/ui"
)

func BackupCreated(backupName string) {
	ui.PrintSuccess("Backup created: %s", backupName)
}

func BackupRestored() {
	ui.PrintSuccess("Backup restored: all your files and secrets have been replaced by the backup. They also have been sent to all members.")
}

func Version() {
	fmt.Printf("keystone-cli version %s\n", constants.Version)
}

func User(currentAccount models.User) {
	fmt.Println(currentAccount.UserID)
}

func Error(err error) bool {
	if err == nil {
		return false
	}

	if kserrors.IsKsError(err) {
		kserr := kserrors.AsKsError(err)
		if kserr == nil {
			return false
		}
		kserr.Print()
	} else {
		ui.PrintError(err.Error())
	}

	return true
}