#!/usr/bin/env bash

VERSION=1.2.6

download() {
    archive=gotop_${VERSION}_${1}.tgz
    curl -LO https://github.com/cjbassi/gotop/releases/download/$VERSION/$archive
    tar xf $archive
    rm $archive
}

arch=$(uname -sm)
case "$arch" in
    Linux\ *64)  download linux_amd64    ;;
    Linux\ *86)  download linux_386      ;;
    Darwin\ *64) download darwin_amd64   ;;
    Darwin\ *86) download darwin_386     ;;
    *)
        echo "No binary found for your system"
        exit 1
        ;;
esac
