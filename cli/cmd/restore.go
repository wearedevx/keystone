package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/cossacklabs/themis/gothemis/cell"
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/archive"
	"github.com/wearedevx/keystone/cli/internal/messages"
	"github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/ui"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore secrets and files from keystone created backup.",
	Long: `Restore secrets and files from keystone created backup.
This will override all the data you have stored locally.`,
	Run: func(cmd *cobra.Command, args []string) {
		argc := len(args)
		if argc == 0 || argc > 1 {
			ui.PrintError(fmt.Sprintf("invalid number of arguments. Expected 1, got %d", argc))
			os.Exit(1)
		}
		if len(ctx.AccessibleEnvironments) < 3 {
			ui.PrintError(fmt.Sprintf("You don't have the permissions to restore a backup."))
			os.Exit(1)
		}

		backupfile := args[0]
		if !utils.FileExists(backupfile) {
			ui.PrintError(fmt.Sprintf("File does not exist : %s", backupfile))
			os.Exit(1)
		}

		password := prompts.StringInput("Password to encrypt backup", "")

		ui.Print(ui.RenderTemplate("confirm files rm",
			`{{ CAREFUL }} You are about to remove the content of .keystone/ which contain all your local secrets and files.
This will override the changes you and other members made since the backup.
It will update other members secrets and files.`, map[string]string{}))
		if !prompts.Confirm("Continue") {
			os.Exit(0)
		}

		if err := os.RemoveAll(ctx.DotKeystonePath()); err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		contents, err := ioutil.ReadFile(backupfile)

		if err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}
		decrypted := decryptBackup(contents, password)

		decryptedPath := fmt.Sprintf("decrypted.tar.gz")

		ioutil.WriteFile(decryptedPath, decrypted, 0644)

		if err := archive.UnGzip(decryptedPath, ctx.Wd); err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		if err := archive.Untar(path.Join(ctx.Wd, ".keystone.tar"), "."); err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		// Remove temp files
		if err := os.Remove(decryptedPath); err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}
		if err := os.Remove(path.Join(ctx.Wd, ".keystone.tar")); err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		var printer = &ui.UiPrinter{}
		ms := messages.NewMessageService(ctx, printer)
		if err := ms.SendEnvironments(ctx.AccessibleEnvironments).Err(); err != nil {
			err.Print()
			os.Exit(1)
			return
		}

		ui.PrintSuccess("Backup restored : all your files and secrets have been replaced by the backup. They also have been sent to all members.")
	},
}

func init() {
	RootCmd.AddCommand(restoreCmd)
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
