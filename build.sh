#!/bin/bash

set -ex

docker run --rm -e GOBIN=/go/bin/ -v "$PWD"/bin:/go/bin/ -v "$PWD":/go/src/github.com/fananchong/go-redis-orm.v2 -w /go/src/github.com/fananchong/go-redis-orm.v2 golang go install ./...
