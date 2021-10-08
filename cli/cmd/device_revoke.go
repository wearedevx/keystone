package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

// deviceCmd represents the device command
var deviceRevokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke access to one of your device.",
	Long:  `Revoke access to one of your device.`,
	Run: func(_ *cobra.Command, _ []string) {
		c, kcErr := client.NewKeystoneClient()
		if kcErr != nil {
			fmt.Println(kcErr)
			os.Exit(1)
		}

		devices, err := c.Devices().GetAll()
		if err != nil {
			handleClientError(err)
		}

		device := prompts.SelectDevice(devices)

		kf := keystonefile.KeystoneFile{}
		kf.Load(ctx.Wd)

		if err = c.Devices().Revoke(device.UID); err != nil {
			handleClientError(err)
		}

		ui.PrintSuccess("Device has been revoked and will no longer be updated with new secrets.")
		ui.Print("If you did this because your account has been compromised, make sure to change your secrets.")

	},
}

func init() {
	deviceCmd.AddCommand(deviceRevokeCmd)

}
