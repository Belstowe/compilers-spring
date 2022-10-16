#!/usr/bin/env bash

./ruster -i "${1}" -o "$(basename $1).ll"
lli "$(basename $1).ll"
