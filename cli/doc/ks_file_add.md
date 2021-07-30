## ks file add

Adds a file to secrets

### Synopsis

Adds a file to secrets

A secret file is a file which have content that can changge
across environments, such as configuration files, credentials,
certificates and so on.

When adding a file, you will be asked for a version of its content
for all known environments.

Examples:
  $ ks file add ./config/config.exs
  
  $ ks file add ./wp-config.php

  $ ks file add ./certs/my-website.cert


```
ks file add [flags]
```

### Options

```
  -h, --help   help for add
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/keystone.yaml)
      --env string      environment to use instead of the current one (default "dev")
  -q, --quiet           make the output machine readable
  -s, --skip            skip prompts and use default
```

### SEE ALSO

* [ks file](ks_file.md)	 - Manages secret files

###### Auto generated by spf13/cobra on 27-Jul-2021