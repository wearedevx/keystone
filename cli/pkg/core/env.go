package core

import (
	"path"

	"github.com/wearedevx/keystone/cli/internal/envfile"
	kserrors "github.com/wearedevx/keystone/cli/internal/errors"
	"github.com/wearedevx/keystone/cli/internal/keystonefile"
	"github.com/wearedevx/keystone/cli/internal/utils"
)

type (
	EnvironmentName string
	SecretValue     string
)

type Secret struct {
	Name      string
	Required  bool
	Values    map[EnvironmentName]SecretValue
	FromCache bool
}

type SecretStrictFlag int

const (
	S_REQUIRED SecretStrictFlag = iota
	S_OPTIONAL
)

// Sets an env variable to keep track of across environments
// [secretName] is the name of the variable to set
// [secretValue] maps environment to the varable value (key is environment name,
// and value, the value of the variable in that environment)
func (ctx *Context) AddSecret(
	secretName string,
	secretValue map[string]string,
	flag SecretStrictFlag,
) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	var err error
	var e *kserrors.Error
	var ksfile keystonefile.KeystoneFile
	// Add new env key to keystone.yaml
	if err = ksfile.
		Load(ctx.Wd).
		SetEnv(secretName, flag == S_REQUIRED).
		Save().
		Err(); err != nil {
		return ctx.setError(kserrors.FailedToUpdateKeystoneFile(err))
	}
	ctx.log.Printf("Added secret %s to keystonefile\n", secretName)

	// Generate .env files in cache for each environment in map
	for env, value := range secretValue {
		e = generateEnvFileInCache(ctx, env, secretName, value)
		if e != nil {
			return ctx.setError(e)
		}
		ctx.log.Printf("Cached %s for env %s as %s\n", secretName, env, value)
	}

	// Copy the new .env for the current environment to .keystone/cache/.env
	currentEnvironment := ctx.CurrentEnvironment()

	if ctx.Err() != nil {
		return ctx
	}

	newDotEnv := ctx.CachedEnvironmentDotEnvPath(currentEnvironment)
	destDotEnv := ctx.CachedDotEnvPath()

	if err = utils.CopyFile(newDotEnv, destDotEnv); err != nil {
		return ctx.setError(kserrors.CopyFailed(newDotEnv, destDotEnv, err))
	}
	ctx.log.Printf("New env file copied from %s to %s\n", newDotEnv, destDotEnv)

	return ctx
}

func generateEnvFileInCache(
	ctx *Context,
	env string,
	secretName string,
	value string,
) (e *kserrors.Error) {
	var err error

	cachePath := ctx.CachedEnvironmentPath(env)
	if !utils.DirExists(cachePath) {
		if err = utils.CreateDirIfNotExist(cachePath); err != nil {
			e = kserrors.CannotCreateDirectory(cachePath, err)

			return e
		}
	}

	envFilePath := path.Join(cachePath, ".env")

	if err = new(envfile.EnvFile).
		Load(envFilePath, nil).
		Set(secretName, value).
		Dump().
		Err(); err != nil {
		e = kserrors.FailedToUpdateDotEnv(envFilePath, err)
		return e
	}

	return nil
}

// Unsets a previously set environment variable
//
// [varname] The variable to unset
// It will be removed in all existing environment.
func (ctx *Context) RemoveSecret(secretName string, purge bool) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	var err error
	var ksfile keystonefile.KeystoneFile
	// Add new env key to keystone.yaml

	if err = ksfile.
		Load(ctx.Wd).
		UnsetEnv(secretName).
		Save().
		Err(); err != nil {
		return ctx.setError(kserrors.FailedToUpdateKeystoneFile(err))
	}
	ctx.log.Printf("Removed %s from keystonefile\n", secretName)

	if purge {
		ctx.purgeSecret(secretName)
	}

	return ctx
}

// purgeSecret removes the values associated to `secretName` from the cache
// of all environments.
// This implies that subsequently sending the environment to other users
// will remove those values for them aswell
func (ctx *Context) purgeSecret(secretName string) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	var err error
	var e *kserrors.Error

	// Update environments' .env files
	environments := ctx.ListEnvironments()

	for _, environment := range environments {
		dir := ctx.CachedEnvironmentPath(environment)
		dotEnvPath := path.Join(dir, ".env")
		dotEnv := new(envfile.EnvFile)

		if err = dotEnv.Load(dotEnvPath, nil).Err(); err != nil {
			return ctx.setError(kserrors.FailedToReadDotEnv(dotEnvPath, err))
		}

		if err = dotEnv.Unset(secretName).Dump().Err(); err != nil {
			return ctx.setError(kserrors.FailedToUpdateDotEnv(dotEnvPath, err))
		}
		ctx.log.Printf(
			"Removed secret value for %s in %s environment\n",
			secretName,
			environment,
		)
	}

	// Copy the new .env for the current environment to .keystone/cache/.env
	currentEnvironment := ctx.CurrentEnvironment()

	if e != nil {
		return ctx.setError(e)
	}

	newDotEnv := ctx.CachedEnvironmentDotEnvPath(currentEnvironment)
	destDotEnv := ctx.CachedDotEnvPath()

	if err = utils.CopyFile(newDotEnv, destDotEnv); err != nil {
		return ctx.setError(kserrors.CopyFailed(newDotEnv, destDotEnv, err))
	}

	return ctx
}

// PurgeSecets Removes from the cache of all environments all secrets that
// are not found in the project’s keystone.yaml
// This implies that sending the environment to other users will remove
// those values for them too
func (ctx *Context) PurgeSecrets() *Context {
	if ctx.Err() != nil {
		return ctx
	}

	var err error
	var e *kserrors.Error
	var ksfile = keystonefile.LoadKeystoneFile(ctx.Wd)
	// Add new env key to keystone.yaml

	if err = ksfile.Err(); err != nil {
		panic(err)
		return ctx.setError(kserrors.FailedToReadKeystoneFile(ksfile.Path, err))
	}

	// Update environments' .env files
	environments := ctx.ListEnvironments()

	for _, environment := range environments {
		dir := ctx.CachedEnvironmentPath(environment)
		dotEnvPath := path.Join(dir, ".env")
		dotEnv := new(envfile.EnvFile)

		if err = dotEnv.Load(dotEnvPath, nil).Err(); err != nil {
			return ctx.setError(kserrors.FailedToReadDotEnv(dotEnvPath, err))
		}

		for secretName := range dotEnv.GetData() {
			if yes, _ := ksfile.HasEnv(secretName); !yes {
				dotEnv.Unset(secretName)
			}
			ctx.log.Printf(
				"Removed secret value for %s in %s environment\n",
				secretName,
				environment,
			)
		}

		if err = dotEnv.Dump().Err(); err != nil {
			return ctx.setError(kserrors.FailedToUpdateDotEnv(dotEnvPath, err))
		}
	}

	// Copy the new .env for the current environment to .keystone/cache/.env
	currentEnvironment := ctx.CurrentEnvironment()

	if e != nil {
		return ctx.setError(e)
	}

	newDotEnv := ctx.CachedEnvironmentDotEnvPath(currentEnvironment)
	destDotEnv := ctx.CachedDotEnvPath()

	if err = utils.CopyFile(newDotEnv, destDotEnv); err != nil {
		return ctx.setError(kserrors.CopyFailed(newDotEnv, destDotEnv, err))
	}

	return ctx
}

// Sets an existing an secret for a given envitronment
// [envName]     name of the target environment
// [secretName]
// [secretValue]
func (ctx *Context) SetSecret(
	envName string,
	secretName string,
	secretValue string,
) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	dotEnvPath := ctx.CachedEnvironmentDotEnvPath(envName)
	dotEnv := new(
		envfile.EnvFile,
	).Load(ctx.CachedEnvironmentDotEnvPath(envName), nil)

	if err := dotEnv.Err(); err != nil {
		return ctx.setError(kserrors.FailedToReadDotEnv(dotEnvPath, err))
	}

	dotEnv.Set(secretName, secretValue).Dump()

	if err := dotEnv.Err(); err != nil {
		return ctx.setError(kserrors.FailedToUpdateDotEnv(dotEnvPath, err))
	}

	ctx.log.Printf("Set secret %s to %s in env %s\n", secretName, secretValue, envName)

	return ctx
}

// UnsetSecret method clears the secret value
func (ctx *Context) UnsetSecret(envName string, secretName string) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	return ctx.SetSecret(envName, secretName, "")
}

// Returns all environments variable
// for the current environment, a map
func (ctx *Context) GetSecrets() map[string]string {
	if ctx.Err() != nil {
		return map[string]string{}
	}

	var err error
	var env map[string]string

	dotEnv := new(envfile.EnvFile).Load(ctx.CachedDotEnvPath(), nil)

	if err = dotEnv.Err(); err != nil {
		ctx.setError(kserrors.FailedToUpdateDotEnv(ctx.CachedDotEnvPath(), err))

		return env
	}

	env = dotEnv.GetData()

	// Allow overring values with a local .env file
	// at the root of the project
	localDotEnvPath := path.Join(ctx.Wd, ".env")
	if utils.FileExists(localDotEnvPath) {
		localDotEnv := new(envfile.EnvFile).Load(localDotEnvPath, nil).GetData()

		for key, value := range localDotEnv {
			env[key] = value
		}
	}

	return env
}

// Returns a secret value for everu environments
func (ctx *Context) GetSecret(secretName string) *Secret {
	secret := new(Secret)

	if ctx.Err() != nil {
		return secret
	}

	var err error
	ksfile := keystonefile.LoadKeystoneFile(ctx.Wd)

	if err = ksfile.Err(); err != nil {
		panic(err)
		ctx.setError(kserrors.FailedToReadKeystoneFile(ksfile.Path, err))
		return secret
	}

	environmentValuesMap := map[string]map[string]string{}
	for _, environment := range ctx.ListEnvironments() {
		dotEnvPath := ctx.CachedEnvironmentDotEnvPath(environment)
		dotEnv := new(envfile.EnvFile).Load(dotEnvPath, nil)

		environmentValuesMap[environment] = dotEnv.GetData()
	}

	for _, envKey := range ksfile.Env {
		name := envKey.Key

		if name == secretName {
			required := envKey.Strict
			values := map[EnvironmentName]SecretValue{}

			for environment, secrets := range environmentValuesMap {
				values[EnvironmentName(environment)] = SecretValue(
					secrets[name],
				)
			}

			secret.Name = name
			secret.Required = required
			secret.Values = values

			break
		}
	}

	return secret
}

// List secret from .keystone/cache, and their value in each environment.
func (ctx *Context) ListSecretsFromCache() []Secret {
	secrets := make([]Secret, 0)

	if ctx.Err() != nil {
		return secrets
	}

	var err error
	ksfile := keystonefile.LoadKeystoneFile(ctx.Wd)

	if err = ksfile.Err(); err != nil {
		panic(err)
		ctx.setError(kserrors.FailedToReadKeystoneFile(ksfile.Path, err))
		return secrets
	}

	environmentValuesMap := map[string]map[string]string{}
	allSecrets := make([]string, 0)

	for _, environment := range ctx.ListEnvironments() {
		dotEnvPath := ctx.CachedEnvironmentDotEnvPath(environment)
		dotEnv := new(envfile.EnvFile).Load(dotEnvPath, nil)

		environmentValuesMap[environment] = dotEnv.GetData()
		for label := range dotEnv.GetData() {
			allSecrets = append(allSecrets, label)
		}
	}

	allSecrets = Uniq(allSecrets)

	for _, envKey := range allSecrets {
		name := envKey
		values := map[EnvironmentName]SecretValue{}

		for environment, secrets := range environmentValuesMap {
			values[EnvironmentName(environment)] = SecretValue(secrets[name])
		}

		secrets = append(secrets, Secret{
			Name:   name,
			Values: values,
		})
	}

	return secrets
}

// Returns secrets from keystone.yaml, and their value in each environment.
func (ctx *Context) ListSecrets() []Secret {
	secrets := make([]Secret, 0)

	if ctx.Err() != nil {
		return secrets
	}

	var err error
	ksfile := keystonefile.LoadKeystoneFile(ctx.Wd)

	if err = ksfile.Err(); err != nil {
		panic(err)
		ctx.setError(
			kserrors.FailedToReadKeystoneFile(
				ksfile.Path,
				err,
			),
		)
		return secrets
	}

	environmentValuesMap := map[string]map[string]string{}
	for _, environment := range ctx.ListEnvironments() {
		dotEnvPath := ctx.CachedEnvironmentDotEnvPath(environment)
		dotEnv := new(envfile.EnvFile).Load(dotEnvPath, nil)

		environmentValuesMap[environment] = dotEnv.GetData()
	}

	for _, envKey := range ksfile.Env {
		name := envKey.Key
		required := envKey.Strict
		values := map[EnvironmentName]SecretValue{}

		for environment, secrets := range environmentValuesMap {
			values[EnvironmentName(environment)] = SecretValue(secrets[name])
		}

		secrets = append(secrets, Secret{
			Name:     name,
			Required: required,
			Values:   values,
		})
	}

	return secrets
}

// Returns a boolean indicating wether the secret `secretName`
// exists in the local files
func (ctx *Context) HasSecret(secretName string) bool {
	haveIt := false

	if ctx.Err() != nil {
		return haveIt
	}

	ksfile := keystonefile.LoadKeystoneFile(ctx.Wd)
	if err := ksfile.Err(); err != nil {
		panic(err)
		ctx.setError(kserrors.FailedToReadKeystoneFile(ksfile.Path, err))
		return haveIt
	}

	for _, envKey := range ksfile.Env {
		if envKey.Key == secretName {
			haveIt = true
			break
		}
	}

	return haveIt
}

// MissingSecretsForEnvironment method returns a list of all the secrets
// that are listed in keystone.yaml and are not in the cache.
// The second value is true if there are missing things.
func (ctx *Context) MissingSecretsForEnvironment(
	environmentName string,
) ([]string, bool) {
	missing := []string{}
	hasMissing := false
	if ctx.Err() != nil {
		return missing, hasMissing
	}

	secrets := ctx.ListSecrets()

	for _, secret := range secrets {
		if secret.Required {
			value, ok := secret.Values[EnvironmentName(environmentName)]

			if !ok || value == "" {
				missing = append(missing, secret.Name)
				hasMissing = true
			}
		}
	}

	return missing, hasMissing
}

// SecretIsRequired method tells if `secretName` is required.
func (ctx *Context) SecretIsRequired(secretName string) bool {
	required := false
	if ctx.Err() != nil {
		return required
	}

	ksfile := keystonefile.LoadKeystoneFile(ctx.Wd)
	if err := ksfile.Err(); err != nil {
		panic(err)
		ctx.setError(kserrors.FailedToReadKeystoneFile(ksfile.Path, err))
		return required
	}

	for _, envKey := range ksfile.Env {
		if envKey.Key == secretName {
			required = envKey.Strict
			break
		}
	}

	return required
}

// MarkSecretRequired method changes the required status of a secret
func (ctx *Context) MarkSecretRequired(
	secretName string,
	required bool,
) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	if err := keystonefile.LoadKeystoneFile(ctx.Wd).
		SetEnv(secretName, required).
		Save().
		Err(); err != nil {
		return ctx.setError(kserrors.FailedToUpdateKeystoneFile(err))
	}

	return ctx
}

// Returns an array of secrets that are in the first list, an not in the second
func FilterSecretsFromCache(
	secretsFromCache []Secret,
	secrets []Secret,
) []Secret {
	secretsFromCacheToDisplay := make([]Secret, 0)

	for _, secretFromCache := range secretsFromCache {
		found := false

		for _, secret := range secrets {
			if secret.Name == secretFromCache.Name {
				found = true
			}
		}

		if !found {
			secretFromCache.FromCache = true
			secretsFromCacheToDisplay = append(
				secretsFromCacheToDisplay,
				secretFromCache,
			)
		}
	}
	return secretsFromCacheToDisplay
}
