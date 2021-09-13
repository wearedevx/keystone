package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/cossacklabs/themis/gothemis/cell"
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/archive"
	"github.com/wearedevx/keystone/cli/ui"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

var password string
var backupName string

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create a backup of your local secrets and files.",
	Long:  `Create a backup of your local secrets and files.`,
	Run: func(cmd *cobra.Command, args []string) {
		if backupName == "" {
			backupName = path.Join(
				ctx.Wd,
				fmt.Sprintf(`keystone-backup-%s-%d.tar.gz`, ctx.GetProjectName(), time.Now().Unix()),
			)
		} else {
			backupName = path.Join(
				ctx.Wd,
				fmt.Sprintf(`%s.tar.gz`, backupName),
			)
		}

		if len(ctx.AccessibleEnvironments) < 3 {
			ui.PrintError(fmt.Sprintf("You don't have the permissions to create a backup."))
			os.Exit(1)
		}
		if password == "" {
			password = prompts.StringInput("Password to encrypt backup", "")
		}

		if err := archive.Tar(ctx.DotKeystonePath(), ctx.Wd); err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		if err := archive.Gzip(path.Join(ctx.Wd, ".keystone.tar"), ctx.Wd); err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		if err := os.Rename(path.Join(ctx.Wd, ".keystone.tar.gz"), backupName); err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		/* #nosec
		 * It is unlikely that BACKUP_NAME is uncontrolled
		 */
		contents, err := ioutil.ReadFile(backupName)

		if err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}
		encrypted := encryptBackup(contents, password)

		/* #nosec
		 * It is unlikely that BACKUP_NAME is uncontrolled
		 */
		if err := ioutil.WriteFile(backupName, encrypted, 0600); err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
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
	if err != nil {
		ui.PrintError(err.Error())
		os.Exit(1)
	}
	encrypted, err := scell.Encrypt([]byte(data), nil)
	if err != nil {
		ui.PrintError(err.Error())
		os.Exit(1)
	}

	return encrypted
}
