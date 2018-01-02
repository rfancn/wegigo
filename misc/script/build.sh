#!/bin/bash

SRC_DIR=/go/src/github.com/rfancn/wegigo
docker run --rm -v "$PWD":$SRC_DIR -w $SRC_DIR golang:1.8 make