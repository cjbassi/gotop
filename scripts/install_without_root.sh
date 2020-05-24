#!/bin/sh
#
# Builds the gotop executable on machines where Go isn't installed,
# or is the wrong version, and you don't have root access to upgrade or
# install Go.
#
# You can run this without cloning the entire gotop repository (the script
# will do this for you.)

set -x

VERSION='1.14.2'     # Go version needed to build
OS='linux'
ARCH='amd64'
BUILDDIR=/tmp/gotop-build
INSTALLDIR=${HOME}/bin

GO_NAME=go${VERSION}.${OS}-${ARCH}

mkdir -p $BUILDDIR
cd $BUILDDIR

curl https://dl.google.com/go/${GO_NAME}.tar.gz --output ./${GO_NAME}.tar.gz

tar -vxzf ${GO_NAME}.tar.gz
rm ${GO_NAME}.tar.gz

PATH=$BUILDDIR/go/bin:$PATH

go env -w GOPATH=$BUILDDIR # otherwise go would create a directory in $HOME

rm -rf ./gotop
git clone https://github.com/xxxserxxx/gotop.git
cd ./gotop
go build -o gotop ./cmd/gotop

go clean -modcache # otherwise $BUILDDIR/pkg would need sudo permissions to remove

mkdir -p $INSTALLDIR
mv gotop ${INSTALLDIR}/gotop

rm -rf $BUILDDIR

printf "gotop installed in ${INSTALLDIR}/gotop\n"

