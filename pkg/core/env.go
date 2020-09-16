package core

import (
	"path"
	"strings"

	. "github.com/wearedevx/keystone/internal/envfile"
	. "github.com/wearedevx/keystone/internal/errors"
	. "github.com/wearedevx/keystone/internal/keystonefile"
	. "github.com/wearedevx/keystone/internal/utils"
)

type EnvironmentName string
type SecretValue string

type Secret struct {
	Name     string
	Required bool
	Values   map[EnvironmentName]SecretValue
}

type SecretStrictFlag int

const (
	S_REQUIRED SecretStrictFlag = iota
	S_OPTIONAL
)

// Sets an env variable to keep track of across environments
// [varname] is the name of the variable to set
// [varvalue] maps environment to the varable value (key is environment name,
//   and value, the value of the variable in that environment)
// TODO: Factorize this plz
func (ctx *Context) AddSecret(secretName string, secretValue map[string]string, flag SecretStrictFlag) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	var err error
	var e *Error
	var ksfile KeystoneFile
	// Add new env key to keystone.yml
	if err = ksfile.Load(ctx.Wd).SetEnv(secretName, flag == S_REQUIRED).Save().Err(); err != nil {
		return ctx.setError(FailedToUpdateKeystoneFile(err))
	}

	// Generate .env files in cache for each environment in map
	for env, value := range secretValue {
		cachePath := path.Join(ctx.cacheDirPath(), env)
		if !DirExists(cachePath) {
			if env != "default" {

				e = EnvironmentDoesntExist(env, strings.Join(ctx.ListEnvironments(), ", "), nil)

				break

			} else {
				if err = CreateDirIfNotExist(cachePath); err != nil {
					e = CannotCreateDirectory(cachePath, err)

					break
				}
			}
		}

		envFilePath := path.Join(cachePath, ".env")

		if err = new(EnvFile).Load(envFilePath).Set(secretName, value).Dump().Err(); err != nil {
			e = FailedToUpdateDotEnv(envFilePath, err)
			break
		}
	}

	if e != nil {
		return ctx.setError(e)
	}

	// Copy the new .env for the current environment to .keystone/cache/.env
	currentEnvironment := ctx.CurrentEnvironment()

	if ctx.Err() != nil {
		return ctx
	}

	newDotEnv := path.Join(ctx.cacheDirPath(), currentEnvironment, ".env")
	destDotEnv := ctx.CachedDotEnvPath()

	if err = CopyFile(newDotEnv, destDotEnv); err != nil {
		return ctx.setError(CopyFailed(newDotEnv, destDotEnv, err))
	}

	return ctx
}

// Unsets a previously set environment variable
//
// [varname] The variable to unset
// It will be removed in all existing environment.
func (ctx *Context) RemoveSecret(secretName string) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	var err error
	var e *Error
	var ksfile KeystoneFile
	// Add new env key to keystone.yml

	if err = ksfile.Load(ctx.Wd).UnsetEnv(secretName).Save().Err(); err != nil {
		return ctx.setError(FailedToUpdateKeystoneFile(err))
	}

	// Update environments' .env files
	environments := ctx.ListEnvironments()

	for _, environment := range environments {
		dir := path.Join(ctx.cacheDirPath(), environment)
		dotEnvPath := path.Join(dir, ".env")

		if err = new(EnvFile).Load(dotEnvPath).Unset(secretName).Dump().Err(); err != nil {
			e = FailedToUpdateDotEnv(dotEnvPath, err)
			break
		}
	}

	// Copy the new .env for the current environment to .keystone/cache/.env
	currentEnvironment := ctx.CurrentEnvironment()

	if e != nil {
		return ctx.setError(e)
	}

	newDotEnv := path.Join(ctx.cacheDirPath(), currentEnvironment, ".env")
	destDotEnv := ctx.CachedDotEnvPath()

	if err = CopyFile(newDotEnv, destDotEnv); err != nil {
		return ctx.setError(CopyFailed(newDotEnv, destDotEnv, err))
	}

	return ctx
}

func (ctx *Context) SetSecret(envName string, secretName string, secretValue string) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	dotEnvPath := path.Join(ctx.cacheDirPath(), envName, ".env")
	dotEnv := new(EnvFile).Load(path.Join(ctx.cacheDirPath(), envName, ".env"))

	if err := dotEnv.Err(); err != nil {
		return ctx.setError(FailedToReadDotEnv(dotEnvPath, err))
	}

	dotEnv.Set(secretName, secretValue).Dump()

	if err := dotEnv.Err(); err != nil {
		return ctx.setError(FailedToUpdateDotEnv(dotEnvPath, err))
	}

	return ctx
}

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

	dotEnv := new(EnvFile).Load(ctx.CachedDotEnvPath())

	if err = dotEnv.Err(); err != nil {
		ctx.setError(FailedToUpdateDotEnv(ctx.CachedDotEnvPath(), err))

		return env
	}

	env = dotEnv.GetData()

	// Allow overring values with a local .env file
	// at the root of the project
	localDotEnvPath := path.Join(ctx.Wd, ".env")
	if FileExists(localDotEnvPath) {
		localDotEnv := new(EnvFile).Load(localDotEnvPath).GetData()

		for key, value := range localDotEnv {
			env[key] = value
		}
	}

	return env
}

func (ctx *Context) GetSecret(secretName string) *Secret {
	secret := new(Secret)

	if ctx.Err() != nil {
		return secret
	}

	var err error
	ksfile := new(KeystoneFile).Load(ctx.Wd)

	if err = ksfile.Err(); err != nil {
		ctx.setError(FailedToReadKeystoneFile(err))
		return secret
	}

	environmentValuesMap := map[string]map[string]string{}
	for _, environment := range ctx.ListEnvironments() {
		dotEnvPath := path.Join(ctx.cacheDirPath(), environment, ".env")
		dotEnv := new(EnvFile).Load(dotEnvPath)

		environmentValuesMap[environment] = dotEnv.GetData()
	}

	for _, envKey := range ksfile.Env {
		name := envKey.Key

		if name == secretName {
			required := envKey.Strict
			values := map[EnvironmentName]SecretValue{}

			for environment, secrets := range environmentValuesMap {
				values[EnvironmentName(environment)] = SecretValue(secrets[name])
			}

			secret.Name = name
			secret.Required = required
			secret.Values = values

			break
		}
	}

	return secret
}

func (ctx *Context) ListSecrets() []Secret {
	secrets := make([]Secret, 0)

	if ctx.Err() != nil {
		return secrets
	}

	var err error
	ksfile := new(KeystoneFile).Load(ctx.Wd)

	if err = ksfile.Err(); err != nil {
		ctx.setError(FailedToReadKeystoneFile(err))
		return secrets
	}

	environmentValuesMap := map[string]map[string]string{}
	for _, environment := range ctx.ListEnvironments() {
		dotEnvPath := path.Join(ctx.cacheDirPath(), environment, ".env")
		dotEnv := new(EnvFile).Load(dotEnvPath)

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

func (ctx *Context) HasSecret(secretName string) bool {
	haveIt := false

	if ctx.Err() != nil {
		return haveIt
	}

	ksfile := new(KeystoneFile).Load(ctx.Wd)
	if err := ksfile.Err(); err != nil {
		ctx.setError(FailedToReadKeystoneFile(err))
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

func (ctx *Context) SecretIsRequired(secretName string) bool {
	required := false
	if ctx.Err() != nil {
		return required
	}

	ksfile := new(KeystoneFile).Load(ctx.Wd)
	if err := ksfile.Err(); err != nil {
		ctx.setError(FailedToReadKeystoneFile(err))
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

func (ctx *Context) MarkSecretRequired(secretName string, required bool) *Context {
	if ctx.Err() != nil {
		return ctx
	}

	if err := new(KeystoneFile).Load(ctx.Wd).SetEnv(secretName, required).Save().Err(); err != nil {
		return ctx.setError(FailedToUpdateKeystoneFile(err))
	}

	return ctx
}
