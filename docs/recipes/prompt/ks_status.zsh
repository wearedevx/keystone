#!/bin/zsh

function keystone_info() {

    ENV=$(ks_prompt env)
    STATUS=$(ks_prompt status)
    s=" "
    if [[ -n $ENV ]]; then
        s+="Ꝅ%{$fg_bold[black]%} $ENV"
        s+="%{$fg_bold[white]%} $STATUS "
    fi
    echo "$s"
}