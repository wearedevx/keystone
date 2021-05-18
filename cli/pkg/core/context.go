package core

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"

	. "github.com/wearedevx/keystone/cli/internal/errors"
	. "github.com/wearedevx/keystone/cli/internal/keystonefile"
	. "github.com/wearedevx/keystone/cli/internal/utils"
)

type Context struct {
	err       *Error
	Wd        string
	TmpDir    string
	ConfigDir string
}

const CTX_INIT = "init"
const CTX_RESOLVE = "resolve"

// Creates a new execution context
//
// When flag equals CTX_INIT, the current working directory is used
// as the project's root.
// When flag equals CTX_RESOLVE, the program tries to find the project's root
// in a parent directory. En error is returned if none is found.
//
// An error will be returned if flag is neither of those values
// TODO: This might need refactiring: too long
func New(flag string) *Context {
	var cwd string
	var err error
	context := new(Context)

	if cwd, err = os.Getwd(); err != nil {
		return context.setError(NoWorkingDirectory(err))
	}

	// Get Wd
	if flag == CTX_INIT {
		context.Wd = cwd
	} else if flag == CTX_RESOLVE {
		wd, err := resolveKeystoneRootDir(cwd)

		if err != nil {
			return context.setError(NotAKeystoneProject(cwd, err))
		} else {
			context.Wd = wd
		}
	} else {
		context.err = UnsupportedFlag(flag, nil)
	}

	// Get a temporary directory
	tmpDir := os.TempDir()
	context.TmpDir = tmpDir

	// Get global configuration path
	var currentUser *user.User
	if currentUser, err = user.Current(); err != nil {
		panic(err)
	}

	configDir := path.Join(currentUser.HomeDir, ".config", "keystone")

	if err = os.MkdirAll(configDir, 0755); err != nil {
		panic(err)
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
	return path.Join(c.dotKeystonePath(), "environment.yml")
}

func (c *Context) rolesFilePath() string {
	return path.Join(c.dotKeystonePath(), "roles.yml")
}

func (c *Context) cacheDirPath() string {
	return path.Join(c.dotKeystonePath(), "cache")
}

func (c *Context) CachedDotEnvPath() string {
	return path.Join(c.cacheDirPath(), ".env")
}

/********************/
/* Public functions */
/********************/

// Remove temporary files
func (context *Context) CleanUp() {
	os.RemoveAll(context.TmpDir)
}

// Accessor for error
func (context *Context) Err() *Error {
	return context.err
}

func (ctx *Context) SetError(err *Error) *Context {
	ctx.err = err

	return ctx
}
func (ctx *Context) setError(err *Error) *Context {
	ctx.err = err

	return ctx
}

// Determines if path matches a keystone managed project root
// path must:
// - be a directory
// - contain a keystone.yml file
func isKeystoneRootDir(path string) bool {
	if !DirExists(path) {
		return false
	}

	return ExistsKeystoneFile(path)
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
			err = fmt.Errorf("Not in a keystone managed project")
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
	ksfile := new(KeystoneFile).Load(c.Wd)
	currentEnvironment := c.CurrentEnvironment()

	for _, env := range ksfile.Environments {
		if env.Name == currentEnvironment {
			return env.EnvironmentID
		}
	}

	return ""
}
