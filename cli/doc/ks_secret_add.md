## ks secret add

Adds a secret to all environments

### Synopsis

Adds a secret to all environments.

Secrets are environment variables which value may vary
across environments, such as 'staging', 'prduction',
and 'development' environments.

The varible name will be added to all such environments,
you will be asked its value for each environment

Example:
  Add an environment variable PORT to all environments
  and set its value to 3000 for the current one.
  $ ks set PORT 3000



```
ks secret add [flags]
```

### Options

```
  -h, --help       help for add
  -o, --optional   mark the secret as optional
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/keystone.yaml)
      --env string      environment to use instead of the current one (default "dev")
  -q, --quiet           make the output machine readable
  -s, --skip            skip prompts and use default
```

### SEE ALSO

* [ks secret](ks_secret.md)	 - Manages secrets

###### Auto generated by spf13/cobra on 27-Jul-2021