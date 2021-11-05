package utils

import (
	"os"
	"path/filepath"
	"strings"
)

type FileInfo struct {
	FullPath string
	Path     string
	Info     os.FileInfo
	IsDir    bool
}

func DirWalk(source string, cb func(info FileInfo) error) error {
	info, err := os.Stat(source)
	if err != nil {
		return err
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			destinationPath := path
			if baseDir != "" {
				destinationPath = filepath.Join(
					baseDir,
					strings.TrimPrefix(path, source),
				)
			}

			return cb(FileInfo{
				FullPath: path,
				Path:     destinationPath,
				Info:     info,
				IsDir:    info.IsDir(),
			})
		})
}
