#!/usr/bin/env bash

VERSION=1.0.1

install() {
    curl -L https://github.com/cjbassi/gotop/releases/download/$VERSION/gotop-$VERSION-${1}.tgz > gotop.tgz
    tar xf gotop.tgz
    rm gotop.tgz
}

arch=$(uname -sm)
case "$arch" in
    Linux\ *64)  install linux_amd64    ;;
    Linux\ *86)  install linux_386      ;;
    Darwin\ *64) install darwin_amd64   ;;
    *)
        echo "No binary found for your system"
        exit 1
        ;;
esac
