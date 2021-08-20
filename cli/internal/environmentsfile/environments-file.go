package environmentsfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	. "github.com/wearedevx/keystone/api/pkg/models"
	. "github.com/wearedevx/keystone/cli/internal/utils"
	"gopkg.in/yaml.v2"
)

type environmentsFileOptions struct {
	Strict bool
}

type Env struct {
	EnvironmentID string `yaml:"id"`
	Name          string `yaml:"name"`
	VersionID     string `yaml:"version_id"`
}

type EnvironmentsFile struct {
	path         string `yaml:"-"`
	err          error  `yaml:"-"`
	Current      string `yaml:"current"`
	Environments []Env  `yaml:"environments"`
}

// Keystone file path for the given context
func environmentsFilePath(dotKeystonePath string) string {
	return path.Join(dotKeystonePath, "environments.yaml")
}

func NewEnvironmentsFile(dotKeystonePath string, updatedEnvironments []Environment) *EnvironmentsFile {
	envs := make([]Env, 0)

	for _, env := range updatedEnvironments {
		envs = append(envs, Env{fmt.Sprint(env.EnvironmentID), env.Name, env.VersionID})
	}
	return &EnvironmentsFile{
		path:         environmentsFilePath(dotKeystonePath),
		err:          nil,
		Current:      "dev",
		Environments: envs,
	}

}

// Checks if current execution context contains a keystone.yaml
func ExistsEnvironmentsFile(dotKeystonePath string) bool {
	return FileExists(environmentsFilePath(dotKeystonePath))
}

// Loads a Keystone from disk
func (file *EnvironmentsFile) Load(dotKeystonePath string) *EnvironmentsFile {
	bytes, err := ioutil.ReadFile(environmentsFilePath(dotKeystonePath))
	// file := newKeystoneFile(context)

	if err != nil {
		file.err = err
	}

	file.path = environmentsFilePath(dotKeystonePath)

	return file.fromYaml(bytes)
}

// Loads contents of yml file into a KeystoneFile struct
func (file *EnvironmentsFile) fromYaml(bytes []byte) *EnvironmentsFile {
	if file.Err() != nil {
		return file
	}

	file.err = yaml.Unmarshal(bytes, &file)

	return file
}

// Turns a EnvironmentsFile into a ynl
func (file *EnvironmentsFile) toYaml() []byte {
	if file.Err() != nil {
		return []byte{}
	}

	bytes, err := yaml.Marshal(file)

	if err != nil {
		file.err = fmt.Errorf("Could not serialize environments file (%w)", err)
	}

	return bytes
}

// Accessor for the EnvironmentsFile's err field
// use for error management
func (file *EnvironmentsFile) Err() error {
	return file.err
}

// Writes the Keystone File to disk
func (file *EnvironmentsFile) Save() *EnvironmentsFile {
	if file.Err() == nil {
		yamlBytes := file.toYaml()

		if err := ioutil.WriteFile(file.path, yamlBytes, 0600); err != nil {
			file.err = fmt.Errorf("Could not write `environments.yaml` (%w)", err)
		}
	}

	return file

}

// Removes the environments file from disk
func (file *EnvironmentsFile) Remove() {
	if file.Err() != nil {
		return
	}

	if err := os.Remove(file.path); err != nil {
		file.err = fmt.Errorf("Could not remove `environments.yaml` (%w)", err)
	}

	return
}

// Adds a variable to the project
// set strict to true if you want to throw an error when it is missing
func (file *EnvironmentsFile) SetVersion(environmentName string, versionID string) *EnvironmentsFile {
	if file.Err() != nil {
		return file
	}

	for i, environment := range file.Environments {
		if environmentName == environment.Name {
			file.Environments[i].VersionID = versionID
		}
	}

	return file
}

func (file *EnvironmentsFile) SetCurrent(environmentName string) *EnvironmentsFile {
	if file.Err() != nil {
		return file
	}

	file.Current = environmentName

	return file
}

func (file *EnvironmentsFile) GetByName(environmentName string) *Env {
	if file.Err() != nil {
		return nil
	}
	for _, env := range file.Environments {
		if env.Name == environmentName {
			return &env

		}
	}
	return nil
}

// Replaces an environment in the environment file with updated data
// If the environment does not exist in the environment file,
// it should be appended to it
func (file *EnvironmentsFile) Replace(environment Environment) *EnvironmentsFile {
	if file.Err() != nil {
		return file
	}

	newEnvironments := make([]Env, 0)

	for _, env := range file.Environments {
		if env.EnvironmentID != environment.EnvironmentID {
			newEnvironments = append(newEnvironments, env)
		}
	}

	newEnvironments = append(newEnvironments, Env{
		EnvironmentID: environment.EnvironmentID,
		Name:          environment.Name,
		VersionID:     environment.VersionID,
	})

	file.Environments = newEnvironments

	return file
}

func (file *EnvironmentsFile) Path() string {
	return file.path
}

// // Removes a variable from the project
// func (file *EnvironmentsFile) UnsetEnv(varname string) *KeystoneFile {
// 	if file.Err() != nil {
// 		return file
// 	}

// 	envs := make([]envKey, 0)

// 	// Filter out previously existing value
// 	for _, env := range file.Env {
// 		if env.Key != varname {
// 			envs = append(envs, env)
// 		}
// 	}

// 	file.Env = eenvironmentFile.Save/ }
