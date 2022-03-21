package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/config"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui/display"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

// deviceCmd represents the device command
var deviceRevokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revokes access to one of your devices",
	Long:  `Revokes access to one of your devices.`,
	Run: func(_ *cobra.Command, _ []string) {
		c, kcErr := client.NewKeystoneClient()
		exitIfErr(kcErr)

		devices, err := c.Devices().GetAll()
		if err != nil {
			handleClientError(err)
			exit(err)

		}

		device := prompts.SelectDevice(devices)

		kf := keystonefile.KeystoneFile{}
		kf.Load(ctx.Wd)
		exitIfErr(kf.Err())

		if err = c.Devices().Revoke(device.UID); err != nil {
			handleClientError(err)
			exit(err)
		}

		exitIfErr(config.RevokeDevice())
		config.Write()

		display.DeviceRevokeSuccess()
	},
}

func init() {
	deviceCmd.AddCommand(deviceRevokeCmd)
}
