package core

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/wearedevx/keystone/api/pkg/models"

	. "github.com/wearedevx/keystone/cli/internal/envfile"
	. "github.com/wearedevx/keystone/cli/internal/environmentsfile"
	. "github.com/wearedevx/keystone/cli/internal/errors"
	. "github.com/wearedevx/keystone/cli/internal/utils"
)

func (ctx *Context) CurrentEnvironment() string {
	if ctx.Err() != nil {
		return ""
	}

	environmentsfile := &EnvironmentsFile{}
	environmentsfile.Load(ctx.dotKeystonePath())

	if err := environmentsfile.Err(); err != nil {
		ctx.setError(CannotReadEnvironment(ctx.environmentFilePath(), err))
	}

	return environmentsfile.Current
}

func (ctx *Context) ListEnvironments() []string {
	if ctx.Err() != nil {
		return []string{}
	}
	envs := make([]string, 0)

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

		environmentsfile := &EnvironmentsFile{}
		if err := environmentsfile.Load(ctx.dotKeystonePath()).SetCurrent(name).Save().Err(); err != nil {
			ctx.setError(FailedToUpdateKeystoneFile(err))
		}

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

func (ctx *Context) UpdateEnvironment(environment models.Environment) *Context {
	environmentFile := new(EnvironmentsFile)

	if err := environmentFile.
		Load(ctx.dotKeystonePath()).
		Replace(environment).
		Save().
		Err(); err != nil {
		ctx.setError(FailedToUpdateDotEnv(environmentFile.Path(), err))
	}

	return ctx

}

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
}

func (ctx *Context) EnvironmentsFromConfig() []Env {
	environmentsfile := new(EnvironmentsFile).Load(ctx.dotKeystonePath())
	return environmentsfile.Environments
}

func (ctx *Context) EnvironmentVersionHasChanged(name string, environmentVersion string) bool {
	currentVersion := ctx.EnvironmentVersionByName(name)
	return currentVersion != environmentVersion
}

func (ctx *Context) LoadEnvironmentsFile() *EnvironmentsFile {
	return new(EnvironmentsFile).Load(ctx.dotKeystonePath())
}
