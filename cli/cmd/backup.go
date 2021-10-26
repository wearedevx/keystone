package cmd

import (
	"encoding/base64"
	"io/ioutil"

	"github.com/cossacklabs/themis/gothemis/cell"
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/archive"
	"github.com/wearedevx/keystone/cli/internal/crypto"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/ui"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

var password string
var backupName string

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Creates a backup of your local secrets and files.",
	Long: `Creates a backup of your local secrets and files.

Since we do not keep a copy of your secrets or files on our servers, 
it can be useful to regularly back them up to a secure location
to prevent losing them all if anything were to happen to your device.`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		if len(ctx.AccessibleEnvironments) < 3 {
			exit(kserrors.BackupDenied(nil))
		}

		backupName = archive.GetBackupPath(
			ctx.Wd,
			ctx.GetProjectName(),
			backupName,
		)

		if password == "" {
			password = prompts.StringInput("Password to encrypt backup", "")
		}

		if err = archive.Archive(
			ctx.DotKeystonePath(),
			ctx.Wd,
			backupName,
		); err != nil {
			err = kserrors.CouldNotCreateArchive(err)
		}
		exitIfErr(err)

		/* #nosec
		 * It is unlikely that BACKUP_NAME is uncontrolled
		 */
		encrypted, err := crypto.EncryptFile(backupName, password)
		if err != nil {
			exit(kserrors.EncryptionFailed(err))
		}

		/* #nosec
		 * It is unlikely that BACKUP_NAME is uncontrolled
		 */
		if err := ioutil.WriteFile(backupName, encrypted, 0600); err != nil {
			exit(kserrors.FailedToWriteBackup(err))
		}

		ui.PrintSuccess("Backup created : %s", backupName)
	},
}

func init() {
	RootCmd.AddCommand(backupCmd)
	backupCmd.Flags().StringVarP(&password, "password", "p", "", "password to encrypt backup with")
	backupCmd.Flags().StringVarP(&backupName, "name", "n", "", "name of the backup file")
}

func encryptBackup(backup []byte, password string) []byte {
	data := base64.StdEncoding.EncodeToString(backup)

	scell, err := cell.SealWithPassphrase(password)
	exitIfErr(err)

	encrypted, err := scell.Encrypt([]byte(data), nil)
	exitIfErr(err)

	return encrypted
}
