package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/ci"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

// ciSendCmd represents the pushCi command
var ciSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Sends environment to a CI service",
	Long: `Sends environment to a CI service.

This command will send all your secrets and files followed by keystone to your CI service.

The CI service must have been setup using: ` + "`" + `ks ci add` + "`" + `
`,
	Example: `# To send the current environment:
ks ci send

# To send a specific environment:
ks ci send --env prod
`,
	Run: func(_ *cobra.Command, _ []string) {
		var environment models.Environment
		ctx.MustHaveEnvironment(currentEnvironment)

		fetch()

		for _, accessibleEnvironment := range ctx.AccessibleEnvironments {
			if accessibleEnvironment.Name == currentEnvironment {
				environment = accessibleEnvironment
			}
		}

		mustNotHaveMissingSecrets(environment)
		mustNotHaveMissingFiles(environment)

		ui.Print(
			"You are about to send the '%s' environment to your CI services.",
			environment,
		)
		if !prompts.Confirm("Continue") {
			os.Exit(0)
		}

		message, err := ctx.PrepareMessagePayload(environment)

		if err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		ciServices, err := ci.ListCiServices(ctx)
		if err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		for _, serviceDef := range ciServices {
			ciService, err := ci.GetCiService(serviceDef.Name, ctx, client.ApiURL)

			if err != nil {
				ui.PrintError(err.Error())
				os.Exit(1)
			}

			ciService.CheckSetup()
			if ciService.Error() != nil {
				ui.PrintError(ciService.Error().Error())
				os.Exit(1)
			}

			ciService.PushSecret(message, currentEnvironment)

			if ciService.Error() != nil {
				ui.PrintError(ciService.Error().Error())
				os.Exit(1)
			}

			ciService.PrintSuccess(currentEnvironment)
		}
	},
}

func init() {
	ciCmd.AddCommand(ciSendCmd)

	ciSendCmd.Flags().StringVar(&serviceName, "with", "", "Ci service name.")
}

func SelectCiService(ctx *core.Context) (ci.CiService, error) {
	var err error

	services, err := ci.ListCiServices(ctx)
	if err != nil {
		return nil, err
	}

	if len(services) == 0 {
		return nil, errors.New("You haven’t set up any CI service yet")
	}

	items := make([]string, len(services), len(services))

	for idx, service := range services {
		items[idx] = service.Name
	}

	if serviceName == "" {
		_, serviceName = prompts.Select(
			"Select a CI service",
			items,
		)
	}

	return ci.GetCiService(serviceName, ctx, client.ApiURL)
}

func mustNotHaveMissingSecrets(environment models.Environment) {
	missing, hasMissing := ctx.MissingSecretsForEnvironment(
		environment.Name,
	)

	if hasMissing {
		error := kserrors.RequiredSecretsAreMissing(missing, environment.Name, nil)
		error.Print()
		os.Exit(1)
	}
}

func mustNotHaveMissingFiles(environment models.Environment) {
	missing, hasMissing := ctx.MissingFilesForEnvironment(
		environment.Name,
	)

	if hasMissing {
		error := kserrors.RequiredFilesAreMissing(missing, environment.Name, nil)
		error.Print()
		os.Exit(1)
	}
}
