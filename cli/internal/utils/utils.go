package utils

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

// func GetEnv(varname string, fallback string) string {
// 	if value, ok := os.LookupEnv(varname); ok {
// 		return value
// 	}

// 	return fallback
// }

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
	info, err := os.Lstat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func DirExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) || info == nil {
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
	}

	return nil
}

// Creates a directory at dirPath if does not exist
func CreateDirIfNotExist(dirPath string) error {

	if !DirExists(dirPath) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("Could not create `%s/` (%w)", dirPath, err)
		}
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

	// Remove destination, because is case of a symlink, os.Create will set empty content to the src of the symlink too!
	err = os.Remove(dst)
	if err != nil {
		return err
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func AppendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

// Check if secret name is SNAKE CAPITAL
func CheckSecretContent(name string) error {
	sampleRegex := regexp.MustCompile("^([A-Z0-9]|_)*$")
	match := sampleRegex.Match([]byte(name))

	if !match {
		return errors.New("Secret " + name + " not allowed. Secret name must be capital snakecase.")
	}
	return nil
}
