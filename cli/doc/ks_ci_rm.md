## ks ci rm

Removes a CI service configuration

### Synopsis

Removes a CI service configuration.

`ks ci send` will no longer send secrets and files to the service.
However, secrets and files sent before calling `ks ci send` will
not be cleaned from the service.

```
ks ci rm [service name] [flags]
```

### Options

```
  -h, --help   help for rm
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/keystone.yaml)
      --env string      environment to use instead of the current one (default "dev")
  -q, --quiet           make the output machine readable
  -s, --skip            skip prompts and use default
```

### SEE ALSO

* [ks ci](ks_ci.md)	 - Manages CI services

###### Auto generated by spf13/cobra on 12-Aug-2021