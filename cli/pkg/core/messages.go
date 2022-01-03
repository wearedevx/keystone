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
	"strings"

	"github.com/udhos/equalfile"
	"github.com/wearedevx/keystone/cli/internal/envfile"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/utils"

	"github.com/wearedevx/keystone/api/pkg/models"
)

type ChangeType string

const (
	ChangeTypeSecretAdd    ChangeType = "secret"
	ChangeTypeSecretChange ChangeType = "change"
	ChangeTypeSecretDelete ChangeType = "delete"
	ChangeTypeFile         ChangeType = "file"
	// This one happens when environment version changed
	// but there is no messages along with it
	ChangeTypeVersion ChangeType = "version"
)

// Change struct describers sercret/file changes
type Change struct {
	Name string     // Secret name or file path
	From string     // Value before, empty if creation
	To   string     // New value, empty if deletion
	Type ChangeType // Type of change (file, secret, version only)
}

// IsSecretAdd method tells if the change is adding a secret
func (c Change) IsSecretAdd() bool {
	return c.Type == ChangeTypeSecretAdd
}

// IsSecretChange method tells if the change is changing a secret
func (c Change) IsSecretChange() bool {
	return c.Type == ChangeTypeSecretChange
}

// IsSecretDelete method tells if the change is removing a secret
func (c Change) IsSecretDelete() bool {
	return c.Type == ChangeTypeSecretDelete
}

// IsFile method tells if the change is about a file
func (c Change) IsFile() bool {
	return c.Type == ChangeTypeFile
}

type Changes []Change

// IsSingleVersionChange method return true if
// the list of changes has single version change, which means
// the environment changed, but we don't have a message for it
func (changes Changes) IsSingleVersionChange() bool {
	if len(changes) == 1 {
		return changes[0].Type == ChangeTypeVersion
	}

	return false
}

// IsEmpty method tells if the list is empty
func (changes Changes) IsEmpty() bool {
	return len(changes) == 0
}

// ChangesByEnvironment struct is a list of changes grouped by environment
type ChangesByEnvironment struct {
	Environments map[string]Changes
}

/// Returns a list of all environments that have a different
/// VersionID but no payload for the user.
func (ce ChangesByEnvironment) ChangedEnvironmentsWithoutPayload() []string {
	result := make([]string, 0)

	for environmentName, changesList := range ce.Environments {
		if changesList.IsSingleVersionChange() {
			result = append(result, environmentName)
		}
	}

	return result
}

/// Writes the contents of messages in the project's cache
/// It will replace secrets value with the new ones,
/// add new secrets and their values
/// remove secrets that have been removed by other user
func (ctx *Context) SaveMessages(
	MessageByEnvironments models.GetMessageByEnvironmentResponse,
) ChangesByEnvironment {
	changes := ChangesByEnvironment{Environments: make(map[string]Changes)}
	cachedLocalSecrets := ctx.ListSecretsFromCache()

	ctx.log.Println("Saving Messages")

	for environmentName, environment := range MessageByEnvironments.Environments {
		// ——— Preparation work ———
		PayloadContent := models.MessagePayload{}

		ctx.log.Printf("-- Environment %s\n", environmentName)

		// If the payload is empty, and the enviroment version has changed,
		// signal the version change.
		switch {
		case ctx.isVersionChangeWithoutPayload(environment):
			changes.Environments[environmentName] = []Change{
				{
					Type: ChangeTypeVersion,
				},
			}
			ctx.log.Printf(
				"\tNew version %s, no payload\n",
				environment.Environment.VersionID,
			)

			continue

		case payloadIsEmpty(environment):
			ctx.log.Println("\tNo payload")
			continue
		}

		// Parse the decrypted message payload
		if err := json.Unmarshal(environment.Message.Payload, &PayloadContent); err != nil {
			ctx.err = kserrors.CouldNotParseMessage(err)
			return changes
		}

		// ——— Handle files ———
		environmentChanges := make([]Change, 0)
		fileChanges := ctx.getFilesChanges(
			PayloadContent.Files,
			environmentName,
		)

		ctx.log.Printf("\tFile Changes: %d\n", len(fileChanges))

		if err := ctx.handleFileChanges(fileChanges, environmentName).Err(); err != nil {
			return changes
		}

		// ——— Handle secrets ———

		// We need the local values for every secret for `environment`
		localSecrets := secretsForEnvironment(
			cachedLocalSecrets,
			environmentName,
		)
		secretChanges := GetSecretsChanges(localSecrets, PayloadContent.Secrets)

		ctx.log.Printf("\tSecrets Changes: %d\n", len(fileChanges))

		if err := ctx.handleSecretChanges(PayloadContent, environmentName).Err(); err != nil {
			return changes
		}

		environmentChanges = append(environmentChanges, fileChanges...)
		environmentChanges = append(environmentChanges, secretChanges...)

		changes.Environments[environmentName] = environmentChanges

		ctx.UpdateEnvironment(environment.Environment)

		ctx.log.Println("-- DONE")
	}

	return changes
}

func secretsForEnvironment(
	cachedLocalSecrets []Secret,
	environmentName string,
) []models.SecretVal {
	localSecrets := make([]models.SecretVal, 0)
	for _, localSecret := range cachedLocalSecrets {
		localSecrets = append(localSecrets, models.SecretVal{
			Label: localSecret.Name,
			Value: string(
				localSecret.Values[EnvironmentName(environmentName)],
			),
		})
	}
	return localSecrets
}

func payloadIsEmpty(response models.GetMessageResponse) bool {
	return len(response.Message.Payload) == 0
}

func (ctx *Context) isVersionChangeWithoutPayload(
	response models.GetMessageResponse,
) bool {
	if ctx.err != nil {
		return false
	}

	if len(response.Message.Payload) == 0 {
		// If the message is empty, but the versionID has changed,
		return ctx.EnvironmentVersionHasChanged(
			response.Environment.Name,
			response.Environment.VersionID,
		)
	}
	return false
}

func (ctx *Context) handleFileChanges(
	fileChanges []Change,
	environmentName string,
) *Context {
	if ctx.err != nil {
		return ctx
	}

	// Only remove cached files that have changed
	for _, fileChange := range fileChanges {
		filesDir := ctx.CachedEnvironmentFilesPath(environmentName)
		filePath := path.Join(filesDir, fileChange.Name)

		err := utils.RemoveFile(filePath)
		if err != nil {
			ctx.err = kserrors.CannotRemoveDirectoryContents(filePath, err)
		}

		ctx.log.Printf("\t\tRemove %s\n", filePath)
	}

	if err := ctx.saveFilesChanges(fileChanges, environmentName); err != nil {
		ctx.err = kserrors.CannotSaveFiles(err.Error(), err)
	}

	return ctx
}

func (ctx *Context) handleSecretChanges(
	PayloadContent models.MessagePayload,
	environmentName string,
) *Context {
	if ctx.err != nil {
		return ctx
	}

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
				ctx.log.Printf("\t\tRemoved %s", key)
			}
		}

		for _, secret := range PayloadContent.Secrets {
			envFile.Set(secret.Label, secret.Value)
			ctx.log.Printf("\t\tSlt %s", secret.Label)
		}

		if err := envFile.Dump().Err(); err != nil {
			ctx.setError(kserrors.FailedToUpdateDotEnv(envFilePath, err))
		}
	}

	return ctx
}

// GetSecretsChanges function
// Returns a list of changes for secrets
func GetSecretsChanges(
	localSecrets []models.SecretVal,
	newSecrets []models.SecretVal,
) (changes []Change) {
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
						Type: ChangeTypeSecretChange,
					})
				}
			}
		}
		// New secret has been fetched
		if !found {
			changes = append(changes, Change{
				Name: secret.Label,
				To:   secret.Value,
				Type: ChangeTypeSecretAdd,
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
				Type: ChangeTypeSecretDelete,
			})
		}
	}

	return changes
}

/// fileHasChanges returns true if the content of file at `pathToExistingFile` is different
// from `candidateContent`, meaning the file contents have changed.
func fileHasChanges(
	pathToExistingFile string,
	candidateContent []byte,
) (sameContent bool, err error) {
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

	sameContent, err = comparator.CompareReader(
		currentFileReader,
		candidateReader,
	)
	if err == nil {
		return !sameContent, nil
	}

	return false, err
}

// getFilesChanges method builds a list of the files that have changed after
// applying the message
func (ctx *Context) getFilesChanges(
	files []models.File,
	environmentName string,
) (changes []Change) {
	changes = make([]Change, 0)

	for _, file := range files {
		fileContent, err := base64.StdEncoding.DecodeString(file.Value)
		if err != nil {
			kserrors.InvalidFileContent(file.Path, err).Print()
			continue
		}

		filePath := path.Join(
			ctx.CachedEnvironmentFilesPath(environmentName),
			file.Path,
		)

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
				Type: ChangeTypeFile,
				Name: file.Path,
				To:   string(fileContent),
			})
		}
	}
	return changes
}

func (ctx *Context) saveFilesChanges(
	changes []Change,
	environmentName string,
) (err error) {
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

		ctx.log.Printf("\t\tSave file %s", cachedFilePath)
	}

	if len(errorList) > 0 {
		errMessage := strings.Join(errorList, "\n")

		return errors.New(errMessage)
	}

	return nil
}

// Return PayloadContent, with secrets and files of current environment.
func (ctx *Context) PrepareMessagePayload(
	environment models.Environment,
) (models.MessagePayload, error) {
	PayloadContent := models.MessagePayload{
		Files:   make([]models.File, 0),
		Secrets: make([]models.SecretVal, 0),
	}

	var err error

	errors := make([]string, 0)

	for _, secret := range ctx.ListSecretsFromCache() {
		PayloadContent.Secrets = append(
			PayloadContent.Secrets,
			models.SecretVal{
				Label: secret.Name,
				Value: string(secret.Values[EnvironmentName(environment.Name)]),
			},
		)
	}

	envCachePath := ctx.CachedEnvironmentFilesPath(environment.Name)

	for _, file := range ctx.ListCachedFilesForEnvironment(environment.Name) {
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

// CompareNewSecretWithChanges method compares recent changes with new
// secret name and values
func (ctx *Context) CompareNewSecretWithChanges(
	secretName string,
	newSecret map[string]string,
	changesByEnvironment ChangesByEnvironment,
) *Context {
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
				environmentValueMapString += fmt.Sprintf(
					"Secret in %s is deleted.\n",
					environment,
				)
			} else {
				environmentValueMapString += fmt.Sprintf("Value in %s is '%s'.\n", environment, value)
			}
		}

		ctx.err = kserrors.SecretHasChanged(
			secretName,
			environmentValueMapString,
			nil,
		)
	}

	return ctx
}

// CompareRemovedSecretWithChanges method compares the removed secrete with
// the latest set of changes
func (ctx *Context) CompareRemovedSecretWithChanges(
	secretName string,
	changesByEnvironment ChangesByEnvironment,
) *Context {
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
			environmentValueMapString += fmt.Sprintf(
				"Value in '%s' is '%s'.\n",
				environment,
				value,
			)
		}

		ctx.err = kserrors.SecretHasChanged(
			secretName,
			environmentValueMapString,
			nil,
		)
	}

	return ctx
}

// CompareNewFileWhithChanges method compares the latest set of changes
// with a new file
func (ctx *Context) CompareNewFileWhithChanges(
	filePath string,
	changesByEnvironment ChangesByEnvironment,
) *Context {
	if ctx.err != nil {
		return ctx
	}

	affectedEnvironments := make([]string, 0)
	for environmentName, changes := range changesByEnvironment.Environments {
		for _, change := range changes {
			if change.Name == filePath {
				affectedEnvironments = append(
					affectedEnvironments,
					environmentName,
				)
			}
		}
	}

	if len(affectedEnvironments) > 0 {
		ctx.err = kserrors.FileHasChanged(
			filePath,
			strings.Join(affectedEnvironments, ","),
			nil,
		)
	}

	return ctx
}
