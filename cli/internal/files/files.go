package files

import (
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/eiannone/keyboard"
	"github.com/wearedevx/keystone/api/pkg/models"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui/display"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

type fileService struct {
	err error
	ctx *core.Context
}

type FileService interface {
	Err() error
	AskContent(
		filePath string,
		environments []models.Environment,
		environmentFileMap map[string][]byte,
		currentContent []byte,
		currentEnvironment string,
	) FileService
	AskToOverrideFilesInCache(fileName string, skipPrompts bool) bool
}

func NewFileService(ctx *core.Context) FileService {
	fs := new(fileService)
	fs.ctx = ctx

	return fs
}

func (fs *fileService) Err() error {
	return fs.err
}

func (fs *fileService) AskContent(
	filePath string,
	environments []models.Environment,
	environmentFileMap map[string][]byte,
	currentContent []byte,
	currentEnvironment string,
) FileService {
	if fs.err != nil {
		return fs
	}

	extension := filepath.Ext(filePath)

	for _, environment := range environments {
		if environment.Name != currentEnvironment {
			display.FileAskForFileContentForEnvironment(
				filePath,
				environment.Name,
			)

			_, _, err := keyboard.GetSingleKey()
			if err != nil {
				fs.err = fmt.Errorf("failed to read user input (%w)", err)
				return fs
			}

			content, err := utils.CaptureInputFromEditor(
				utils.GetPreferredEditorFromEnvironment,
				extension,
				string(currentContent),
			)
			if err != nil {
				fs.err = fmt.Errorf(
					"failed to get content from editor (%w)",
					err,
				)
				return fs
			}

			environmentFileMap[environment.Name] = content
		}
	}
	return fs
}

func (fs *fileService) AskToOverrideFilesInCache(fileName string, skipPrompts bool) bool {
	if fs.err != nil {
		return false
	}

	files := fs.ctx.ListFilesFromCache()
	var found keystonefile.FileKey
	for _, file := range files {
		if file.Path == fileName {
			found = file
		}
	}
	if !reflect.ValueOf(found).IsZero() {
		display.FileContentsForEnvironments(
			fileName,
			fs.ctx.AccessibleEnvironments,
			fs.ctx.GetFileContents,
		)

		override := false

		if !skipPrompts {
			override = prompts.ConfirmOverrideFileContents()
		}

		return !override
	}

	return false
}
