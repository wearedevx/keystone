package core

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	. "github.com/wearedevx/keystone/cli/internal/envfile"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	. "github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/pkg/client"

	"github.com/wearedevx/keystone/api/pkg/models"
	"gopkg.in/yaml.v2"
)

func (ctx *Context) SaveMessages(MessageByEnvironments models.GetMessageByEnvironmentResponse) (models.GetMessageByEnvironmentResponse, error) {

	for environmentName, environment := range MessageByEnvironments.Environments {
		var PayloadContent = models.MessagePayload{}

		if err := yaml.Unmarshal(environment.Message.Payload, &PayloadContent); err != nil {
			panic(err)
		}

		// Remove content of cache directory to ensure old files are deleted
		RemoveContents(ctx.CachedEnvironmentPath(environmentName))

		for _, file := range PayloadContent.Files {
			fileContent, _ := base64.StdEncoding.DecodeString(file.Value)
			CreateFileIfNotExists(path.Join(ctx.CachedEnvironmentFilesPath(environmentName), file.Path), string(fileContent))
		}

		envFilePath := ctx.CachedEnvironmentDotEnvPath(environmentName)

		for _, secret := range PayloadContent.Secrets {
			if err := new(EnvFile).Load(envFilePath).Set(secret.Label, secret.Value).Dump().Err(); err != nil {
				err = kserrors.FailedToUpdateDotEnv(envFilePath, err)
				// fmt.Println(err.Error())
			}
			// CreateFileIfNotExists(path.Join(ctx.cacheDirPath(), environmentName, file.Path), string(fileContent))
		}
	}

	return MessageByEnvironments, nil
}

// Return PayloadContent, with secrets and files of current environment.
func (ctx *Context) PrepareMessagePayload(environment models.Environment) (models.MessagePayload, error) {
	var PayloadContent = models.MessagePayload{
		Files:   make([]models.File, 0),
		Secrets: make([]models.SecretVal, 0),
	}

	var err error

	errors := make([]string, 0)

	for _, secret := range ctx.ListSecrets() {
		PayloadContent.Secrets = append(PayloadContent.Secrets, models.SecretVal{
			Label: secret.Name,
			Value: string(secret.Values[EnvironmentName(environment.Name)]),
		})
	}

	cachePath := ctx.cacheDirPath()
	envCachePath := path.Join(cachePath, environment.Name)

	for _, file := range ctx.ListFiles() {
		filePath := path.Join(envCachePath, file.Path)
		fileContent, err := ioutil.ReadFile(filePath)

		if err != nil {
			errors = append(errors, err.Error())
		}

		PayloadContent.Files = append(PayloadContent.Files, models.File{
			Path:  file.Path,
			Value: base64.StdEncoding.EncodeToString(fileContent),
		})
	}

	if len(errors) > 0 {
		err = fmt.Errorf(strings.Join(errors, "\n"))
	}

	return PayloadContent, err
}

func (ctx *Context) FetchNewMessages(result *models.GetMessageByEnvironmentResponse) error {
	// Get keystone key.
	c, kcErr := client.NewKeystoneClient()

	if kcErr != nil {
		kcErr.Print()
		os.Exit(1)
	}

	projectID := ctx.GetProjectID()

	*result, _ = c.Messages().GetMessages(projectID)

	return nil
}

func (ctx *Context) WriteNewMessages(messagesByEnvironments models.GetMessageByEnvironmentResponse) error {
	c, kcErr := client.NewKeystoneClient()

	if kcErr != nil {
		kcErr.Print()
		os.Exit(1)
	}

	ctx.SaveMessages(messagesByEnvironments)

	if err := ctx.Err(); err != nil {
		err.Print()
		return err
	}
	for environmentName, environment := range messagesByEnvironments.Environments {
		messageID := environment.Message.ID
		if messageID != 0 {
			fmt.Println("Environment", environmentName, "updated ✔")
			response, _ := c.Messages().DeleteMessage(environment.Message.ID)
			if !response.Success {
				fmt.Println("Can't delete message", response.Error)
			}
		} else {
			environmentChanged := ctx.EnvironmentVersionHasChanged(environmentName, environment.VersionID)
			if environmentChanged {
				fmt.Println("Environment", environmentName, "has changed but no message available. Ask someone to push their secret ⨯")
			} else {
				fmt.Println("Environment", environmentName, "up to date ✔")
			}
		}
	}
	return nil
}

func (ctx *Context) CompareNewSecretWithMessages(secretName string, newSecret map[string]string, fetchedSecrets models.GetMessageByEnvironmentResponse, localSecrets []Secret) *kserrors.Error {

	// where are stored changed values
	environmentValueMap := make(map[string]string)

	for environmentName, message := range fetchedSecrets.Environments {
		var PayloadContent = models.MessagePayload{}
		if err := yaml.Unmarshal(message.Message.Payload, &PayloadContent); err != nil {
			panic(err)
		}

		for _, secret := range PayloadContent.Secrets {
			// Get Secret we want to change in messages
			if secret.Label == secretName {
				// Get local value to see if it has changed in fetchedSecrets
				for _, localSecret := range localSecrets {
					if secretName == localSecret.Name {
						// Compare local value with value from message
						localSecretValue := string(localSecret.Values[EnvironmentName(environmentName)])
						if secret.Value != localSecretValue {
							environmentValueMap[environmentName] = secret.Value
						}
					}
				}
			}
		}
	}

	if len(environmentValueMap) > 0 {
		environmentValueMapString := ""
		for environment, value := range environmentValueMap {
			environmentValueMapString += fmt.Sprintf("Value in '%s' is '%s'.\n", environment, value)
		}
		return kserrors.SecretHasChanged(secretName, environmentValueMapString, nil)
	}
	return nil
}
