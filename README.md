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

### Building

The easiest way is to
```
go get github.com/xxxserxxx/gotop/cmd/gotop
```

To create the cross-compile builds, there's a `make.sh` script; it has a lot of dependencies and has only been tested on my computer. When it works, it creates archives for numerous OSes & architectures. There's no testing for whether dependencies are available; it assumes they are and will fail in strange ways when they aren't.

- bash
- Go
- zip
- nfpm (for deb & rpm)
- docker (for darwin)

It is *just* smart enough to not rebuild things when it doesn't have to, and it tries to keep the darwin docker container around so it's not building from scratch every time. There are no guarantees.

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
4. Legal widget names are: cpu, disk, mem, temp, batt, net, procs, power
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

### Metrics

gotop can export widget data as Prometheus metrics. This allows users to take
snapshots of the current state of a machine running gotop, or to query gotop
remotely.

All metrics are in the `gotop` namespace, and are tagged with
`goto_<widget>_<enum>`. Metrics are only exported for widgets
that are enabled, and are updated with the same frequency as the configured
update interval.  Most widgets are exported as Prometheus gauges.

Metrics are disabled by default, and must be enabled with the `--export` flag.
The flag takes an interface port in the idiomatic Go format of
`<addy>:<port>`; a common pattern is `-x :2112`. There is **no security**
on this feature; I recommend that you run this bound to a localhost interface,
e.g. `127.0.0.1:7653`, and if you want to access this remotely, run it behind
a proxy that provides SSL and authentication such as
[Caddy](https://caddyserver.com).

Once enabled, any widgets that are enabled will appear in the HTTP payload of
a call to `http://<addy>:<port>/metrics`. For example,

```
➜  ~ curl -s http://localhost:2112/metrics | egrep -e '^gotop'
gotop_battery_0 0.6387792286668692
gotop_cpu_0 12.871287128721228
gotop_cpu_1 11.000000000001364
gotop_disk_:dev:nvme0n1p1 0.63
gotop_memory_main 49.932259713701434
gotop_memory_swap 0
gotop_net_recv 129461
gotop_net_sent 218525
gotop_temp_coretemp_core0 37
gotop_temp_coretemp_core1 37
```

Disk metrics are reformatted to replace `/` with `:` which makes them legal
Prometheus names:

```
➜  ~ curl -s http://localhost:2112/metrics | egrep -e '^gotop_disk' | tr ':' '/'
gotop_disk_/dev/nvme0n1p1 0.63
```

This feature satisfies a ticket request to provide a "snapshot" for comparison
with a known state, but it is also foundational for a future feature where
widgets can be configured with virtual devices fed by data from remote gotop
instances. The objective for that feature is to allow monitoring of multiple
remote VMs without having to have a wall of gotops running on a large monitor.

### CLI Options

`-c`, `--color=NAME` Set a colorscheme.  
`-m`, `--minimal` Only show CPU, Mem and Process widgets.  (DEPRECATED for `-l minimal`)  
`-r`, `--rate=RATE` Number of times per second to update CPU and Mem widgets [default: 1].  
`-V`, `--version` Print version and exit.  
`-p`, `--percpu` Show each CPU in the CPU widget.  
`-a`, `--averagecpu` Show average CPU in the CPU widget.  
`-f`, `--fahrenheit` Show temperatures in fahrenheit.  
`-s`, `--statusbar` Show a statusbar with the time.  
`-b`, `--battery` Show battery level widget (`minimal` turns off). [preview](./assets/screenshots/battery.png)  (DEPRECATED for `-l battery`)  
`-i`, `--interface=NAME` Select network interface. Several interfaces can be defined using comma separated values. Interfaces can also be ignored by prefixing the interface with `!` [default: all].  
`-l`, `--layout=NAME` Choose a layout. gotop searches for a file by NAME in \$XDG_CONFIG_HOME/gotop, then relative to the current path. "-" reads a layout from stdin, allowing for simple, one-off layouts such as `echo net | gotop -l -`  
`-x`, `--export=PORT` Enable metrics for export on the specified port. This feature is disabled by default.


## Built With

- [gizak/termui](https://github.com/gizak/termui)
- [nsf/termbox](https://github.com/nsf/termbox-go)
- [exrook/drawille-go](https://github.com/exrook/drawille-go)
- [shirou/gopsutil](https://github.com/shirou/gopsutil)
- [goreleaser/nfpm](https://github.com/goreleaser/nfpm)
- [distatus/battery](https://github.com/distatus/battery)
- [prometheus/client_golang](https://github.com/prometheus/client_golang)
