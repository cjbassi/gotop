#!/bin/bash

VERSION=v1.0

download() {
    curl -L https://github.com/cjbassi/gotop/releases/download/$VERSION/${1} > /usr/bin/gotop
    chmod +x /usr/bin/gotop
}

arch=$(uname -sm)
case "$arch" in
  Linux\ *64)   download gotop-linux_amd64  ;;
esac
