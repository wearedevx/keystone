## ks member set-role

Sets the role for a member

### Synopsis

Sets the role for a member.
If no role argument is provided, it will be prompted.

Roles determine access rights to environments.

```
ks member set-role <member id> [role] [flags]
```

### Examples

```
# Set the role directly
ks member set-role john@gitlab devops

# Set the role with a prompt
ks member set-role sandra@github
```

### Options

```
  -h, --help   help for set-role
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

* [ks member](ks_member.md)	 - Manages members

###### Auto generated by spf13/cobra on 3-Jan-2022
