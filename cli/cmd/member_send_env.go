package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/wearedevx/keystone/api/pkg/models"
	kerrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/messages"
	"github.com/wearedevx/keystone/cli/ui"

	// . "github.com/wearedevx/keystone/cli/ui"

	"github.com/spf13/cobra"
)

var member string
var allEnv = false

// initCmd represents the init command
var sendEnvCmd = &cobra.Command{
	Use:   "send-env <member id>",
	Short: "Sends secrets and files from current environment to member.",
	Long: `Sends secrets and files from current environment to member.

If a member hasn't received secrets and files last time someone sent an update,
it can be done again with this command.
`,
	Example: `ks member send-env john@gitlab`,

	Args: func(_ *cobra.Command, args []string) error {
		r := regexp.MustCompile(`[\w-_.]+@(gitlab|github)`)
		argc := len(args)

		if argc == 0 || argc > 1 {
			return fmt.Errorf("invalid number of arguments. Expected 1, got %d", argc)
		}

		member = args[0]

		if !r.Match([]byte(member)) {
			return fmt.Errorf("invalid member id: %s", member)
		}

		return nil
	},
	Run: func(_ *cobra.Command, _ []string) {
		var err *kerrors.Error

		ctx.MustHaveEnvironment(currentEnvironment)

		var printer = &ui.UiPrinter{}
		ms := messages.NewMessageService(ctx, printer)
		ms.GetMessages()

		if err = ms.Err(); err != nil {
			err.Print()
			os.Exit(1)
		}

		environments := make([]models.Environment, 0)

		if allEnv {
			environments = ctx.AccessibleEnvironments
		} else {
			environments = append(environments, models.Environment{Name: currentEnvironment})

		}

		for i, env := range environments {
			localEnvironment := ctx.LoadEnvironmentsFile().GetByName(env.Name)

			environments[i] = models.Environment{
				Name:          localEnvironment.Name,
				VersionID:     localEnvironment.VersionID,
				EnvironmentID: localEnvironment.EnvironmentID,
			}
		}

		if err = ms.SendEnvironmentsToOneMember(environments, member).Err(); err != nil {
			err.Print()
			return
		}

	},
}

func init() {
	memberCmd.AddCommand(sendEnvCmd)

	sendEnvCmd.Flags().StringVar(&member, "all", "a", "Member to send env to.")
	sendEnvCmd.Flags().BoolVarP(&allEnv, "all-env", "", false, "Send secrets from all environments to member.")
}
