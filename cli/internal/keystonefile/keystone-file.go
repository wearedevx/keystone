package keystonefile

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/wearedevx/keystone/api/pkg/models"
	"gopkg.in/yaml.v2"

	"github.com/wearedevx/keystone/cli/internal/utils"
)

type EnvKey struct {
	Key    string
	Strict bool
}

type FileKey struct {
	Path      string
	Strict    bool
	FromCache bool `yaml:"-"`
}

type keystoneFileOptions struct {
	Strict bool
}

// Represents the contents of the keystone.yaml file
type KeystoneFile struct {
	Path        string `yaml:"-"`
	err         error  `yaml:"-"`
	ProjectId   string `yaml:"project_id"`
	ProjectName string `yaml:"name"`
	Env         []EnvKey
	Files       []FileKey
	Options     keystoneFileOptions
	CiServices  []CiService `yaml:"ci_services"`
}

var ksf *KeystoneFile
var loaded bool = false

type CiService struct {
	Name    string            `yaml:"name"`
	Type    string            `yaml:"type"`
	Options map[string]string `yaml:"options"`
}

// Keystone file path for the given context
func keystoneFilePath(wd string) string {
	return path.Join(wd, "keystone.yaml")
}

// NewKeystoneFile function returns a new instance of a KeyatoneFile
func NewKeystoneFile(wd string, project models.Project) *KeystoneFile {
	return &KeystoneFile{
		Path:        keystoneFilePath(wd),
		err:         nil,
		ProjectId:   project.UUID,
		ProjectName: project.Name,
		Env:         make([]EnvKey, 0),
		Files:       make([]FileKey, 0),
		Options: keystoneFileOptions{
			Strict: false,
		},
	}
}

func LoadKeystoneFile(wd string) *KeystoneFile {
	if loaded && ksf != nil {
		return ksf
	} else {
		k := new(KeystoneFile).Load(wd)
		if k.err != nil {
			panic(k.err)
		} else {
			loaded = true
			ksf = k
		}

		return ksf
	}
}

// Checks if current execution context contains a keystone.yaml
func ExistsKeystoneFile(wd string) bool {
	return utils.FileExists(keystoneFilePath(wd))
}

// Loads a Keystone from disk
func (file *KeystoneFile) Load(wd string) *KeystoneFile {
	var bytes []byte
	var err error

	file.Path = keystoneFilePath(wd)
	/* #nosec
	 * We generate the file path, and its content is about to be parsed
	 */
	bytes, err = ioutil.ReadFile(file.Path)
	if err != nil {
		panic(err)
		// file.err = err
		//
		// return file
	}

	return file.fromYaml(bytes)
}

// Loads contents of yml file into a KeystoneFile struct
func (file *KeystoneFile) fromYaml(bytes []byte) *KeystoneFile {
	if file.Err() != nil {
		return file
	}

	err := yaml.Unmarshal(bytes, &file)
	if err != nil {
		file.err = fmt.Errorf("Parsing error: %w", err)
	}

	return file
}

// Turns a KeystoneFile into a ynl
func (file *KeystoneFile) toYaml() []byte {
	if file.Err() != nil {
		return []byte{}
	}

	bytes, err := yaml.Marshal(file)
	if err != nil {
		file.err = fmt.Errorf("could not serialize keystone file (%w)", err)
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

		if err := ioutil.WriteFile(file.Path, yamlBytes, 0o600); err != nil {
			file.err = fmt.Errorf("could not write `keystone.yaml` (%w)", err)
		}
	}

	return file
}

// Removes the keystone file from disk
func (file *KeystoneFile) Remove() *KeystoneFile {
	if file.Err() != nil {
		return file
	}

	if err := os.Remove(file.Path); err != nil {
		file.err = fmt.Errorf("could not remove `keystone.yaml` (%w)", err)
	}

	return file
}

// Adds a variable to the project
// set strict to true if you want to throw an error when it is missing
func (file *KeystoneFile) SetEnv(varname string, strict bool) *KeystoneFile {
	if file.Err() != nil {
		return file
	}

	file.UnsetEnv(varname) // avoid duplicates

	file.Env = append(file.Env, EnvKey{
		Key:    varname,
		Strict: strict,
	})

	return file
}

// HasEnv method returns true the environment variable `varname` exists
// in the keystone file
func (file *KeystoneFile) HasEnv(varname string) (hasIt bool, strict bool) {
	if file.Err() != nil {
		return false, false
	}

	for _, ek := range file.Env {
		if ek.Key == varname {
			return true, ek.Strict
		}
	}

	return false, false
}

// Removes a variable from the project
func (file *KeystoneFile) UnsetEnv(varname string) *KeystoneFile {
	if file.Err() != nil {
		return file
	}

	envs := make([]EnvKey, 0)

	// Filter out previously existing value
	for _, env := range file.Env {
		if env.Key != varname {
			envs = append(envs, env)
		}
	}

	file.Env = envs

	return file
}

// AddFile method adds a tracked file to the keystone file
func (file *KeystoneFile) AddFile(filekey FileKey) *KeystoneFile {
	if file.Err() != nil {
		return file
	}

	file.RemoveFile(filekey.Path)

	file.Files = append(file.Files, filekey)

	return file
}

// RemoveFile method removes a file from the keysotne file
func (file *KeystoneFile) RemoveFile(filepath string) *KeystoneFile {
	if file.Err() != nil {
		return file
	}

	files := make([]FileKey, 0)

	for _, f := range file.Files {
		if f.Path != filepath {
			files = append(files, f)
		}
	}

	file.Files = files

	return file
}

// SetFileRequired method marks a file a required in the keystone file
func (file *KeystoneFile) SetFileRequired(
	filepath string,
	required bool,
) *KeystoneFile {
	if file.Err() != nil {
		return file
	}

	for index, f := range file.Files {
		if f.Path == filepath {
			file.Files[index].Strict = required
		}
	}

	return file
}

// AddCiService method adds a CI service to the keystone file
func (file *KeystoneFile) AddCiService(ciService CiService) *KeystoneFile {
	if file.Err() != nil {
		return file
	}

	file.RemoveCiService(ciService.Name)
	file.CiServices = append(file.CiServices, ciService)

	return file
}

// RemoveCiService method removes a CI service from the keystone file
func (file *KeystoneFile) RemoveCiService(serviceName string) *KeystoneFile {
	if file.Err() != nil {
		return file
	}

	services := make([]CiService, 0)

	for _, service := range file.CiServices {
		if service.Name != serviceName {
			services = append(services, service)
		}
	}

	file.CiServices = services

	return file
}

// GetCiService method returns the CI service named `serviceName`
func (file *KeystoneFile) GetCiService(serviceName string) CiService {
	var ciService CiService
	for _, service := range file.CiServices {
		if service.Name == serviceName {
			ciService = service
		}
	}
	return ciService
}
