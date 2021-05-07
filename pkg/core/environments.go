package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/wearedevx/keystone/internal/config"
	. "github.com/wearedevx/keystone/internal/envfile"
	. "github.com/wearedevx/keystone/internal/errors"
	. "github.com/wearedevx/keystone/internal/keystonefile"
	. "github.com/wearedevx/keystone/internal/utils"
	"github.com/wearedevx/keystone/pkg/client"
)

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
		newEnvDir := path.Join(ctx.cacheDirPath(), name)
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
		envDir := path.Join(ctx.cacheDirPath(), name)
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
		dotEnvFilePath := path.Join(ctx.cacheDirPath(), name, ".env")
		currentDotEnvFilePath := path.Join(ctx.cacheDirPath(), ".env")

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
		dotEnvPath := path.Join(ctx.cacheDirPath(), name, ".env")

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
		dotEnvPath := path.Join(ctx.cacheDirPath(), envName, ".env")

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

	return DirExists(path.Join(ctx.cacheDirPath(), name))
}

func (ctx *Context) MustHaveEnvironment(name string) {
	if !ctx.HasEnvironment(name) {
		EnvironmentDoesntExist(name, strings.Join(ctx.ListEnvironments(), ", "), nil).Print()
		os.Exit(0)
	}
}

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

func (ctx *Context) Fetch(environment string) {

	currentAccount, _ := config.GetCurrentAccount()
	token := config.GetAuthToken()
	userID := currentAccount["user_id"]
	ksClient := client.NewKeystoneClient(userID, token)

	// Get env hash from config
	localEnvironmentVersion := ctx.EnvironmentVersion()
	environmentID := ctx.EnvironmentID()

	// Request: Get env hash from remote
	results, err := ksClient.GetMessages(environmentID, localEnvironmentVersion)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(results.Messages)
	if results.VersionID == localEnvironmentVersion {
		return
	}

	fmt.Println(results.VersionID)

	// 204: no hash set for env
	//    -> Set new hash for env
	// 200: hash and new message
	//    -> Set new hash for env
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
func (ctx *Context) EnvironmentID() string {
	environments := ctx.EnvironmentsFromConfig()
	currentEnvironment := ctx.CurrentEnvironment()

	for _, e := range environments {
		if e.Name == currentEnvironment {
			return e.EnvironmentID
		}
	}
	return ""
}

func (ctx *Context) EnvironmentsFromConfig() []Env {

	ksfile := new(KeystoneFile).Load(ctx.Wd)
	return ksfile.Environments
}
