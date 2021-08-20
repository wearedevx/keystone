package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/udhos/equalfile"
	"github.com/wearedevx/keystone/api/pkg/models"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	. "github.com/wearedevx/keystone/cli/internal/gitignorehelper"
	. "github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/internal/utils"
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
		ctx.setError(kserrors.FailedToReadKeystoneFile(err))
		return make([]FileKey, 0)
	}

	return ksfile.Files
}

func (ctx *Context) ListFilesFromCache() []FileKey {
	if ctx.Err() != nil {
		return make([]FileKey, 0)
	}

	filesFromCache := make([]string, 0)

	for _, envname := range ctx.ListEnvironments() {
		cachePath := ctx.CachedEnvironmentFilesPath(envname)

		err := filepath.Walk(cachePath,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				// skip directories
				if info.IsDir() {
					return nil
				}

				fileRelativePath := strings.ReplaceAll(path, cachePath, "")
				regexp, err := regexp.Compile(`^\/`)

				fileRelativePath = regexp.ReplaceAllString(fileRelativePath, "")

				if len(fileRelativePath) > 0 {
					filesFromCache = append(filesFromCache, fileRelativePath)
				}
				return nil
			})

		if err != nil {
			ctx.setError(kserrors.FailedToReadKeystoneFile(err))
			return make([]FileKey, 0)
		}
	}

	filesFromCache = Uniq(filesFromCache)

	fileKey := make([]FileKey, 0)
	for _, f := range filesFromCache {
		newFileKey := FileKey{
			Path:      f,
			Strict:    false,
			FromCache: true,
		}
		fileKey = append(fileKey, newFileKey)
	}

	return fileKey
}

func (ctx *Context) AddFile(file FileKey, envContentMap map[string][]byte) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	// Add file path to the keystone file
	if err := new(KeystoneFile).Load(ctx.Wd).AddFile(file).Save().Err(); err != nil {
		return ctx.setError(kserrors.FailedToUpdateKeystoneFile(err))
	}

	environments := ctx.ListEnvironments()
	current := ctx.CurrentEnvironment()

	// Use current content for current environment.
	dest := path.Join(ctx.CachedEnvironmentFilesPath(current), file.Path)
	dir := filepath.Dir(dest)
	os.MkdirAll(dir, 0o755)

	if err := CopyFile(file.Path, dest); err != nil {
		return ctx.setError(kserrors.CopyFailed(file.Path, dest, err))
	}

	// Set content for every other environment
	for _, environment := range environments {
		dest := path.Join(ctx.CachedEnvironmentFilesPath(environment), file.Path)
		parentDir := filepath.Dir(dest) + string(os.PathSeparator)

		if err := os.MkdirAll(parentDir, 0o755); err != nil {
			ctx.setError(kserrors.CannotCreateDirectory(parentDir, err))
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

func (ctx *Context) SetFile(filePath string, content []byte) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	// Add file path to the keystone file
	ksfile := new(KeystoneFile).Load(ctx.Wd)
	if err := ksfile.Err(); err != nil {
		return ctx.setError(kserrors.FailedToReadKeystoneFile(err))
	}

	currentEnvironment := ctx.CurrentEnvironment()
	dest := path.Join(ctx.CachedEnvironmentFilesPath(currentEnvironment), filePath)

	destFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_TRUNC, 0o644)
	if err == nil {
		defer destFile.Close()

		_, err = destFile.Write(content)
	}

	if err != nil {
		println(fmt.Sprintf("Failed to write %s (%s)", dest, err.Error()))
		os.Exit(1)
	}

	return ctx
}

// LocallyModifiedFiles returns the list of file whose local content are
// different than the version in cache for the given environment
// (e.g. modified by the user)
func (ctx *Context) LocallyModifiedFiles(envname string) []FileKey {
	if ctx.Err() != nil {
		return []FileKey{}
	}

	files := ctx.ListFiles()
	modified := make([]FileKey, len(files))

	for index, fileKey := range files {
		if ctx.IsFileModified(fileKey.Path, envname) {
			modified[index] = fileKey
		}
		if ctx.Err() != nil {
			return []FileKey{}
		}
	}

	return modified
}

func (ctx *Context) IsFileModified(filePath, environment string) (isModified bool) {
	if ctx.Err() != nil {
		return false
	}
	var localPath, cachedPath string
	var localReader, cachedReader *os.File
	var err error

	localPath = path.Join(ctx.Wd, filePath)
	cachedPath = path.Join(ctx.CachedEnvironmentFilesPath(environment), filePath)

	localReader, err = os.Open(localPath)
	if err != nil {
		// ctx.setError(kserrors.CannotCopyFile(localPath, cachedPath, err))
		return false

	}
	cachedReader, err = os.Open(cachedPath)
	if err != nil {
		ctx.setError(kserrors.CannotCopyFile(localPath, cachedPath, err))
		return false
	}

	comparator := equalfile.New(nil, equalfile.Options{})
	sameContent, err := comparator.CompareReader(localReader, cachedReader)

	if err != nil {
		ctx.setError(kserrors.CannotCopyFile(localPath, cachedPath, err))
		return false
	}

	localReader.Close()
	cachedReader.Close()

	return !sameContent
}

// FilesUseEnvironment creates copies of files found in the project’s
// keystone.yaml file, from the environment `targeEnvironment` in cache.
func (ctx *Context) FilesUseEnvironment(
	envname string,
	targetEnvironment string,
	forceCopy bool,
) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	ksfile := new(KeystoneFile).Load(ctx.Wd)
	if err := ksfile.Err(); err != nil {
		return ctx.setError(kserrors.FailedToReadKeystoneFile(err))
	}

	cachePath := ctx.CachedEnvironmentFilesPath(envname)
	files := ksfile.Files

	for _, file := range files {
		cachedFilePath := path.Join(cachePath, file.Path)
		localPath := path.Join(ctx.Wd, file.Path)

		if !FileExists(cachedFilePath) {
			if file.Strict {
				return ctx.setError(kserrors.FileNotInEnvironment(file.Path, envname, nil))
			}
			fmt.Fprintf(os.Stderr, "File \"%s\" not in environment\n", file.Path)
		}

		if ctx.IsFileModified(file.Path, envname) &&
			forceCopy == CTX_KEEP_LOCAL_FILES {
			fmt.Fprintln(os.Stderr, ui.RenderTemplate("modified file",
				`{{ "Warning!" | yellow }} File '{{ .Path }}' has been locally modified. To discard local changes, run 'ks file reset {{ .Path }}'`,
				file,
			))
		} else {
			if FileExists(localPath) {
				os.Remove(localPath)
			}

			parentDir := filepath.Dir(localPath)

			if err := os.MkdirAll(parentDir, 0700); err != nil {
				return ctx.setError(kserrors.CannotCopyFile(file.Path, cachedFilePath, err))
			}

			if err := utils.CopyFile(cachedFilePath, localPath); err != nil {
				return ctx.setError(kserrors.CannotCopyFile(file.Path, cachedFilePath, err))
			}
		}
	}

	return ctx
}

func (ctx *Context) RemoveFile(filePath string, force bool, purge bool, accessibleEnvironments []models.Environment) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	ksfile := new(KeystoneFile).Load(ctx.Wd)

	if err := ksfile.Err(); err != nil {
		return ctx.setError(kserrors.FailedToReadKeystoneFile(err))
	}

	filteredFiles := make([]FileKey, 0)
	found := false
	for _, file := range ksfile.Files {
		if file.Path != filePath {
			filteredFiles = append(filteredFiles, file)
		} else {
			found = true
		}
	}
	if !found {
		err := errors.New("The file is not added to keystone.")
		return ctx.setError(kserrors.CannotRemoveFile(filePath, err))

	}

	ksfile.Files = filteredFiles
	ksfile.Save()

	if err := ksfile.Err(); err != nil {
		return ctx.setError(kserrors.FailedToUpdateKeystoneFile(err))
	}

	dest := path.Join(ctx.Wd, filePath)

	if force {
		fmt.Println("Force remove file on filesystem.")
		os.Remove(dest)
	} else {

		currentEnvironment := ctx.CurrentEnvironment()
		currentCached := path.Join(ctx.CachedEnvironmentFilesPath(currentEnvironment), filePath)

		// Remove destination, because is case of a symlink, os.Create will set empty content to the src of the symlink too!
		os.Remove(dest)

		if err := CopyFile(currentCached, dest); err != nil {
			return ctx.setError(kserrors.CopyFailed(currentCached, dest, err))
		}
	}

	if purge {
		for _, environment := range accessibleEnvironments {
			cachedFilePath := path.Join(ctx.CachedEnvironmentFilesPath(environment.Name), filePath)

			if FileExists(cachedFilePath) {
				os.Remove(cachedFilePath)
			}
		}
	}

	GitUnignore(ctx.Wd, filePath)

	return ctx
}

// Returns a boolean indicating wether the file `fileName`
// exists in the local files
func (ctx *Context) HasFile(fileName string) bool {
	haveIt := false

	if ctx.Err() != nil {
		return haveIt
	}

	ksfile := new(KeystoneFile).Load(ctx.Wd)
	if err := ksfile.Err(); err != nil {
		ctx.setError(kserrors.FailedToReadKeystoneFile(err))
		return haveIt
	}

	for _, fileKey := range ksfile.Files {
		if fileKey.Path == fileName {
			haveIt = true
			break
		}
	}

	return haveIt
}

func (ctx *Context) MarkFileRequired(
	filePath string,
	required bool,
) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	if err := new(KeystoneFile).
		Load(ctx.Wd).
		SetFileRequired(filePath, required).
		Save().
		Err(); err != nil {
		return ctx.setError(kserrors.FailedToUpdateKeystoneFile(err))
	}

	return ctx
}

// GetFileContents returns the file contents for the given envsrionment
// as a slice of bytes.
// It returns an error if reading the file fails (Pemission denied, no exists…)
// or if the file is empty (content length equals 0)
func (ctx *Context) GetFileContents(
	fileName string, environmentName string,
) (contents []byte, err error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	cachePath := ctx.CachedEnvironmentFilesPath(environmentName)
	filePath := path.Join(cachePath, fileName)

	contents, err = ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	if len(contents) == 0 {
		return nil, fmt.Errorf("No contents")
	}

	return contents, err
}

// MissingFilesForEnvironment return a list of missing files for the given
// environment.
// The second returned value indicates whether the list contains something
// or not.
func (ctx *Context) MissingFilesForEnvironment(
	environmentName string,
) ([]string, bool) {
	missing := []string{}
	hasMissing := false

	files := ctx.ListFiles()

	for _, file := range files {
		if file.Strict {
			if _, err := ctx.GetFileContents(file.Path, environmentName); err != nil {
				hasMissing = true
				missing = append(missing, file.Path)
			}
		}
	}

	return missing, hasMissing
}
