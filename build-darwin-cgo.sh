#!/usr/bin/env bash

xgo --targets="darwin/386,darwin/amd64" $PWD
mv gotop-darwin-10.6-386 dist/darwin_386/gotop
mv gotop-darwin-10.6-amd64 dist/darwin_amd64/gotop
