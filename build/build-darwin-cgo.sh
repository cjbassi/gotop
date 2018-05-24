#!/usr/bin/env bash

version=$(go run main.go -v)

xgo --targets="darwin/386,darwin/amd64" .
tar czf gotop_$version_darwin_386.tgz gotop-darwin-10.6-386
tar czf gotop_$version_darwin_amd64.tgz gotop-darwin-10.6-amd64
