package core

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/config"
	"github.com/wearedevx/keystone/cli/internal/environmentsfile"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	ksfile "github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/pkg/constants"
)

type Context struct {
	err                    *kserrors.Error
	Wd                     string
	TmpDir                 string
	ConfigDir              string
	AccessibleEnvironments []models.Environment
}

const (
	CTX_INIT    = "init"
	CTX_RESOLVE = "resolve"
)

const (
	CTX_OVERWRITE_LOCAL_FILES = true
	CTX_KEEP_LOCAL_FILES      = !CTX_OVERWRITE_LOCAL_FILES
)

// Creates a new execution context
//
// When flag equals CTX_INIT, the current working directory is used
// as the project's root.
// When flag equals CTX_RESOLVE, the program tries to find the project's root
// in a parent directory. En error is returned if none is found.
//
// An error will be returned if flag is neither of those values
func New(flag string) *Context {
	var cwd string
	var err error
	context := new(Context)

	if cwd, err = os.Getwd(); err != nil {
		return context.setError(kserrors.NoWorkingDirectory(err))
	}

	// Get Wd
	switch flag {
	case CTX_INIT:
		context.Wd = cwd
	case CTX_RESOLVE:
		wd, err := resolveKeystoneRootDir(cwd)

		if err != nil {
			return context.setError(kserrors.NotAKeystoneProject(cwd, err))
		} else {
			context.Wd = wd
		}
	default:
		context.err = kserrors.UnsupportedFlag(flag, nil)
	}

	// Get a temporary directory
	tmpDir := os.TempDir()
	context.TmpDir = tmpDir

	// Get global configuration path
	configDir, err := config.ConfigDir()
	if err != nil {
		println("failed get the config directory")
		os.Exit(1)
	}

	context.ConfigDir = configDir

	return context
}

/**************************/
/* Private path utilities */
/**************************/

const dotKeystone string = ".keystone"

func (c *Context) dotKeystonePath() string {
	return path.Join(c.Wd, ".keystone")
}

func (c *Context) environmentFilePath() string {
	return path.Join(c.dotKeystonePath(), "environments.yaml")
}

func (c *Context) rolesFilePath() string {
	return path.Join(c.dotKeystonePath(), "roles.yaml")
}

func (c *Context) cacheDirPath() string {
	return path.Join(c.dotKeystonePath(), "cache")
}

func (c *Context) CachedDotEnvPath() string {
	return path.Join(c.cacheDirPath(), ".env")
}

func (c *Context) CachedEnvironmentPath(environmentName string) string {
	if environmentName == "" {
		environmentName = string(constants.DEV)
	}

	p := path.Join(c.cacheDirPath(), environmentName)
	c.mustEnvironmentNameBeValid(environmentName)

	if !c.fileBelongsToContext(p) {
		kserrors.
			EnvironmentDoesntExist(environmentName, "dev, staging, prod", nil).
			Print()
		os.Exit(1)
	}

	return path.Join(c.cacheDirPath(), environmentName)
}

func (c *Context) CachedEnvironmentDotEnvPath(environmentName string) string {
	return path.Join(c.CachedEnvironmentPath(environmentName), ".env")
}

func (c *Context) CachedEnvironmentFilesPath(environmentName string) string {
	return path.Join(c.CachedEnvironmentPath(environmentName), "files")
}

/********************/
/* Public functions */
/********************/

// Remove temporary files
func (context *Context) CleanUp() {
	err := os.RemoveAll(context.TmpDir)
	if err != nil {
		panic(err)
	}
}

// Accessor for error
func (context *Context) Err() *kserrors.Error {
	return context.err
}

func (ctx *Context) SetError(err *kserrors.Error) *Context {
	ctx.err = err

	return ctx
}

func (ctx *Context) setError(err *kserrors.Error) *Context {
	ctx.err = err

	return ctx
}

// Determines if path matches a keystone managed project root
// path must:
// - be a directory
// - contain a keystone.yaml file
func isKeystoneRootDir(path string) bool {
	if !utils.DirExists(path) {
		return false
	}

	return ksfile.ExistsKeystoneFile(path)
}

// Looks for a keystone managed project root in a parent directory
// [cwd] is the directory to start the search
// Returns the path of the first project root it finds
// Returns an error if project root could be found
func resolveKeystoneRootDir(cwd string) (string, error) {
	candidate := cwd
	found := false
	var err error

	for !found && err == nil {
		if isKeystoneRootDir(candidate) {
			found = true
		} else if candidate == "" || candidate == "/" {
			err = fmt.Errorf("not in a keystone managed project")
			break
		} else {
			candidate = filepath.Dir(candidate)
		}
	}

	return candidate, err
}

func (c *Context) currentEnvironmentCachePath() string {
	envCachePath := c.cacheDirPath()
	currentEnvironment := c.CurrentEnvironment()
	return path.Join(envCachePath, currentEnvironment)
}

func (c *Context) getCurrentEnvironmentId() string {
	environmentsfile := new(environmentsfile.EnvironmentsFile).Load(c.Wd)
	currentEnvironment := c.CurrentEnvironment()

	for _, env := range environmentsfile.Environments {
		if env.Name == currentEnvironment {
			return env.EnvironmentID
		}
	}

	return ""
}

func (c *Context) DotKeystonePath() string {
	return path.Join(c.Wd, ".keystone")
}
