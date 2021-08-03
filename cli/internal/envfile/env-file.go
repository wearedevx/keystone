package envfile

import (
	"fmt"
	"io"
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
	opts LoadOptions
}

// Options for loading the .env file
type LoadOptions struct {
	// Don’t unescape double quotes and other special characters on read.
	// If loaded with this option, dumping will panic.
	DontUnescapeChars bool
}

func DefaultLoadOptions() LoadOptions {
	return LoadOptions{}
}

// Loads a .env file from disk
func (f *EnvFile) Load(path string, opts *LoadOptions) *EnvFile {
	f.path = path
	f.data = make(map[string]string)
	if opts == nil {
		f.opts = DefaultLoadOptions()
	} else {
		f.opts = *opts
	}

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

	envFile, err := readFile(path, f.opts)

	for key, value := range envFile {
		f.data[key] = value
	}

	// for scanner.Scan() {
	// 	if scanner.Err() == nil {
	// 		line := scanner.Text()
	// 		if !strings.HasPrefix(line, "#") && len(line) > 0 {
	// 			if strings.HasSuffix(line, `"`) {
	// 				key, value2 := split(line)
	// 				value = value2

	// 				f.data[key] = value
	// 			} else {
	// 				 value =

	// 			// key, value := split(line)

	// 			// f.data[key] = value
	// 		}
	// 	} else {
	// 		return f.SetError("Failed to read `%s` (%w)", path, err)
	// 	}
	// }

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
// Panics if the EnvFile was loaded with the DontUnescapeChars option
func (f *EnvFile) Dump() *EnvFile {
	if f.Err() != nil {
		return f
	}

	if f.opts.DontUnescapeChars {
		panic("Writing .env file with unescaped chars should not happen")
	}

	var sb strings.Builder

	for key, value := range f.data {
		trimed := strings.Trim(value, " \n\r\t")
		escaped := doubleQuoteEscape(trimed)

		sb.WriteString(fmt.Sprintf("%s=\"%s\"\n", key, escaped))
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

// Parse reads an env file from io.Reader, returning a map of keys and values.
func Parse(r io.Reader, opts LoadOptions) (map[string]string, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return UnmarshalBytes(data, opts)
}

// Read all env (with same file loading semantics as Load) but return values as
// a map rather than automatically writing values into env
func Read(filename string, opts LoadOptions) (envMap map[string]string, err error) {
	envMap = make(map[string]string)

	individualEnvMap, individualErr := readFile(filename, opts)

	if individualErr != nil {
		err = individualErr
		return // return early on a spazout
	}

	for key, value := range individualEnvMap {
		envMap[key] = value
	}

	return
}

// UnmarshalBytes parses env file from byte slice of chars, returning a map of keys and values.
func UnmarshalBytes(src []byte, opts LoadOptions) (map[string]string, error) {
	out := make(map[string]string)
	err := parseBytes(src, out, opts)
	return out, err
}

func readFile(filename string, opts LoadOptions) (envMap map[string]string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	return Parse(file, opts)
}
