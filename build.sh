#!/bin/bash

set -ex

BIN_DIR=$PWD/bin/
export GOBIN=$BIN_DIR
go install ./tools/...
./bin/redis2go --input_dir=./example/redis_def --output_dir=./example --package=main
go install ./example/...

if [ "$1" = "publish" ]; then
    go get github.com/bradfitz/goimports
    docker build -t redis2go .
    docker tag redis2go:latest fananchong/redis2go:latest
    docker push fananchong/redis2go:latest
fi
