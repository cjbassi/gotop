#!/usr/bin/env bash


get_latest_release_version() {
  curl --silent "https://api.github.com/repos/$1/releases/latest" | # Get latest release from GitHub api
    grep '"tag_name":' |                                            # Get tag line
    sed -E 's/.*"([^"]+)".*/\1/'                                    # Pluck JSON value
}


download() {
    archive=gotop_${version}_${1}.tgz
    curl -LO https://github.com/cjbassi/gotop/releases/download/$version/$archive
    tar xf $archive
    rm $archive
}


arch=$(uname -sm)
version=$(get_latest_release_version 'cjbassi/gotop')

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
        echo "\
No binary found for your system.
Feel free to request that we prebuild one that works on your system."
        exit 1
        ;;
esac
