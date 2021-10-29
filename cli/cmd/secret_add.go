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
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/environments"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/internal/utils"
	core "github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui"
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

		err = utils.CheckSecretContent(secretName)
		exitIfErr(err)

		ctx.MustHaveEnvironment(currentEnvironment)

		if yes, values := checkSecretAlreadyInCache(secretName); yes {
			// TODO: that printing part should be in the ui packag
			ui.Print(`The secret already exist. Values are:`)
			for env, value := range values {
				ui.Print(`%s: %s`, env, value)
			}

			useCache = !prompts.ConfirmOverrideSecretValue(skipPrompts)
		}

		if !useCache {
			es := environments.NewEnvironmentService(ctx)
			exitIfErr(es.Err())

			environmentValueMap := setValuesForEnvironments(
				secretName,
				secretValue,
				ctx.AccessibleEnvironments,
			)

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

		ui.PrintSuccess("Variable '%s' is set for %d environment(s)", secretName, len(ctx.AccessibleEnvironments))
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
	secretAddCmd.Flags().BoolVarP(&addOptional, "optional", "o", false, "mark the secret as optional")
}

func setValuesForEnvironments(secretName string, secretValue string, accessibleEnvironments []models.Environment) map[string]string {

	environmentValueMap := make(map[string]string)
	// Ask value for each env
	if !skipPrompts {
		ui.Print(ui.RenderTemplate("ask new value for environment", `
Enter a values for {{ . }}:`, secretName))

		for _, environment := range accessibleEnvironments {

			// multiline edit
			if strings.Contains(secretValue, "\n") {
				var defaultContent strings.Builder

				defaultContent.WriteString(secretValue)
				defaultContent.WriteRune('\n')
				defaultContent.WriteRune('\n')
				defaultContent.WriteRune('\n')
				defaultContent.WriteString("# Enter value for secret ")
				defaultContent.WriteString(secretName)
				defaultContent.WriteString(" on environment ")
				defaultContent.WriteString(environment.Name)

				result, err := utils.CaptureInputFromEditor(
					utils.GetPreferredEditorFromEnvironment,
					"",
					defaultContent.String(),
				)
				stringResult := string(result)

				// remove blank line and comment from secret
				stringResult = regexp.MustCompile(`#.*$`).ReplaceAllString(strings.TrimSpace(stringResult), "")
				stringResult = regexp.MustCompile(`[\t\r\n]+`).ReplaceAllString(strings.TrimSpace(stringResult), "\n")

				if err != nil {
					if err.Error() != "^C" {
						ui.PrintError(err.Error())
						os.Exit(1)
					}
					os.Exit(0)
				}

				environmentValueMap[environment.Name] = strings.Trim(string(stringResult), " ")
			} else {
				environmentValueMap[environment.Name] = prompts.StringInput(
					environment.Name,
					secretValue,
				)
			}
		}

	} else {
		for _, environment := range accessibleEnvironments {
			environmentValueMap[environment.Name] = strings.Trim(secretValue, " ")
		}

	}

	return environmentValueMap
}

// TODO: should be a core function
func checkSecretAlreadyInCache(secretName string) (inCache bool, _ map[core.EnvironmentName]core.SecretValue) {
	var found core.Secret
	values := make(map[core.EnvironmentName]core.SecretValue)
	secrets := ctx.ListSecretsFromCache()

	for _, secret := range secrets {
		if secret.Name == secretName {
			found = secret
		}
	}

	inCache = !reflect.ValueOf(found).IsZero()

	if inCache {
		values = found.Values
	}

	return inCache, values
}
