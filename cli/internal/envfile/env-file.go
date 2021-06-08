package envfile

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/wearedevx/keystone/cli/internal/utils"
)

// An EnvFile represents a .env file data
type EnvFile struct {
	// err is set everytime an error occurs while working with the envFile
	err error
	// path to the .env file on the filesystem
	path string
	// .env variables
	data map[string]string
}

// Loads a .env file from disk
func (f *EnvFile) Load(path string) *EnvFile {
	f.path = path
	f.data = make(map[string]string)

	err := utils.CreateFileIfNotExists(path, "")
	if err != nil {
		f.SetError("Failed to create file %s, %+v", f.path, err)
		return f
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0644)

	if err != nil {
		return f.SetError("Failed to open `%s` (%w)", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if scanner.Err() == nil {
			line := scanner.Text()
			if !strings.HasPrefix(line, "#") && len(line) > 0 {
				key, value := split(line)

				f.data[key] = value
			}
		} else {
			return f.SetError("Failed to read `%s` (%w)", path, err)
		}
	}

	return f
}

// Error accessor used to check if
// an erro occurred
func (f *EnvFile) Err() error {
	return f.err
}

// Set the internal error field
func (f *EnvFile) SetError(message string, args ...interface{}) *EnvFile {
	f.err = fmt.Errorf(message, args...)
	return f
}

// Writes the .env to the disk
func (f *EnvFile) Dump() *EnvFile {
	if f.Err() != nil {
		return f
	}

	var sb strings.Builder

	for key, value := range f.data {
		sb.WriteString(key)
		sb.WriteRune('=')
		sb.WriteString(value)
		sb.WriteRune('\n')
	}

	contents := sb.String()

	if err := ioutil.WriteFile(f.path, []byte(contents), 0o644); err != nil {
		f.err = fmt.Errorf("Failed to write `%s` (%w)", f.path, err)
	}

	return f
}

// Looks for the value associated to key.
// If the .env file does not contain the key,
// the second returned value is set to false.
func (f *EnvFile) Lookup(key string) (string, bool) {
	if f.Err() == nil {
		if val, ok := f.data[key]; ok {
			return val, true
		}
	}
	return "", false
}

// Returns all the key/value pairs found in the .env file
// in a map
func (f *EnvFile) GetData() map[string]string {
	if f.Err() == nil {
		return f.data
	}

	return make(map[string]string)
}

func (f *EnvFile) SetData(data map[string]string) *EnvFile {
	if f.Err() == nil {
		f.data = data
	}

	return f
}

// Adds a key-value pair to the .env file
func (f *EnvFile) Set(key string, value string) *EnvFile {
	if f.Err() == nil {
		f.data[key] = value
	}

	return f
}

// Removes a key-value pair to the .env file
func (f *EnvFile) Unset(key string) *EnvFile {
	if f.Err() == nil {
		delete(f.data, key)
	}

	return f
}

// Utility: slits on equal
func split(s string) (string, string) {

	slice := strings.Split(s, "=")
	return slice[0], slice[1]
}
