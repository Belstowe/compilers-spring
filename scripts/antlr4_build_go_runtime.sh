#!/usr/bin/env bash

readonly antlr_ver=4.10.1

go get -d github.com/antlr/antlr4/runtime/Go/antlr@$antlr_ver

java org.antlr.v4.Tool -Dlanguage=Go build/grammar/RustLexer.g4