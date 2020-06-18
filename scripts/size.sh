#!/bin/bash
#
# size.sh is used to bisect the repository and find changes that negatively
# impacted the gotop binary size. It does this by building gotop and exiting
# successfully if the binary size is under a defined amount.
# 
# Example:
# ```
# git bisect start
# git bisect bad master
# git bisect good 755037d211cc8e58e9ce43ee74a95a3036053dee
# git bisect run ./size
# ```

GOODSIZE=6000000

# Caleb's directory structure was different from the current structure, so
# we have to find the main package first.
pt=$(dirname $(find . -name main.go))
# Give the executable a unique-ish name
fn=gotop_$(git rev-list -1 HEAD)
go build -o $fn $pt
sz=$(ls -l $fn | awk '{print $5}')
git checkout -- .
[[ $sz -gt $GOODSIZE ]] && exit 1
exit 0
