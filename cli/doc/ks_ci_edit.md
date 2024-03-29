## ks ci edit

Configures an existing CI service integration

### Synopsis

Configures an existing CI service integration.

Use this command to modify CI service specific settings,
like API key and project name.

```
ks ci edit [ci service name] [flags]
```

### Examples

```
ks ci edit

# To avoid the prompt
ks ci edit my-gitub-ci-service
```

### Options

```
  -h, --help   help for edit
```

### Options inherited from parent commands

```
  -c, --config string   config file (default is $HOME/.config/keystone/keystone.yaml)
      --debug           debug output
      --env string      environment to use instead of the current one
  -q, --quiet           make the output machine readable
  -s, --skip            skip prompts and use default
```

### SEE ALSO

* [ks ci](ks_ci.md)	 - Manages CI services

###### Auto generated by spf13/cobra on 3-Jan-2022
