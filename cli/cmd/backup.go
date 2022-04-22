package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/archive"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/ui/display"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

var password string
var backupName string
var short bool

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Creates a backup of your local secrets and files.",
	Long: `Creates a backup of your local secrets and files.

Since we do not keep a copy of your secrets or files on our servers, 
it can be useful to regularly back them up to a secure location
to prevent losing them all if anything were to happen to your device.`,
	Run: func(_ *cobra.Command, _ []string) {
		var err error

		backupName = archive.GetBackupPath(
			ctx.Wd,
			ctx.GetProjectName(),
			backupName,
		)

		if password == "" {
			password = prompts.PasswordToEncrypt()
		}

		if err = archive.ArchiveWithPassphrase(
			ctx.DotKeystonePath(),
			backupName,
			password,
		); err != nil {
			exit(
				kserrors.CouldNotCreateArchive(err),
			)
		}

		display.BackupCreated(backupName, short)
	},
}

func init() {
	RootCmd.AddCommand(backupCmd)
	backupCmd.Flags().StringVarP(&password, "password", "p", "", "password to encrypt backup with")
	backupCmd.Flags().StringVarP(&backupName, "name", "n", "", "name of the backup file")
	backupCmd.Flags().BoolVar(&short, "short", false, "short output, for use in scrpits")
}
