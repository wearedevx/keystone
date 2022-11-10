/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/internal/secrets"
	"github.com/wearedevx/keystone/cli/internal/utils"
	core "github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui/display"
	"github.com/wearedevx/keystone/cli/ui/prompts"

	"github.com/spf13/cobra"
)

var addOptional bool = false

// secretAddCmd represents the set command
var secretAddCmd = &cobra.Command{
	Use:   "add <secret name> [secret value]",
	Short: "Adds a secret to all environments",
	Long: `Adds a secret to all environments.

Secrets are environment variables which value may vary
across environments, such as 'staging', 'prod',
and 'dev' environments.

The variable name will be added to all such environments,
you will be asked its value for each environment
`,
	Example: `# Add a secret ` + "`" + `PORT` + "`" + ` to all environments
# and set its value to 3000 for the current one.
ks secret add PORT 3000

# Add a secret ` + "`" + `PORT` + "`" + `without setting a default:
ks secret add PORT`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(_ *cobra.Command, args []string) {
		var err error
		var useCache bool

		secretName := args[0]
		secretValue := ""

		if len(args) == 2 {
			secretValue = args[1]
		}

		// Fetch messages first, because, if something have changed,
		// and we get the messages after, we throw an error and the user
		// loses its input
		mustFetchMessages()

		err = utils.CheckSecretContent(secretName)
		exitIfErr(err)

		ctx.MustHaveEnvironment(currentEnvironment)
		secretService := secrets.NewSecretService(ctx)

		if yes, values := secretService.IsSecretInCache(secretName); yes {
			display.SecretAlreadyExitsts(values)

			useCache = !prompts.ConfirmOverrideSecretValue(skipPrompts)
		}

		if !useCache {
			environmentValueMap, setErr := secretService.SetValuesForEnvironments(
				secretName,
				secretValue,
				ctx.AccessibleEnvironments,
				skipPrompts,
			)
			exitIfErr(setErr)

			changes, messageService := mustFetchMessages()
			flag := core.S_REQUIRED

			if addOptional {
				flag = core.S_OPTIONAL
			}

			exitIfErr(ctx.
				CompareNewSecretWithChanges(
					secretName,
					environmentValueMap,
					changes,
				).
				AddSecret(secretName, environmentValueMap, flag).
				Err())

			exitIfErr(messageService.
				SendEnvironments(ctx.AccessibleEnvironments).
				Err())
		} else {
			var ksfile keystonefile.KeystoneFile
			// Add new env key to keystone.yaml
			err = ksfile.
				Load(ctx.Wd).
				SetEnv(secretName, true).
				Save().
				Err()
			if err != nil {
				err = kserrors.FailedToUpdateKeystoneFile(err)
			}

			exit(err)
		}

		display.SecretIsSetForEnvironment(
			secretName,
			len(ctx.AccessibleEnvironments),
		)
	},
}

func init() {
	secretsCmd.AddCommand(secretAddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	secretAddCmd.Flags().
		BoolVarP(&addOptional, "optional", "o", false, "mark the secret as optional")
}
