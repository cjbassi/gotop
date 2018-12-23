<div align="center">

<img src="./assets/logo.png" width="20%" />
<br><br>

Another terminal based graphical activity monitor, inspired by [gtop](https://github.com/aksakalli/gtop) and [vtop](https://github.com/MrRio/vtop), this time written in [Go](https://golang.org/)!

<img src="./assets/demo.gif" />
<img src="./assets/minimal.png" width="96%" />

</div>

## Installation

Working and tested on Linux, FreeBSD and OSX. Windows support is planned.

### Using Git

Clone the repo and then run [scripts/download.sh](./scripts/download.sh) to download the correct binary for your system from the [releases tab](https://github.com/cjbassi/gotop/releases):

```bash
git clone --depth 1 https://github.com/cjbassi/gotop /tmp/gotop
/tmp/gotop/scripts/download.sh
```

Then move `gotop` into your \$PATH somewhere.

### Arch Linux

Install `gotop-bin`, `gotop`, or `gotop-git` from the AUR.

### FreeBSD

```
pkg install gotop
```

### Homebrew

```
brew tap cjbassi/gotop
brew install gotop
```

### Source

**Note**: Requires Go 1.11+.

```bash
go get github.com/cjbassi/gotop
```

## Usage

### Keybinds

- Quit: `q` or `<C-c>`
- Process navigation
  - `k` and `<Up>`: up
  - `j` and `<Down`: down
  - `<C-u>`: half page up
  - `<C-d>`: half page down
  - `<C-b>`: full page up
  - `<C-f>`: full page down
  - `gg` and `<Home>`: jump to top
  - `G` and `<End>`: jump to bottom
- Process actions:
  - `<Tab>`: toggle process grouping
  - `dd`: kill selected process or group of processes
- Process sorting
  - `c`: CPU
  - `m`: Mem
  - `p`: PID
- CPU and Mem graph scaling:
  - `h`: scale in
  - `l`: scale out
- `?`: toggles keybind help menu

### Mouse

- click to select process
- mouse wheel to scroll through processes

### Colorschemes

gotop ships with a few colorschemes which can be set with the `-c` flag followed by the name of one. You can find all the colorschemes in the [colorschemes folder](./colorschemes).

To make a custom colorscheme, check out the [template](./colorschemes/template.go) for instructions and then use [default.json](./colorschemes/default.json) as a starter. Then put the file at `~/.config/gotop/{name}.json` on Linux or `~/Library/Application Support/gotop/{name}.json` on OSX and load it with `gotop -c {name}`. Colorschemes PR's are welcome!

### CLI Options

`-c`, `--color=NAME` Set a colorscheme.  
`-m`, `--minimal` Only show CPU, Mem and Process widgets.  
`-r`, `--rate=RATE` Number of times per second to update CPU and Mem widgets [default: 1].  
`-v`, `--version` Show version.  
`-p`, `--percpu` Show each CPU in the CPU widget.  
`-a`, `--averagecpu` Show average CPU in the CPU widget.

## Building deb/rpms

To build dep/rpms using [nfpm](https://github.com/goreleaser/nfpm):

```bash
make all
```

This will place the built packages into the `dist` folder.

## Credits

- [mdnazmulhasan27771](https://github.com/mdnazmulhasan27771) for the [logo](./assets/logo.png)
- [f1337](https://github.com/f1337) for helping port gotop to OSX

## Built With

- [cjbassi/termui](https://github.com/cjbassi/termui)
  - [drawille-go](https://github.com/exrook/drawille-go)
  - [termbox](https://github.com/nsf/termbox-go)
- [gopsutil](https://github.com/shirou/gopsutil)
- [goreleaser](https://github.com/goreleaser/goreleaser)

## Stargazers over time

[![Stargazers over time](https://starcharts.herokuapp.com/cjbassi/gotop.svg)](https://starcharts.herokuapp.com/cjbassi/gotop)
