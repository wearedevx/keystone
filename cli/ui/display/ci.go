package display

import (
	"fmt"

	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/ui"
)

/// Displays the list of CI services configurations
func CiConfigurations(services []keystonefile.CiService) {
	if len(services) != 0 {
		ui.Print(ui.RenderTemplate("ci list", `
CI Services:{{ range $service := .Services }} 
 - {{ $service.Name }} ({{ $service.Type }}){{end}}`, struct {
			Services []keystonefile.CiService
		}{
			Services: services,
		}))
	}
}

func CiAdded() {
	ui.PrintSuccess("CI service added successfully")
}

func CiSecretsRemoved(environmentName string) {
	ui.PrintSuccess(
		fmt.Sprintf(
			"Secrets successfully removed from CI service, environment %s.",
			environmentName,
		),
	)
}

func CiNoSecretsForEnvironment(environmentName string) {
	ui.PrintSuccess(
		fmt.Sprintf(
			"No secret found for environment %s in CI service",
			environmentName,
		),
	)
}

func CiServiceSetupSuccessfully() {
	ui.PrintSuccess("CI service setup successfully")
}

func CiServiceRemoved(serviceName string) {
	ui.PrintSuccess("CI service '%s' successfully removed", serviceName)
}

func CiSecretSent(serviceName, environmentName string) {
	ui.PrintSuccess(
		fmt.Sprintf(
			`Secrets successfully sent to %s CI service, environment %s.
See https://github.com/wearedevx/keystone-action to use them.`,
			serviceName,
			environmentName,
		),
	)
}
