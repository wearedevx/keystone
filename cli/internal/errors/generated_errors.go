package errors

var helpTexts map[string]string = map[string]string{
	"InitFailed": `
{{ ERROR }} {{ .Name | red }}
This happened because: {{ .Cause }}
`,
	"ServiceNotAvailable": `
{{ ERROR }} {{ .Name | red }}
`,
	"InvalidConnectionToken": `
{{ ERROR }} {{ .Name | red }}
Your current connection token has probably expired.

Try to login again:
  $ ks login
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
	"AlreadyKeystoneProject": `
{{ ERROR }} {{ .Name | red }}
You are trying to create a Keystone project but there already are keystone files in your current directory.
Please remove the .keystone directory and keystone.yaml file beforehand.
`,
	"DeviceNotRegistered": `
{{ ERROR }} {{ .Name | red }}
This device is not registered to your account or its access has been revoked.
To register it, please logout, then login again.
`,
	"BadDeviceName": `
{{ ERROR }} {{ .Name | red }}
Device names must be alphanumeric with ., -, _
`,
	"CannotSaveConfig": `
{{ ERROR }} {{ .Name | red }}
You have been successfully logged in, but the configuration file could
not be written.

This happened because: {{ .Cause }}
`,
	"FailedToReadKeystoneFile": `
{{ ERROR }} {{ .Name | red }}
The keystone.yaml file exists, but it might not be readable or writable.
Its content may also be corrupted and may not be parsable.

This happened because: {{ .Cause }}
`,
	"FailedToUpdateKeystoneFile": `
{{ ERROR }} {{ .Name | red }}
The keystone.yaml file exists, but it might not be readable or writable.
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
	"RoleDoesNotExist": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .RoleName | red }} {{- "'" | red }}
`,
	"RoleNeedsUpgrade": `
{{ ERROR }} {{ .Name | red }}
You are not allowed to set roles other than admin for a free organization.

To upgrade your plan:
  $ ks orga upgrade
`,
	"ProjectDoesntExist": `
{{ ERROR }} {{ "Project" | red }} {{ .ProjectName | red }} {{ "Does Not Exist" | red }}
Project in your keystone.yaml does not exist or your are not part of it.

If you have this configuration from a project member, ask them to add you in the keystone project.
`,
	"OrganizationNotPaid": `
{{ ERROR }} {{ .Name | red }} 
The project belongs to an organization on a free plan
and some members of the project are not admin.
Roles feature is only available for paid organizations,
this prevents you from sharing secrets with each other.

To unlock the situation, you can set the role to admin
for each member of the project:
  $ ks member set-role <member id>

Or, upgrade your organization plan:
  $ ks orga upgrade
`,
	"NameDoesNotMatch": `
{{ ERROR }} {{ .Name | red }}
`,
	"CouldNotRemoveLocalFiles": `
{{ ERROR }} {{ .Name | red }}
The project has been destroyed on our servers,
but some locale files could not be removed.

Check file system permissions and remove the following files manually:
  - keystone.yaml
  - .keystone/

This happened because: {{ .Cause }}
`,
	"EnvironmentDoesntExist": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .Environment | red }} {{- "'" | red }}
Available environments are: {{ .Available }}

To use another environment:
  $ ks env switch {{ .Environment }}
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
	"PermissionDenied": `
{{ ERROR }} {{ .Name | red }} 
You do not have the rights to change the '{{ .Environment }}' environment.
`,
	"CannotGetEnvironmentKeys": `
{{ ERROR }} {{ .Name | red }}
Public keys for the '{{ .Environment }}' could not be retrieved.

This happened because: {{ .Cause }}
`,
	"YouHaveLocallyModifiedFiles": `
{{ ERROR }} {{ .Name | red }}
{{ range $file := .Files }}  - {{ $file }}
{{ end }}

If you want to make those changes permanent for the '{{ .Environment }}'
and send them all members:
  $ ks file set <filepath>

If you want to discard those changes:
  $ ks file reset [filepath]...
`,
	"SecretDoesNotExist": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .Secret | red }} {{- "'" | red }}

To list secrets:
  $ ks secret

To add a {{ .Secret }} secret to all environments:
  $ ks secret add {{ .Secret }} <secret-value>
`,
	"SecretRequired": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .Secret | red }} {{- "'" | red }}
You are trying to either unset '{{ .Secret }}', or to set it to a blank value,
but is required.

If you are sure this secret can be unset or blank, you can mark as optionl with:
  $ ks secret optional {{ .Secret }}
`,
	"SecretHasChanged": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .Secret | red }} {{- "'" | red }}
You are trying to set a value for '{{ .Secret }}', but a new value has been set by another member.
If you want to override their value, try again.

{{ .Values }}
`,
	"RequiredSecretsAreMissing": `
{{ ERROR }} {{ .Name | red }} 
You are trying to send the environment '{{ .EnvironmentName }}' to a CI/CD service, but some required secrets are missing their value:
{{ range $secretName := .MissingSecrets }}
  - {{ $secretName }}
{{ end }}

You may set value for those secrets with:
  $ ks --env {{ .EnvironmentName }} secret set <SECRET_NAME> <SECRET_VALUE>

Or make them optional using:
  $ ks secret optional <SECRET_NAME>
`,
	"FileDoesNotExist": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .FileName | red }} {{- "'" | red }}

To list files:
  $ ks file

To add {{ .FileName }} to all environments:
  $ ks file add {{ .FileName }}
`,
	"RequiredFilesAreMissing": `
{{ ERROR }} {{ .Name | red }}
You are trying to send the environment '{{ .EnvironmentName }}' to a CI/CD service, but some required files are missing or empty:
{{ range $filePath := .MissingFiles }}
 - {{ $filePath }}
{{ end }}

You may set the content for those files with:
  $ ks file add <FILE_PATH>

Or make them optional using:
  $ ks file optional <FILE_PATH>
`,
	"FileNotInWorkingDirectory": `
{{ ERROR }} {{ .Name | red }}
The file you are trying to add ({{ .FilePath }}) does not belong
to the project's current working directory :
    {{ .Wd }}

Only files belonging to {{ .Wd }} or its subdirectories can be added.
`,
	"EnvironmentsHaveChanged": `
{{ ERROR }} {{ .Name | red }}
We couldn't find data for the following environments: '{{ .EnvironmentsName }}',
but a new value has been set by another member.
Ask someone to use 'ks env send' to make newa data available to you.
`,
	"FileHasChanged": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{ .FilePath | red }}
You are trying to update the file '{{ .FilePath }}',
but another member has changed its content.
If you want to override their changes, try again.

Affected environments: {{ .AffectedEnvironments }}
`,
	"CannotAddFile": `
{{ ERROR }} {{ .Name | red }} {{- ": \"" | red }} {{- .Path | red }} {{- "\"" | red }}
This happened because: {{ .Cause }}
`,
	"CannotSetFile": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .Path | red }} {{- "'" | red }}
You tried to set the content of '{{ .Path }}', but it could not be read or found.
Make sure the file exists, and has been added to the project using:
  $ ks file add {{ .Path }}

This happened because: {{ .Cause }}
`,
	"CannotRemoveFile": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .Path | red }} {{- "'" | red }}
This happened because: {{ .Cause }}
`,
	"CannotCopyFile": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .Path | red }} {{- "'" | red }}
A copy to {{ .CachePath }} could not be created at {{ .Path }}.

This happened because: {{ .Cause }}
`,
	"FileNotInEnvironment": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .Path | red }} {{- "'" | red }}
No version of '{{ .Path }}' was found for the '{{ .Environment }}' environment.

To get the latest variables and files for '{{ .Environment }}':
  $ ks file && ks secret
`,
	"CannotCreateDirectory": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .Path | red }} {{- "'" | red }}
This happened because: {{ .Cause }}
`,
	"CannotRemoveDirectoryContents": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .Path | red }} {{- "'" | red }}
This happened because: {{ .Cause }}
`,
	"CannotSaveFiles": `
{{ ERROR }} {{ .Name | red }} {{- ": '" | red }} {{- .FileList | red }} {{- "'" | red }}
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
	"CannotFindProjectID": `
{{ ERROR }} {{ .Name | red }}
Keystone.yaml may be malformated
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
  $ ks invite <email>
`,
	"CannotAddMembers": `
{{ ERROR }} {{ .Name | red }}

This happened because: {{ .Cause }}
`,
	"CannotRemoveMembers": `
{{ ERROR }} {{ .Name | red }}

This happened because: {{ .Cause }}
`,
	"MemberHasNoAccessToEnv": `
{{ ERROR }} {{ .Name | red }}

This happend because: {{ .Cause }}
`,
	"CouldNotDecryptMessages": `
{{ ERROR }} {{ .Name | red }}
{{ .Message }}

This happened because: {{ .Cause }}
`,
	"CouldNotEncryptMessages": `
{{ ERROR }} {{ .Name | red }}

This happened because: {{ .Cause }}
`,
	"EncryptionFailed": `
{{ ERROR }} {{ .Name | red }}

This happened because: {{ .Cause }}
`,
	"CouldNotParseMessage": `
{{ ERROR }} {{ .Name | red }}

This happened because: {{ .Cause }}
`,
	"PayloadErrors": `
{{ ERROR }} {{ .Name | red }}

This happened because: {{ .Cause }}
`,
	"InvalidFileContent": `
{{ ERROR }} {{ .Name | red }}
The file '{{ .Path }}' content could not be decoded from a base64 string

This happened because: {{ .Cause }}
`,
	"FailedCheckingChanges": `
{{ ERROR }} {{ .Name | red }}
Error were encountered for file: '{{ .Path }}'

This happened because: {{ .Cause }}
`,
	"FeatureRequiresToUpgrade": `
{{ ERROR }} {{ .Name | red }}
To take full advantage of all the features Keystone has to offer,
such as roles and logs, you should upgrade your plan, using:
  $ ks orga upgrade
`,
	"AlreadySubscribed": `
{{ ERROR }} {{ .Name | red }}
You already have access to all of Keystone features!
`,
	"CannotUpgrade": `
{{ ERROR }} {{ .Name | red }}
Communication with the billing service failed,
try again later.
`,
	"ManagementInaccessible": `
{{ ERROR }} {{ .Name | red }}
Communication with the billing service failed,
try again later.
`,
	"OrganizationDoesNotExist": `
{{ ERROR }} {{ .Name | red }}
To create a new organization, use:
  $ ks orga add <ORGANIZATION_NAME>

Orgnaization names must be unique
`,
	"BadOrganizationName": `
{{ ERROR }} {{ .Name | red }}
Organization names must be alphanumeric with ., -, _
`,
	"OrganizationNameAlreadyTaken": `

`,
	"MustOwnTheOrganization": `
{{ ERROR }} {{ .Name | red }}
You tried to perform an operation that requires ownership
of the organization.

You should ask the owner of said organization to perform it for you.
`,
	"YouDoNotOwnTheOrganization": `
{{ ERROR }} {{ .Name | red }} {{ .OrganizationName | red }}
To see organization you own, use:
  $ ks orga

To create a new organization, use:
  $ ks orga add {{ .OrganizationName }}
`,
	"CouldntSendInvite": `
{{ ERROR }} {{ .Name | red }}

This happened because: {{ .Cause }}
`,
	"CouldntSetRole": `
{{ ERROR }} {{ .Name | red }}

This happened because: {{ .Cause }}
`,
	"BackupDenied": `
{{ ERROR }} {{ .Name | red }}
You are not allowed to create backups.
`,
	"RestoreDenied": `
{{ ERROR }} {{ .Name | red }}
You are not allowed to restore backups.
`,
	"CouldNotCreateArchive": `
{{ ERROR }} {{ .Name | red }}

This happened because: {{ .Cause }}
`,
	"FailedToReadBackup": `
{{ ERROR }} {{ .Name | red }}

This happened because: {{ .Cause }}
`,
	"FailedToWriteBackup": `
{{ ERROR }} {{ .Name | red }}

This happened because: {{ .Cause }}
`,
	"BackupNotSetUp": `
{{ ERROR }} {{ .Name | red }}

This functionnality requires backups to have been setup.
You can do it using:
  $ ks backup --setup
`,
	"NoBackup": `
{{ ERROR }} {{ .Name | red }}
No backup could be found for the project "{{ .ProjectName }}" in
"{{ .BackupDirPath }}".

You can create one manually with:
  $ ks backup

Or change the backup directory with:
  $ ks backup --setup
`,
	"RestoreFailed": `
{{ ERROR }} {{ .Name | red }}

This happened because: {{ .Cause }}
`,
	"BackupFailed": `
{{ ERROR }} {{ .Name | red }}

This happened because: {{ .Cause }}
`,
	"NoCIServices": `
{{ ERROR }} {{ .Name | red }}
You haven't configured any CI services yet.

To add a new service, try:
  $ ks ci add
`,
	"CiServiceAlreadyExists": `
{{ ERROR }} {{ .Name | red }} {{- ": " | red }} {{- .ServiceName | red }}
If you wish to modify it, try:
  $ ks ci edit {{ .ServiceName }}

If you wish to add a new service, pick another name, then
  $ ks ci add <other-service-name>
`,
	"NoSuchService": `
{{ ERROR }} {{ .Name | red }} {{ .ServiceName | red }}
`,
	"CouldNotAddService": `
{{ ERROR }} {{ .Name | red }} {{ .ServiceName | red }}

This happened because: {{ .Cause }}
`,
	"CouldNotCleanService": `
{{ ERROR }} {{ .Name | red }}

This happened because: {{ .Cause }}
`,
	"CouldNotChangeService": `
{{ ERROR }} {{ .Name | red }} {{ .ServiceName | red }}

This happened because: {{ .Cause }}
`,
	"CouldNotRemoveService": `
{{ ERROR }} {{ .Name | red }}

This happened because: {{ .Cause }}
`,
	"MissingCIInformation": `
{{ ERROR }} {{ .Name | red }}
To edit the service, try:
  $ ks ci edit {{ .ServiceName }}
`,
	"CouldNotSendToCIService": `
{{ ERROR }} {{ .Name | red }}

This happened because: {{ .Cause }}
`,
}

func InitFailed(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Init Failed", helpTexts["InitFailed"], meta, cause)
}

func ServiceNotAvailable(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Service not available", helpTexts["ServiceNotAvailable"], meta, cause)
}

func InvalidConnectionToken(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Invalid Connection Token", helpTexts["InvalidConnectionToken"], meta, cause)
}

func NotAKeystoneProject(path string, cause error) *Error {
	meta := map[string]interface{}{
		"Path": string(path),
	}
	return NewError("Not A Keystone Project", helpTexts["NotAKeystoneProject"], meta, cause)
}

func NoWorkingDirectory(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("No Working Directory", helpTexts["NoWorkingDirectory"], meta, cause)
}

func UnsupportedFlag(flag string, cause error) *Error {
	meta := map[string]interface{}{
		"Flag": string(flag),
	}
	return NewError("Unsupported Flag", helpTexts["UnsupportedFlag"], meta, cause)
}

func AlreadyKeystoneProject(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Already a Keystone project", helpTexts["AlreadyKeystoneProject"], meta, cause)
}

func DeviceNotRegistered(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Device not registered", helpTexts["DeviceNotRegistered"], meta, cause)
}

func BadDeviceName(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Bad Device Name", helpTexts["BadDeviceName"], meta, cause)
}

func CannotSaveConfig(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Cannot Save Config", helpTexts["CannotSaveConfig"], meta, cause)
}

func FailedToReadKeystoneFile(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Failed To Read Keystone File", helpTexts["FailedToReadKeystoneFile"], meta, cause)
}

func FailedToUpdateKeystoneFile(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Failed To Update Keystone File", helpTexts["FailedToUpdateKeystoneFile"], meta, cause)
}

func FailedToUpdateDotEnv(path string, cause error) *Error {
	meta := map[string]interface{}{
		"Path": string(path),
	}
	return NewError("Failed To Update .env", helpTexts["FailedToUpdateDotEnv"], meta, cause)
}

func FailedToReadDotEnv(path string, cause error) *Error {
	meta := map[string]interface{}{
		"Path": string(path),
	}
	return NewError("Failed To Read .env", helpTexts["FailedToReadDotEnv"], meta, cause)
}

func RoleDoesNotExist(rolename string, cause error) *Error {
	meta := map[string]interface{}{
		"RoleName": string(rolename),
	}
	return NewError("Role Not Available", helpTexts["RoleDoesNotExist"], meta, cause)
}

func RoleNeedsUpgrade(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Needs Upgrade", helpTexts["RoleNeedsUpgrade"], meta, cause)
}

func ProjectDoesntExist(projectname string, projectid string, cause error) *Error {
	meta := map[string]interface{}{
		"ProjectName": string(projectname),
		"ProjectId":   string(projectid),
	}
	return NewError("Project Does Not Exist", helpTexts["ProjectDoesntExist"], meta, cause)
}

func OrganizationNotPaid(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Upgrade to a Paid Plan Is Required", helpTexts["OrganizationNotPaid"], meta, cause)
}

func NameDoesNotMatch(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Name Does Not Match", helpTexts["NameDoesNotMatch"], meta, cause)
}

func CouldNotRemoveLocalFiles(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Could Not Remove Local Files", helpTexts["CouldNotRemoveLocalFiles"], meta, cause)
}

func EnvironmentDoesntExist(environment string, available string, cause error) *Error {
	meta := map[string]interface{}{
		"Environment": string(environment),
		"Available":   string(available),
	}
	return NewError("Environment Does Not Exist", helpTexts["EnvironmentDoesntExist"], meta, cause)
}

func FailedToSetCurrentEnvironment(environment string, path string, cause error) *Error {
	meta := map[string]interface{}{
		"Environment": string(environment),
		"Path":        string(path),
	}
	return NewError("Failed To Set Current Environment", helpTexts["FailedToSetCurrentEnvironment"], meta, cause)
}

func CannotReadEnvironment(path string, cause error) *Error {
	meta := map[string]interface{}{
		"Path": string(path),
	}
	return NewError("Cannot Read Environment", helpTexts["CannotReadEnvironment"], meta, cause)
}

func PermissionDenied(environment string, cause error) *Error {
	meta := map[string]interface{}{
		"Environment": string(environment),
	}
	return NewError("Permission Denied", helpTexts["PermissionDenied"], meta, cause)
}

func CannotGetEnvironmentKeys(environment string, cause error) *Error {
	meta := map[string]interface{}{
		"Environment": string(environment),
	}
	return NewError("Cannot Get Environment Puplic Keys", helpTexts["CannotGetEnvironmentKeys"], meta, cause)
}

func YouHaveLocallyModifiedFiles(environment string, files []string, cause error) *Error {
	meta := map[string]interface{}{
		"Environment": string(environment),
		"Files":       []string(files),
	}
	return NewError("You Have Locally Modified Files", helpTexts["YouHaveLocallyModifiedFiles"], meta, cause)
}

func SecretDoesNotExist(secret string, cause error) *Error {
	meta := map[string]interface{}{
		"Secret": string(secret),
	}
	return NewError("Secret Doesn't Exist", helpTexts["SecretDoesNotExist"], meta, cause)
}

func SecretRequired(secret string, cause error) *Error {
	meta := map[string]interface{}{
		"Secret": string(secret),
	}
	return NewError("Secret Required", helpTexts["SecretRequired"], meta, cause)
}

func SecretHasChanged(secret string, values string, cause error) *Error {
	meta := map[string]interface{}{
		"Secret": string(secret),
		"Values": string(values),
	}
	return NewError("Secret has changed", helpTexts["SecretHasChanged"], meta, cause)
}

func RequiredSecretsAreMissing(missingsecrets []string, environmentname string, cause error) *Error {
	meta := map[string]interface{}{
		"MissingSecrets":  []string(missingsecrets),
		"EnvironmentName": string(environmentname),
	}
	return NewError("Required Secrets Are Missing", helpTexts["RequiredSecretsAreMissing"], meta, cause)
}

func FileDoesNotExist(filename string, cause error) *Error {
	meta := map[string]interface{}{
		"FileName": string(filename),
	}
	return NewError("File Doesn't Exist", helpTexts["FileDoesNotExist"], meta, cause)
}

func RequiredFilesAreMissing(missingfiles []string, environmentname string, cause error) *Error {
	meta := map[string]interface{}{
		"MissingFiles":    []string(missingfiles),
		"EnvironmentName": string(environmentname),
	}
	return NewError("Required Files Are Missing", helpTexts["RequiredFilesAreMissing"], meta, cause)
}

func FileNotInWorkingDirectory(filepath string, wd string, cause error) *Error {
	meta := map[string]interface{}{
		"FilePath": string(filepath),
		"Wd":       string(wd),
	}
	return NewError("File Not In Working Directory", helpTexts["FileNotInWorkingDirectory"], meta, cause)
}

func EnvironmentsHaveChanged(environmentsname string, cause error) *Error {
	meta := map[string]interface{}{
		"EnvironmentsName": string(environmentsname),
	}
	return NewError("Messages expired", helpTexts["EnvironmentsHaveChanged"], meta, cause)
}

func FileHasChanged(filepath string, affectedenvironments string, cause error) *Error {
	meta := map[string]interface{}{
		"FilePath":             string(filepath),
		"AffectedEnvironments": string(affectedenvironments),
	}
	return NewError("File has changed", helpTexts["FileHasChanged"], meta, cause)
}

func CannotAddFile(path string, cause error) *Error {
	meta := map[string]interface{}{
		"Path": string(path),
	}
	return NewError("Cannot Add File", helpTexts["CannotAddFile"], meta, cause)
}

func CannotSetFile(path string, cause error) *Error {
	meta := map[string]interface{}{
		"Path": string(path),
	}
	return NewError("Cannot Set File", helpTexts["CannotSetFile"], meta, cause)
}

func CannotRemoveFile(path string, cause error) *Error {
	meta := map[string]interface{}{
		"Path": string(path),
	}
	return NewError("Cannot Remove File", helpTexts["CannotRemoveFile"], meta, cause)
}

func CannotCopyFile(path string, cachepath string, cause error) *Error {
	meta := map[string]interface{}{
		"Path":      string(path),
		"CachePath": string(cachepath),
	}
	return NewError("Cannot Copy File", helpTexts["CannotCopyFile"], meta, cause)
}

func FileNotInEnvironment(path string, environment string, cause error) *Error {
	meta := map[string]interface{}{
		"Path":        string(path),
		"Environment": string(environment),
	}
	return NewError("File Not Found For Environment", helpTexts["FileNotInEnvironment"], meta, cause)
}

func CannotCreateDirectory(path string, cause error) *Error {
	meta := map[string]interface{}{
		"Path": string(path),
	}
	return NewError("Cannot Create Directory", helpTexts["CannotCreateDirectory"], meta, cause)
}

func CannotRemoveDirectoryContents(path string, cause error) *Error {
	meta := map[string]interface{}{
		"Path": string(path),
	}
	return NewError("Cannot Remove Directory Contents", helpTexts["CannotRemoveDirectoryContents"], meta, cause)
}

func CannotSaveFiles(filelist string, cause error) *Error {
	meta := map[string]interface{}{
		"FileList": string(filelist),
	}
	return NewError("Cannot Save Files", helpTexts["CannotSaveFiles"], meta, cause)
}

func CannotRemoveDirectory(path string, cause error) *Error {
	meta := map[string]interface{}{
		"Path": string(path),
	}
	return NewError("Cannot Remove Directory", helpTexts["CannotRemoveDirectory"], meta, cause)
}

func CopyFailed(source string, destination string, cause error) *Error {
	meta := map[string]interface{}{
		"Source":      string(source),
		"Destination": string(destination),
	}
	return NewError("Copy failed", helpTexts["CopyFailed"], meta, cause)
}

func MustBeLoggedIn(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("You must be logged in", helpTexts["MustBeLoggedIn"], meta, cause)
}

func CannotFindProjectID(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Cannot find project ID in config file", helpTexts["CannotFindProjectID"], meta, cause)
}

func UnkownError(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Unkown Error", helpTexts["UnkownError"], meta, cause)
}

func UsersDontExist(message string, cause error) *Error {
	meta := map[string]interface{}{
		"Message": string(message),
	}
	return NewError("Users Don't Exist", helpTexts["UsersDontExist"], meta, cause)
}

func CannotAddMembers(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Cannot Add Members", helpTexts["CannotAddMembers"], meta, cause)
}

func CannotRemoveMembers(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Cannot Remove Members", helpTexts["CannotRemoveMembers"], meta, cause)
}

func MemberHasNoAccessToEnv(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Member has no access to environment", helpTexts["MemberHasNoAccessToEnv"], meta, cause)
}

func CouldNotDecryptMessages(message string, cause error) *Error {
	meta := map[string]interface{}{
		"Message": string(message),
	}
	return NewError("Could not decrypt messages", helpTexts["CouldNotDecryptMessages"], meta, cause)
}

func CouldNotEncryptMessages(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Could not encrypt messages", helpTexts["CouldNotEncryptMessages"], meta, cause)
}

func EncryptionFailed(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Encryption Failed", helpTexts["EncryptionFailed"], meta, cause)
}

func CouldNotParseMessage(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Could not parse message", helpTexts["CouldNotParseMessage"], meta, cause)
}

func PayloadErrors(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Errors occured while preparing the payload", helpTexts["PayloadErrors"], meta, cause)
}

func InvalidFileContent(path string, cause error) *Error {
	meta := map[string]interface{}{
		"Path": string(path),
	}
	return NewError("Invalid file content", helpTexts["InvalidFileContent"], meta, cause)
}

func FailedCheckingChanges(path string, cause error) *Error {
	meta := map[string]interface{}{
		"Path": string(path),
	}
	return NewError("Failed While Checking for Changes", helpTexts["FailedCheckingChanges"], meta, cause)
}

func FeatureRequiresToUpgrade(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("This Feature Requires to Upgrade", helpTexts["FeatureRequiresToUpgrade"], meta, cause)
}

func AlreadySubscribed(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("The Organization Has Already Been Upgraded", helpTexts["AlreadySubscribed"], meta, cause)
}

func CannotUpgrade(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Cannot Upgrade", helpTexts["CannotUpgrade"], meta, cause)
}

func ManagementInaccessible(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Management Inaccessible", helpTexts["ManagementInaccessible"], meta, cause)
}

func OrganizationDoesNotExist(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Organizaiton Does Not Exist", helpTexts["OrganizationDoesNotExist"], meta, cause)
}

func BadOrganizationName(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Bad Organization Name", helpTexts["BadOrganizationName"], meta, cause)
}

func OrganizationNameAlreadyTaken(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Organization Name Already Taken", helpTexts["OrganizationNameAlreadyTaken"], meta, cause)
}

func MustOwnTheOrganization(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("You Must Own the Organization", helpTexts["MustOwnTheOrganization"], meta, cause)
}

func YouDoNotOwnTheOrganization(organizationname string, cause error) *Error {
	meta := map[string]interface{}{
		"OrganizationName": string(organizationname),
	}
	return NewError("You Do Not Own An Organization Named", helpTexts["YouDoNotOwnTheOrganization"], meta, cause)
}

func CouldntSendInvite(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Couldn't Send Invite", helpTexts["CouldntSendInvite"], meta, cause)
}

func CouldntSetRole(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Couldn't Set Role", helpTexts["CouldntSetRole"], meta, cause)
}

func BackupDenied(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Permission Denied", helpTexts["BackupDenied"], meta, cause)
}

func RestoreDenied(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Permission Denied", helpTexts["RestoreDenied"], meta, cause)
}

func CouldNotCreateArchive(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Could Not Create Archive", helpTexts["CouldNotCreateArchive"], meta, cause)
}

func FailedToReadBackup(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Failed To Read Backup", helpTexts["FailedToReadBackup"], meta, cause)
}

func FailedToWriteBackup(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Failed To Write Backup", helpTexts["FailedToWriteBackup"], meta, cause)
}

func BackupNotSetUp(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Backup required", helpTexts["BackupNotSetUp"], meta, cause)
}

func NoBackup(projectname string, backupdirpath string, cause error) *Error {
	meta := map[string]interface{}{
		"ProjectName":   string(projectname),
		"BackupDirPath": string(backupdirpath),
	}
	return NewError("No Backup For Current Project", helpTexts["NoBackup"], meta, cause)
}

func RestoreFailed(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Backup Restoration Failed", helpTexts["RestoreFailed"], meta, cause)
}

func BackupFailed(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Backup Failed", helpTexts["BackupFailed"], meta, cause)
}

func NoCIServices(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("No CI Services", helpTexts["NoCIServices"], meta, cause)
}

func CiServiceAlreadyExists(servicename string, cause error) *Error {
	meta := map[string]interface{}{
		"ServiceName": string(servicename),
	}
	return NewError("A CI Service Already Exists", helpTexts["CiServiceAlreadyExists"], meta, cause)
}

func NoSuchService(servicename string, cause error) *Error {
	meta := map[string]interface{}{
		"ServiceName": string(servicename),
	}
	return NewError("No Such Service", helpTexts["NoSuchService"], meta, cause)
}

func CouldNotAddService(servicename string, cause error) *Error {
	meta := map[string]interface{}{
		"ServiceName": string(servicename),
	}
	return NewError("Could Not Add Service", helpTexts["CouldNotAddService"], meta, cause)
}

func CouldNotCleanService(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Could Not Clean Service", helpTexts["CouldNotCleanService"], meta, cause)
}

func CouldNotChangeService(servicename string, cause error) *Error {
	meta := map[string]interface{}{
		"ServiceName": string(servicename),
	}
	return NewError("Could Not Change Service", helpTexts["CouldNotChangeService"], meta, cause)
}

func CouldNotRemoveService(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Could Not Remove Service", helpTexts["CouldNotRemoveService"], meta, cause)
}

func MissingCIInformation(servicename string, cause error) *Error {
	meta := map[string]interface{}{
		"ServiceName": string(servicename),
	}
	return NewError("Missing Information for CI Service", helpTexts["MissingCIInformation"], meta, cause)
}

func CouldNotSendToCIService(cause error) *Error {
	meta := map[string]interface{}{}

	return NewError("Could Not Send to CI Service", helpTexts["CouldNotSendToCIService"], meta, cause)
}
