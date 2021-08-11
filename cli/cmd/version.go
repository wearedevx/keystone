package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/wearedevx/keystone/cli/pkg/constants"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Displays the current CLI version",
	Long:    "Displays the current CLI version.",
	Example: "ks version",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("keystone-cli version %s\n", constants.Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
