#!/usr/bin/env sh

./build/bin/llir-graph --input $1 --dot-def-use | dot -Tsvg | tee $2
