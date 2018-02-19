#!/bin/bash

VERSION=v1.0

download() {
    curl -L https://github.com/cjbassi/gotop/releases/download/$VERSION/${1} > /usr/bin/gotop
    chmod +x /usr/bin/gotop
}

uninstall() {
    rm /usr/bin/gotop
}

for opt in "$@"; do
  case $opt in
    --uninstall)
      uninstall
      exit 0
      ;;
    *)
      echo "unknown option: $opt"
      exit 1
      ;;
  esac
done

arch=$(uname -sm)
case "$arch" in
  Linux\ *64)   download gotop-linux_amd64  ;;
esac
