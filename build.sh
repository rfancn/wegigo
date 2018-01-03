#!/bin/bash

BASE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

SRC_DIR=/go/src/github.com/rfancn/wegigo
docker run --rm -v "$BASE_DIR":$SRC_DIR -w $SRC_DIR golang:1.8 make