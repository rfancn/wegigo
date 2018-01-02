#!/bin/bash

docker run --rm -v "$PWD":/go/src/github.com/rfancn/wegigo -w /go/src/github.com/rfancn/wegigo golang:1.8 make