<div align="center">

<img src="https://github.com/cjbassi/gotop/blob/master/assets/logo.png" width="20%" />
<br><br>

Another terminal based graphical activity monitor, inspired by [gtop](https://github.com/aksakalli/gtop) and [vtop](https://github.com/MrRio/vtop), this time written in [Go](https://golang.org/)!

<img src="https://github.com/cjbassi/gotop/blob/master/assets/demo.gif" />
<img src="https://github.com/cjbassi/gotop/blob/master/assets/minimal.png" width="96%" />

</div>


## Installation

Working and tested on Linux and OSX, with Windows support being worked on.


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

* [mdnazmulhasan27771](https://github.com/mdnazmulhasan27771) for the [logo](https://github.com/cjbassi/gotop/blob/master/assets/logo.png)
* [f1337](https://github.com/f1337) for helping port gotop to OSX


## Built With

* [My termui fork](https://github.com/cjbassi/termui)
    * [drawille-go](https://github.com/exrook/drawille-go)
    * [termbox](https://github.com/nsf/termbox-go)
* [gopsutil](https://github.com/shirou/gopsutil)
* [goreleaser](https://github.com/goreleaser/goreleaser)


## Stargazers over time

[![Stargazers over time](https://starcharts.herokuapp.com/cjbassi/gotop.svg)](https://starcharts.herokuapp.com/cjbassi/gotop)
