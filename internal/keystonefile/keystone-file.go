package keystonefile

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	uuid "github.com/satori/go.uuid"
	. "github.com/wearedevx/keystone/internal/utils"
	"gopkg.in/yaml.v2"
)

type envKey struct {
	Key    string
	Strict bool
}

type keystoneFileOptions struct {
	Strict bool
}

// Represents the contents of the keystone.yml file
type KeystoneFile struct {
	path      string `yaml:"-"`
	err       error  `yaml:"-"`
	ProjectId string `yaml:"project_id"`
	Env       []envKey
	Files     []string
	Options   keystoneFileOptions
}

// Keystone file path for the given context
func keystoneFilePath(wd string) string {

	return path.Join(wd, "keystone.yml")
}

//
func NewKeystoneFile(wd string) *KeystoneFile {
	project_id := uuid.NewV4()

	return &KeystoneFile{
		path:      keystoneFilePath(wd),
		err:       nil,
		ProjectId: project_id.String(),
		Env:       make([]envKey, 0),
		Files:     make([]string, 0),
		Options: keystoneFileOptions{
			Strict: false,
		},
	}
}

// Checks if current execution context contains a keystone.yml
func ExistsKeystoneFile(wd string) bool {
	return FileExists(keystoneFilePath(wd))
}

// Loads a Keystone from disk
func (file *KeystoneFile) Load(wd string) *KeystoneFile {
	bytes, err := ioutil.ReadFile(keystoneFilePath(wd))
	// file := newKeystoneFile(context)

	if err != nil {
		file.err = err
	}

	file.path = keystoneFilePath(wd)

	return file.fromYaml(bytes)
}

// Loads contents of yml file into a KeystoneFile struct
func (file *KeystoneFile) fromYaml(bytes []byte) *KeystoneFile {
	if file.Err() != nil {
		return file
	}

	file.err = yaml.Unmarshal(bytes, &file)

	return file
}

// Turns a KeystoneFile into a ynl
func (file *KeystoneFile) toYaml() []byte {
	if file.Err() != nil {
		return []byte{}
	}

	bytes, err := yaml.Marshal(file)

	if err != nil {
		file.err = fmt.Errorf("Could not serialize keystone file (%w)", err)
	}

	return bytes
}

// Accessor for the KeystoneFile's err field
// use for error management
func (file *KeystoneFile) Err() error {
	return file.err
}

// Writes the Keystone File to disk
func (file *KeystoneFile) Save() *KeystoneFile {
	if file.Err() == nil {
		yamlBytes := file.toYaml()

		if err := ioutil.WriteFile(file.path, yamlBytes, 0644); err != nil {
			file.err = fmt.Errorf("Could not write `keystone.yml` (%w)", err)
		}
	}

	return file

}

// Removes the keystone file from disk
func (file *KeystoneFile) Remove() {
	if file.Err() != nil {
		return
	}

	if err := os.Remove(file.path); err != nil {
		file.err = fmt.Errorf("Could not remove `keystone.yml` (%w)", err)
	}

	return
}

// Adds a variable to the project
// set strict to true if you want to throw an error when it is missing
func (file *KeystoneFile) SetEnv(varname string, strict bool) *KeystoneFile {
	if file.Err() != nil {
		return file
	}

	file.UnsetEnv(varname) // avoid duplicates

	file.Env = append(file.Env, envKey{
		Key:    varname,
		Strict: strict,
	})

	return file
}

// Removes a variable from the project
func (file *KeystoneFile) UnsetEnv(varname string) *KeystoneFile {
	if file.Err() != nil {
		return file
	}

	envs := make([]envKey, 0)

	// Filter out previously existing value
	for _, env := range file.Env {
		if env.Key != varname {
			envs = append(envs, env)
		}
	}

	file.Env = envs

	return file
}

func (file *KeystoneFile) AddFile(filepath string) *KeystoneFile {
	if file.Err() != nil {
		return file
	}

	file.RemoveFile(filepath)

	file.Files = append(file.Files, filepath)

	return file
}

func (file *KeystoneFile) RemoveFile(filepath string) *KeystoneFile {
	if file.Err() != nil {
		return file
	}

	files := make([]string, 0)

	for _, f := range file.Files {
		if f != filepath {
			files = append(files, f)
		}
	}

	file.Files = files

	return file
}