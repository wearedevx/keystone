package display

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/ui"
)

// DeviceList function Prints a list of devices
func DeviceList(devices []models.Device) {
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

// DeviceRevokeSuccess function Message when device is revoked
func DeviceRevokeSuccess() {
	ui.PrintSuccess(
		"Device has been revoked and will no longer be updated with new secrets.",
	)
	ui.Print(
		"If you did this because your account has been compromised, make sure to change your secrets.",
	)
}

func formatDevice(device models.Device) string {
	lastUsedAtString := ""

	if device.LastUsedAt.IsZero() {
		lastUsedAtString = "never used"
	} else {
		lastUsedAtString = fmt.Sprintf("last used on %s", device.LastUsedAt.Format("2006/01/02"))
	}

	return fmt.Sprintf(
		"  - %s, %s, created on %s\n",
		device.Name,
		lastUsedAtString,
		device.CreatedAt.Format("2006/01/02"),
	)
}
