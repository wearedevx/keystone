package core

import (
	"encoding/base64"
	"fmt"
	"path"

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
		RemoveContents(path.Join(ctx.cacheDirPath(), environmentName))

		for _, file := range PayloadContent.Files {
			fileContent, _ := base64.StdEncoding.DecodeString(file.Value)
			CreateFileIfNotExists(path.Join(ctx.cacheDirPath(), environmentName, file.Path), string(fileContent))
		}

		envFilePath := path.Join(ctx.cacheDirPath(), environmentName, ".env")

		for _, secret := range PayloadContent.Secrets {

			if err := new(EnvFile).Load(envFilePath).Set(secret.Label, secret.Value).Dump().Err(); err != nil {
				err = FailedToUpdateDotEnv(envFilePath, err)
				fmt.Println(err.Error())
			}
			// CreateFileIfNotExists(path.Join(ctx.cacheDirPath(), environmentName, file.Path), string(fileContent))
		}
	}

	return MessageByEnvironments, nil
}
