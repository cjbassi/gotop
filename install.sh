#!/usr/bin/env bash

VERSION=1.0.0

print_error() {
    echo "No binary found for your architecture. If your architecture is compatible with a binary"
    echo "that's already on GitHub, you can manually download and install it and open an issue"
    echo "saying so. Otherwise, create an issue requesting binaries to be build for your"
    echo "architecture and you can build from source in the meantime if you like."
}

install() {
    curl -L https://github.com/cjbassi/gotop/releases/download/$VERSION/gotop-${1}.tgz > /tmp/gotop.tgz
    tar xf /tmp/gotop.tgz -C /usr/bin
    rm /tmp/gotop.tgz
}

update() {
    cur_version=$(gotop -v 2>/dev/null)
    if [[ $? != 0 ]]; then
        download
    fi
    if (( "${cur_version//.}" < "${VERSION//.}" )); then
        download
    fi
}

arch=$(uname -sm)
case "$arch" in
    Linux\ x86_64)  install linux_amd64 ;;
    *)
        print_error
        exit 1
        ;;
esac
