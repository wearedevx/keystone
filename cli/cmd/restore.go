package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/archive"
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

		backupfile := args[0]
		if !utils.FileExists(backupfile) {
			ui.PrintError(fmt.Sprintf("File does not exist : %s", backupfile))
			os.Exit(1)
		}

		if !prompts.Confirm("Sure you want to remove the content of .keystone/ with your backup") {
			os.Exit(0)
		}

		if err := os.RemoveAll(ctx.DotKeystonePath()); err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		if err := archive.UnGzip(backupfile, ctx.Wd); err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		if err := archive.Untar(path.Join(ctx.Wd, ".keystone.tar"), "."); err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		ui.PrintSuccess("Backup restored : all your local files and secrets have been replaced by the backup.")
	},
}

func init() {
	RootCmd.AddCommand(restoreCmd)
}
