package core

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/udhos/equalfile"
	"github.com/wearedevx/keystone/cli/internal/envfile"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/utils"
	"github.com/wearedevx/keystone/cli/ui"

	"github.com/wearedevx/keystone/api/pkg/models"
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

func (ctx *Context) SaveMessages(MessageByEnvironments models.GetMessageByEnvironmentResponse) ChangesByEnvironment {
	changes := ChangesByEnvironment{Environments: make(map[string][]Change)}

	for environmentName, environment := range MessageByEnvironments.Environments {
		var PayloadContent = models.MessagePayload{}

		localSecrets := make([]models.SecretVal, 0)
		for _, localSecret := range ctx.ListSecretsFromCache() {
			localSecrets = append(localSecrets, models.SecretVal{
				Label: localSecret.Name,
				Value: string(localSecret.Values[EnvironmentName(environmentName)]),
			})
		}
		if len(environment.Message.Payload) == 0 {
			continue
		}
		if err := json.Unmarshal(environment.Message.Payload, &PayloadContent); err != nil {
			ctx.err = kserrors.CouldNotParseMessage(err)
			return changes

		}

		environmentChanges := make([]Change, 0)
		fileChanges := ctx.getFilesChanges(PayloadContent.Files, environmentName)

		if len(fileChanges) > 0 {
			filesDir := ctx.CachedEnvironmentFilesPath(environmentName)
			if err := utils.RemoveContents(filesDir); err != nil {
				ctx.err = kserrors.CannotRemoveDirectoryContents(filesDir, err)
				return changes
			}
		}

		err := ctx.saveFilesChanges(fileChanges, environmentName)
		if err != nil {
			ctx.err = kserrors.CannotSaveFiles(err.Error(), err)
			return changes
		}

		secretChanges := GetSecretsChanges(localSecrets, PayloadContent.Secrets)
		// NOTE: if there are changes, the .env file gets rewritten, therefore
		// there is no need to delete it

		environmentChanges = append(environmentChanges, fileChanges...)
		environmentChanges = append(environmentChanges, secretChanges...)

		if len(PayloadContent.Secrets) > 0 {
			envFilePath := ctx.CachedEnvironmentDotEnvPath(environmentName)
			envFile := new(envfile.EnvFile)
			envFile.Load(envFilePath, nil)

			for key := range envFile.GetData() {
				found := false

				for _, secret := range PayloadContent.Secrets {
					if key == secret.Label {
						found = true
						break
					}
				}

				if !found {
					envFile.Unset(key)
				}
			}

			for _, secret := range PayloadContent.Secrets {
				envFile.Set(secret.Label, secret.Value)
			}

			if err := envFile.Dump().Err(); err != nil {
				ctx.setError(kserrors.FailedToUpdateDotEnv(envFilePath, err))
			}
		}

		changes.Environments[environmentName] = environmentChanges
		ctx.UpdateEnvironment(environment.Environment)
	}

	return changes
}

func GetSecretsChanges(localSecrets []models.SecretVal, newSecrets []models.SecretVal) (changes []Change) {
	for _, secret := range newSecrets {
		// Get Secret we want to change in messages
		// Get local value to see if it has changed in fetchedSecrets
		found := false
		for _, localSecret := range localSecrets {
			if secret.Label == localSecret.Label {
				found = true
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
		// New secret has been fetched
		if !found {
			changes = append(changes, Change{
				Name: secret.Label,
				From: "",
				To:   secret.Value,
				Type: "secret",
			})

		}
	}

	// Check for secret that has been deleted
	for _, localSecret := range localSecrets {
		found := false
		for _, newSecret := range newSecrets {
			if newSecret.Label == localSecret.Label {
				found = true
			}
		}

		if !found {
			changes = append(changes, Change{
				Name: localSecret.Label,
				From: localSecret.Value,
				To:   "",
				Type: "secret",
			})

		}
	}

	return changes
}

/// fileHasChanges returns true if the content of file at `pathToExistingFile` is different
// from `candidateContent`, meaning the file contents have changed.
func fileHasChanges(pathToExistingFile string, candidateContent []byte) (sameContent bool, err error) {
	candidateReader := bytes.NewReader(candidateContent)
	/* #nosec
	 * pathToExistingFile must be checked befor call
	 * to ensure that it belongs de ctx.Wd
	 */
	currentFileReader, err := os.Open(pathToExistingFile)

	if err != nil {
		if os.IsNotExist(err) {
			// Not really an error, just create the file
			return true, nil
		}

		return sameContent, err
	}

	comparator := equalfile.New(nil, equalfile.Options{})

	sameContent, err = comparator.CompareReader(currentFileReader, candidateReader)
	if err == nil {
		return !sameContent, nil
	}

	return false, err
}

func (ctx *Context) getFilesChanges(files []models.File, environmentName string) (changes []Change) {
	changes = make([]Change, 0)

	for _, file := range files {
		fileContent, err := base64.StdEncoding.DecodeString(file.Value)
		if err != nil {
			kserrors.InvalidFileContent(file.Path, err).Print()
			continue
		}

		filePath := path.Join(ctx.CachedEnvironmentFilesPath(environmentName), file.Path)

		if !ctx.fileBelongsToContext(filePath) {
			ctx.err = kserrors.FileNotInWorkingDirectory(filePath, ctx.Wd, nil)
			return []Change{}
		}

		fileHasChanges, err := fileHasChanges(filePath, fileContent)
		if err != nil {
			kserrors.FailedCheckingChanges(filePath, err).Print()
			continue
		}

		if fileHasChanges {
			changes = append(changes, Change{
				Type: "file",
				Name: file.Path,
				To:   string(fileContent),
			})
		}
	}
	return changes
}

func (ctx *Context) saveFilesChanges(changes []Change, environmentName string) (err error) {
	errorList := make([]string, 0)
	cacheDir := ctx.CachedEnvironmentFilesPath(environmentName)

	for _, change := range changes {
		cachedFilePath := path.Join(cacheDir, change.Name)

		err = utils.CreateFileIfNotExists(cachedFilePath, change.To)
		if err != nil {
			errorList = append(errorList, err.Error())
			continue
		}

		if err = ioutil.WriteFile(cachedFilePath, []byte(change.To), 0o600); err != nil {
			errorList = append(errorList, err.Error())
			continue
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

	for _, secret := range ctx.ListSecretsFromCache() {
		PayloadContent.Secrets = append(PayloadContent.Secrets, models.SecretVal{
			Label: secret.Name,
			Value: string(secret.Values[EnvironmentName(environment.Name)]),
		})
	}

	envCachePath := ctx.CachedEnvironmentFilesPath(environment.Name)

	for _, file := range ctx.ListFiles() {
		filePath := path.Join(envCachePath, file.Path)
		/* #nosec */
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

// DEPRECATED
func (ctx *Context) WriteNewMessages(messagesByEnvironments models.GetMessageByEnvironmentResponse) (ChangesByEnvironment, *kserrors.Error) {
	changes := ctx.SaveMessages(messagesByEnvironments)

	if err := ctx.Err(); err != nil {
		err.Print()
		return changes, kserrors.UnkownError(err)
	}

	changedEnvironments := make([]string, 0)

	// FIXME: Should not this be an application-wide constant ?
	var envList []string = []string{"dev", "staging", "prod"}

	for _, environmentName := range envList {
		environment, ok := messagesByEnvironments.Environments[environmentName]
		if !ok {
			continue
		}

		messageID := environment.Message.ID

		if messageID != 0 {
			// IF changes detected
			if len(changes.Environments[environmentName]) > 0 {
				ui.Print("Environment " + environmentName + ": " + strconv.Itoa(len(changes.Environments[environmentName])) + " secret(s) changed")

				for _, change := range changes.Environments[environmentName] {
					// No previous cotent => secret is new
					if len(change.From) == 0 {
						ui.Print("++ " + change.Name + " : " + change.To)
					} else {
						ui.Print("   " + change.Name + " : " + change.From + " ↦ " + change.To)
					}

				}
			} else {
				ui.Print("Environment " + environmentName + " up to date ✔")
			}

			if err := ctx.Err(); err != nil {
				err.Print()
				return changes, ctx.err
			}
		} else {
			environmentChanged := ctx.EnvironmentVersionHasChanged(environmentName, environment.Environment.VersionID)

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

func (ctx *Context) CompareNewSecretWithChanges(secretName string, newSecret map[string]string, changesByEnvironment ChangesByEnvironment) *Context {
	if ctx.err != nil {
		return ctx
	}

	environmentValueMap := make(map[string]string)

	for environmentName, changes := range changesByEnvironment.Environments {
		for _, change := range changes {
			if change.Name == secretName {
				environmentValueMap[environmentName] = change.To
			}
		}
	}

	if len(environmentValueMap) > 0 {
		environmentValueMapString := ""

		for environment, value := range environmentValueMap {
			if len(value) == 0 {
				environmentValueMapString += fmt.Sprintf("Secret in %s is deleted.\n", environment)
			} else {
				environmentValueMapString += fmt.Sprintf("Value in %s is '%s'.\n", environment, value)
			}
		}

		ctx.err = kserrors.SecretHasChanged(secretName, environmentValueMapString, nil)
	}

	return ctx
}

func (ctx *Context) CompareRemovedSecretWithChanges(secretName string, changesByEnvironment ChangesByEnvironment) *Context {
	if ctx.err != nil {
		return ctx
	}

	environmentValueMap := make(map[string]string)

	for environmentName, changes := range changesByEnvironment.Environments {
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

		ctx.err = kserrors.SecretHasChanged(secretName, environmentValueMapString, nil)
	}

	return ctx
}

func (ctx *Context) CompareNewFileWhithChanges(filePath string, changesByEnvironment ChangesByEnvironment) *Context {
	if ctx.err != nil {
		return ctx
	}

	affectedEnvironments := make([]string, 0)
	for environmentName, changes := range changesByEnvironment.Environments {
		for _, change := range changes {
			if change.Name == filePath {
				affectedEnvironments = append(affectedEnvironments, environmentName)
			}
		}
	}

	if len(affectedEnvironments) > 0 {
		ctx.err = kserrors.FileHasChanged(filePath, strings.Join(affectedEnvironments, ","), nil)
	}

	return ctx

}
