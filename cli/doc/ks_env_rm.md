## ks env rm

Removes and environment

### Synopsis

Permanently removes an environment.

Every secret and tracked file content will be lost.
This is permanent and cannot be undone.

Example:
  $ ks env rm temp


```
ks env rm [flags]
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

* [ks env](ks_env.md)	 - Manage environments

###### Auto generated by spf13/cobra on 2-Dec-2020
