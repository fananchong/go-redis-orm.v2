#!/bin/bash

set -ex


# update
./legacy.sh

cp -f ./Godeps.json.template ./Godeps.json
cd ..
rm -rf ./vendor
docker run --rm -e GOPATH=/go/:/temp/ -v /temp/:/temp/ -v "$PWD":/go/src/github.com/fananchong/go-redis-orm.v2 -w /go/src/github.com/fananchong/go-redis-orm.v2 fananchong/godep save ./...

