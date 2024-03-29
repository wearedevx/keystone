package utils

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
	info, err := os.Lstat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return false
	}
	return !info.IsDir()
}

// DirExists function returns true if `filename` exists and is a directory
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
		dir := filepath.Dir(filePath)

		if err := os.MkdirAll(dir, 0o700); err != nil {
			return fmt.Errorf("could not create directory `%s` (%w)", dir, err)
		}

		if err := ioutil.WriteFile(filePath, []byte(defaultContent), 0o600); err != nil {
			return fmt.Errorf("could not create `%s` (%w)", filePath, err)
		}
	}

	return nil
}

// Creates a directory at dirPath if does not exist
func CreateDirIfNotExist(dirPath string) error {
	if !DirExists(dirPath) {
		if err := os.MkdirAll(dirPath, 0o700); err != nil {
			return fmt.Errorf("could not create `%s/` (%w)", dirPath, err)
		}
	}

	return nil
}

type Closer interface {
	Close() error
}

// Close function closes anything with a `Close()` method
func Close(r Closer) {
	if err := r.Close(); err != nil {
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			panic(err)
		}
	}
}

// CopyFile function copies a file from `src` to `dst`
func CopyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	/* #nosec */
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer Close(source)

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer Close(destination)

	_, err = io.Copy(destination, source)
	return err
}

// RemoveFile function removes a file
func RemoveFile(f string) error {
	err := os.Remove(f)

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to remove file %s: %w", f, err)
	}

	return nil
}

// RemoveContents function remove all contents of a directory
func RemoveContents(dir string) error {
	/* #nosec
	 * as long as contents of dir are only removed
	 */
	d, err := os.Open(dir)
	if err != nil {
		return err
	}

	defer Close(d)

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

// AppendIfMissing function appends a string to a slice, if it is not already there
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
		return errors.New(
			"Secret " + name + " not allowed. Secret name must be capital snakecase.",
		)
	}
	return nil
}

const doubleQuoteSpecialChars = "\\\n\r\"!$`"

// DoubleQuoteEscape function
func DoubleQuoteEscape(line string) string {
	for _, c := range doubleQuoteSpecialChars {
		toReplace := "\\" + string(c)
		if c == '\n' {
			toReplace = `\n`
		}
		if c == '\r' {
			toReplace = `\r`
		}
		line = strings.ReplaceAll(line, string(c), toReplace)
	}
	return line
}
