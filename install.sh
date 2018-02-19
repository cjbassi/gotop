#!/bin/bash

VERSION=1.0.0

download() {
    arch=$(uname -sm)
    case "$arch" in
        Linux\ *64)   exe=gotop-linux_amd64  ;;
    esac
    curl -L https://github.com/cjbassi/gotop/releases/download/$VERSION/$exe > /usr/bin/gotop
    chmod +x /usr/bin/gotop
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

uninstall() {
    rm /usr/bin/gotop 2>/dev/null
}

for opt in "$@"; do
  case $opt in
    --uninstall)
      uninstall
      exit 0
      ;;
    --update)
      update
      exit 0
      ;;
    *)
      echo "unknown option: $opt"
      exit 1
      ;;
  esac
done

download
