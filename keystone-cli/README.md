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
@keystone.sh/cli/1.0.0 linux-x64 node-v10.19.0
$ ks --help [COMMAND]
USAGE
  $ ks COMMAND
...
```
<!-- usagestop -->

# Commands

<!-- commands -->
* [`ks autocomplete [SHELL]`](#ks-autocomplete-shell)
* [`ks diff FILEPATH`](#ks-diff-filepath)
* [`ks env add ENV`](#ks-env-add-env)
* [`ks env checkout ENV`](#ks-env-checkout-env)
* [`ks env config`](#ks-env-config)
* [`ks env list`](#ks-env-list)
* [`ks env reset`](#ks-env-reset)
* [`ks env rm ENV`](#ks-env-rm-env)
* [`ks help [COMMAND]`](#ks-help-command)
* [`ks init [PROJECT_NAME]`](#ks-init-project_name)
* [`ks list`](#ks-list)
* [`ks login [BLOCKSTACK_ID]`](#ks-login-blockstack_id)
* [`ks logout`](#ks-logout)
* [`ks member add BLOCKSTACKID EMAIL`](#ks-member-add-blockstackid-email)
* [`ks member invite [EMAILS]`](#ks-member-invite-emails)
* [`ks member list`](#ks-member-list)
* [`ks member rm`](#ks-member-rm)
* [`ks project config`](#ks-project-config)
* [`ks project list`](#ks-project-list)
* [`ks project rm PROJECT`](#ks-project-rm-project)
* [`ks pull`](#ks-pull)
* [`ks push [FILEPATH]`](#ks-push-filepath)
* [`ks rm FILEPATHS`](#ks-rm-filepaths)
* [`ks status`](#ks-status)
* [`ks token create ENV`](#ks-token-create-env)
* [`ks token revoke ENV`](#ks-token-revoke-env)
* [`ks whoami`](#ks-whoami)

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

## `ks diff FILEPATH`

output a diff of the changes you made to a file

```
USAGE
  $ ks diff FILEPATH

ARGUMENTS
  FILEPATH  path to your file.

EXAMPLE
  $ ks diff path/to/file
```

_See code: [src/commands/diff.js](https://github.com/wearedevx/keystone/blob/v1.0.0/src/commands/diff.js)_

## `ks env add ENV`

add a new environment to the project

```
USAGE
  $ ks env add ENV

ARGUMENTS
  ENV  environment name

EXAMPLE
  $ ks env add ENV_NAME
```

_See code: [src/commands/env/add.js](https://github.com/wearedevx/keystone/blob/v1.0.0/src/commands/env/add.js)_

## `ks env checkout ENV`

switch environment and pull files

```
USAGE
  $ ks env checkout ENV

ARGUMENTS
  ENV  environment name

EXAMPLE
  $ ks env checkout ENV_NAME
```

_See code: [src/commands/env/checkout.js](https://github.com/wearedevx/keystone/blob/v1.0.0/src/commands/env/checkout.js)_

## `ks env config`

manage members role in project environments

```
USAGE
  $ ks env config

DESCRIPTION
  roles are the followings :
     reader: can only read files from the the environment and pull them locally
     contributor: can read, write and add new files to the environement
     admin: all the above plus ask people to join the project

EXAMPLE
  $ ks env config
```

_See code: [src/commands/env/config.js](https://github.com/wearedevx/keystone/blob/v1.0.0/src/commands/env/config.js)_

## `ks env list`

List environments

```
USAGE
  $ ks env list

EXAMPLE
  $ ks env list
```

_See code: [src/commands/env/list.js](https://github.com/wearedevx/keystone/blob/v1.0.0/src/commands/env/list.js)_

## `ks env reset`

reset changes you made locally in tracked files

```
USAGE
  $ ks env reset

EXAMPLE
  $ ks env reset
```

_See code: [src/commands/env/reset.js](https://github.com/wearedevx/keystone/blob/v1.0.0/src/commands/env/reset.js)_

## `ks env rm ENV`

remove an environment

```
USAGE
  $ ks env rm ENV

ARGUMENTS
  ENV  environment name

EXAMPLE
  $ ks env rm ENV_NAME
```

_See code: [src/commands/env/rm.js](https://github.com/wearedevx/keystone/blob/v1.0.0/src/commands/env/rm.js)_

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

create Keystone config file

```
USAGE
  $ ks init [PROJECT_NAME]

ARGUMENTS
  PROJECT_NAME  your project name

EXAMPLE
  $ ks init project_name
```

_See code: [src/commands/init.js](https://github.com/wearedevx/keystone/blob/v1.0.0/src/commands/init.js)_

## `ks list`

list files tracked for your current environment

```
USAGE
  $ ks list

OPTIONS
  -a, --all  list every files in your gaia hub

EXAMPLE
  $ ks list
```

_See code: [src/commands/list.js](https://github.com/wearedevx/keystone/blob/v1.0.0/src/commands/list.js)_

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

_See code: [src/commands/login.js](https://github.com/wearedevx/keystone/blob/v1.0.0/src/commands/login.js)_

## `ks logout`

Logs you out of your account and erase your session from this computer.

```
USAGE
  $ ks logout

EXAMPLE
  $ ks logout
```

_See code: [src/commands/logout.js](https://github.com/wearedevx/keystone/blob/v1.0.0/src/commands/logout.js)_

## `ks member add BLOCKSTACKID EMAIL`

add a member to a project.

```
USAGE
  $ ks member add BLOCKSTACKID EMAIL

ARGUMENTS
  BLOCKSTACKID  blockstack_id to add
  EMAIL         email associated to an invitation

DESCRIPTION
  adding a member give them access to the project.
  the member should have accepted your invitation for this to work

  you can add the member to an environment with : $ ks env config

EXAMPLE
  $ ks member add example.id.blockstack example@mail.com #add a user to a project
```

_See code: [src/commands/member/add.js](https://github.com/wearedevx/keystone/blob/v1.0.0/src/commands/member/add.js)_

## `ks member invite [EMAILS]`

Invites one or more people by email to a project.

```
USAGE
  $ ks member invite [EMAILS]

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
  $ ks member invite friend@example.com #Send an invitation to friend@example.com as a reader on the project
  $ ks member invite friend@example.com friend2@example.com --role=admin #Invite as admin on the project
  $ ks member invite friend@example.com friend2@example.com --removal #Removes the invitations for friend and friend2
```

_See code: [src/commands/member/invite.js](https://github.com/wearedevx/keystone/blob/v1.0.0/src/commands/member/invite.js)_

## `ks member list`

list members from current environment or project

```
USAGE
  $ ks member list

EXAMPLE
  $ ks member list
```

_See code: [src/commands/member/list.js](https://github.com/wearedevx/keystone/blob/v1.0.0/src/commands/member/list.js)_

## `ks member rm`

remove one or more users

```
USAGE
  $ ks member rm

EXAMPLE
  $ ks member rm nickname1.id.blockstack nickname2.id.blockstack
```

_See code: [src/commands/member/rm.js](https://github.com/wearedevx/keystone/blob/v1.0.0/src/commands/member/rm.js)_

## `ks project config`

manage members role in the project

```
USAGE
  $ ks project config

DESCRIPTION
  roles are the followings :
     reader: can't do anything regarding the project itself
     contributor: can add or remove environments
     administrator: can add or remove environments, add and remove users, change users roles

EXAMPLE
  $ ks project config
```

_See code: [src/commands/project/config.js](https://github.com/wearedevx/keystone/blob/v1.0.0/src/commands/project/config.js)_

## `ks project list`

List projects in user workspace

```
USAGE
  $ ks project list

EXAMPLE
  $ ks project list
```

_See code: [src/commands/project/list.js](https://github.com/wearedevx/keystone/blob/v1.0.0/src/commands/project/list.js)_

## `ks project rm PROJECT`

remove project and its files from your storage

```
USAGE
  $ ks project rm PROJECT

ARGUMENTS
  PROJECT  project name (with uuid)

EXAMPLE
  $ ks project rm PROJECT_NAME
```

_See code: [src/commands/project/rm.js](https://github.com/wearedevx/keystone/blob/v1.0.0/src/commands/project/rm.js)_

## `ks pull`

fetch files for current environment. Write them locally.

```
USAGE
  $ ks pull

OPTIONS
  -f, --force  Overwrite any changes made locally

DESCRIPTION
  Once pulled files can be one of the three states  
     - updated : The file has been updated because someone else pushed a newer version
     - auto-merged : The file was modified and has been merged with someone else's changes
     - conflicted : The file has been modified and some lines are in conflict with someone else's changes. You should 
  fix the conflicts and push your changes

EXAMPLE
  $ ks pull
```

_See code: [src/commands/pull.js](https://github.com/wearedevx/keystone/blob/v1.0.0/src/commands/pull.js)_

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

_See code: [src/commands/push.js](https://github.com/wearedevx/keystone/blob/v1.0.0/src/commands/push.js)_

## `ks rm FILEPATHS`

Deletes one or more files.

```
USAGE
  $ ks rm FILEPATHS

ARGUMENTS
  FILEPATHS  Path to your file. Accepts a glob pattern

DESCRIPTION
  If you're an administrator or a contributor, the files will be removed for everyone.
  If you're a reader on the environment, you can't delete any files.

EXAMPLE
  $ ks rm path/to/file
```

_See code: [src/commands/rm.js](https://github.com/wearedevx/keystone/blob/v1.0.0/src/commands/rm.js)_

## `ks status`

shows the status of tracked files

```
USAGE
  $ ks status

EXAMPLE
  $ ks status
```

_See code: [src/commands/status.js](https://github.com/wearedevx/keystone/blob/v1.0.0/src/commands/status.js)_

## `ks token create ENV`

give access to your files with a non-blockstack user

```
USAGE
  $ ks token create ENV

ARGUMENTS
  ENV  environment you want the token to be created on

DESCRIPTION
  generate a token
  the token should be set in the system environment
  it allow the user to only run $ ks pull in order to pull locally files from the selected env

EXAMPLE
  $ ks token create ENV_NAME
```

_See code: [src/commands/token/create.js](https://github.com/wearedevx/keystone/blob/v1.0.0/src/commands/token/create.js)_

## `ks token revoke ENV`

revoke access to your files with a non-blockstack user

```
USAGE
  $ ks token revoke ENV

ARGUMENTS
  ENV  environment you want the token to be revoked on

EXAMPLE
  $ ks token revoke ENV_NAME
```

_See code: [src/commands/token/revoke.js](https://github.com/wearedevx/keystone/blob/v1.0.0/src/commands/token/revoke.js)_

## `ks whoami`

Shows the blockstack id of the currently logged in user

```
USAGE
  $ ks whoami

EXAMPLE
  $ ks whoami
```

_See code: [src/commands/whoami.js](https://github.com/wearedevx/keystone/blob/v1.0.0/src/commands/whoami.js)_
<!-- commandsstop -->
