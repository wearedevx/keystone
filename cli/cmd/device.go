package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/cli/internal/config"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/client/auth"
	"github.com/wearedevx/keystone/cli/ui"
)

// deviceCmd represents the device command
var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "List all devices linked to your account.",
	Long:  `List all devices linked to your account.`,
	Run: func(cmd *cobra.Command, args []string) {
		c, kcErr := client.NewKeystoneClient()

		if kcErr != nil {
			fmt.Println(kcErr)
			os.Exit(1)
		}

		kf := keystonefile.KeystoneFile{}
		kf.Load(ctx.Wd)

		devices, err := c.Devices().GetAll()

		if err != nil {
			if errors.Is(err, auth.ErrorUnauthorized) {
				config.Logout()
				kserrors.InvalidConnectionToken(err).Print()
			} else {
				kserrors.UnkownError(err).Print()
			}
			os.Exit(1)
		}
		ui.Print("You have %d device(s) registered for this account :", len(devices))

		fmt.Println()
		for _, device := range devices {
			fmt.Printf("  - %s, created at %s\n", device.Name, device.CreatedAt.Format("2006/01/02"))
		}
		fmt.Println()
		fmt.Println("To revoke access to one of these devices, use :\n  $ ks device revoke <device>")

	},
}

func init() {
	RootCmd.AddCommand(deviceCmd)

}
