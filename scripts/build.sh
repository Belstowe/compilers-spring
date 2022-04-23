#!/usr/bin/env bash

readonly antlr_ver=4.10.1
readonly antlr=antlr-$antlr_ver-complete.jar
readonly antlr_path=$PWD/$antlr
readonly output=$PWD/parser

mkdir -p $output

if [ ! -f "$antlr_path" ]; then
    curl -O https://www.antlr.org/download/$antlr
fi

go get -d github.com/antlr/antlr4/runtime/Go/antlr@$antlr_ver
(cd build/grammar && java -jar $antlr_path -o $output -visitor -no-listener -Dlanguage=Go *.g4)