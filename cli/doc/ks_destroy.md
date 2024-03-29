## ks destroy

Destroys the whole Keystone project

### Synopsis

Destroys the whole Keystone project.

The project will be deleted, members won’t be able to send nor receive
updates about it. 

All secrets and files managed by Keystone *WILL BE LOST*.
It is highly recommended that you backup everything up beforehand.

This is irreversible.


```
ks destroy [flags]
```

### Examples

```
ks destroy
```

### Options

```
  -h, --help   help for destroy
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

* [ks](ks.md)	 - A safe system for developers to store, share and use secrets.

###### Auto generated by spf13/cobra on 3-Jan-2022
