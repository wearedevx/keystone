// +build test

package config

import (
	"os"
	"path"

	"github.com/wearedevx/keystone/cli/internal/utils"
)

// ConfigDir function returns the path to the config dir,
// and creates it if it doesn't exist
func ConfigDir() (dirpath string, err error) {
	dirpath = path.Join(os.Getenv("HOME"), ".config", "keystone")

	if err = utils.CreateDirIfNotExist(dirpath); err != nil {
		return "", err
	}

	return dirpath, nil
}

// ConfigPath function returns the path to the config file.
func ConfigPath() (configPath string, err error) {
	configDirPath, err := ConfigDir()
	if err != nil {
		return "", err
	}

	configPath = path.Join(configDirPath, "keystone.yaml")

	return configPath, nil
}
