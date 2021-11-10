package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/archive"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/messages"
	"github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/ui/display"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore <path to archive>",
	Short: "Restores secrets and files from keystone created backup",
	Long: `Restores secrets and files from keystone created backup.
This will override all the data you have stored locally.`,
	Example: "ks restore keystone-backup-project-163492022.tar.gz",
	Args:    cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		if len(ctx.AccessibleEnvironments) < 3 {
			exit(kserrors.RestoreDenied(nil))
		}

		backupfile := args[0]
		if !utils.FileExists(backupfile) {
			exit(kserrors.FileDoesNotExist(backupfile, nil))
		}

		if password == "" {
			password = prompts.PasswordToDecrypt()
		}

		extractTarget := ctx.Wd

		if !skipPrompts {
			if !prompts.ConfirmDotKeystonDirRemoval() {
				fmt.Println("srietn")

				exit(nil)
			}
		}

		exitIfErr(os.RemoveAll(ctx.DotKeystonePath()))

		exitIfErr(
			archive.ExtractWithPassphrase(
				backupfile,
				extractTarget,
				password,
			),
		)

		ms := messages.NewMessageService(ctx)
		exitIfErr(
			ms.SendEnvironments(ctx.AccessibleEnvironments).Err(),
		)

		display.BackupRestored()
	},
}

func init() {
	RootCmd.AddCommand(restoreCmd)
	restoreCmd.Flags().StringVarP(&password, "password", "p", "", "password to encrypt backup with")
}
