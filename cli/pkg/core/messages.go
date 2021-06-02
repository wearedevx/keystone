package core

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/udhos/equalfile"
	"github.com/wearedevx/keystone/api/pkg/models"
	. "github.com/wearedevx/keystone/cli/internal/envfile"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	. "github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/pkg/client"
	"github.com/wearedevx/keystone/cli/ui"
	"gopkg.in/yaml.v2"
)

type Change struct {
	Name string
	From string
	To   string
	Type string // secret || file
}

type ChangesByEnvironment struct {
	Environments map[string][]Change
}

func (ctx *Context) SaveMessages(MessageByEnvironments models.GetMessageByEnvironmentResponse) (ChangesByEnvironment, error) {
	changes := ChangesByEnvironment{Environments: make(map[string][]Change)}

	for environmentName, environment := range MessageByEnvironments.Environments {
		var PayloadContent = models.MessagePayload{}
		envFilePath := ctx.CachedEnvironmentDotEnvPath(environmentName)

		localSecrets := make([]models.SecretVal, 0)
		for _, localSecret := range ctx.ListSecrets() {
			localSecrets = append(localSecrets, models.SecretVal{
				Label: localSecret.Name,
				Value: string(localSecret.Values[EnvironmentName(environmentName)]),
			})
		}

		if err := yaml.Unmarshal(environment.Message.Payload, &PayloadContent); err != nil {
			// TODO: Error Handling
			panic(err)
		}

		// Remove content of cache directory to ensure old files are deleted
		RemoveContents(ctx.CachedEnvironmentPath(environmentName))

		environmentChanges := make([]Change, 0)
		fileChanges, err := ctx.getFilesChanges(PayloadContent.Files, environmentName)
		if err != nil {
			// TODO: Error Handling
			panic(err)
		}

		err = ctx.saveFilesChanges(fileChanges)
		if err != nil {
			// TODO: Error Handling
			panic(err)
		}

		secretChanges := GetSecretsChanges(localSecrets, PayloadContent.Secrets)

		environmentChanges = append(environmentChanges, fileChanges...)
		environmentChanges = append(environmentChanges, secretChanges...)

		for _, secret := range PayloadContent.Secrets {
			if err := new(EnvFile).
				Load(envFilePath).
				Set(secret.Label, secret.Value).
				Dump().
				Err(); err != nil {
				err = kserrors.FailedToUpdateDotEnv(envFilePath, err)
			}
		}

		changes.Environments[environmentName] = environmentChanges
	}

	return changes, nil
}

func GetSecretsChanges(localSecrets []models.SecretVal, newSecrets []models.SecretVal) (changes []Change) {
	for _, secret := range newSecrets {
		// Get Secret we want to change in messages
		// Get local value to see if it has changed in fetchedSecrets
		for _, localSecret := range localSecrets {
			if secret.Label == localSecret.Label {
				// Compare local value with value from message
				if secret.Value != localSecret.Value {
					changes = append(changes, Change{
						Name: secret.Label,
						From: localSecret.Value,
						To:   secret.Value,
						Type: "secret",
					})
				}
			}
		}
	}

	return changes
}

/// fileHasChanges returns true if the content of file at `pathToExistingFile` is different
// from `candidateContent`, meaning the file contents have changed.
func fileHasChanges(pathToExistingFile string, candidateContent []byte) (sameContent bool, err error) {
	candidateReader := bytes.NewReader(candidateContent)
	currentFileReader, err := os.Open(pathToExistingFile)

	if err != nil {
		return sameContent, err
	}

	comparator := equalfile.New(nil, equalfile.Options{})

	sameContent, err = comparator.CompareReader(currentFileReader, candidateReader)
	if err == nil {
		return !sameContent, nil
	}

	return false, err
}

func (ctx *Context) getFilesChanges(files []models.File, environmentName string) (changes []Change, err error) {
	changes = make([]Change, 0)

	for _, file := range files {
		fileContent, err := base64.StdEncoding.DecodeString(file.Value)
		if err != nil {
			//TODO: Prettify this error
			fmt.Println("ERROR decoding base64 decrypted file content", err.Error())
			continue
		}

		filePath := path.Join(ctx.CachedEnvironmentFilesPath(environmentName), file.Path)

		fileHasChanges, err := fileHasChanges(filePath, fileContent)
		if err != nil {
			//TODO: Prettify this error
			fmt.Println("ERROR checking for file changes: ", err.Error())
		}

		if fileHasChanges {
			changes = append(changes, Change{
				Type: "file",
				Name: filePath,
				To:   string(fileContent),
			})
		}
	}
	return changes, err
}

func (ctx *Context) saveFilesChanges(changes []Change) (err error) {
	errorList := make([]string, 0)

	for _, change := range changes {
		err = CreateFileIfNotExists(change.Name, change.To)
		if err != nil {
			errorList = append(errorList, err.Error())
		}
	}

	if len(errorList) > 0 {
		errMessage := strings.Join(errorList, "\n")

		return errors.New(errMessage)
	}

	return nil
}

// Return PayloadContent, with secrets and files of current environment.
func (ctx *Context) PrepareMessagePayload(environment models.Environment) (models.MessagePayload, error) {
	var PayloadContent = models.MessagePayload{
		Files:   make([]models.File, 0),
		Secrets: make([]models.SecretVal, 0),
	}

	var err error

	errors := make([]string, 0)

	for _, secret := range ctx.ListSecrets() {
		PayloadContent.Secrets = append(PayloadContent.Secrets, models.SecretVal{
			Label: secret.Name,
			Value: string(secret.Values[EnvironmentName(environment.Name)]),
		})
	}

	envCachePath := ctx.CachedEnvironmentFilesPath(environment.Name)

	for _, file := range ctx.ListFiles() {
		filePath := path.Join(envCachePath, file.Path)
		fileContent, err := ioutil.ReadFile(filePath)

		if err != nil {
			errors = append(errors, err.Error())
		}

		PayloadContent.Files = append(PayloadContent.Files, models.File{
			Path:  file.Path,
			Value: base64.StdEncoding.EncodeToString(fileContent),
		})
	}

	if len(errors) > 0 {
		err = fmt.Errorf(strings.Join(errors, "\n"))
	}

	return PayloadContent, err
}

func (ctx *Context) FetchNewMessages(result *models.GetMessageByEnvironmentResponse) error {
	// Get keystone key.
	c, kcErr := client.NewKeystoneClient()

	if kcErr != nil {
		kcErr.Print()
		os.Exit(1)
	}

	projectID := ctx.GetProjectID()

	r, err := c.Messages().GetMessages(projectID)
	if err == nil {
		*result = r
	}

	return err
}

func (ctx *Context) WriteNewMessages(messagesByEnvironments models.GetMessageByEnvironmentResponse) (ChangesByEnvironment, *kserrors.Error) {
	c, kcErr := client.NewKeystoneClient()

	if kcErr != nil {
		kcErr.Print()
		os.Exit(1)
	}

	changes, _ := ctx.SaveMessages(messagesByEnvironments)

	if err := ctx.Err(); err != nil {
		err.Print()
		return changes, kserrors.UnkownError(err)
	}

	changedEnvironments := make([]string, 0)

	for environmentName, environment := range messagesByEnvironments.Environments {
		messageID := environment.Message.ID

		if messageID != 0 {
			// IF changes detected
			if len(changes.Environments[environmentName]) > 0 {
				ui.Print("Environment " + environmentName + ": " + strconv.Itoa(len(changes.Environments[environmentName])) + " secret(s) or file(s) changed")
				for _, change := range changes.Environments[environmentName] {
					if change.Type == "secret" {
						ui.Print("secret " + change.Name + ": " + change.From + " ↦ " + change.To)
					} else {
						ui.Print("file " + change.Name + " changed")
					}
				}
			} else {
				ui.Print("Environment " + environmentName + " up to date ✔")
			}

			response, _ := c.Messages().DeleteMessage(environment.Message.ID)

			if !response.Success {
				ui.Print("Can't delete message " + response.Error)
			} else {
				ctx.SetEnvironmentVersion(environmentName, environment.VersionID)
			}
		} else {
			environmentChanged := ctx.EnvironmentVersionHasChanged(environmentName, environment.VersionID)

			if environmentChanged {
				ui.Print("Environment " + environmentName + " has changed but no message available. Ask someone to push their secret ⨯")
				changedEnvironments = append(changedEnvironments, environmentName)
			} else {
				ui.Print("Environment " + environmentName + " up to date ✔")
			}
		}
	}

	if len(changedEnvironments) > 0 {
		return changes, kserrors.EnvironmentsHaveChanged(strings.Join(changedEnvironments, ", "), nil)
	}

	return changes, nil
}

func (ctx *Context) CompareNewSecretWithChanges(secretName string, newSecret map[string]string, changesByEnvironment ChangesByEnvironment) *kserrors.Error {
	// where are stored changed values
	environmentValueMap := make(map[string]string)

	for environmentName, changes := range changesByEnvironment.Environments {
		// var PayloadContent = models.MessagePayload{}
		// if err := yaml.Unmarshal(message.Message.Payload, &PayloadContent); err != nil {
		// 	panic(err)
		// }
		for _, change := range changes {

			if change.Name == secretName {
				environmentValueMap[environmentName] = change.To
			}
		}

	}

	if len(environmentValueMap) > 0 {
		environmentValueMapString := ""
		for environment, value := range environmentValueMap {
			environmentValueMapString += fmt.Sprintf("Value in '%s' is '%s'.\n", environment, value)
		}
		return kserrors.SecretHasChanged(secretName, environmentValueMapString, nil)
	}
	return nil
}
