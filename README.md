# Keystone

[![Coverage Status](https://coveralls.io/repos/github/wearedevx/keystone/badge.svg?branch=master)](https://coveralls.io/github/wearedevx/keystone?branch=master)
  
**Secrets synced and safe.**   

Sync your environment variables across team members, environments and codebase versions without leaving your terminal.

## Installation
### Linux
Use snap to install this package.
```sh
snap install keystone-cli
```
By default, you can call keystone command with `keystone-cli.ks`. You might want to create an alias using [snap aliases](https://snapcraft.io/docs/commands-and-aliases).
```sh
snap alias keystone-cli.ks ks
```

### macOS (via homebrew)
Install the Keystone tap
```sh
brew tap wearedevx/keystone
```

Install the latest stable version
```sh
brew install wearedevx/keystone/keystone
```

You can also install the development version with
```sh
brew install wearedevx/keystone/keystone-develop
```
And to update the development version
```sh
brew reinstall wearedevx/keystone/keystone-develop
```

### Usage
To start using Keystone you will need to login with [`ks login`](https://github.com/wearedevx/keystone/blob/master/cli/doc/ks_login.md), using your GitHub or GitLab account.
If your project is not keystone-managed yet, bootstrap it with [`ks init <YOUR_PROJECT_NAME>`](https://github.com/wearedevx/keystone/blob/master/cli/doc/ks_init.md).  
  
To start managing secrets and files, and access all of Keystone’s features, refer to the [complete CLI documentation](https://github.com/wearedevx/keystone/blob/master/cli/doc/ks.md)
