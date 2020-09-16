package utils

import (
	"fmt"
	"os"
	"path"

	uuid "github.com/satori/go.uuid"
)

func CreateTestDir() (string, error) {
	testDir := path.Join(".", "tests", uuid.NewV4().String())
	err := os.MkdirAll(testDir, 0755)

	if err != nil {
		return testDir, err
	}

	return testDir, nil
}

func CleanTestDir(testDir string) error {
	err := os.RemoveAll(testDir)

	if err != nil {
		fmt.Printf("%+v\n", err)
		return fmt.Errorf("Error Cleaning Test Dir: `%s` (%w)", testDir, err)
	}

	return err
}
