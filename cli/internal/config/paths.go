// +build !test

package config

import (
	"os"
	"path"

	"github.com/wearedevx/keystone/cli/internal/utils"
)

func ConfigDir() (dirpath string, err error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	if err = utils.CreateDirIfNotExist(userConfigDir); err != nil {
		return "", err
	}

	dirpath = path.Join(userConfigDir, "keystone")

	if err = utils.CreateDirIfNotExist(dirpath); err != nil {
		return "", err
	}

	return dirpath, nil
}

func ConfigPath() (configPath string, err error) {
	configDirPath, err := ConfigDir()
	if err != nil {
		return "", err
	}

	configPath = path.Join(configDirPath, "keystone.yaml")

	return configPath, nil
}
