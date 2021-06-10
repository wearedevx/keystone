package core

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	. "github.com/wearedevx/keystone/cli/internal/errors"
	. "github.com/wearedevx/keystone/cli/internal/gitignorehelper"
	. "github.com/wearedevx/keystone/cli/internal/keystonefile"
	. "github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/ui"
)

type FileStrictFlag int

const (
	F_REQUIRED FileStrictFlag = iota
	F_OPTIONAL
)

func (ctx *Context) ListFiles() []FileKey {
	if ctx.Err() != nil {
		return make([]FileKey, 0)
	}

	ksfile := new(KeystoneFile).Load(ctx.Wd)

	if err := ksfile.Err(); err != nil {
		ctx.setError(FailedToReadKeystoneFile(err))
		return make([]FileKey, 0)
	}

	return ksfile.Files
}

func (ctx *Context) AddFile(file FileKey, envContentMap map[string][]byte) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	// Add file path to the keystone file
	if err := new(KeystoneFile).Load(ctx.Wd).AddFile(file).Save().Err(); err != nil {
		return ctx.setError(FailedToUpdateKeystoneFile(err))
	}

	environments := ctx.ListEnvironments()
	current := ctx.CurrentEnvironment()

	// Use current content for current environment.
	dest := path.Join(ctx.CachedEnvironmentFilesPath(current), file.Path)
	dir := filepath.Dir(dest)
	os.MkdirAll(dir, 0o755)

	if err := CopyFile(file.Path, dest); err != nil {
		return ctx.setError(CopyFailed(file.Path, dest, err))
	}

	// Set content for every other environment
	for _, environment := range environments {
		dest := path.Join(ctx.CachedEnvironmentFilesPath(environment), file.Path)
		parentDir := filepath.Dir(dest) + string(os.PathSeparator)

		if err := os.MkdirAll(parentDir, 0o755); err != nil {
			ctx.setError(CannotCreateDirectory(parentDir, err))
			return ctx
		}

		destFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0o644)
		if err == nil {
			defer destFile.Close()

			_, err = destFile.Write(envContentMap[environment])
		}

		if err != nil {
			println(fmt.Sprintf("Failed to write %s (%s)", dest, err.Error()))
			os.Exit(1)
		}
	}

	return ctx
}

// FilesUseEnvironment creates symlinks for files found in the project’s
// keystone.yml file, pointing them to the environment `envname` in cache.
func (ctx *Context) FilesUseEnvironment(envname string) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	ksfile := new(KeystoneFile).Load(ctx.Wd)
	if err := ksfile.Err(); err != nil {
		return ctx.setError(FailedToReadKeystoneFile(err))
	}

	cachePath := ctx.CachedEnvironmentFilesPath(envname)
	files := ksfile.Files

	for _, file := range files {
		cachedFilePath := path.Join(cachePath, file.Path)
		linkPath := path.Join(ctx.Wd, file.Path)

		if FileExists(linkPath) {
			os.Remove(linkPath)
		}

		if !FileExists(cachedFilePath) {
			if file.Strict {
				return ctx.setError(FileNotInEnvironment(file.Path, envname, nil))
			}
			ui.Print("File \"%s\" not in envrironment", file.Path)
		}

		parentDir := filepath.Dir(linkPath)

		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return ctx.setError(CannotLinkFile(file.Path, cachedFilePath, err))
		}

		if err := os.Symlink(cachedFilePath, linkPath); err != nil {
			return ctx.setError(CannotLinkFile(file.Path, cachedFilePath, err))
		}
	}

	return ctx
}

func (ctx *Context) RemoveFile(filePath string, force bool) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	ksfile := new(KeystoneFile).Load(ctx.Wd)

	if err := ksfile.Err(); err != nil {
		return ctx.setError(FailedToReadKeystoneFile(err))
	}

	filteredFiles := make([]FileKey, 0)
	for _, file := range ksfile.Files {
		if file.Path != filePath {
			filteredFiles = append(filteredFiles, file)
		}
	}

	ksfile.Files = filteredFiles
	ksfile.Save()

	if err := ksfile.Err(); err != nil {
		return ctx.setError(FailedToUpdateKeystoneFile(err))
	}

	environments := ctx.ListEnvironments()
	currentEnvironment := ctx.CurrentEnvironment()

	currentCached := path.Join(ctx.CachedEnvironmentFilesPath(currentEnvironment), filePath)
	dest := path.Join(ctx.Wd, filePath)

	if force {
		fmt.Println("Force remove file on filesystem.")
		os.Remove(dest)
	} else {
		fmt.Println("Keep file on filesystem.")
	}

	CopyFile(currentCached, dest)

	for _, environment := range environments {
		cachedFilePath := path.Join(ctx.CachedEnvironmentFilesPath(environment), filePath)

		if FileExists(cachedFilePath) {
			os.Remove(cachedFilePath)
		}
	}

	GitUnignore(ctx.Wd, filePath)

	return ctx
}
