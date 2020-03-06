# Keystone CLI

Open-source decentralized system for developers to store, share and use secrets safely

<!-- toc -->
* [Keystone CLI](#keystone-cli)
* [Usage](#usage)
* [Commands](#commands)
<!-- tocstop -->

# Usage

<!-- usage -->
```sh-session
$ npm install -g @keystone.sh/cli
$ ks COMMAND
running command...
$ ks (-v|--version|version)
@keystone.sh/cli/0.0.35 linux-x64 node-v10.19.0
$ ks --help [COMMAND]
USAGE
  $ ks COMMAND
...
```
<!-- usagestop -->

# Commands

<!-- commands -->
* [`ks add BLOCKSTACKID EMAIL`](#ks-add-blockstackid-email)
* [`ks autocomplete [SHELL]`](#ks-autocomplete-shell)
* [`ks cat PATH`](#ks-cat-path)
* [`ks delete [FILEPATHS]`](#ks-delete-filepaths)
* [`ks diff FILEPATH`](#ks-diff-filepath)
* [`ks env [ACTION] [ENV]`](#ks-env-action-env)
* [`ks help [COMMAND]`](#ks-help-command)
* [`ks init [PROJECT_NAME]`](#ks-init-project_name)
* [`ks invite [EMAILS]`](#ks-invite-emails)
* [`ks list TYPE`](#ks-list-type)
* [`ks login [BLOCKSTACK_ID]`](#ks-login-blockstack_id)
* [`ks logout`](#ks-logout)
* [`ks project ACTION`](#ks-project-action)
* [`ks pull`](#ks-pull)
* [`ks push [FILEPATH]`](#ks-push-filepath)
* [`ks remove`](#ks-remove)
* [`ks share ENV`](#ks-share-env)
* [`ks status`](#ks-status)
* [`ks whoami`](#ks-whoami)

## `ks add BLOCKSTACKID EMAIL`

Add a member to a project.

```
USAGE
  $ ks add BLOCKSTACKID EMAIL

ARGUMENTS
  BLOCKSTACKID  Blockstack_id to add
  EMAIL         email associated to an invitation

DESCRIPTION
  Adding a member give them access to the project.
  The member should have accepted your invitation for this to work

  You  can add the member to an environment with : $ ks env config

EXAMPLE
  $ ks add example.id.blockstack example@mail.com #add a user to a project
```

_See code: [src/commands/add.js](https://github.com/wearedevx/keystone/blob/v0.0.35/src/commands/add.js)_

## `ks autocomplete [SHELL]`

display autocomplete installation instructions

```
USAGE
  $ ks autocomplete [SHELL]

ARGUMENTS
  SHELL  shell type

OPTIONS
  -r, --refresh-cache  Refresh cache (ignores displaying instructions)

EXAMPLES
  $ ks autocomplete
  $ ks autocomplete bash
  $ ks autocomplete zsh
  $ ks autocomplete --refresh-cache
```

_See code: [@oclif/plugin-autocomplete](https://github.com/oclif/plugin-autocomplete/blob/v0.1.5/src/commands/autocomplete/index.ts)_

## `ks cat PATH`

Output a remote file.

```
USAGE
  $ ks cat PATH

ARGUMENTS
  PATH  path to your file

EXAMPLE
  $ ks cat path/to/file
```

_See code: [src/commands/cat.js](https://github.com/wearedevx/keystone/blob/v0.0.35/src/commands/cat.js)_

## `ks delete [FILEPATHS]`

Deletes one or more files.

```
USAGE
  $ ks delete [FILEPATHS]

ARGUMENTS
  FILEPATHS  Path to your file. Accepts a glob pattern

OPTIONS
  -p, --project=project  Use this flag to completely delete all files of a project from your storage.

DESCRIPTION
  If you're an administrator or a contributor, the files will be removed for everyone.
  If you're a reader on the environment, you can't delete any files.

EXAMPLES
  $ ks delete path/to/file
  $ ks delete -p project_name/2b6a10c6-ea91-48b1-b340-a7504326961e
```

_See code: [src/commands/delete.js](https://github.com/wearedevx/keystone/blob/v0.0.35/src/commands/delete.js)_

## `ks diff FILEPATH`

Output a diff of the changes you made to a file

```
USAGE
  $ ks diff FILEPATH

ARGUMENTS
  FILEPATH  Path to your file.

EXAMPLE
  $ ks diff path/to/file
```

_See code: [src/commands/diff.js](https://github.com/wearedevx/keystone/blob/v0.0.35/src/commands/diff.js)_

## `ks env [ACTION] [ENV]`

Manage environments.

```
USAGE
  $ ks env [ACTION] [ENV]

ARGUMENTS
  ACTION
      - config
           Change users role for each environment.

         - new 
           Create a new environment

         - remove 
           Remove an environment

  ENV
      Set working env

DESCRIPTION
  You need to be administrator in the project in order to access the command.

  You can change the role set by using the role flag. You have 3 choices:
  - reader: can only read files from the the environment and pull them locally
  - contributor: can read, write and add new files to the environement
  - admin: all the above plus ask people to join the project

EXAMPLES
  $ ks env config
  $ ks env new ENV_NAME
  $ ks env remove ENV_NAME
```

_See code: [src/commands/env.js](https://github.com/wearedevx/keystone/blob/v0.0.35/src/commands/env.js)_

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

_See code: [@oclif/plugin-help](https://github.com/oclif/plugin-help/blob/v2.2.3/src/commands/help.ts)_

## `ks init [PROJECT_NAME]`

Create Keystone config file

```
USAGE
  $ ks init [PROJECT_NAME]

ARGUMENTS
  PROJECT_NAME  Your project name

EXAMPLE
  $ ks init project_name
```

_See code: [src/commands/init.js](https://github.com/wearedevx/keystone/blob/v0.0.35/src/commands/init.js)_

## `ks invite [EMAILS]`

Invites one or more people by email to a project.

```
USAGE
  $ ks invite [EMAILS]

ARGUMENTS
  EMAILS  Emails for invitations to be sent

OPTIONS
  -r, --role=reader|contributor|admin  [default: reader] Assigns a role
  --check                              Check your pending invitations
  --removal                            Deletes an invitation

DESCRIPTION
  By default, people you invite are readers.
  You can change the role set by using the role flag. You have 3 choices:
  - reader: cannot do anything project wide. Need to be added to an environment to pull files
  - contributor: can add and remove environments from the project
  - admin: all the above plus invite and add users to the project

EXAMPLES
  $ ks invite friend@example.com #Send an invitation to friend@example.com as a reader on the project
  $ ks invite friend@example.com friend2@example.com --role=admin #Invite as admin on the project
  $ ks invite friend@example.com friend2@example.com --removal #Removes the invitations for friend and friend2
```

_See code: [src/commands/invite.js](https://github.com/wearedevx/keystone/blob/v0.0.35/src/commands/invite.js)_

## `ks list TYPE`

Lists projects, environments, members and files

```
USAGE
  $ ks list TYPE

ARGUMENTS
  TYPE  What do you want to list (projects, environments, members or files)

OPTIONS
  -a, --all  For files listing, list every files in your gaia hub. For members, list files from project, instead of the
             environment.

EXAMPLES
  $ ks list members
  $ ks list members --all
  $ ks list projects
  $ ks list environments
  $ ks list files
```

_See code: [src/commands/list.js](https://github.com/wearedevx/keystone/blob/v0.0.35/src/commands/list.js)_

## `ks login [BLOCKSTACK_ID]`

Logs into your account with Blockstack or creates a new one

```
USAGE
  $ ks login [BLOCKSTACK_ID]

ARGUMENTS
  BLOCKSTACK_ID  Your blockstack id

EXAMPLE
  $ ks login nickname.id.blockstack
```

_See code: [src/commands/login.js](https://github.com/wearedevx/keystone/blob/v0.0.35/src/commands/login.js)_

## `ks logout`

Logs you out of your account and erase your session from this computer.

```
USAGE
  $ ks logout

EXAMPLE
  $ ks logout
```

_See code: [src/commands/logout.js](https://github.com/wearedevx/keystone/blob/v0.0.35/src/commands/logout.js)_

## `ks project ACTION`

Manage users role in the project.

```
USAGE
  $ ks project ACTION

ARGUMENTS
  ACTION  Configure project members

DESCRIPTION
  You can change the role set by using the role flag. You have 3 choices:
  - reader: can't do anything regarding the project itself.
  - contributor: can add or remove environments.
  - administrator: can add or remove environments, add and remove users, change users roles.

EXAMPLE
  $ ks project config
```

_See code: [src/commands/project.js](https://github.com/wearedevx/keystone/blob/v0.0.35/src/commands/project.js)_

## `ks pull`

Fetch files for current environment. Write them locally.

```
USAGE
  $ ks pull

OPTIONS
  -f, --force  Overwrite any changes made locally

DESCRIPTION
  Once pulled files can be one of the three states :
     - updated : The file has been updated because someone else pushed a newer version
     - auto-merged : The file was modified and has been merged with someone else's changes
     - conflicted : The file has been modified and some lines are in conflict with someone else's changes. You should 
  fix the conflicts and push your changes

EXAMPLE
  $ ks pull
```

_See code: [src/commands/pull.js](https://github.com/wearedevx/keystone/blob/v0.0.35/src/commands/pull.js)_

## `ks push [FILEPATH]`

Push a file to a project.

```
USAGE
  $ ks push [FILEPATH]

ARGUMENTS
  FILEPATH  Path to your file. Accepts a glob pattern

EXAMPLES
  $ ks push path/to/my/file
  $ ks push
```

_See code: [src/commands/push.js](https://github.com/wearedevx/keystone/blob/v0.0.35/src/commands/push.js)_

## `ks remove`

Remove one or more users.

```
USAGE
  $ ks remove

OPTIONS
  -u, --users=users  List of user you want to remove. Separated by space.

DESCRIPTION
  ...
  If you are an administrator, you can remove a user from a project.

EXAMPLE
  $ ks remove nickname1.id.blockstack nickname2.id.blockstack
```

_See code: [src/commands/remove.js](https://github.com/wearedevx/keystone/blob/v0.0.35/src/commands/remove.js)_

## `ks share ENV`

Share your files with a non-blockstack user

```
USAGE
  $ ks share ENV

ARGUMENTS
  ENV  Environment you want the user to be created on.

DESCRIPTION
  Generate a token. 
  The token should be set in the system environment of any user.
  This user will be able to run only $ ks pull in order to pull locally files from the selected env.

EXAMPLE
  $ ks share ENV_NAME
```

_See code: [src/commands/share.js](https://github.com/wearedevx/keystone/blob/v0.0.35/src/commands/share.js)_

## `ks status`

Shows the status of tracked files

```
USAGE
  $ ks status

EXAMPLE
  $ ks status
```

_See code: [src/commands/status.js](https://github.com/wearedevx/keystone/blob/v0.0.35/src/commands/status.js)_

## `ks whoami`

Shows the blockstack id of the currently logged in user

```
USAGE
  $ ks whoami

EXAMPLE
  $ ks whoami
```

_See code: [src/commands/whoami.js](https://github.com/wearedevx/keystone/blob/v0.0.35/src/commands/whoami.js)_
<!-- commandsstop -->
