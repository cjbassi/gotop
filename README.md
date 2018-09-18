<div align="center">

<img src="https://github.com/cjbassi/gotop/blob/master/assets/logo.png" width="20%" />
<br><br>

Another terminal based graphical activity monitor, inspired by [gtop](https://github.com/aksakalli/gtop) and [vtop](https://github.com/MrRio/vtop), this time written in [Go](https://golang.org/)!

<img src="https://github.com/cjbassi/gotop/blob/master/assets/demo.gif" />
<img src="https://github.com/cjbassi/gotop/blob/master/assets/minimal.png" width="96%" />

</div>

## Installation

Working and tested on Linux and OSX. Windows support is planned.

### Using Git

Clone the repo and then run [scripts/download.sh](https://github.com/cjbassi/gotop/blob/master/scripts/download.sh) to download the correct binary for your system from the [releases tab](https://github.com/cjbassi/gotop/releases):

```sh
git clone --depth 1 https://github.com/cjbassi/gotop /tmp/gotop
/tmp/gotop/scripts/download.sh
```

Then move `gotop` into your $PATH somewhere.

### Arch Linux

Install the `gotop-bin` package from the AUR.

### Homebrew

```
brew tap cjbassi/gotop
brew install gotop
```

### Source

```sh
go get github.com/cjbassi/gotop
```

## Usage

### Keybinds

- Quit: `q` or `<C-c>`
- Process Navigation:
  - `<up>`/`<down>` and `j`/`k`: up and down
  - `<C-d>` and `<C-u>`: up and down half a page
  - `<C-f>` and `<C-b>`: up and down a full page
  - `gg` and `G`: jump to top and bottom
- Process Sorting:
  - `c`: CPU
  - `m`: Mem
  - `p`: PID
- `<tab>`: toggle process grouping
- `dd`: kill the selected process or process group
- `h` and `l`: zoom in and out of CPU and Mem graphs
- `?`: toggles keybind help menu

### Mouse

- click to select process
- mouse wheel to scroll through processes

### Colorschemes

gotop ships with a few colorschemes which can be set with the `-c` flag followed by the name of one. You can find all the colorschemes in [colorschemes](https://github.com/cjbassi/gotop/tree/master/colorschemes).

To make a custom colorscheme, check out the [template](https://github.com/cjbassi/gotop/blob/master/colorschemes/template.go) for instructions and then use [default.json](https://github.com/cjbassi/gotop/blob/master/colorschemes/default.json) as a starter. Then you can put the file at `~/.config/gotop/{name}.json` and load it with `gotop -c {name}`. Colorschemes PR's are welcome!

### CLI Options

`-c`, `--color=NAME` Set a colorscheme.  
`-m`, `--minimal` Only show CPU, Mem and Process widgets.  
`-r`, `--rate=RATE` Number of times per second to update CPU and Mem widgets [default: 1].  
`-v`, `--version` Show version.  
`-p`, `--percpu` Show each CPU in the CPU widget.  
`-a`, `--averagecpu` Show average CPU in the CPU widget.

## Credits

- [mdnazmulhasan27771](https://github.com/mdnazmulhasan27771) for the [logo](https://github.com/cjbassi/gotop/blob/master/assets/logo.png)
- [f1337](https://github.com/f1337) for helping port gotop to OSX

## Built With

- [cjbassi/termui](https://github.com/cjbassi/termui)
  - [drawille-go](https://github.com/exrook/drawille-go)
  - [termbox](https://github.com/nsf/termbox-go)
- [gopsutil](https://github.com/shirou/gopsutil)
- [goreleaser](https://github.com/goreleaser/goreleaser)

## Stargazers over time

[![Stargazers over time](https://starcharts.herokuapp.com/cjbassi/gotop.svg)](https://starcharts.herokuapp.com/cjbassi/gotop)
