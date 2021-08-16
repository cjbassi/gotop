<div align="center">

<a href="./assets/logo">
    <img src="./assets/logo/logo.png" width="20%" />
</a>
<br><br>

Another terminal based graphical activity monitor, inspired by [gtop](https://github.com/aksakalli/gtop) and [vtop](https://github.com/MrRio/vtop), this time written in [Go](https://golang.org/)!

Join us in [\#gotop:matrix.org](https://app.element.io/#/room/#gotop:matrix.org) ![](https://img.shields.io/matrix/gotop:matrix.org) ([matrix clients](https://matrix.to/#/#gotop:matrix.org)).

![](https://github.com/xxxserxxx/gotop/workflows/Build%20Go%20binaries/badge.svg)
![](https://img.shields.io/github/v/release/xxxserxxx/gotop)
![](https://img.shields.io/github/release-date/xxxserxxx/gotop)
![](https://img.shields.io/librariesio/github/xxxserxxx/gotop)
![](https://img.shields.io/github/commit-activity/m/xxxserxxx/gotop)

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
- **Source**: This requires Go >= 1.16. `go install github.com/xxxserxxx/gotop/v4/cmd/gotop@latest`

### Extensions update

Extensions have proven problematic; go plugins are not usable in real-world cases, and the solution I had running for a while was hacky, at best. Consequently, extensions have been moved into the main code base for now.

-  nvidia support: use the `--nvidia` flag to enable. You must have the `nvidia- smi` package installed, and gotop must be able to find the `nvidia-smi` executable, for this to work.
-  remote: allows gotop to pull sensor data from applications exporting Prometheus metrics, including remote gotop instances themselves.

### Console Users

gotop requires a font that has braille and block character Unicode code points; some distributions do not provide this.  In the gotop repository is a `pcf` font that has these points, and setting this font may improve how gotop renders in your console.  To use this, run these commands:

```shell
curl -O -L https://raw.githubusercontent.com/xxxserxxx/gotop/master/fonts/Lat15-VGA16-braille.psf
setfont Lat15-VGA16-braille.psf
```

### Building

This is the download & compile approach.

gotop requires Go 1.16 or later to build, as it relies on the embed feature released with 1.16; a library it uses, lingo, uses both embed and the `io/fs` package.  For a version of gotop that builds with earlier versions, check out one of the tags prior to v4.2.0.

```shell
git clone https://github.com/xxxserxxx/gotop.git
cd gotop
# This ugly SOB gets a usable version from the git tag list
VERS="$(git tag -l --sort=-v:refname | sed 's/v\([^-].*\)/\1/g' | head -1 | tr -d '-' ).$(git describe --long --tags | sed 's/\([^-].*\)-\([0-9]*\)-\(g.*\)/r\2.\3/g' | tr -d '-')"
DAT=$(date +%Y%m%dT%H%M%S)
go build -o gotop \
	-ldflags "-X main.Version=v${VERS} -X main.BuildDate=${DAT}" \
	./cmd/gotop
```

If you want to compact the executable as much as possible on Linux, change the `ldflags` line to this:

```
-ldflags "-X main.Version=v${VERS} -X main.BuildDate=${DAT} -extldflags '-s -w'" \
```

Now move the `gotop` executable to somewhere in your `$PATH`.

If Go is not installed or is the wrong version, and you don't have root access or don't want to upgrade Go, a script is provided to download Go and the gotop sources, compile gotop, and then clean up. See `scripts/install_without_root.sh`.

#### go generate

With Go 1.16, it is no longer necessary to call `go generate`. Translations and Apple SMC tags are embedded with `go:embed`.

## Usage

Run with `-h` to get an extensive list of command line arguments.  Many of these can be configured by creating a configuration file; see the next section for more information.  Key bindings can be viewed while gotop is running by pressing the `?` key, or they can be printed out by using the `--list keys` command.

In addition to the key bindings, the mouse can be used to control the process list:

- click to select process
- mouse wheel to scroll through processes

For more information on other topics, see:

- [Layouts](https://github.com/xxxserxxx/gotop/blob/master/docs/layouts.md)
- [Configuration](https://github.com/xxxserxxx/gotop/blob/master/docs/configuration.md)
- [Color schemes](https://github.com/xxxserxxx/gotop/blob/master/docs/colorschemes.md)
- [Device filtering](https://github.com/xxxserxxx/gotop/blob/master/docs/devices.md)
- [Extensions](https://github.com/xxxserxxx/gotop/blob/master/docs/extensions.md)

Monitoring remote machines
--------------------------

gotop can monitor gotops running on remote machines and display (some of the)
metrics within a single instance. gotop expects to be behind a proxy, or within
a secure intranet, so while it's not exactly hard to set up, it's also not
trivial. An example set-up is explained in the 
[Remote Monitoring](https://github.com/xxxserxxx/gotop/blob/master/docs/remote-monitoring.md) 
document.

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
- [lingo](https://github.com/xxxserxxx/lingo) is forked from [jdkeke142's](https://github.com/jdkeke142/lingo-toml) lingo, which was in turn forked from [kortemy's](https://github.com/kortemy/lingo) original project.


## History

**ca. 2020-01-25** The original author of gotop started a new tool in Rust, called [ytop](https://github.com/cjbassi/ytop), and deprecated his Go version.  This repository is a fork of original gotop project with a new maintainer to keep the project alive and growing.  An objective of the fork is to maintain a small, focused core while providing a path to extend functionality for less universal use cases; examples of this is sensor support for NVidia graphics cards, and for aggregating data from remote gotop instances.

## Alternatives

I obviously think gotop is the Bee's Knees, but there are many alternatives. Many of these have been around for years. All of them are terminal-based tools.

- Grandpa [top](http://sourceforge.net/projects/unixtop/). Written 36 years ago, C, installed by default on almost every Unix descendant.
- [bashtop](https://github.com/aristocratos/bashtop), in pure bash! Beautiful and space efficient, and [deserves special comment](docs/bashtop.md).
- [bpytop](https://github.com/aristocratos/bpytop), @aristocratos, the author of bashtop, rewrote it in Python in mid-2020; it's the same beautiful interface, and a very nice alternative.
- [htop](https://hisham.hm/htop/). A prettier top. Similar functionality. 16 years old!
- [atop](https://www.atoptool.nl/). Detailed process-focused inspection with a table-like view. Been around for 9 long years.
- [iftop](http://www.ex-parrot.com/~pdw/iftop/), a top for network connections.  More than just data transfer, iftop will show what interfaces are connecting to what IP addresses. Requires root access to run.
- [iotop](http://guichaz.free.fr/iotop/), top for disk access. Tells you *which* processes are writing to and from disk space, and how much. Also requires root access to run.
- [nmon](http://nmon.sourceforge.net) a dashboard style top; widgets can be dynamically enabled and disabled, pure ASCII rendering, so it doesn't rely on fancy character sets to draw bars. 
- [ytop](https://github.com/cjbassi/ytop), a rewrite of gotop (ca. 3.0) in Rust.  Same great UI, different programming language.
- [slabtop](https://gitlab.com/procps-ng/procps), part of procps-ng, looks like top but provides kernel slab cache information! Requires root.
- [systemd-cgtop](https://www.github.com/systemd/systemd), comes with systemd (odds are your system uses systemd, so this is already installed), provides a resource use view of control groups -- basically, which services are using what resources. Does *not* require root to run.
- [virt-top](https://libvirt.org/) top for virtualized containers (VMs, like QEMU).
- [ctop](https://bcicen.github.io/ctop/) top for containers (LXC, like docker)


### A comment on clones

In a chat room I heard someone refer to gotop as "another one of those fancy language rewrites people do."  I'm not the original author of gotop, so it's easy to not take offense, but I'm going on record as saying that I disagree with that sentiment: I think these rewrites are valuable, useful, and healthy to the community. They increase software diversity at very little [cost to users](https://en.wikipedia.org/wiki/Information_overload), and are a sort of evolutionary mechanism: as people do rewrites, some are worse, but some are better, and users benefit.  Rewrites provide options, which fight against [monocultures](https://github.com). As importantly, most developers are really only fluent in a couple of programming languages. We all have *familiarity* with a dozen, and may even have extensive experience with a half-dozen, but if you don't constantly use a language, you tend to forget the extended library APIs, your development environment isn't tuned, you're rusty with using the tool sets, and you may have forgotten a lot of the language peculiarities and gotchas. The barrier to entry for contributing to a software project -- to simply finding and fixing a bug -- in a language you're not intimate with can be very high. It gets much worse if the project owner is a stickler for a particular style.  So I believe that gotop's original author's decision to rewrite his project in Rust is a net positive. He probably made fewer design mistakes in ytop (we always do, on the second rewrite), and Rust developers -- who may have hesitated learning or brushing up on Go to submit an improvement -- have another project to which they can contribute.

Diversity is good. Don't knock the free stuff.


## Stargazers over time

[![Stargazers over time](https://starcharts.herokuapp.com/xxxserxxx/gotop.svg)](https://starcharts.herokuapp.com/xxxserxxx/gotop)
