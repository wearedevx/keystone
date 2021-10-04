## ks ci update

Updates a CI service integration

### Synopsis

Updates CI service integration.

Use this command to modify CI service specific settings
like API key and project name.

```
ks ci update [ci service name] [flags]
```

### Examples

```
ks ci update

# To avoid the prompt
ks ci update my-gitub-ci-service
```

### Options

```
  -h, --help   help for update
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/keystone/keystone.yaml)
      --env string      environment to use instead of the current one (default "dev")
  -q, --quiet           make the output machine readable
  -s, --skip            skip prompts and use default
```

### SEE ALSO

* [ks ci](ks_ci.md)	 - Manages CI services

###### Auto generated by spf13/cobra on 7-Sep-2021