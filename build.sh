#!/bin/bash

set -ex

export GOBIN=$PWD/bin
go install -race ./tools/...
./bin/redis2go --input_dir=./example/redis_def --output_dir=./example --package=main
go install -race ./example/...

if [ "$1" = "publish" ]; then
    docker build -t redis2go .
    docker tag redis2go:latest fananchong/redis2go:latest
    docker push fananchong/redis2go:latest
fi
