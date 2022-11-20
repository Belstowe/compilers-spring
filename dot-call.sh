#!/usr/bin/env sh

./build/bin/llir-graph --input $1 --dot-callgraph | dot -Tsvg | tee $2
