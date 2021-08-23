package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui"
)

// deviceCmd represents the device command
var deviceRevokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke access to one of your device.",
	Long:  `Revoke access to one of your device.`,
	Run: func(cmd *cobra.Command, args []string) {

		argc := len(args)
		if argc == 0 || argc > 1 {
			ui.PrintError(fmt.Sprintf("invalid number of arguments. Expected 1, got %d", argc))
			os.Exit(1)
		}
		deviceName := args[0]

		c, kcErr := client.NewKeystoneClient()

		if kcErr != nil {
			fmt.Println(kcErr)
			os.Exit(1)
		}

		kf := keystonefile.KeystoneFile{}
		kf.Load(ctx.Wd)

		err := c.Devices().Revoke(deviceName)

		// ui.PrintError(err)
		if err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}
		ui.PrintSuccess("Device has been revoked and will no longer be updated with new secrets.")

	},
}

func init() {
	deviceCmd.AddCommand(deviceRevokeCmd)

}
