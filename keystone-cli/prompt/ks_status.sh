#!/bin/sh

keystone_info() {

    ENV=$(ks_prompt env)
    STATUS=$(ks_prompt status)
    s=" "
    if [[ -n $ENV ]]; then 
	s+="Ꝅ $ENV"
	s+=" $STATUS "
    fi
    echo "$s"
}

