package core

import (
	"path"

	. "github.com/wearedevx/keystone/internal/errors"
	. "github.com/wearedevx/keystone/internal/gitignorehelper"
	. "github.com/wearedevx/keystone/internal/keystonefile"
	. "github.com/wearedevx/keystone/internal/utils"

	"github.com/wearedevx/keystone/internal/models"
)

// Initialize the projects directory structure.
//
// In the execution context's working directory,
// it creates:
// - keystone.yml
// - .keystone/
// - .keystone/environment
// - .keystone/cache/
// - .keystone/cache/.env
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
			if !ExistsKeystoneFile(ctx.Wd) {
				return NewKeystoneFile(ctx.Wd, project).Save().Err()
			}
			return nil
		},
		func() error {
			return CreateDirIfNotExist(ctx.dotKeystonePath())
		},
		func() error {
			return CreateFileIfNotExists(ctx.environmentFilePath(), "default")
		},
		func() error {
			return CreateDirIfNotExist(ctx.cacheDirPath())
		},
		func() error {
			return CreateFileIfNotExists(ctx.CachedDotEnvPath(), "")
		},
		func() error {
			return CreateDirIfNotExist(path.Join(ctx.cacheDirPath(), "default"))
		},
		func() error {
			return GitIgnore(ctx.Wd, dotKeystone)
		},
	}

	for _, op := range ops {
		if err = op(); err != nil {
			return ctx.setError(InitFailed(err))
		}
	}

	return ctx
}
