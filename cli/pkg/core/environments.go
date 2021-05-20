package core

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/config"
	. "github.com/wearedevx/keystone/cli/internal/envfile"
	. "github.com/wearedevx/keystone/cli/internal/errors"
	. "github.com/wearedevx/keystone/cli/internal/keystonefile"
	. "github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/pkg/client"
)

func (ctx *Context) CurrentEnvironment() string {
	if ctx.Err() != nil {
		return ""
	}

	bytes, err := ioutil.ReadFile(ctx.environmentFilePath())

	if err != nil {
		ctx.setError(CannotReadEnvironment(ctx.environmentFilePath(), err))
	}

	return strings.Trim(string(bytes), "\n")
}

func (ctx *Context) ListEnvironments() []string {
	if ctx.Err() != nil {
		return []string{}
	}
	envs := make([]string, 0)
	envs = append(envs, "default")

	cacheDir := ctx.cacheDirPath()
	contents, err := ioutil.ReadDir(cacheDir)

	if err != nil {
		ctx.setError(UnkownError(err))
		return envs
	}

	for _, file := range contents {
		if !file.IsDir() {
			continue
		}

		envname := file.Name()
		contained := false

		for _, e := range envs {
			if e == envname {
				contained = true
				break
			}
		}

		if !contained {
			envs = append(envs, file.Name())
		}
	}

	return envs
}

func (ctx *Context) CreateEnvironment(name string) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	if !ctx.HasEnvironment(name) {
		newEnvDir := ctx.CachedEnvironmentPath(name)
		err := os.MkdirAll(newEnvDir, 0o755)

		if err != nil {
			ctx.setError(CannotCreateDirectory(newEnvDir, err))
		}
	} else {
		ctx.setError(EnvironmentAlreadyExists(name, nil))
	}

	return ctx
}

func (ctx *Context) RemoveEnvironment(name string) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	if current := ctx.CurrentEnvironment(); current == name {
		return ctx.setError(CannotRemoveCurrentEnvironment(name, nil))
	}

	if ctx.HasEnvironment(name) {
		envDir := ctx.CachedEnvironmentPath(name)
		err := os.RemoveAll(envDir)

		if err != nil {
			return ctx.setError(CannotRemoveDirectory(envDir, err))
		}
	} else {
		return ctx.setError(EnvironmentDoesntExist(name, strings.Join(ctx.ListEnvironments(), ", "), nil))
	}

	return ctx
}

func (ctx *Context) SetCurrent(name string) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	if ctx.HasEnvironment(name) {
		dotEnvFilePath := ctx.CachedEnvironmentDotEnvPath(name)
		currentDotEnvFilePath := ctx.CachedDotEnvPath()

		err := CopyFile(dotEnvFilePath, currentDotEnvFilePath)

		if err != nil {
			return ctx.setError(CopyFailed(dotEnvFilePath, currentDotEnvFilePath, err))
		}

		err = ioutil.WriteFile(ctx.environmentFilePath(), []byte(name), 0o644)

		if err != nil {
			return ctx.setError(FailedToSetCurrentEnvironment(name, ctx.environmentFilePath(), err))
		}

		ctx.FilesUseEnvironment(name)

	} else {
		return ctx.setError(EnvironmentDoesntExist(name, strings.Join(ctx.ListEnvironments(), ", "), nil))
	}

	return ctx
}

func (ctx *Context) SetAllSecrets(name string, secrets map[string]string) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	if ctx.HasEnvironment(name) {
		dotEnvPath := ctx.CachedEnvironmentDotEnvPath(name)

		if err := new(EnvFile).Load(dotEnvPath).SetData(secrets).Dump().Err(); err != nil {
			return ctx.setError(FailedToUpdateDotEnv(dotEnvPath, err))
		}

	} else {
		return ctx.setError(EnvironmentDoesntExist(name, strings.Join(ctx.ListEnvironments(), ", "), nil))
	}

	return ctx
}

func (ctx *Context) GetAllSecrets(envName string) map[string]string {
	emptyMap := map[string]string{}

	if ctx.Err() != nil {
		return emptyMap
	}

	if ctx.HasEnvironment(envName) {
		dotEnvPath := ctx.CachedEnvironmentDotEnvPath(envName)

		envFile := new(EnvFile).Load(dotEnvPath)

		if err := envFile.Err(); err != nil {
			ctx.setError(FailedToReadDotEnv(dotEnvPath, err))
			return emptyMap
		}

		return envFile.GetData()
	} else {
		ctx.setError(EnvironmentDoesntExist(envName, strings.Join(ctx.ListEnvironments(), ", "), nil))
	}

	return emptyMap
}

func (ctx *Context) HasEnvironment(name string) bool {
	if ctx.Err() != nil {
		return false
	}

	return DirExists(ctx.CachedEnvironmentPath(name))
}

func (ctx *Context) MustHaveEnvironment(name string) {
	if !ctx.HasEnvironment(name) {
		EnvironmentDoesntExist(name, strings.Join(ctx.ListEnvironments(), ", "), nil).Print()
		os.Exit(0)
	}
}

func (ctx *Context) EnvironmentVersion() string {
	environments := ctx.EnvironmentsFromConfig()
	currentEnvironment := ctx.CurrentEnvironment()

	for _, e := range environments {
		if e.Name == currentEnvironment {
			return e.VersionID
		}
	}
	return ""
}

func (ctx *Context) EnvironmentVersionByName(name string) string {
	environments := ctx.EnvironmentsFromConfig()

	for _, e := range environments {
		if e.Name == name {
			return e.VersionID
		}
	}
	return ""
}

func (ctx *Context) EnvironmentID() string {
	return ctx.getCurrentEnvironmentId()
	// environments := ctx.EnvironmentsFromConfig()
	// currentEnvironment := ctx.CurrentEnvironment()

	// for _, e := range environments {
	// 	if e.Name == currentEnvironment {
	// 		return e.EnvironmentID
	// 	}
	// }
	// return ""
}

func (ctx *Context) EnvironmentsFromConfig() []Env {

	ksfile := new(KeystoneFile).Load(ctx.Wd)
	return ksfile.Environments
}

// Push current environnement.
// Post to server []MessageToWritePayload.
func (ctx *Context) PushEnv() error {
	// var result client.GenericResponse

	// Get public keyrs
	c, kcErr := client.NewKeystoneClient()

	if kcErr != nil {
		return kcErr
	}

	environmentId := ctx.EnvironmentID()

	userPublicKeys, err := c.Users().GetEnvironmentPublicKeys(environmentId)

	if err != nil {
		return err
	}

	PayloadContent, err := ctx.PrepareMessagePayload()

	if err != nil {
		return err
	}

	messagesToWrite := models.MessagesToWritePayload{
		Messages: make([]models.MessageToWritePayload, 0),
	}

	account, _ := config.GetCurrentAccount()

	// Create one message per user
	for _, userPublicKey := range userPublicKeys.Keys {
		if userPublicKey.UserID != account.UserID {
			// TODO: encrypt payload with recipient public key
			// crypto.EncryptForUser()
			var payload string
			PayloadContent.Serialize(&payload)

			RecipientID, _ := strconv.ParseUint(userPublicKey.UserID, 10, 64)
			RecipientIDUint := uint(RecipientID)

			messagesToWrite.Messages = append(messagesToWrite.Messages, models.MessageToWritePayload{
				Payload:       []byte(payload),
				UserID:        userPublicKey.UserID,
				RecipientID:   RecipientIDUint,
				EnvironmentID: environmentId,
			})
		}
	}

	_, err = c.Messages().SendMessages(environmentId, messagesToWrite)

	if err != nil {
		return err
	}

	// TODO
	// Set new version id

	return nil
}

func (ctx *Context) EnvironmentVersionHasChanged(name string, environmentID string) bool {
	currentVersion := ctx.EnvironmentVersionByName(name)
	return currentVersion != environmentID
}
