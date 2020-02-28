#!/bin/bash

export VERSION=$(go run ./cmd/gotop -V)

rm -f build.log

# Set up some common environment variables
function out() {
	export GOOS=$1
	export GOARCH=$2
	OUT=build/gotop_${VERSION}_${GOOS}_${GOARCH}
	if [[ -n $3 ]]; then
		export GOARM=$3
		OUT=${OUT}_v${GOARM}
	fi
	D=build/gotop
	if [[ $GOOS == "windows" ]]; then
		D=${D}.exe
		AR=${OUT}.zip
	else
		AR=${OUT}.tgz
	fi
}
# Clean up environment variables
function uset() {
	unset GOOS GOARCH GOARM CGO_ENABLED OUT D AR
}
# Compile something that can be cross-compiled without docker
function c() {
	out $1 $2 $3
	go build -o $D ./cmd/gotop >> build.log 2>&1 
	if [[ $? -ne 0 ]]; then
		printenv | grep GO >> build.log
		echo "############### FAILED ###############"
		echo FAILED COMPILE $OUT
	fi
}
# Zip up something that's been compiled
function z() {
	out $1 $2 $3
	if [[ -e $AR ]]; then
		echo SKIP $AR
		return
	fi
	cd build
	if [[ $GOOS == "windows" ]]; then
		zip -q $(basename $AR) $(basename $D)
	else
		tar -czf $(basename $AR) $(basename $D)
	fi
	if [[ $? -ne 0 ]]; then
		echo "############### FAILED ###############"
		echo FAILED ARCHIVE $AR
		cd ..
		return
	fi
	cd ..
	echo BUILT $AR
}
# Do c(), then z(), and then clean up.
function candz() {
	unset OUT
	out $1 $2 $3
	local AR=${OUT}.zip
	if [[ -e $AR ]]; then
		echo SKIP $AR
		return
	fi
	c $1 $2 $3
	z $1 $2 $3
	rm -f $D
	uset
}
# Build the deb and rpm archives
function nfpmAR() {
	out $1 $2 $3
	sed -i "s/arch: .*/arch: \"${GOARCH}\"/" build/nfpm.yml
	for y in rpm deb; do
		local AR=build/gotop_${VERSION}_linux_${GOARCH}.${y}
		if [[ -e $AR ]]; then
			echo SKIP $AR
		else
			echo Building $AR
			nfpm pkg -t ${AR} -f build/nfpm.yml
		fi
	done
}
# Cross-compile the darwin executable.  This requires docker.
function cdarwinz() {
	if [[ ! -d darwin ]]; then
		git clone . darwin
		cd darwin
	else
		cd darwin
		git checkout -- .
		git pull
	fi
	export CGO_ENABLED=1
	if [[ `docker images -q gotopdarwin` == "" ]]; then
		docker run -it --name osxcross --mount type=bind,source="$(pwd)",target=/mnt dockercore/golang-cross /bin/bash /mnt/build/osx.sh
		local ID=`docker commit osxcross dockercore/golang-cross | cut -d : -f 2`
		docker tag $ID gotopdarwin
	else
		docker run -it --name osxcross --mount type=bind,source="$(pwd)",target=/mnt gotopdarwin /bin/bash /mnt/build/osx.sh
	fi
	cd ..
	mv darwin/build/gotop*darwin*.tgz build/
	docker rm osxcross
	uset
}

##########################################################
# Iterate through the OS/ARCH permutations and build an
# archive for each, using the previously defined functions.
if [[ $1 == "darwin" ]]; then
	export MACOSX_DEPLOYMENT_TARGET=10.10.0 
	export CC=o64-clang 
	export CXX=o64-clang++ 
	export CGO_ENABLED=1 
	c darwin 386
	z darwin 386
	c darwin amd64
	z darwin amd64
else
	candz linux arm64

	for x in 5 6 7; do
		candz linux arm $x
	done

	for x in 386 amd64; do
		c linux $x
		z linux $x
		nfpmAR linux $x
		rm -f $D
		uset

		candz windows $x
		candz freebsd $x

		# TODO Preliminary OpenBSD support [#112] [#117] [#118]
		# candz openbsd $x
	done
	cdarwinz
fi
