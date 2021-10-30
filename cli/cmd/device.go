package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui/display"
)

// deviceCmd represents the device command
var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Manages devices",
	Long: `Manages devices.

Used without arguments, lists devices this account has registered.
`,
	Example: "ks device",
	Run: func(_ *cobra.Command, _ []string) {
		c, kcErr := client.NewKeystoneClient()
		exitIfErr(kcErr)

		kf := keystonefile.KeystoneFile{}
		kf.Load(ctx.Wd)

		devices, err := c.Devices().GetAll()
		if err != nil {
			handleClientError(err)
			exit(err)
		}

		display.DeviceList(devices)
	},
}

func init() {
	RootCmd.AddCommand(deviceCmd)

}
