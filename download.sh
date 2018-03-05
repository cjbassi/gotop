#!/usr/bin/env bash

VERSION=1.1.1

download() {
    curl -L https://github.com/cjbassi/gotop/releases/download/$VERSION/gotop-$VERSION-${1}.tgz > gotop.tgz
    tar xf gotop.tgz
    rm gotop.tgz
}

arch=$(uname -sm)
case "$arch" in
    Linux\ *64)  download linux_amd64    ;;
    Linux\ *86)  download linux_386      ;;
    Darwin\ *64) download darwin_amd64   ;;
    *)
        echo "No binary found for your system"
        exit 1
        ;;
esac
