local resetColor="%{$reset_color%}"
local prefix="⭓"
local dir="%{$fg_bold[black]%}%c$resetColor$resetColor"

PROMPT='%{$fg_bold[black]%}$prefix $(ssh_connection)$dir$(keystone_info)%{$fg[white]%}$(git_prompt_info)$resetColor'
MODE_INDICATOR="%{$fg_bold[white]%}%{$bg[black]%} NORMAL %{$reset_color%}"

function ssh_connection() {
    if [[ -n $SSH_CONNECTION ]]; then
        echo "%{$fg_bold[black]%}☰ "
            fi
}

ZSH_THEME_GIT_PROMPT_PREFIX="Ꝩ %{$fg_bold[black]%}"
ZSH_THEME_GIT_PROMPT_SUFFIX="$resetColor "
ZSH_THEME_GIT_PROMPT_DIRTY="%{$fg_bold[magenta]%} ⨯"
ZSH_THEME_GIT_PROMPT_CLEAN="%{$fg_bold[green]%} ✔"