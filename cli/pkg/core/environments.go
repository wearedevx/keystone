package core

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/wearedevx/keystone/api/pkg/models"

	. "github.com/wearedevx/keystone/cli/internal/envfile"
	. "github.com/wearedevx/keystone/cli/internal/environmentsfile"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	. "github.com/wearedevx/keystone/cli/internal/utils"
)

func (ctx *Context) CurrentEnvironment() string {
	if ctx.Err() != nil {
		return ""
	}

	environmentsfile := &EnvironmentsFile{}
	environmentsfile.Load(ctx.dotKeystonePath())

	if err := environmentsfile.Err(); err != nil {
		ctx.setError(kserrors.CannotReadEnvironment(ctx.environmentFilePath(), err))
	}

	if environmentsfile.Current != "" {
		ctx.mustEnvironmentNameBeValid(environmentsfile.Current)
	}

	return environmentsfile.Current
}

func (ctx *Context) mustEnvironmentNameBeValid(name string) {
	valid := false

	switch name {
	case "dev":
		valid = true
	case "staging":
		valid = true
	case "prod":
		valid = true
	}

	if !valid {
		kserrors.EnvironmentDoesntExist(
			name,
			"dev, staging, prod",
			nil,
		).Print()

		os.Exit(1)
	}
}

// ListEnvironments lists all environmnts
// present on disk
func (ctx *Context) ListEnvironments() []string {
	if ctx.Err() != nil {
		return []string{}
	}
	envs := make([]string, 0)

	cacheDir := ctx.cacheDirPath()
	contents, err := ioutil.ReadDir(cacheDir)

	if err != nil {
		ctx.setError(kserrors.UnkownError(err))
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
			ctx.mustEnvironmentNameBeValid(envname)
			envs = append(envs, envname)
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
		err := os.MkdirAll(newEnvDir, 0o700)

		if err != nil {
			ctx.setError(kserrors.CannotCreateDirectory(newEnvDir, err))
		}
	} else {
		ctx.setError(kserrors.EnvironmentAlreadyExists(name, nil))
	}

	return ctx
}

func (ctx *Context) RemoveEnvironment(name string) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	if current := ctx.CurrentEnvironment(); current == name {
		return ctx.setError(kserrors.CannotRemoveCurrentEnvironment(name, nil))
	}

	if ctx.HasEnvironment(name) {
		envDir := ctx.CachedEnvironmentPath(name)
		err := os.RemoveAll(envDir)

		if err != nil {
			return ctx.setError(kserrors.CannotRemoveDirectory(envDir, err))
		}
	} else {
		return ctx.setError(kserrors.EnvironmentDoesntExist(name, strings.Join(ctx.ListEnvironments(), ", "), nil))
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
			return ctx.setError(kserrors.CopyFailed(dotEnvFilePath, currentDotEnvFilePath, err))
		}

		environmentsfile := &EnvironmentsFile{}
		if err := environmentsfile.Load(ctx.dotKeystonePath()).SetCurrent(name).Save().Err(); err != nil {
			ctx.setError(kserrors.FailedToUpdateKeystoneFile(err))
		}

		if err != nil {
			return ctx.setError(kserrors.FailedToSetCurrentEnvironment(name, ctx.environmentFilePath(), err))
		}

		ctx.FilesUseEnvironment(ctx.CurrentEnvironment(), name, CTX_KEEP_LOCAL_FILES)

	} else {
		return ctx.setError(kserrors.EnvironmentDoesntExist(name, strings.Join(ctx.ListEnvironments(), ", "), nil))
	}

	return ctx
}

func (ctx *Context) SetAllSecrets(name string, secrets map[string]string) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	if ctx.HasEnvironment(name) {
		dotEnvPath := ctx.CachedEnvironmentDotEnvPath(name)

		if err := new(EnvFile).Load(dotEnvPath, nil).SetData(secrets).Dump().Err(); err != nil {
			return ctx.setError(kserrors.FailedToUpdateDotEnv(dotEnvPath, err))
		}

	} else {
		return ctx.setError(kserrors.EnvironmentDoesntExist(name, strings.Join(ctx.ListEnvironments(), ", "), nil))
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

		envFile := new(EnvFile).Load(dotEnvPath, nil)

		if err := envFile.Err(); err != nil {
			ctx.setError(kserrors.FailedToReadDotEnv(dotEnvPath, err))
			return emptyMap
		}

		return envFile.GetData()
	} else {
		ctx.setError(kserrors.EnvironmentDoesntExist(envName, strings.Join(ctx.ListEnvironments(), ", "), nil))
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
		kserrors.EnvironmentDoesntExist(name, strings.Join(ctx.ListEnvironments(), ", "), nil).Print()
		os.Exit(1)
	}
}

func (ctx *Context) MustHaveAccessToEnvironment(environmentName string) *Context {
	for _, accessible := range ctx.AccessibleEnvironments {
		if accessible.Name == environmentName {
			return ctx
		}
	}

	kserrors.PermissionDenied(environmentName, nil).Print()
	os.Exit(1)

	return ctx
}

func (ctx *Context) UpdateEnvironment(environment models.Environment) *Context {
	environmentFile := new(EnvironmentsFile)

	if err := environmentFile.
		Load(ctx.dotKeystonePath()).
		Replace(environment).
		Save().
		Err(); err != nil {
		ctx.setError(kserrors.FailedToUpdateDotEnv(environmentFile.Path(), err))
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

func (ctx *Context) RemoveForbiddenEnvironments(accessibleEnvironments []models.Environment) {
	accessibleEnvironmentsNames := make([]string, 0)

	for _, accessibleEnvironment := range accessibleEnvironments {
		accessibleEnvironmentsNames = append(accessibleEnvironmentsNames, accessibleEnvironment.Name)
	}

	for _, localEnvironment := range ctx.ListEnvironments() {
		// If environment is not accessible, remove directory in cache
		if !Contains(accessibleEnvironmentsNames, localEnvironment) {
			ctx.RemoveEnvironment(localEnvironment)
		}

	}

}
