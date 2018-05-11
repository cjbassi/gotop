<div align="center">

<img src="https://github.com/cjbassi/gotop/blob/master/assets/logo.png" width="20%" />
<br><br>

Another terminal based graphical activity monitor, inspired by [gtop](https://github.com/aksakalli/gtop) and [vtop](https://github.com/MrRio/vtop), this time written in [Go](https://golang.org/)!

<img src="https://github.com/cjbassi/gotop/blob/master/assets/demo.gif" />
<img src="https://github.com/cjbassi/gotop/blob/master/assets/minimal.png" width="96%" />

</div>


## Installation

Only working and tested on Linux. OSX is no longer supported due to issues with gopsutil, but that is currently being worked on. Windows support is also in the works.

Go programs compile to a single binary and there are currently prebuilt ones for 32/64bit Linux and ARM Linux. Use one of the following methods to get the binary installed on your machine:

### Using Git

Clone the repo then run [scripts/download.sh](https://github.com/cjbassi/gotop/blob/master/scripts/download.sh) to download the correct binary:

```sh
git clone --depth 1 https://github.com/cjbassi/gotop /tmp/gotop
/tmp/gotop/scripts/download.sh
```

Then move `gotop` into your $PATH somewhere.

### Arch Linux

Install the `gotop-bin` package from the AUR.

### Source

```sh
go get github.com/cjbassi/gotop
```

### Docker

```
docker run -it --rm cjbassi/gotop
```

Note: Process list doesn't work when running with Docker.


## Usage

### Keybinds

* Quit: `q` or `<C-c>`
* Process Navigation:
    * `<up>`/`<down>` and `j`/`k`: up and down
    * `<C-d>` and `<C-u>`: up and down half a page
    * `<C-f>` and `<C-b>`: up and down a full page
    * `gg` and `G`: jump to top and bottom
* Process Sorting:
    * `c`: CPU
    * `m`: Mem
    * `p`: PID
* `<tab>`: toggle process grouping
* `dd`: kill the selected process or process group
* `h` and `l`: zoom in and out of CPU and Mem graphs
* `?`: toggles keybind help menu

### Mouse

* click to select process
* mouse wheel to scroll through processes


### Colorschemes

A different Colorscheme can be set with the `-c` flag followed its name.
You can find different ones in [src/colorschemes](https://github.com/cjbassi/gotop/tree/master/src/colorschemes).
Feel free to add a new one.
You can use 256 colors, bold, underline, and reverse.
You can see the template and get more info [here](https://github.com/cjbassi/gotop/blob/master/src/colorschemes/template.go)
and see the default colorscheme as an example [here](https://github.com/cjbassi/gotop/blob/master/src/colorschemes/default.go).

### CLI Options

`-m`, `--minimal`         Only show CPU, Mem and Process widgets.  
`-r`, `--rate=RATE`       Number of times per second to update CPU and Mem widgets [default: 1].


## Credits

* [Logo](https://github.com/cjbassi/gotop/blob/master/assets/logo.png) by [mdnazmulhasan27771](https://github.com/mdnazmulhasan27771)


## Built With

* [My termui fork](https://github.com/cjbassi/termui)
    * [drawille-go](https://github.com/exrook/drawille-go)
    * [termbox](https://github.com/nsf/termbox-go)
* [gopsutil](https://github.com/shirou/gopsutil)
* [goreleaser](https://github.com/goreleaser/goreleaser)


## Stargazers over time

[![Stargazers over time](https://starcharts.herokuapp.com/cjbassi/gotop.svg)](https://starcharts.herokuapp.com/cjbassi/gotop)
