## ks file set

Updates a file’s content for the current environment

### Synopsis

Updates a file’s content for the current environment.

Changes the content of a file without altering other environments.
The local version of the file will be used.


```
ks file set <path to a file> [flags]
```

### Examples

```
ks file set ./config.php

# Change the content of ./config.php for the 'staging' environment:
ks --env staging file set ./config.php

```

### Options

```
  -h, --help   help for set
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

* [ks file](ks_file.md)	 - Manages secret files

###### Auto generated by spf13/cobra on 3-Jan-2022
