package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/archive"
	"github.com/wearedevx/keystone/cli/ui"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create a backup of your local secrets and files.",
	Long:  `Create a backup of your local secrets and files.`,
	Run: func(cmd *cobra.Command, args []string) {
		BACKUP_NAME := fmt.Sprintf(`./keystone-backup-%d.tar.gz`, time.Now().Unix())

		if err := archive.Tar(ctx.DotKeystonePath(), "./"); err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		if err := archive.Gzip("./.keystone.tar", "./"); err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		os.Rename("./.keystone.tar.gz", BACKUP_NAME)

		ui.PrintSuccess("Backup created : %s", BACKUP_NAME)
	},
}

func init() {
	RootCmd.AddCommand(backupCmd)

}
