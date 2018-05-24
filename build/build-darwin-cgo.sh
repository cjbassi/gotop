#!/usr/bin/env bash

version=$(go run main.go -v)

xgo --targets="darwin/386,darwin/amd64" .

mv gotop-darwin-10.6-386 gotop
tar czf gotop_${version}_darwin_386.tgz gotop
rm -f gotop

mv gotop-darwin-10.6-amd64 gotop
tar czf gotop_${version}_darwin_amd64.tgz gotop
rm -f gotop
