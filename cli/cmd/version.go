package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/wearedevx/keystone/cli/pkg/constants"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("keystone-cli version %s\n", constants.Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
