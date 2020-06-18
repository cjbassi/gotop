#!/bin/bash

go test -covermode=count -coverprofile=docs/coverage.out ./...
go tool cover -html=docs/coverage.out -o docs/quality.html

revive -config .revive.toml ./... > docs/lint.txt

abcgo -path . -sort -no-test | sed 's%^/home/ser/workspace/gotop.d/gotop/%%' > docs/abcgo.txt
