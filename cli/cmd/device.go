package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui"
)

// deviceCmd represents the device command
var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "List all devices linked to your account.",
	Long:  `List all devices linked to your account.`,
	Run: func(_ *cobra.Command, _ []string) {
		c, kcErr := client.NewKeystoneClient()
		if kcErr != nil {
			fmt.Println(kcErr)
			os.Exit(1)
		}

		kf := keystonefile.KeystoneFile{}
		kf.Load(ctx.Wd)

		devices, err := c.Devices().GetAll()
		if err != nil {
			handleClientError(err)
		}

		printDeviceList(devices)

	},
}

func printDeviceList(devices []models.Device) {
	devStrings := []string{}

	for _, device := range devices {
		devStrings = append(devStrings, formatDevice(device))
	}

	ui.Print(
		ui.RenderTemplate(
			"device-list",
			`You have {{ .Len }} device(s) registered for this account:

{{ .Devices }}

To revoke access to one of these devices, use:
  $ ks device revoke
`,
			map[string]string{
				"Devices": strings.Join(devStrings, "\n"),
				"Len":     strconv.Itoa(len(devices)),
			},
		),
	)
}

func formatDevice(device models.Device) string {
	lastUsedAtString := ""

	if device.LastUsedAt.IsZero() {
		lastUsedAtString = "never used"
	} else {
		lastUsedAtString = fmt.Sprintf("last used on %s", device.CreatedAt.Format("2006/01/02"))
	}

	return fmt.Sprintf(
		"  - %s, %s, created at %s\n",
		device.Name,
		lastUsedAtString,
		device.CreatedAt.Format("2006/01/02"),
	)
}

func init() {
	RootCmd.AddCommand(deviceCmd)

}
