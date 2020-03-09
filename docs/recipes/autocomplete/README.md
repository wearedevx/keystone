## Autocomplete

Bash and zsh users can benefit from autocompletion with ks command.

<p align="center">
    <img src="autocomplete.gif" height="210"/>
</p>

#### Just run one of the two commands, depending on your shell.

```bash
# BASH
$ printf "$(ks autocomplete:script bash)" | awk '/KS_AC/ {print $0}' >> ~/.bashrc; source ~/.bashrc

# ZSH
$ printf "$(ks autocomplete:script zsh)" | awk '/KS_AC/ {print $0}' >> ~/.zshrc; source ~/.zshrc
```

#### Test it out, e.g.:

```bash
$ ks <TAB>                 # Command completion
$ ks command --<TAB>       # Flag completion
```
