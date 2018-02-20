# Maintainer: Caleb Bassi <calebjbassi@gmail.com>

pkgname=gotop
pkgver=1.0.1
pkgrel=1
pkgdesc="A TUI graphical activity monitor inspired by gtop"
arch=("x86_64" "i686")
url="https://github.com/cjbassi/gotop"
license=("AGPLv3")
provides=("gotop")

case "$CARCH" in
    x86_64)
        _arch=amd64
        ;;
    i686)
        _arch=386
        ;;
esac

source=("https://github.com/cjbassi/gotop/releases/download/$pkgver/gotop-$pkgver-linux_$_arch.tgz")
md5sums=("SKIP")

package() {
    mkdir -p "$pkgdir/usr/bin"
    mv $srcdir/gotop $pkgdir/usr/bin
}
