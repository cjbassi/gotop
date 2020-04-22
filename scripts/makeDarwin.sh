#!/bin/bash

GOOS=darwin
GOARCH=amd64
CGO_ENABLED=1
MACOSX_DEPLOYMENT_TARGET=10.10.0
CC=o64-clang
CXX=o64-clang++

export GOOS GOARCH CGO_ENABLED MACOSX_DEPLOYMENT_TARGET CC CXX
go build -o gotop.darwin ./cmd/gotop 
