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

// CiAdded function Message on CI service successfully added
func CiAdded() {
	ui.PrintSuccess("CI service added successfully")
}

// CiSecretsRemoved function Message on secrets successfully removed
// from the CI service
func CiSecretsRemoved(environmentName string) {
	ui.PrintSuccess(
		fmt.Sprintf(
			"Secrets successfully removed from CI service, environment %s.",
			environmentName,
		),
	)
}

// CiNoSecretsForEnvironment function Message when there are no secrets
// for environment in the CI service
func CiNoSecretsForEnvironment(environmentName string) {
	ui.PrintSuccess(
		fmt.Sprintf(
			"No secret found for environment %s in CI service",
			environmentName,
		),
	)
}

// CiServiceSetupSuccessfully function Message when CI setup happened successfully
func CiServiceSetupSuccessfully() {
	ui.PrintSuccess("CI service setup successfully")
}

// CiServiceRemoved function Message when CI service removal happened successfully
func CiServiceRemoved(serviceName string) {
	ui.PrintSuccess("CI service '%s' successfully removed", serviceName)
}

// CiSecretSent function Message when secrets were sent successfully
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
