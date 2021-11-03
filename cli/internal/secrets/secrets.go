package secrets

import (
	"fmt"
	"gitEnterValuee/cli/ui/display"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui/cli/ui/display"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

type SecretService struct {
	err error
	ctx *core.Context
}

func NewSecretService(ctx *core.Context) *SecretService {
	return &SecretService{
		err: nil,
		ctx: ctx,
	}
}

func (ss *SecretService) Err() error {
	return ss.err
}

func (ss *SecretService) IsSecretInCache(secretName string) (inCache bool, _ map[core.EnvironmentName]core.SecretValue) {
	{
		var found core.Secret
		values := make(map[core.EnvironmentName]core.SecretValue)
		secrets := ss.ctx.ListSecretsFromCache()

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
}

func (ss *SecretService) SetValuesForEnvironments(
	secretName, secretValue string,
	environments []models.Environment,
	skipPrompts bool,
) (map[string]string, error) {
	{
		environmentValueMap := make(map[string]string)
		// Ask value for each env
		if !skipPrompts {
			display.EnterValue(secretName)

			for _, environment := range environments {
				// multiline edit
				if strings.Contains(secretValue, "\n") {
					result, err := utils.CaptureInputFromEditor(
						utils.GetPreferredEditorFromEnvironment,
						"",
						defaultContentForSetValues(
							secretName,
							environment.Name,
						),
					)
					stringResult := trimCommentsFromSecretInput(string(result))

					if err != nil {
						if err.Error() != "^C" {
							return nil, err
						}
						os.Exit(0)
					}

					environmentValueMap[environment.Name] = stringResult
				} else {
					newValue := prompts.StringInput(
						environment.Name,
						secretValue,
					)
					environmentValueMap[environment.Name] = strings.TrimSpace(newValue)
				}
			}

		} else {
			for _, environment := range environments {
				environmentValueMap[environment.Name] = strings.TrimSpace(secretValue)
			}
		}

		return environmentValueMap, nil
	}
}

func defaultContentForSetValues(secretName, environmentName string) string {
	return fmt.Sprintf(
		`%s


# Enter value for secret %s on environment %s`,
		secretName, secretName, environmentName)
}

func trimCommentsFromSecretInput(input string) string {
	input = regexp.MustCompile(`#.*$`).
		ReplaceAllString(strings.TrimSpace(input), "")

	input = regexp.MustCompile(`[\t\r\n]+`).
		ReplaceAllString(strings.TrimSpace(input), "\n")

	return strings.TrimSpace(string(input))
}
