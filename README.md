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

See the [mini-blog](/xxxserxxx/gotop/wiki/blog) for updates on the build status, and the [change log](/xxxserxxx/gotop/blob/master/CHANGELOG.md) for release updates.


<img src="./assets/screenshots/demo.gif" />
<img src="./assets/screenshots/kitchensink.gif" />

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

### Console Users Note

gotop requires a font that has braille and block character Unicode code points; some distributions do not provide this.  In the gotop repository is a `pcf` font that has these points, and setting this font may improve how gotop renders in your console.  To use this, run these commands:

```shell
$ curl -O -L https://raw.githubusercontent.com/xxxserxxx/gotop/master/fonts/Lat15-VGA16-braille.psf
$ setfont Lat15-VGA16-braille.psf
```

### Building

This is the download & compile approach.

gotop should build with most versions of Go.  If you have a version other than 1.14 installed, remove the `go` line at the end of `go.mod`.

```
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

## Config file

Most command-line settings can be persisted into a configuration file. The config file is named `gotop.conf` and can be located in several places. The first place gotop will look is in the current directory; after this, the locations depend on the OS and distribution. On Linux using XDG, for instance, the home location of `~/.config/gotop/gotop.conf` is the second location. The last location is a system-wide global location, such as `/etc/gotop/gotop.conf`. The `-h` help command will print out all of the locations, in order. Command-line options override values in any config files, and only the first config file found is loaded.

A configuration file can be created using the `--write-config` command-line argument. This will try to place the config file in the home config directory (the second location), but if it's unable to do so it'll write a file to the current directory.

Config file changes can be made by combining command-line arguments with `--write-config`. For example, to persist the `solarized` theme, call:

```
gotop -c solarized --write-config
```

### Colorschemes

gotop ships with a few colorschemes which can be set with the `-c` flag followed by the name of one. You can find all the colorschemes in the [colorschemes folder](./colorschemes).

To make a custom colorscheme, check out the [template](./colorschemes/template.go) for instructions and then use [default.json](./colorschemes/default.json) as a starter. Then put the file at `~/.config/gotop/<name>.json` and load it with `gotop -c <name>`. Colorschemes PR's are welcome!

To list all built-in color schemes, call:

```
gotop --list colorschemes
```

### Layouts

gotop can parse and render layouts from a specification file.  The format is
intentionally simple.  The amount of nesting levels is limited.  Some examples
are in the `layouts` directory; you can try each of these with, e.g.,
`gotop --layout-file layouts/procs`.  If you stick your layouts in
`$XDG_CONFIG_HOME/gotop`, you can reference them on the command line with the
`-l` argument, e.g. `gotop -l procs`.

The syntax for each widget in a row is:
```
(rowspan:)?widget(/weight)?
```
and these are separated by spaces.

1. Each line is a row
2. Empty lines are skipped
3. Spaces are compressed (so you can do limited visual formatting)
4. Legal widget names are: cpu, disk, mem, temp, batt, net, procs
5. Widget names are not case sensitive
4. The simplest row is a single widget, by name, e.g. `cpu`
5. **Weights**
    1. Widgets with no weights have a weight of 1.
    2. If multiple widgets are put on a row with no weights, they will all have
       the same width.
    3. Weights are integers
    4. A widget will have a width proportional to its weight divided by the
       total weight count of the row. E.g.,

       ```
       cpu      net
       disk/2   mem/4
       ```

       The first row will have two widgets: the CPU and network widgets; each
       will be 50% of the total width wide.  The second row will have two
       widgets: disk and memory; the first will be 2/6 ~= 33% wide, and the
       second will be 5/7 ~= 67% wide (or, memory will be twice as wide as disk).
9.  If prefixed by a number and colon, the widget will span that number of
    rows downward. E.g.

    ```
    mem   2:cpu
    net
    ```

    Here, memory and network will be in the same row as CPU, one over the other,
    and each half as high as CPU; it'll look like this:

    ```
     +------+------+
     | Mem  |      |
     +------+ CPU  |
     | Net  |      |
     +------+------+
    ```
     
10. Negative, 0, or non-integer weights will be recorded as "1".  Same for row
    spans. 
11. Unrecognized widget names will cause the application to abort.                          
12. In rows with multi-row spanning widgets **and** weights, weights in
    lower rows are ignored.  Put the weight on the widgets in that row, not
    in later (spanned) rows.
13. Widgets are filled in top down, left-to-right order.
14. The larges row span in a row defines the top-level row span; all smaller
    row spans constitude sub-rows in the row. For example, `cpu mem/3 net/5`
    means that net/5 will be 5 rows tall overall, and mem will compose 3 of
    them. If following rows do not have enough widgets to fill the gaps,
    spacers will be used.

Yes, you're clever enough to break the layout algorithm, but if you try to
build massive edifices, you're in for disappointment.

To list all built-in color schemes, call:

```
gotop --list layouts
```

### Device filtering

Some devices have quite a number of data points; on OSX, for instance, there are dozens of temperature readings. These can be filtered through a configuration file.  There is no command-line argument for this filter.

The list will grow, but for now the only device that supports filtering is the temperature widget.  The configuration entry is called `temperature`, and it contains an exact-match list of comma-separated values with no spaces.  To see the list of valid values, run gotop with the `--list devices` command.  Gotop will print out the type of device and the legal values.  For example, on Linux:

```
$ gotop --list devices
Temperatures:
        acpitz
        nvme_composite
        nvme_sensor1
        nvme_sensor2
        pch_cannonlake
        coretemp_packageid0
        coretemp_core0
        coretemp_core1
        coretemp_core2
        coretemp_core3
        ath10k_hwmon
```
You might then add the following line to the config file.  First, find where gotop looks for config files:
```
$ gotop -h | tail -n 6
Colorschemes & layouts that are not built-in are searched for (in order) in:
/home/USER/workspace/gotop.d/gotop, /home/USER/.config/gotop, /etc/xdg/gotop
The first path in this list is always the cwd. The config file
'gotop.config' can also reside in one of these directories.

Log files are stored in /home/ser/.cache/gotop/errors.log
```
So you might use `/home/YOU/.config/gotop.conf`, and add (or modify) this line:
```
temperatures=acpitz,coretemp_core0,ath10k_hwmon
```
This will cause the temp widget to show only four of the eleven temps.

### CLI Options

Run `gotop -h` to see the list of all command line options.

## More screen shots

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

The original author of gotop started a new tool in Rust, called [ytop](https://github.com/cjbassi/ytop).  This repository is a fork of original gotop project with a new maintainer.

## Stargazers over time

[![Stargazers over time](https://starcharts.herokuapp.com/xxxserxxx/gotop.svg)](https://starcharts.herokuapp.com/xxxserxxx/gotop)
