#!/usr/bin/env bash

VERSION=1.0.0

arch=$(uname -sm)

case "$arch" in
    Linux\ *64)   arch=linux_amd64   ;;
esac

curl -L https://github.com/cjbassi/gotop/releases/download/$VERSION/gotop-$arch.tgz > /tmp/gotop.tgz
tar xf /tmp/gotop.tgz -C /usr/bin
rm /tmp/gotop.tgz

update() {
    cur_version=$(gotop -v 2>/dev/null)
    if [[ $? != 0 ]]; then
        download
    fi
    if (( "${cur_version//.}" < "${VERSION//.}" )); then
        download
    fi
}
