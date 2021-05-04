package errors

var helpTexts map[string]string = map[string]string{
	"InitFailed": `
{{ ERROR }} {{ .Name | red }}
This happened because: {{ .Cause }}

`,
	"NotAKeystoneProject": `
{{ ERROR }} {{ .Name | red }}
It seems you are not in a keystone project.

Neither the current directory ({{ .Path }}), nor any of its parent,
have a keystone.yaml file.
If this is a new project, start with:
  $ ks init <your-project-name>

`,
	"NoWorkingDirectory": `
{{ ERROR }} {{ .Name | red }}
A current working directory could not be determined.

This happened because: {{ .Cause }}

`,
	"UnsupportedFlag": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .Flag | red }} {{- "'" | red }}

`,
	"FailedToReadKeystoneFile": `
{{ ERROR }} {{ .Name | red }}
The keystone.yml file exists, but it might not be readable or writable.
Its content may also be corrupted and may not be parsable.

This happened because: {{ .Cause }}

`,
	"FailedToUpdateKeystoneFile": `
{{ ERROR }} {{ .Name | red }}
The keystone.yml file exists, but it might not be readable or writable.
Its content may also be corrupted and may not be parsable.

This happened because: {{ .Cause }}

`,
	"FailedToUpdateDotEnv": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .Path | red }} {{- "'" | red }}
This happened because: {{ .Cause }}

`,
	"FailedToReadDotEnv": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .Path | red }} {{- "'" | red }}
This happened because: {{ .Cause }}

`,
	"FailedToReadRolesFile": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .Path | red }} {{- "'" | red }}
This happened because: {{ .Cause }}

`,
	"RoleDoesNotExist": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .RoleName | red }} {{- "'" | red }}
Available roles are: {{ .Available }}

You can manage roles for the current project by editing the roles file:
  .keystone/roles.yml

`,
	"EnvironmentDoesntExist": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .Environment | red }} {{- "'" | red }}
Available environments are: {{ .Available }}

To create a new environment:
  $ ks env new {{ .Environment }}

`,
	"EnvironmentAlreadyExists": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .Environment | red }} {{- "'" | red }}
You tried to create an environment with the name '{{ .Environment }}',
but your project already have one with that name.

To use the '{{ .Environment }}':
  $ ks env {{ .Environment }}

`,
	"FailedToSetCurrentEnvironment": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .Environment | red }} {{- "'" | red }}
The file at '{{ .Path }}' could not be written.

This happened because: {{ .Cause }}

`,
	"CannotReadEnvironment": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .Path | red }} {{- "'" | red }}
This happened because: {{ .Cause }}

`,
	"CannotRemoveCurrentEnvironment": `
{{ ERROR }} {{ .Name | red }}
You are trying to remove the '{{ .Environment }}' environment,
but it is currently in use.

Change to another environment:
  $ ks env default

`,
	"SecretDoesNotExist": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .Secret | red }} {{- "'" | red }}

To list secrets:
  $ ks secret

To add a {{ .Secret }} secret to all environments:
  $ ks secret {{ .Secret }} <secret-value>

`,
	"SecretRequired": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .Secret | red }} {{- "'" | red }}
You are trying to either unset '{{ .Secret }}', or to set it to a blank value,
but is required.

`,
	"CannotAddFile": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .Path | red }} {{- "'" | red }}
This happened because: {{ .Cause }}

`,
	"CannotRemoveFile": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .Path | red }} {{- "'" | red }}
This happened because: {{ .Cause }}

`,
	"CannotLinkFile": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .Path | red }} {{- "'" | red }}
The symlink to {{ .CachePath }} could not be created at {{ .Path }}.

This happened because: {{ .Cause }}

`,
	"FileNotInEnvironment": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .Path | red }} {{- "'" | red }}
No version of '{{ .Path }}' was found for the '{{ .Environment }}' environment.

To get the latest variables and files for '{{ .Environment }}':
  $ ks --env {{ .Environment }} pull

`,
	"CannotCreateDirectory": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .Path | red }} {{- "'" | red }}
This happened because: {{ .Cause }}

`,
	"CannotRemoveDirectory": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .Path | red }} {{- "'" | red }}
This happened because: {{ .Cause }}

`,
	"CopyFailed": `
{{ ERROR }} {{ .Name | red }}
Trying to copy '{{ .Source }}' to '{{ .Destination }}'

This happened because: {{ .Cause }}

`,
	"MustBeLoggedIn": `
{{ ERROR }} {{ .Name | red }}
You must be logged to execute this command.

Please run:
  $ ks login

`,
	"UnkownError": `
{{ ERROR }} {{ .Name | red }}
Ouch! We didn't think this could happen.

Feel free to file an issue at https://github.com/wearedevx/keystone
Sorry for the inconvenience

This happened because: {{ .Cause }}

`,
	"UsersDontExist": `
{{ ERROR }} {{ .Name | red }}
{{ .Message }}

You can invite those users to Keystone using
  $ ks invite <emai>

`,
	"CannotAddMembers": `
{{ ERROR }} {{ .Name | red }}
This happened because: {{ .Cause }}

`,
	"CannotRemoveMembers": `
{{ ERROR }} {{ .Name | red }}
This happened because: {{ .Cause }}

`,
}

func InitFailed(cause error) *Error {
	meta := map[string]string{}

	return NewError("Init Failed", helpTexts["InitFailed"], meta, cause)
}

func NotAKeystoneProject(path string, cause error) *Error {
	meta := map[string]string{
		"Path": string(path),
	}
	return NewError("Not A Keystone Project", helpTexts["NotAKeystoneProject"], meta, cause)
}

func NoWorkingDirectory(cause error) *Error {
	meta := map[string]string{}

	return NewError("No Working Directory", helpTexts["NoWorkingDirectory"], meta, cause)
}

func UnsupportedFlag(flag string, cause error) *Error {
	meta := map[string]string{
		"Flag": string(flag),
	}
	return NewError("Unsupported Flag", helpTexts["UnsupportedFlag"], meta, cause)
}

func FailedToReadKeystoneFile(cause error) *Error {
	meta := map[string]string{}

	return NewError("Failed To Read Keystone File", helpTexts["FailedToReadKeystoneFile"], meta, cause)
}

func FailedToUpdateKeystoneFile(cause error) *Error {
	meta := map[string]string{}

	return NewError("Failed To Update Keystone File", helpTexts["FailedToUpdateKeystoneFile"], meta, cause)
}

func FailedToUpdateDotEnv(path string, cause error) *Error {
	meta := map[string]string{
		"Path": string(path),
	}
	return NewError("Failed To Update .env", helpTexts["FailedToUpdateDotEnv"], meta, cause)
}

func FailedToReadDotEnv(path string, cause error) *Error {
	meta := map[string]string{
		"Path": string(path),
	}
	return NewError("Failed To Read .env", helpTexts["FailedToReadDotEnv"], meta, cause)
}

func FailedToReadRolesFile(path string, cause error) *Error {
	meta := map[string]string{
		"Path": string(path),
	}
	return NewError("Failed To Read Roles File", helpTexts["FailedToReadRolesFile"], meta, cause)
}

func RoleDoesNotExist(rolename string, available string, cause error) *Error {
	meta := map[string]string{
		"RoleName":  string(rolename),
		"Available": string(available),
	}
	return NewError("Role Does Not Exist", helpTexts["RoleDoesNotExist"], meta, cause)
}

func EnvironmentDoesntExist(environment string, available string, cause error) *Error {
	meta := map[string]string{
		"Environment": string(environment),
		"Available":   string(available),
	}
	return NewError("Environment Doesn't Exist", helpTexts["EnvironmentDoesntExist"], meta, cause)
}

func EnvironmentAlreadyExists(environment string, cause error) *Error {
	meta := map[string]string{
		"Environment": string(environment),
	}
	return NewError("Environment Already Exists", helpTexts["EnvironmentAlreadyExists"], meta, cause)
}

func FailedToSetCurrentEnvironment(environment string, path string, cause error) *Error {
	meta := map[string]string{
		"Environment": string(environment),
		"Path":        string(path),
	}
	return NewError("Failed To Set Current Environment", helpTexts["FailedToSetCurrentEnvironment"], meta, cause)
}

func CannotReadEnvironment(path string, cause error) *Error {
	meta := map[string]string{
		"Path": string(path),
	}
	return NewError("Cannot Read Environment", helpTexts["CannotReadEnvironment"], meta, cause)
}

func CannotRemoveCurrentEnvironment(environment string, cause error) *Error {
	meta := map[string]string{
		"Environment": string(environment),
	}
	return NewError("Cannot Remove Current Environment", helpTexts["CannotRemoveCurrentEnvironment"], meta, cause)
}

func SecretDoesNotExist(secret string, cause error) *Error {
	meta := map[string]string{
		"Secret": string(secret),
	}
	return NewError("Secret Doesn't Exist", helpTexts["SecretDoesNotExist"], meta, cause)
}

func SecretRequired(secret string, cause error) *Error {
	meta := map[string]string{
		"Secret": string(secret),
	}
	return NewError("Secret Required", helpTexts["SecretRequired"], meta, cause)
}

func CannotAddFile(path string, cause error) *Error {
	meta := map[string]string{
		"Path": string(path),
	}
	return NewError("Cannot Add File", helpTexts["CannotAddFile"], meta, cause)
}

func CannotRemoveFile(path string, cause error) *Error {
	meta := map[string]string{
		"Path": string(path),
	}
	return NewError("Cannot Remove File", helpTexts["CannotRemoveFile"], meta, cause)
}

func CannotLinkFile(path string, cachepath string, cause error) *Error {
	meta := map[string]string{
		"Path":      string(path),
		"CachePath": string(cachepath),
	}
	return NewError("Cannot Link File", helpTexts["CannotLinkFile"], meta, cause)
}

func FileNotInEnvironment(path string, environment string, cause error) *Error {
	meta := map[string]string{
		"Path":        string(path),
		"Environment": string(environment),
	}
	return NewError("File Not Found For Environment", helpTexts["FileNotInEnvironment"], meta, cause)
}

func CannotCreateDirectory(path string, cause error) *Error {
	meta := map[string]string{
		"Path": string(path),
	}
	return NewError("Cannot Create Directory", helpTexts["CannotCreateDirectory"], meta, cause)
}

func CannotRemoveDirectory(path string, cause error) *Error {
	meta := map[string]string{
		"Path": string(path),
	}
	return NewError("Cannot Remove Directory", helpTexts["CannotRemoveDirectory"], meta, cause)
}

func CopyFailed(source string, destination string, cause error) *Error {
	meta := map[string]string{
		"Source":      string(source),
		"Destination": string(destination),
	}
	return NewError("Copy failed", helpTexts["CopyFailed"], meta, cause)
}

func MustBeLoggedIn(cause error) *Error {
	meta := map[string]string{}

	return NewError("You must be logged in", helpTexts["MustBeLoggedIn"], meta, cause)
}

func UnkownError(cause error) *Error {
	meta := map[string]string{}

	return NewError("Unkown Error", helpTexts["UnkownError"], meta, cause)
}

func UsersDontExist(message string, cause error) *Error {
	meta := map[string]string{
		"Message": string(message),
	}
	return NewError("Users Don't Exist", helpTexts["UsersDontExist"], meta, cause)
}

func CannotAddMembers(cause error) *Error {
	meta := map[string]string{}

	return NewError("Cannot Add Members", helpTexts["CannotAddMembers"], meta, cause)
}

func CannotRemoveMembers(cause error) *Error {
	meta := map[string]string{}

	return NewError("Cannot Remove Members", helpTexts["CannotRemoveMembers"], meta, cause)
}
