## ks env

Manage environments

### Synopsis

Manage environments.

Displays a list of available environments:
  $ ks env
    
   * default
     staging
	 prod

With an argument name, activates the anvironments:
  $ ks env staging


```
ks env [flags]
```

### Options

```
  -h, --help   help for env
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/keystone.yaml)
      --env string      environment to use instead of the current one
  -q, --quiet           make the output machine readable
```

### SEE ALSO

* [ks](ks.md)	 - A safe system for developers to store, share and use secrets.
* [ks env new](ks_env_new.md)	 - Creates a new environment
* [ks env rm](ks_env_rm.md)	 - Removes and environment

###### Auto generated by spf13/cobra on 2-Dec-2020
