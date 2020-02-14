<div align="center">

<a href="./assets/logo">
    <img src="./assets/logo/logo.png" width="20%" />
</a>
<br><br>

Another terminal based graphical activity monitor, inspired by [gtop](https://github.com/aksakalli/gtop) and [vtop](https://github.com/MrRio/vtop), this time written in [Go](https://golang.org/)!

The original author of gotop has re-implemented the application in Rust, as [ytop](https://github.com/cjbassi/ytop).  This is a fork of original gotop project with a new maintainer.

<img src="./assets/demos/demo.gif" />
<img src="./assets/screenshots/minimal.png" width="96%" />

</div>

## Installation

Working and tested on Linux, FreeBSD and macOS. Windows support is planned. OpenBSD works with some caveats.

### Source

```bash
go get github.com/xxxserxxx/gotop/...
```

### Prebuilt binaries

**Note**: Doesn't require Go.

Visit [here](https://github.com/xxxserxxx/gotop/releases) with your web browser and download a version that works for you.

Unzip it and then move `gotop` into your `$PATH` somewhere.  If you're on a Debian or Redhat derivative, you can download an `.rpm` or `.deb` to install.

## Usage

### Keybinds

- Quit: `q` or `<C-c>`
- Process navigation:
  - `k` and `<Up>`: up
  - `j` and `<Down>`: down
  - `<C-u>`: half page up
  - `<C-d>`: half page down
  - `<C-b>`: full page up
  - `<C-f>`: full page down
  - `gg` and `<Home>`: jump to top
  - `G` and `<End>`: jump to bottom
- Process actions:
  - `<Tab>`: toggle process grouping
  - `dd`: kill selected process or group of processes with SIGTERM
  - `d3`: kill selected process or group of processes with SIGQUIT
  - `d9`: kill selected process or group of processes with SIGKILL
- Process sorting
  - `c`: CPU
  - `m`: Mem
  - `p`: PID
- Process filtering:
  - `/`: start editing filter
  - (while editing):
    - `<Enter>` accept filter
    - `<C-c>` and `<Escape>`: clear filter
- CPU and Mem graph scaling:
  - `h`: scale in
  - `l`: scale out
- `?`: toggles keybind help menu

### Mouse

- click to select process
- mouse wheel to scroll through processes

### Colorschemes

gotop ships with a few colorschemes which can be set with the `-c` flag followed by the name of one. You can find all the colorschemes in the [colorschemes folder](./colorschemes).

To make a custom colorscheme, check out the [template](./colorschemes/template.go) for instructions and then use [default.json](./colorschemes/default.json) as a starter. Then put the file at `~/.config/gotop/<name>.json` and load it with `gotop -c <name>`. Colorschemes PR's are welcome!

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
4. The simplest row is a single widget, by name, e.g.
   ```
   cpu
   ```
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
9. If prefixed by a number and colon, the widget will span that number of
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

### CLI Options

`-c`, `--color=NAME` Set a colorscheme.  
`-m`, `--minimal` Only show CPU, Mem and Process widgets.  
`-r`, `--rate=RATE` Number of times per second to update CPU and Mem widgets [default: 1].  
`-V`, `--version` Print version and exit.  
`-p`, `--percpu` Show each CPU in the CPU widget.  
`-a`, `--averagecpu` Show average CPU in the CPU widget.  
`-f`, `--fahrenheit` Show temperatures in fahrenheit.  
`-s`, `--statusbar` Show a statusbar with the time.  
`-b`, `--battery` Show battery level widget (`minimal` turns off). [preview](./assets/screenshots/battery.png)  
`-i`, `--interface=NAME` Select network interface [default: all].
`-l`, `--layout=NAME` Choose a layout. gotop searches for a file by NAME in \$XDG_CONFIG_HOME/gotop, then relative to the current path. "-" reads a layout from stdin, allowing for simple, one-off layouts such as `echo net | gotop -l -`

Several interfaces can be defined using comma separated values.

Interfaces can also be ignored using `!`

## Built With

- [gizak/termui](https://github.com/gizak/termui)
- [nsf/termbox](https://github.com/nsf/termbox-go)
- [exrook/drawille-go](https://github.com/exrook/drawille-go)
- [shirou/gopsutil](https://github.com/shirou/gopsutil)
- [goreleaser/nfpm](https://github.com/goreleaser/nfpm)
- [distatus/battery](https://github.com/distatus/battery)
