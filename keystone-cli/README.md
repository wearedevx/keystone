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
@keystone.sh/cli/0.0.3 darwin-x64 node-v10.14.0
$ ks --help [COMMAND]
USAGE
  $ ks COMMAND
...
```
<!-- usagestop -->

# Commands

<!-- commands -->
* [`ks add BLOCKSTACKID EMAIL`](#ks-add-blockstackid-email)
* [`ks cat PATH [DECRYPT]`](#ks-cat-path-decrypt)
* [`ks check`](#ks-check)
* [`ks delete [FILEPATHS]`](#ks-delete-filepaths)
* [`ks env [ACTION] [ENV]`](#ks-env-action-env)
* [`ks help [COMMAND]`](#ks-help-command)
* [`ks init [PROJECT_NAME]`](#ks-init-project_name)
* [`ks invite [EMAILS]`](#ks-invite-emails)
* [`ks list TYPE`](#ks-list-type)
* [`ks login [BLOCKSTACK_ID]`](#ks-login-blockstack_id)
* [`ks logout`](#ks-logout)
* [`ks open`](#ks-open)
* [`ks project [ACTION]`](#ks-project-action)
* [`ks pull`](#ks-pull)
* [`ks push [FILEPATH]`](#ks-push-filepath)
* [`ks remove`](#ks-remove)
* [`ks reset`](#ks-reset)
* [`ks share ACTION`](#ks-share-action)
* [`ks whoami`](#ks-whoami)
* [`ks wp`](#ks-wp)

## `ks add BLOCKSTACKID EMAIL`

Add a member to a project.

```
USAGE
  $ ks add BLOCKSTACKID EMAIL

ARGUMENTS
  BLOCKSTACKID  Blockstack_id to add
  EMAIL         email associated to an invitation

OPTIONS
  -c, --config=config    Set the path to the blockstack session file
  -p, --project=project  Set the project
  --removal              Deletes an invitation

DESCRIPTION
  The member should have accepted your invitation for this to work

EXAMPLE
  $ ks add example.id.blockstack example@mail.com #add a user to a project
```

_See code: [src/commands/add.js](https://github.com/wearedevx/keystone/blob/v0.0.3/src/commands/add.js)_

## `ks cat PATH [DECRYPT]`

Output a remote file.

```
USAGE
  $ ks cat PATH [DECRYPT]

ARGUMENTS
  PATH     path to your file
  DECRYPT  [default: true] should decrypt

OPTIONS
  -c, --config=config    Set the path to the blockstack session file
  -d, --debug            cat file with full path
  -o, --origin=origin    from origin
  -p, --project=project  Set the project
  --[no-]decrypt         Indiciate to decrypt or not
  --[no-]json            Indiciate to parse json or not
  --removal              Deletes an invitation

EXAMPLE
  $ ks cat my-file
```

_See code: [src/commands/cat.js](https://github.com/wearedevx/keystone/blob/v0.0.3/src/commands/cat.js)_

## `ks check`

```
USAGE
  $ ks check

OPTIONS
  -c, --config=config    Set the path to the blockstack session file
  -p, --project=project  Set the project
```

_See code: [src/commands/check.js](https://github.com/wearedevx/keystone/blob/v0.0.3/src/commands/check.js)_

## `ks delete [FILEPATHS]`

Deletes one or more files.

```
USAGE
  $ ks delete [FILEPATHS]

ARGUMENTS
  FILEPATHS  Path to your file. Accepts a glob pattern

OPTIONS
  -c, --config=config    Set the path to the blockstack session file
  -p, --project=project  Set the project

DESCRIPTION
  If you're an administrator or a contributor, the files will be removed for everyone.
  If you're a reader on the project, you can't delete any files.
```

_See code: [src/commands/delete.js](https://github.com/wearedevx/keystone/blob/v0.0.3/src/commands/delete.js)_

## `ks env [ACTION] [ENV]`

Manage environments.

```
USAGE
  $ ks env [ACTION] [ENV]

ARGUMENTS
  ACTION  Configure add or remove an environment
  ENV     Set working env

OPTIONS
  -n, --name=name  Enviroment name

EXAMPLE
  $ ks env remove --name dev
```

_See code: [src/commands/env.js](https://github.com/wearedevx/keystone/blob/v0.0.3/src/commands/env.js)_

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

OPTIONS
  -c, --config=config    Set the path to the blockstack session file
  -p, --project=project  Set the project

EXAMPLE
  $ ks init project-name
```

_See code: [src/commands/init.js](https://github.com/wearedevx/keystone/blob/v0.0.3/src/commands/init.js)_

## `ks invite [EMAILS]`

Invites one or more people by email to a project.

```
USAGE
  $ ks invite [EMAILS]

ARGUMENTS
  EMAILS  Emails for invitations to be sent

OPTIONS
  -c, --config=config                  Set the path to the blockstack session file
  -p, --project=project                Set the project
  -r, --role=reader|contributor|admin  [default: reader] Assigns a role
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
  $ ks invite friend@example.com friend2@example.com --role=admin #Invite as admin on the project
  $ ks invite friend@example.com friend2@example.com --removal #Removes the invitations for friend and friend2
```

_See code: [src/commands/invite.js](https://github.com/wearedevx/keystone/blob/v0.0.3/src/commands/invite.js)_

## `ks list TYPE`

Lists projects, environments, members and files

```
USAGE
  $ ks list TYPE

ARGUMENTS
  TYPE  What do you want to list (projects, environments, members or files)

OPTIONS
  -a, --all              List all elements
  -c, --config=config    Set the path to the blockstack session file
  -p, --project=project  Set the project

EXAMPLE
  $ ks list members
```

_See code: [src/commands/list.js](https://github.com/wearedevx/keystone/blob/v0.0.3/src/commands/list.js)_

## `ks login [BLOCKSTACK_ID]`

Logs into your account with Blockstack or creates a new one

```
USAGE
  $ ks login [BLOCKSTACK_ID]

ARGUMENTS
  BLOCKSTACK_ID  Your blockstack id
```

_See code: [src/commands/login.js](https://github.com/wearedevx/keystone/blob/v0.0.3/src/commands/login.js)_

## `ks logout`

Logs you out of your account and erase your session from this computer.

```
USAGE
  $ ks logout

EXAMPLE
  $ ks logout
```

_See code: [src/commands/logout.js](https://github.com/wearedevx/keystone/blob/v0.0.3/src/commands/logout.js)_

## `ks open`

```
USAGE
  $ ks open

OPTIONS
  -c, --config=config    Set the path to the blockstack session file
  -p, --project=project  Set the project
```

_See code: [src/commands/open.js](https://github.com/wearedevx/keystone/blob/v0.0.3/src/commands/open.js)_

## `ks project [ACTION]`

Manage project.

```
USAGE
  $ ks project [ACTION]

ARGUMENTS
  ACTION  Configure project members

OPTIONS
  -c, --config=config    Set the path to the blockstack session file
  -p, --project=project  Set the project

EXAMPLE
  $ ks env config
```

_See code: [src/commands/project.js](https://github.com/wearedevx/keystone/blob/v0.0.3/src/commands/project.js)_

## `ks pull`

Fetch files for current environment.

```
USAGE
  $ ks pull

OPTIONS
  -f, --force  Overwrite any changes
```

_See code: [src/commands/pull.js](https://github.com/wearedevx/keystone/blob/v0.0.3/src/commands/pull.js)_

## `ks push [FILEPATH]`

Push a file to a project.

```
USAGE
  $ ks push [FILEPATH]

ARGUMENTS
  FILEPATH  Path to your file. Accepts a glob pattern

OPTIONS
  -c, --config=config    Set the path to the blockstack session file
  -e, --encrypt=encrypt  * DEBUG ONLY * encrypt the file with given blockstackid
  -p, --path=path        * DEBUG ONLY * push the file to the given path
  -p, --project=project  Set the project
```

_See code: [src/commands/push.js](https://github.com/wearedevx/keystone/blob/v0.0.3/src/commands/push.js)_

## `ks remove`

```
USAGE
  $ ks remove

OPTIONS
  -c, --config=config    Set the path to the blockstack session file
  -p, --project=project  Set the project
```

_See code: [src/commands/remove.js](https://github.com/wearedevx/keystone/blob/v0.0.3/src/commands/remove.js)_

## `ks reset`

Remove everything but your public.key file.

```
USAGE
  $ ks reset

OPTIONS
  -c, --config=config    Set the path to the blockstack session file
  -p, --project=project  Set the project

EXAMPLE
  $ ks reset
```

_See code: [src/commands/reset.js](https://github.com/wearedevx/keystone/blob/v0.0.3/src/commands/reset.js)_

## `ks share ACTION`

Share your file file with a non-blockstack user

```
USAGE
  $ ks share ACTION

ARGUMENTS
  ACTION  new || pull. Create a new shared user or pull files based on keystone-link.json file.

OPTIONS
  -e, --env=env    Env you want to create the user in.
  -l, --link=link  Path to your link file.

EXAMPLE
  $ ks share
```

_See code: [src/commands/share.js](https://github.com/wearedevx/keystone/blob/v0.0.3/src/commands/share.js)_

## `ks whoami`

Shows the blockstack id of the currently logged in user

```
USAGE
  $ ks whoami

OPTIONS
  -c, --config=config    Set the path to the blockstack session file
  -p, --project=project  Set the project

EXAMPLE
  $ ks whoami
```

_See code: [src/commands/whoami.js](https://github.com/wearedevx/keystone/blob/v0.0.3/src/commands/whoami.js)_

## `ks wp`

Print project name in current workspace

```
USAGE
  $ ks wp

OPTIONS
  -c, --config=config    Set the path to the blockstack session file
  -p, --project=project  Set the project

EXAMPLE
  $ ks wp
```

_See code: [src/commands/wp.js](https://github.com/wearedevx/keystone/blob/v0.0.3/src/commands/wp.js)_
<!-- commandsstop -->
