## ks secrets

Manage secrets

### Synopsis

Manage secrets.

Used without arguments, displays a table of secrets.

```
ks secrets [flags]
```

### Options

```
  -h, --help   help for secrets
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/keystone.yaml)
      --env string      environment to use instead of the current one
  -q, --quiet           make the output machine readable
```

### SEE ALSO

* [ks](ks.md)	 - A safe system for developers to store, share and use secrets.
* [ks secrets add](ks_secrets_add.md)	 - Adds a secret to all environments
* [ks secrets optional](ks_secrets_optional.md)	 - Marks a secret as optional
* [ks secrets require](ks_secrets_require.md)	 - Marks a secret as required
* [ks secrets rm](ks_secrets_rm.md)	 - Removes a secret from all environments
* [ks secrets set](ks_secrets_set.md)	 - Updates a secret's value for the current environment
* [ks secrets unset](ks_secrets_unset.md)	 - Clears a secret for the current environment

###### Auto generated by spf13/cobra on 2-Dec-2020
