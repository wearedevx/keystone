package utils

import (
	"fmt"
	"os"
	"path"

	uuid "github.com/satori/go.uuid"
)

func CreateTestDir() (string, error) {
	testDir := path.Join(os.TempDir(), "tests", uuid.NewV4().String())
	err := os.MkdirAll(testDir, 0o700)
	if err != nil {
		return testDir, err
	}

	return testDir, nil
}

func CleanTestDir(testDir string) error {
	err := os.RemoveAll(testDir)
	if err != nil {
		return fmt.Errorf("error cleaning test dir: `%s` (%w)", testDir, err)
	}

	return err
}
