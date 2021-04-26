package core

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	. "github.com/wearedevx/keystone/internal/errors"
	. "github.com/wearedevx/keystone/internal/gitignorehelper"
	. "github.com/wearedevx/keystone/internal/keystonefile"
	. "github.com/wearedevx/keystone/internal/utils"
)

func (ctx *Context) ListFiles() []string {
	if ctx.Err() != nil {
		return make([]string, 0)
	}

	ksfile := new(KeystoneFile).Load(ctx.Wd)

	if err := ksfile.Err(); err != nil {
		ctx.setError(FailedToReadKeystoneFile(err))
		return make([]string, 0)
	}

	return ksfile.Files
}

func (ctx *Context) AddFile(filePath string, envContentMap map[string][]byte) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	// Add file path to the keystone file
	if err := new(KeystoneFile).Load(ctx.Wd).AddFile(filePath).Save().Err(); err != nil {
		return ctx.setError(FailedToUpdateKeystoneFile(err))
	}

	environments := ctx.ListEnvironments()
	current := ctx.CurrentEnvironment()

	// Use current content for current environment.
	dest := path.Join(ctx.cacheDirPath(), current, filePath)
	dir := filepath.Dir(dest)
	os.MkdirAll(dir, 0o755)

	if err := CopyFile(filePath, dest); err != nil {
		return ctx.setError(CopyFailed(filePath, dest, err))
	}

	// Set content for every other environment
	for _, environment := range environments {
		// current is already set
		if environment == current {
			continue
		}

		dest := path.Join(ctx.cacheDirPath(), environment, filePath)
		parentDir := filepath.Dir(dest) + "/"

		if err := os.MkdirAll(parentDir, 0o755); err != nil {
			ctx.setError(CannotCreateDirectory(parentDir, err))
			panic(err)
		}

		destFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0o644)

		if err == nil {
			defer destFile.Close()

			destFile.Write(envContentMap[environment])
		} else {
			panic(err)
		}

	}

	return ctx
}

func (ctx *Context) FilesUseEnvironment(envname string) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	ksfile := new(KeystoneFile).Load(ctx.Wd)

	if err := ksfile.Err(); err != nil {
		return ctx.setError(FailedToReadKeystoneFile(err))
	}

	files := ksfile.Files

	for _, file := range files {
		cachedFilePath := path.Join(ctx.cacheDirPath(), envname, file)
		linkPath := path.Join(ctx.Wd, file)

		if FileExists(linkPath) {
			os.Remove(linkPath)
		}

		if !FileExists(cachedFilePath) {
			return ctx.setError(FileNotInEnvironment(file, envname, nil))
		}

		parentDir := filepath.Dir(linkPath)

		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return ctx.setError(CannotLinkFile(file, cachedFilePath, err))
		}

		if err := os.Symlink(cachedFilePath, linkPath); err != nil {
			return ctx.setError(CannotLinkFile(file, cachedFilePath, err))
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

	filteredFiles := make([]string, 0)
	for _, file := range ksfile.Files {
		if file != filePath {
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

	currentCached := path.Join(ctx.cacheDirPath(), currentEnvironment, filePath)
	dest := path.Join(ctx.Wd, filePath)

	if force {
		fmt.Println("Force remove file on filesystem.")
		os.Remove(dest)
	} else {
		fmt.Println("Keep file on filesystem.")
	}

	CopyFile(currentCached, dest)

	for _, environment := range environments {
		cachedFilePath := path.Join(ctx.cacheDirPath(), environment, filePath)

		// relativePathFile := strings.Replace(cachedFilePath, ctx.Wd, "", 1)
		// fmt.Println("- File to delete: ", "."+relativePathFile)

		if FileExists(cachedFilePath) {
			os.Remove(cachedFilePath)
		}
	}

	GitUnignore(ctx.Wd, filePath)

	return ctx
}
