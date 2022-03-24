package backup

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/wearedevx/keystone/cli/internal/archive"
	"github.com/wearedevx/keystone/cli/internal/config"
	"github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/pkg/core"
	"github.com/wearedevx/keystone/cli/ui/prompts"
)

var (
	ErrorBackupNotFound = errors.New("backup not found")
	ErrorBackupFailed   = errors.New("backup failed")
	ErrorRestoreFailed  = errors.New("restore failed")
)

type backupService struct {
	log *log.Logger
	ctx *core.Context
}

type BackupService interface {
	Setup() error
	IsSetup() bool
	Backup() (string, string, error)
	Restore() error
}

func NewBackupService(ctx *core.Context) BackupService {
	if nobackp := os.Getenv("NOBACKUP"); nobackp == "true" {
		return &stubBackupService{}
	}

	return &backupService{
		log: log.New(log.Writer(), "[Backup] ", 0),
		ctx: ctx,
	}
}

func (b *backupService) Setup() error {
	if prompts.ConfirmCreateBackupStrategy(b.IsSetup()) {
		backupPath := prompts.BackupPath()
		config.SetBackupStrategy(true, backupPath)

		b.log.Printf("Backups will be written to `%s`\n", backupPath)
	} else {
		config.SetBackupStrategy(false, "")

		b.log.Println("There will be no backups")
	}

	config.Write()

	return nil
}

func (b *backupService) IsSetup() bool {
	return config.IsBackupStrategySetup()
}

func (b *backupService) Backup() (string, string, error) {
	doBackup, backupPath := config.GetBackupStrategy()
	projectName := b.ctx.GetProjectName()
	backupFile := fmt.Sprintf("%s.tar.gz", projectName)
	destination := filepath.Join(backupPath, backupFile)

	if doBackup {
		if !utils.DirExists(backupPath) {
			err := os.MkdirAll(backupPath, 0o755)
			if err != nil {
				return "", "", fmt.Errorf("Failed to create backup directory %w", err)
			}
		}

		if utils.DirExists(backupPath) {

			err := archive.Archive(b.ctx.DotKeystonePath(), destination)
			if err != nil {
				return "", "", fmt.Errorf("Failed to create backup archive %w", err)
			}
		}
	}

	return destination, backupFile, nil
}

func (b *backupService) Restore() error {
	doBackup, backupDirPath := config.GetBackupStrategy()

	if doBackup {
		projectName := b.ctx.GetProjectName()
		backupPath := path.Join(backupDirPath, fmt.Sprintf("%s.tar.gz", projectName))

		if utils.FileExists(backupPath) {
			file, err := os.OpenFile(backupPath, os.O_RDONLY, 0o644)
			defer file.Close()
			if err != nil {
				b.log.Printf("Error: %s\n", err.Error())
				return err
			}

			err = archive.Extract(file, b.ctx.DotKeystonePath())
			if err != nil {
				b.log.Printf("Error: %s\n", err.Error())
				return err
			}
		} else {
			return ErrorBackupNotFound
		}
	}

	return nil
}
