#!/usr/bin/env bash

VERSION=1.0.0

arch=$(uname -sm)

case "$arch" in
    Linux\ *64)   exe=linux_amd64   ;;
esac

curl -L https://github.com/cjbassi/gotop/releases/download/$VERSION/gotop-$exe > /usr/bin/gotop
chmod +x /usr/bin/gotop

update() {
    cur_version=$(gotop -v 2>/dev/null)
    if [[ $? != 0 ]]; then
        download
    fi
    if (( "${cur_version//.}" < "${VERSION//.}" )); then
        download
    fi
}
