## ks secrets require

Marks a secret as required

### Synopsis

Marks a secret as required.

Secrets marked as required cannot be unset or set to blank value.
If they are, 'ks source' will exit with a non-zero exit code.


```
ks secrets require [flags]
```

### Options

```
  -h, --help   help for require
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/keystone.yaml)
      --env string      environment to use instead of the current one
  -q, --quiet           make the output machine readable
```

### SEE ALSO

* [ks secrets](ks_secrets.md)	 - Manage secrets

###### Auto generated by spf13/cobra on 2-Dec-2020
