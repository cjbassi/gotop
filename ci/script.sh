#!/usr/bin/env bash

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
    VERSION=$(go run main.go -v) # used by nfpm
    docker run --rm \
        -v $PWD:/tmp/pkg \
        -e VERSION=${VERSION} \
        goreleaser/nfpm pkg \
            --config /tmp/pkg/ci/nfpm.yml \
            --target /tmp/pkg/dist/${FILE}.deb || ERROR=true
    docker run --rm \
        -v $PWD:/tmp/pkg \
        -e VERSION=${VERSION} \
        goreleaser/nfpm pkg \
            --config /tmp/pkg/ci/nfpm.yml \
            --target /tmp/pkg/dist/${FILE}.rpm || ERROR=true
fi

if [ ${ERROR} == "true" ]; then
    exit 1
fi
