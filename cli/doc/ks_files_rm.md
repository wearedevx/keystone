## ks file rm

Removes a file from secrets

### Synopsis

Removes a file from secrets.

The file will no longer be gitignored and its content
will no longer be updated when changing environment.

The content of the file for other environments will be lost.
This is permanent, and cannot be undone.

Example:
  $ ks file rm config/old-test-config.php

```
ks file rm [flags]
```

### Options

```
  -h, --help   help for rm
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/keystone.yaml)
      --env string      environment to use instead of the current one
  -q, --quiet           make the output machine readable
```

### SEE ALSO

* [ks file](ks_files.md)	 - Manage secret files

###### Auto generated by spf13/cobra on 2-Dec-2020