#!/usr/bin/env bash

readonly antlr=antlr-4.10.1-complete.jar
readonly workdir=~/.local/share/applications

curl -O https://www.antlr.org/download/$antlr -o $workdir/$antlr
export CLASSPATH=".:$workdir/$antlr:$CLASSPATH"