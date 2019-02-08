#!/usr/bin/env bash

function main {
    # Check if any command failed
    ERROR=false

    GOARCH=${_GOARCH}
    GOOS=${_GOOS}

    if [[ ! ${GOARCH} ]]; then
        exit
    fi

    env GOOS=${GOOS} GOARCH=${GOARCH} GOARM=${GOARM} go build -o ${NAME} || ERROR=true

    mkdir -p dist

    if [[ ${GOARCH} == "arm64" ]]; then
        FILE=${NAME}_${TRAVIS_BRANCH}_${GOOS}_arm8
    else
        FILE=${NAME}_${TRAVIS_BRANCH}_${GOOS}_${GOARCH}${GOARM}
    fi

    tar -czf dist/${FILE}.tgz ${NAME} || ERROR=true

    if [[ ${GOOS} == "linux" && ${GOARCH} == "amd64" ]]; then
        make all || ERROR=true
        rm dist/gotop
    fi

    if [ ${ERROR} == "true" ]; then
        exit 1
    fi
}

main
