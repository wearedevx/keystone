package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/ui/display"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Displays the current CLI version",
	Long:    "Displays the current CLI version.",
	Example: "ks version",
	Run: func(_ *cobra.Command, _ []string) {
		display.Version()
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
