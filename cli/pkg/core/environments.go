package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/config"
	. "github.com/wearedevx/keystone/cli/internal/envfile"
	. "github.com/wearedevx/keystone/cli/internal/environmentsfile"
	. "github.com/wearedevx/keystone/cli/internal/errors"
	. "github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/pkg/client"
)

func (ctx *Context) CurrentEnvironment() string {
	if ctx.Err() != nil {
		return ""
	}

	bytes := make([]byte, 0)
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

// func (ctx *Context) Fetch(environment string) {

// 	currentAccount, _ := config.GetCurrentAccount()
// 	token := config.GetAuthToken()
// 	userID := currentAccount["user_id"]
// 	ksClient := client.NewKeystoneClient(userID, token)

// 	// Get env hash from config

// 	// Request: Get env hash from remote
// 	results, err := ksClient.GetMessages(environmentID, localEnvironmentVersion)

// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	fmt.Println(results.Messages)
// 	if results.VersionID == localEnvironmentVersion {
// 		return
// 	}

// 	fmt.Println(results.VersionID)

// 	// 204: no hash set for env
// 	//    -> Set new hash for env
// 	// 200: hash and new message
// 	//    -> Set new hash for env
// }

func (ctx *Context) SetEnvironmentVersion(name string, version_id string) string {
	environments := ctx.EnvironmentsFromConfig()

	for _, e := range environments {
		if e.Name == name {
			return e.VersionID
		}
	}
	return ""
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

	environmentsfile := new(EnvironmentsFile).Load(ctx.Wd)
	return environmentsfile.Environments
}

// Push current environnement.
// Post to server []MessageToWritePayload.
func (ctx *Context) PushEnv(environments []models.Environment) error {
	// var result client.GenericResponse

	// Get public keyrs
	c, kcErr := client.NewKeystoneClient()

	if kcErr != nil {
		return kcErr
	}

	messagesToWrite := models.MessagesToWritePayload{
		Messages: make([]models.MessageToWritePayload, 0),
	}
	for _, environment := range environments {
		environmentId := environment.EnvironmentID

		userPublicKeys, err := c.Users().GetEnvironmentPublicKeys(environmentId)

		if err != nil {
			return err
		}

		PayloadContent, err := ctx.PrepareMessagePayload(environment)

		if err != nil {
			return err
		}

		var currentUser map[string]string

		if currentUser, _ = config.GetCurrentAccount(); err != nil {
			panic(err)
		}
		// Create one message per user
		for _, userPublicKey := range userPublicKeys.Keys {

			// Dont't send message to current user
			if userPublicKey.UserUID != currentUser["user_id"] {

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

	}

	result, pushErr := c.Messages().SendMessages(messagesToWrite)
	if pushErr != nil {
		return pushErr
	}

	// TODO
	// Set new version id
	fmt.Println(result.Environments)
	for _, environment := range result.Environments {
		fmt.Println(environment.VersionID)
	}

	return nil
}

func (ctx *Context) EnvironmentVersionHasChanged(name string, environmentID string) bool {
	currentVersion := ctx.EnvironmentVersionByName(name)
	return currentVersion != environmentID
}

func (ctx *Context) GetAccessibleEnvironments() []models.Environment {
	// Get public keyrs
	c, kcErr := client.NewKeystoneClient()

	if kcErr != nil {
		ctx.setError(kcErr)
		return make([]models.Environment, 0)
	}

	projectID := ctx.GetProjectID()

	accessibleEnvironments, err := c.Project(projectID).GetAccessibleEnvironments()
	if err != nil {
		ctx.setError(UnkownError(err))
	}

	return accessibleEnvironments

}
