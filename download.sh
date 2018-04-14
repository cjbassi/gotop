#!/usr/bin/env bash

VERSION=1.2.12

download() {
    archive=gotop_${VERSION}_${1}.tgz
    curl -LO https://github.com/cjbassi/gotop/releases/download/$VERSION/$archive
    tar xf $archive
    rm $archive
}

arch=$(uname -sm)
case "$arch" in
    # order matters
    Darwin\ *64)        download darwin_amd64   ;;
    Darwin\ *86)        download darwin_386     ;;
    Linux\ armv5*)      download linux_arm5     ;;
    Linux\ armv6*)      download linux_arm6     ;;
    Linux\ armv7*)      download linux_arm7     ;;
    Linux\ armv8*)      download linux_arm8     ;;
    Linux\ aarch64*)    download linux_arm8     ;;
    Linux\ *64)         download linux_amd64    ;;
    Linux\ *86)         download linux_386      ;;
    *)
        echo "No binary found for your system"
        exit 1
        ;;
esac
