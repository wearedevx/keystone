package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func DirExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// Creates a file with defaultContent at filePath if it doesn't exist
func CreateFileIfNotExists(filePath string, defaultContent string) error {
	if !FileExists(filePath) {
		if err := ioutil.WriteFile(filePath, []byte(defaultContent), 0644); err != nil {
			return fmt.Errorf("Could not create `%s` (%w)", filePath, err)
		}
		fmt.Printf("Creating File: %s\n", filePath)
	}

	return nil
}

// Creates a directory at dirPath if does not exist
func CreateDirIfNotExist(dirPath string) error {

	if !DirExists(dirPath) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("Could not create `%s/` (%w)", dirPath, err)
		}
		fmt.Printf("Creating Directory: %s\n", dirPath)
	}

	return nil
}

func CopyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}
