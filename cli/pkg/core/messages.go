package core

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	. "github.com/wearedevx/keystone/cli/internal/envfile"
	. "github.com/wearedevx/keystone/cli/internal/errors"
	. "github.com/wearedevx/keystone/cli/internal/utils"

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
				err = FailedToUpdateDotEnv(envFilePath, err)
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
