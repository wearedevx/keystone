package core

import (
	"github.com/wearedevx/keystone/cli/internal/environmentsfile"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/gitignorehelper"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/pkg/constants"

	"github.com/wearedevx/keystone/api/pkg/models"
)

// Initialize the projects directory structure.
//
// In the execution context's working directory,
// it creates:
// - keystone.yaml
// - .keystone/
// - .keystone/environment
// - .keystone/cache/
// - .keystone/cache/.env
// - .keystone/cache/dev/
// - .keystone/cache/ci/
// - .keystone/cache/staging/
// - .keystone/cache/prod/
//
// It adds .keystone to .gitignore, creating
// it if does not exist
func (ctx *Context) Init(project models.Project) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	var err error

	// Operations are declared in an array
	// and executed in a for loop to ease
	// erro handling
	ops := []func() error{
		func() error {
			if !keystonefile.ExistsKeystoneFile(ctx.Wd) {
				return keystonefile.NewKeystoneFile(ctx.Wd, project).
					Save().
					Err()
			}
			return nil
		},
		func() error {
			return utils.CreateDirIfNotExist(ctx.dotKeystonePath())
		},
		func() error {
			if !environmentsfile.ExistsEnvironmentsFile(ctx.dotKeystonePath()) {
				return environmentsfile.NewEnvironmentsFile(ctx.dotKeystonePath(), project.Environments).
					Save().
					Err()
			}
			return nil
		},
		func() error {
			return utils.CreateDirIfNotExist(ctx.cacheDirPath())
		},
		func() error {
			return utils.CreateFileIfNotExists(ctx.CachedDotEnvPath(), "")
		},
		func() error {
			return utils.CreateDirIfNotExist(ctx.CachedEnvironmentPath(string(constants.DEV)))
		},
		func() error {
			return utils.CreateDirIfNotExist(
				ctx.CachedEnvironmentPath(string(constants.STAGING)),
			)
		},
		func() error {
			return utils.CreateDirIfNotExist(ctx.CachedEnvironmentPath(string(constants.PROD)))
		},
		func() error {
			return utils.CreateFileIfNotExists(
				ctx.CachedEnvironmentDotEnvPath(string(constants.DEV)),
				"",
			)
		},
		func() error {
			return utils.CreateFileIfNotExists(
				ctx.CachedEnvironmentDotEnvPath(string(constants.STAGING)),
				"",
			)
		},
		func() error {
			return utils.CreateFileIfNotExists(
				ctx.CachedEnvironmentDotEnvPath(string(constants.PROD)),
				"",
			)
		},
		func() error {
			return utils.CreateDirIfNotExist(
				ctx.CachedEnvironmentFilesPath(string(constants.DEV)),
			)
		},
		func() error {
			return utils.CreateDirIfNotExist(
				ctx.CachedEnvironmentFilesPath(string(constants.STAGING)),
			)
		},
		func() error {
			return utils.CreateDirIfNotExist(
				ctx.CachedEnvironmentFilesPath(string(constants.PROD)),
			)
		},
		func() error {
			return gitignorehelper.GitIgnore(ctx.Wd, dotKeystone)
		},
	}

	for _, op := range ops {
		if err = op(); err != nil {
			return ctx.setError(kserrors.InitFailed(err))
		}
	}

	return ctx
}
