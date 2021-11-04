package core

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/wearedevx/keystone/api/pkg/models"

	"github.com/wearedevx/keystone/cli/internal/envfile"
	"github.com/wearedevx/keystone/cli/internal/environmentsfile"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/pkg/constants"
)

// CurrentEnvironment method returns the current environment name
func (ctx *Context) CurrentEnvironment() string {
	if ctx.Err() != nil {
		return ""
	}

	environmentsfile := &environmentsfile.EnvironmentsFile{
		Current: string(constants.DEV),
	}
	environmentsfile.Load(ctx.dotKeystonePath())

	if err := environmentsfile.Err(); err != nil {
		ctx.setError(
			kserrors.CannotReadEnvironment(ctx.environmentFilePath(), err),
		)
	}

	return environmentsfile.Current
}

func (ctx *Context) mustEnvironmentNameBeValid(name string) {
	valid := false

	switch constants.EnvName(name) {
	case constants.DEV:
		valid = true
	case constants.STAGING:
		valid = true
	case constants.PROD:
		valid = true
	}

	if !valid {
		kserrors.EnvironmentDoesntExist(
			name,
			constants.EnvList.String(),
			nil,
		).Print()

		os.Exit(1)
	}
}

// ListEnvironments method returns a list of environments that can be found
// in the cache
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

		if !contained && envname != "" {
			ctx.mustEnvironmentNameBeValid(envname)
			envs = append(envs, envname)
		}
	}

	return envs
}

// RemoveEnvironment method
// Deprecated
func (ctx *Context) RemoveEnvironment(name string) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	if current := ctx.CurrentEnvironment(); current == name {
		panic("cannot remove current envrionment")
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

// SetCurrent method changes the current environment
func (ctx *Context) SetCurrent(name string) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	if ctx.HasEnvironment(name) {
		dotEnvFilePath := ctx.CachedEnvironmentDotEnvPath(name)
		currentDotEnvFilePath := ctx.CachedDotEnvPath()

		err := utils.CopyFile(dotEnvFilePath, currentDotEnvFilePath)
		if err != nil {
			return ctx.setError(
				kserrors.CopyFailed(dotEnvFilePath, currentDotEnvFilePath, err),
			)
		}

		environmentsfile := &environmentsfile.EnvironmentsFile{}
		if err := environmentsfile.Load(ctx.dotKeystonePath()).SetCurrent(name).Save().Err(); err != nil {
			ctx.setError(kserrors.FailedToUpdateKeystoneFile(err))
		}

		if err != nil {
			return ctx.setError(
				kserrors.FailedToSetCurrentEnvironment(
					name,
					ctx.environmentFilePath(),
					err,
				),
			)
		}

		ctx.FilesUseEnvironment(
			ctx.CurrentEnvironment(),
			name,
			CTX_OVERWRITE_LOCAL_FILES,
		)

	} else {
		return ctx.setError(kserrors.EnvironmentDoesntExist(name, strings.Join(ctx.ListEnvironments(), ", "), nil))
	}

	return ctx
}

// SetAllSecrets method sets all the secrets at once in the cache
func (ctx *Context) SetAllSecrets(
	name string,
	secrets map[string]string,
) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	if ctx.HasEnvironment(name) {
		dotEnvPath := ctx.CachedEnvironmentDotEnvPath(name)

		if err := new(envfile.EnvFile).Load(dotEnvPath, nil).SetData(secrets).Dump().Err(); err != nil {
			return ctx.setError(kserrors.FailedToUpdateDotEnv(dotEnvPath, err))
		}

	} else {
		return ctx.setError(kserrors.EnvironmentDoesntExist(name, strings.Join(ctx.ListEnvironments(), ", "), nil))
	}

	return ctx
}

// GetAllSecrets method returns all the secrets and thei value for the given environment
func (ctx *Context) GetAllSecrets(envName string) map[string]string {
	emptyMap := map[string]string{}

	if ctx.Err() != nil {
		return emptyMap
	}

	if ctx.HasEnvironment(envName) {
		dotEnvPath := ctx.CachedEnvironmentDotEnvPath(envName)

		envFile := new(envfile.EnvFile).Load(dotEnvPath, nil)

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

// HasEnvironment method returns true environment exists
func (ctx *Context) HasEnvironment(name string) bool {
	if ctx.Err() != nil {
		return false
	}

	return utils.DirExists(ctx.CachedEnvironmentPath(name))
}

// MustHaveEnvironment method exits with an error if the environment
// does not exist
func (ctx *Context) MustHaveEnvironment(name string) {
	if !ctx.HasEnvironment(name) {
		kserrors.EnvironmentDoesntExist(name, strings.Join(ctx.ListEnvironments(), ", "), nil).
			Print()
		os.Exit(1)
	}
}

// MustHaveAccessToEnvironment method exits with an error if the user
// does not have read access to the environment
func (ctx *Context) MustHaveAccessToEnvironment(
	environmentName string,
) *Context {
	for _, accessible := range ctx.AccessibleEnvironments {
		if accessible.Name == environmentName {
			return ctx
		}
	}

	kserrors.PermissionDenied(environmentName, nil).Print()
	os.Exit(1)

	return ctx
}

// UpdateEnvironment method updates environment info (e.g. versionID)
func (ctx *Context) UpdateEnvironment(environment models.Environment) *Context {
	environmentFile := new(environmentsfile.EnvironmentsFile)

	if err := environmentFile.
		Load(ctx.dotKeystonePath()).
		Replace(environment).
		Save().
		Err(); err != nil {
		ctx.setError(kserrors.FailedToUpdateDotEnv(environmentFile.Path(), err))
	}

	return ctx
}

// EnvironmentVersion method returns the local version of the current
// environment
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

// EnvironmentVersionByName method returns the local version of the
// environment named `name`
func (ctx *Context) EnvironmentVersionByName(name string) string {
	environments := ctx.EnvironmentsFromConfig()

	for _, e := range environments {
		if e.Name == name {
			return e.VersionID
		}
	}
	return ""
}

// EnvironmentID method returns the environmentID of the current environment
func (ctx *Context) EnvironmentID() string {
	return ctx.getCurrentEnvironmentId()
}

// EnvironmentsFromConfig method reads the environmentfile
func (ctx *Context) EnvironmentsFromConfig() []environmentsfile.Env {
	environmentsfile := new(
		environmentsfile.EnvironmentsFile,
	).Load(ctx.dotKeystonePath())
	return environmentsfile.Environments
}

// EnvironmentVersionHasChanged method indiciates whether the versions differ
func (ctx *Context) EnvironmentVersionHasChanged(
	name string,
	environmentVersion string,
) bool {
	currentVersion := ctx.EnvironmentVersionByName(name)
	return currentVersion != environmentVersion
}

// LoadEnvironmentsFile method
func (ctx *Context) LoadEnvironmentsFile() *environmentsfile.EnvironmentsFile {
	return new(environmentsfile.EnvironmentsFile).Load(ctx.dotKeystonePath())
}

// RemoveForbiddenEnvironments method removes the environments information
// the current user may not access.
func (ctx *Context) RemoveForbiddenEnvironments(
	accessibleEnvironments []models.Environment,
) {
	accessibleEnvironmentsNames := make([]string, 0)

	for _, accessibleEnvironment := range accessibleEnvironments {
		accessibleEnvironmentsNames = append(
			accessibleEnvironmentsNames,
			accessibleEnvironment.Name,
		)
	}

	for _, localEnvironment := range ctx.ListEnvironments() {
		// If environment is not accessible, remove directory in cache
		if !Contains(accessibleEnvironmentsNames, localEnvironment) {
			ctx.RemoveEnvironment(localEnvironment)
		}
	}
}
