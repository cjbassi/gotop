#!/usr/bin/env bash

VERSION=$(go run gotop.go -v)

env GOOS=darwin GOARCH=amd64 go build
tar czf gotop-$VERSION-darwin_amd64.tgz gotop
rm gotop

env GOOS=linux GOARCH=386 go build
tar czf gotop-$VERSION-linux_386.tgz gotop
rm gotop

env GOOS=linux GOARCH=amd64 go build
tar czf gotop-$VERSION-linux_amd64.tgz gotop
rm gotop
