# Keystone, the go version

## Installation
### Linux
Use snap to install this package.
```
$ snap install keystone-cli
```
By default, you can call keystone command with `keystone-cli.ks`. You might want to create an alias using [snap aliases](https://snapcraft.io/docs/commands-and-aliases).
```
$ snap alias keystone-cli.ks ks
```

### macOS (via homebrew)
Install the Keystone tap
```
$ brew tap wearedevx/keystone
```

Install the latest stable version
```
$ brew install wearedevx/keystone/keystone
```

You can also install the development version with
```
$ brew install wearedevx/keystone/keystone-development
```
And to update the development version
```
$ brew reinstall wearedevx/keystone/keystone-development
```

### Usage
To start using Keystone you will need to login with [`ks login`](https://github.com/wearedevx/keystone/blob/master/cli/doc/ks_login.md).  
If your project is not keystone-managed yet, bootstrap it with [`ks init <YOUR_PROJECT_NAME>`](https://github.com/wearedevx/keystone/blob/master/cil/doc/ks_init.md).  
  
To start managing secrets and files, and access all of Keystone’s features, refer to the [complete CLI documentation](https://github.com/wearedevx/keystone/blob/master/cli/doc/ks.md)
