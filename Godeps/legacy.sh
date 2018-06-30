#!/bin/bash

set -ex

git config http.proxy http://127.0.0.1:8123
git config https.proxy https://127.0.0.1:8123

export GOPATH=/temp

go get -u -d github.com/gomodule/redigo/redis
go get -u -d github.com/FZambia/sentinel
go get -u -d github.com/mna/redisc
go get -u -d github.com/fananchong/goredis

unset GOPATH

git config --unset http.proxy
git config --unset https.proxy
