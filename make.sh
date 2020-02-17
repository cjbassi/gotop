#!/bin/bash

export VERSION=$(go run ./cmd/gotop -V)

rm -f build.log

function candz() {
	export GOOS=$1
	export GOARCH=$2
	OUT=build/gotop_${VERSION}_${GOOS}_${GOARCH}
	if [[ -n $3 ]]; then
		export GOARM=$3
		OUT=${OUT}_v${GOARM}
	fi
	OUT=${OUT}.zip
	if [[ -e $OUT ]]; then
		echo SKIP $OUT
		return
	fi
	D=build/gotop
	if [[ $GOOS == "windows" ]]; then
		D=${D}.exe
	fi
	go build -o $D ./cmd/gotop >> build.log 2>&1 
	unset GOOS GOARCH GOARM CGO_ENABLED
	if [[ $? -ne 0 ]]; then
		printenv | grep GO >> build.log
		echo "############### FAILED ###############" >> build.log
		echo >> build.log
		echo >> build.log
		echo FAILED $OUT
		return
	fi
	cd build
	zip $(basename $OUT) $(basename $D) >> ../build.log 2>&1
	cd ..
	rm -f $D
	if [[ $? -ne 0 ]]; then
		echo "############### FAILED ###############" >> build.log
		echo >> build.log
		echo >> build.log
		echo FAILED $OUT
		return
	fi
	echo BUILT $OUT
}

candz linux arm64
for x in 5 6 7; do
	candz linux arm $x
done
for x in 386 amd64; do
	candz linux $x

	sed -i "s/arch: .*/arch: \"${x}\"/" build/nfpm.yml
	for y in rpm deb; do
		OUT=build/gotop_${VERSION}_linux_${x}.${y}
		if [[ -e $OUT ]]; then
			echo SKIP $OUT
		else
			echo Building $OUT
			nfpm pkg -t ${OUT} -f build/nfpm.yml
		fi
	done

	candz windows $x
	candz freebsd $x

	export CGO_ENABLED=1
	candz darwin $x
	candz openbsd $x
	unset CGO_ENABLED
done

rm -f build/gotop
