secret
======

Store and manage your secret on the blockchain with Blockstack

[![oclif](https://img.shields.io/badge/cli-oclif-brightgreen.svg)](https://oclif.io)
[![Version](https://img.shields.io/npm/v/secret.svg)](https://npmjs.org/package/secret)
[![Downloads/week](https://img.shields.io/npm/dw/secret.svg)](https://npmjs.org/package/secret)
[![License](https://img.shields.io/npm/l/secret.svg)](https://github.com/http://github.com/samuelroy/secret/secret/blob/master/package.json)

<!-- toc -->
* [Usage](#usage)
* [Commands](#commands)
<!-- tocstop -->
# Usage
<!-- usage -->
```sh-session
$ npm install -g keystone-cli
$ ks COMMAND
running command...
$ ks (-v|--version|version)
keystone-cli/0.0.1 darwin-x64 node-v10.14.0
$ ks --help [COMMAND]
USAGE
  $ ks COMMAND
...
```
<!-- usagestop -->
# Commands
<!-- commands -->
* [`ks add ID EMAIL ROLE`](#ks-add-id-email-role)
* [`ks cat PATH [DECRYPT]`](#ks-cat-path-decrypt)
* [`ks check`](#ks-check)
* [`ks delete [FILES]`](#ks-delete-files)
* [`ks fetch [FILES]`](#ks-fetch-files)
* [`ks help [COMMAND]`](#ks-help-command)
* [`ks init PROJECT_NAME`](#ks-init-project_name)
* [`ks invite [EMAILS]`](#ks-invite-emails)
* [`ks list [PROJECT_NAME]`](#ks-list-project_name)
* [`ks login [BLOCKSTACK_ID]`](#ks-login-blockstack_id)
* [`ks logout`](#ks-logout)
* [`ks new PROJECT_NAME`](#ks-new-project_name)
* [`ks push FILEPATH`](#ks-push-filepath)
* [`ks remove`](#ks-remove)
* [`ks rm PATH`](#ks-rm-path)
* [`ks whoami`](#ks-whoami)

## `ks add ID EMAIL ROLE`

Add a member to a project. -- debug function

```
USAGE
  $ ks add ID EMAIL ROLE

ARGUMENTS
  ID     Blockstack_id to add
  EMAIL  email associated to an invitation
  ROLE   [default: reader] role to add

OPTIONS
  -p, --project=project                Set the project
  -r, --role=reader|contributor|admin  [default: reader] Assigns a role
  --removal                            Deletes an invitation

DESCRIPTION
  The member should have accepted your invitation for this to work

EXAMPLE
  $ ks add example.id.blockstack example@mail.com contributor --project=keystone #add a user to a project
```

_See code: [src/commands/add.js](https://github.com/keystone.sh/blob/v0.0.1/src/commands/add.js)_

## `ks cat PATH [DECRYPT]`

Output a remote file.

```
USAGE
  $ ks cat PATH [DECRYPT]

ARGUMENTS
  PATH     path to your file
  DECRYPT  [default: true] should decrypt

OPTIONS
  -p, --project=project  Set the project

EXAMPLE
  $ ks cat my-file
```

_See code: [src/commands/cat.js](https://github.com/keystone.sh/blob/v0.0.1/src/commands/cat.js)_

## `ks check`

Check if a blockstack id has a Keystone application public key.

```
USAGE
  $ ks check

OPTIONS
  -p, --project=project  Set the project

EXAMPLE
  $ ks check example.id.blockstack
```

_See code: [src/commands/check.js](https://github.com/keystone.sh/blob/v0.0.1/src/commands/check.js)_

## `ks delete [FILES]`

Deletes one or more files.

```
USAGE
  $ ks delete [FILES]

ARGUMENTS
  FILES  List of files to delete

OPTIONS
  -p, --project=project  Set the project
  -t, --tags=tags        [default: ] Looks for files with one of the tags set

DESCRIPTION
  If you're an administrator or a contributor, the files will be removed for everyone.
  If you're a reader on the project, you can't delete any files.

EXAMPLES
  $ ks delete config/.env #deletes config/.env from your project set in .ksconfig
  $ ks delete #deletes all files from your project set in .ksconfig
  $ ks delete --tags=dev #deletes all files tagged with dev
  $ ks delete --project=my-project #deletes all files from my-project
```

_See code: [src/commands/delete.js](https://github.com/keystone.sh/blob/v0.0.1/src/commands/delete.js)_

## `ks fetch [FILES]`

Fetch one or more files.

```
USAGE
  $ ks fetch [FILES]

ARGUMENTS
  FILES  List of files to fetch

OPTIONS
  -d, --directory=directory  Set the destination folder
  -p, --project=project      Set the project
  -t, --tags=tags            Looks for files with one of the tags set

EXAMPLES
  $ ks fetch #fetch all files from the project set in .ksconfig
  $ ks fetch my-file #fetch my-file from the project set in .ksconfig
  $ ks fetch --tags=dev #fetch all files tagged with dev from the project set in .ksconfig
  $ ks fetch --project=my-project #fetch all files from my-project
  $ ks fetch my-file --directory=config/ #fetch my-file and copy to directory config/
```

_See code: [src/commands/fetch.js](https://github.com/keystone.sh/blob/v0.0.1/src/commands/fetch.js)_

## `ks help [COMMAND]`

display help for ks

```
USAGE
  $ ks help [COMMAND]

ARGUMENTS
  COMMAND  command to show help for

OPTIONS
  --all  see all commands in CLI
```

_See code: [@oclif/plugin-help](https://github.com/oclif/plugin-help/blob/v2.1.6/src/commands/help.ts)_

## `ks init PROJECT_NAME`

Set a default project and create a keystone config file in your folder

```
USAGE
  $ ks init PROJECT_NAME

ARGUMENTS
  PROJECT_NAME  Your project name

OPTIONS
  -d, --directory=directory  [default: ./] Default directory. Relative path from your root folder.

EXAMPLE
  $ ks init default-project
```

_See code: [src/commands/init.js](https://github.com/keystone.sh/blob/v0.0.1/src/commands/init.js)_

## `ks invite [EMAILS]`

Invites one or more people by email to a project.

```
USAGE
  $ ks invite [EMAILS]

ARGUMENTS
  EMAILS  Emails for invitations to be sent

OPTIONS
  -p, --project=project                Set the project
  -r, --role=reader|contributor|admin  [default: reader] Assigns a role
  --accept                             Accept an invitation
  --check                              Check your pending invitations
  --removal                            Deletes an invitation

DESCRIPTION
  By default, people you invite are readers. 
  You can change the role set by using the role flag. You have 3 choices:
  - reader: can only read files from the project.
  - contributor: can read, write and add new files to the project
  - admin: all the above plus ask people to join the project

EXAMPLES
  $ ks invite friend@example.com #Send an invitation to friend@example.com as a reader on the project
  $ ks invite friend@example.com --project=my-project #Send an invitation to friend@example for the project my-project
  $ ks invite friend@example.com friend2@example.com --role=admin #Invite as admin on the project
  $ ks invite friend@example.com friend2@example.com --removal #Removes the invitations for friend and friend2
```

_See code: [src/commands/invite.js](https://github.com/keystone.sh/blob/v0.0.1/src/commands/invite.js)_

## `ks list [PROJECT_NAME]`

Lists projects and files

```
USAGE
  $ ks list [PROJECT_NAME]

ARGUMENTS
  PROJECT_NAME  Your project name

OPTIONS
  -k, --keystone  Debug function to show all files managed under Keystone

EXAMPLES
  $ ks list
  $ ks list --project=my-project
```

_See code: [src/commands/list.js](https://github.com/keystone.sh/blob/v0.0.1/src/commands/list.js)_

## `ks login [BLOCKSTACK_ID]`

Logs into your account with Blockstack or creates a new one

```
USAGE
  $ ks login [BLOCKSTACK_ID]

ARGUMENTS
  BLOCKSTACK_ID  Your blockstack id
```

_See code: [src/commands/login.js](https://github.com/keystone.sh/blob/v0.0.1/src/commands/login.js)_

## `ks logout`

Logs you out of your account and erase your session from this computer.

```
USAGE
  $ ks logout

EXAMPLE
  $ ks logout
```

_See code: [src/commands/logout.js](https://github.com/keystone.sh/blob/v0.0.1/src/commands/logout.js)_

## `ks new PROJECT_NAME`

Creates a new project.

```
USAGE
  $ ks new PROJECT_NAME

ARGUMENTS
  PROJECT_NAME  Your project name

EXAMPLE
  $ ks new my-project
```

_See code: [src/commands/new.js](https://github.com/keystone.sh/blob/v0.0.1/src/commands/new.js)_

## `ks push FILEPATH`

Push a file to a project.

```
USAGE
  $ ks push FILEPATH

ARGUMENTS
  FILEPATH  Path to your file(s). Accepts a glob pattern

OPTIONS
  -p, --project=project  project to push files to
  -t, --tags=tags        assigns one or more tags to your file

EXAMPLES
  $ ks push my-file
  $ ks push my-file --project=my-project
  $ ks push my-file --tags=dev,prod
```

_See code: [src/commands/push.js](https://github.com/keystone.sh/blob/v0.0.1/src/commands/push.js)_

## `ks remove`

Remove a project.

```
USAGE
  $ ks remove

OPTIONS
  -f, --force=force      Force the deletion. Beware, it might let some files from the project on your storage.
  -p, --project=project  Set the project

DESCRIPTION
  ...
  If you're an administrator, the project will be removed for everyone.

  If you're a contributor or a reader, you will be removed from the project.

EXAMPLES
  $ ks remove #remove your project set in .ksconfig and all its files
  $ ks remove --project=my-project #remove your project called my-project
```

_See code: [src/commands/remove.js](https://github.com/keystone.sh/blob/v0.0.1/src/commands/remove.js)_

## `ks rm PATH`

Remove a remote file.

```
USAGE
  $ ks rm PATH

ARGUMENTS
  PATH  path to your file

OPTIONS
  -p, --project=project  Set the project

EXAMPLE
  $ ks rm my-file
```

_See code: [src/commands/rm.js](https://github.com/keystone.sh/blob/v0.0.1/src/commands/rm.js)_

## `ks whoami`

Shows the blockstack id of the currently logged in user

```
USAGE
  $ ks whoami

EXAMPLE
  $ ks whoami
```

_See code: [src/commands/whoami.js](https://github.com/keystone.sh/blob/v0.0.1/src/commands/whoami.js)_
<!-- commandsstop -->
