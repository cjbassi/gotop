<div align="center">

<a href="./assets/logo">
    <img src="./assets/logo/logo.png" width="20%" />
</a>
<br><br>

Another terminal based graphical activity monitor, inspired by [gtop](https://github.com/aksakalli/gtop) and [vtop](https://github.com/MrRio/vtop), this time written in [Go](https://golang.org/)!

Join us in [\#gotop:matrix.org](https://riot.im/app/#/room/#gotop:matrix.org) ([matrix clients](https://matrix.to/#/#gotop:matrix.org)).

![](https://github.com/xxxserxxx/gotop/workflows/Build%20Go%20binaries/badge.svg)
![](https://github.com/xxxserxxx/gotop/workflows/Create%20pre-release/badge.svg)

![](https://raw.githubusercontent.com/xxxserxxx/gotop/master/docs/release.svg)

See the [mini-blog](https://github.com/xxxserxxx/gotop/wiki/Micro-Blog) for updates on the build status, and the [change log](/CHANGELOG.md) for release updates.

<img src="./assets/screenshots/demo.gif" />

</div>

## Installation

Working and tested on Linux, FreeBSD and MacOS. Windows binaries are provided, but have limited testing. OpenBSD works with some caveats; cross-compiling is difficult and binaries are not provided.

If you install gotop by hand, or you download or create new layouts or colorschemes, you will need to put the layout files where gotop can find them.  To see the list of directories gotop looks for files, run `gotop -h`.  The first directory is always the directory from which gotop is run.

-  **Arch**: Install from AUR, e.g. `yay -S gotop-bin`. There is also `gotop` and `gotop-git`
-  **Gentoo**: gotop is available on [guru](https://gitweb.gentoo.org/repo/proj/guru.git) overlay. 
    ```shell
    sudo layman -a guru
    sudo emerge gotop
    ```
- **OSX**: gotop is in *homebrew-core*.  `brew install gotop`.  Make sure to uninstall and untap any previous installations or taps.
- **Prebuilt binaries**: Binaries for most systems can be downloaded from [the github releases page](https://github.com/xxxserxxx/gotop/releases). RPM and DEB packages are also provided.
- **Source**: This requires Go >= 1.14. `go get -u github.com/xxxserxxx/gotop/cmd/gotop`

### Extension builds

An evolving mechanism in gotop are extensions. This is designed to allow gotop to support feature sets that are not universally needed without blowing up the application for average users with usused features.  Examples are support for specific hardware sets like video cards, or things that are just obviously not a core objective of the application, like remote server monitoring.

The path to these extensions is a tool called [gotop-builder](https://github.com/xxxserxxx/gotop-builder). It is easy to use and depends only on having Go installed.  You can read more about it on the project page, where you can also find binaries for Linux that have *all* extensions built in. If you want less than an all-inclusive build, or one for a different OS/architecture, you can use gotop-builder itself to create your own.


### Console Users

gotop requires a font that has braille and block character Unicode code points; some distributions do not provide this.  In the gotop repository is a `pcf` font that has these points, and setting this font may improve how gotop renders in your console.  To use this, run these commands:

```shell
curl -O -L https://raw.githubusercontent.com/xxxserxxx/gotop/master/fonts/Lat15-VGA16-braille.psf
setfont Lat15-VGA16-braille.psf
```

### Building

This is the download & compile approach.

gotop should build with most versions of Go.  If you have a version other than 1.14 installed, remove the `go` line at the end of `go.mod` and it should work.

```shell
git clone https://github.com/xxxserxxx/gotop.git
cd gotop
sed -i '/^go/d' go.mod          # Do this if you have go != 1.14
go build -o gotop ./cmd/gotop
```

Move `gotop` to somewhere in your `$PATH`.

If Go is not installed or is the wrong version, and you don't have root access or don't want to upgrade Go, a script is provided to download Go and the gotop sources, compile gotop, and then clean up. See `scripts/install_without_root.sh`.

## Usage

Run with `-h` to get an extensive list of command line arguments.  Many of these can be configured by creating a configuration file; see the next section for more information.  Key bindings can be viewed while gotop is running by pressing the `?` key, or they can be printed out by using the `--list keys` command.

In addition to the key bindings, the mouse can be used to control the process list:

- click to select process
- mouse wheel to scroll through processes

For more information on other topics, see:

- [Layouts](docs/layouts.md)
- [Configuration](docs/configuration.md)
- [Color schemes](docs/colorschemes.md)
- [Device filtering](docs/devices.md)
- [Extensions](docs/extensions.md)


## More screen shots

#### '-l kitchensink' + colorscheme
<img src="./assets/screenshots/kitchensink.gif" />

#### "-l battery"
<img src="./assets/screenshots/battery.png" />

#### "-l minimal"
<img src="./assets/screenshots/minimal.png" />

#### Custom (layouts/procs)
<img src="./assets/screenshots/procs.png" />

## Built With

- [gizak/termui](https://github.com/gizak/termui)
- [nsf/termbox](https://github.com/nsf/termbox-go)
- [exrook/drawille-go](https://github.com/exrook/drawille-go)
- [shirou/gopsutil](https://github.com/shirou/gopsutil)
- [goreleaser/nfpm](https://github.com/goreleaser/nfpm)
- [distatus/battery](https://github.com/distatus/battery)
- [VictoriaMetrics/metrics](https://github.com/VictoriaMetrics/metrics) Check this out! The API is clean, elegant, introduces many fewer indirect dependencies than the Prometheus client, and adds 50% less size to binaries.

## History

**ca. 2020-01-25** The original author of gotop started a new tool in Rust, called [ytop](https://github.com/cjbassi/ytop), and deprecated his Go version.  This repository is a fork of original gotop project with a new maintainer to keep the project alive and growing.  An objective of the fork is to maintain a small, focused core while providing a path to extend functionality for less universal use cases; examples of this is sensor support for NVidia graphics cards, and for aggregating data from remote gotop instances.

## Stargazers over time

[![Stargazers over time](https://starcharts.herokuapp.com/xxxserxxx/gotop.svg)](https://starcharts.herokuapp.com/xxxserxxx/gotop)
