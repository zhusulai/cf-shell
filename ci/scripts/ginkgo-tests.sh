#!/usr/bin/env bash
set -e

go get github.com/onsi/ginkgo/ginkgo
go install github.com/onsi/ginkgo/ginkgo

mkdir -p $GOPATH/src/github.com/zhusulai
cp -r cf-shell $GOPATH/src/github.com/zhusulai

cd $GOPATH/src/github.com/zhusulai/cf-shell/cfcli

ginkgo

