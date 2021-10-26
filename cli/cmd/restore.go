package cmd

import (
	"encoding/base64"
	"os"

	"github.com/cossacklabs/themis/gothemis/cell"
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/archive"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/messages"
	"github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/ui"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restores secrets and files from keystone created backup",
	Long: `Restores secrets and files from keystone created backup.
This will override all the data you have stored locally.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(ctx.AccessibleEnvironments) < 3 {
			exit(kserrors.RestoreDenied(nil))
		}

		backupfile := args[0]
		if !utils.FileExists(backupfile) {
			exit(kserrors.FileDoesNotExist(backupfile, nil))
		}

		if password == "" {
			password = prompts.StringInput("Password to decrypt backup", "")
		}

		extractTarget := ctx.Wd

		if !skipPrompts {
			ui.Print(ui.RenderTemplate("confirm files rm",
				`{{ CAREFUL }} You are about to remove the content of .keystone/ which contain all your local secrets and files.
This will override the changes you and other members made since the backup.
It will update other members secrets and files.`, map[string]string{}))
			if !prompts.Confirm("Continue") {
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

		ui.PrintSuccess("Backup restored: all your files and secrets have been replaced by the backup. They also have been sent to all members.")
	},
}

func init() {
	RootCmd.AddCommand(restoreCmd)
	restoreCmd.Flags().StringVarP(&password, "password", "p", "", "password to encrypt backup with")
}

func decryptBackup(backup []byte, password string) []byte {

	scell, err := cell.SealWithPassphrase(password)
	if err != nil {
		ui.PrintError(err.Error())
		os.Exit(1)
	}
	decrypted, err := scell.Decrypt([]byte(backup), nil)
	if err != nil {
		ui.PrintError("Cannot decrypt backup with this password.")
		os.Exit(1)
	}

	data, err := base64.StdEncoding.DecodeString(string(decrypted))
	return data
}
