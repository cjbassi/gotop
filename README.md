# gotop

Another TUI graphical activity monitor, inspired by [gtop](https://github.com/aksakalli/gtop) and [vtop](https://github.com/MrRio/vtop), this time written in [Go](https://golang.org/)!  
Built with [gopsutil](https://github.com/shirou/gopsutil), [drawille-go](https://github.com/exrook/drawille-go), and a modified version of [termui](https://github.com/gizak/termui).

<img src="https://github.com/cjbassi/gotop/blob/master/media/demo.gif" />
<img src="https://github.com/cjbassi/gotop/blob/master/media/minimal.png" />


## Installation

### Binaries

Binaries are currently available for 32/64bit Linux and 64bit OSX.

To download the latest binary for your system from GitHub, you can run the [download](https://github.com/cjbassi/gotop/blob/master/download.sh) script:

```
sh -c "$(curl https://raw.githubusercontent.com/cjbassi/gotop/master/download.sh)"
```

Then move `gotop` into your $PATH somewhere.


### Arch Linux

Alternatively, if you're on Arch Linux, you can install the `gotop` package from the AUR.

### Source

```
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

A different Colorscheme can be set with the `-c` flag followed its name. You can find them in the `colorschemes` folder.
Feel free to add a new one. You can use 256 colors, bold, underline, and reverse. You can see the template and get more info [here](https://github.com/cjbassi/gotop/blob/master/colorschemes/template.go) and see the default colorscheme as an example [here](https://github.com/cjbassi/gotop/blob/master/colorschemes/default.go).


### CLI Options

`-m`, `--minimal`         Only show CPU, Mem and Process widgets.  
`-r`, `--rate=RATE`       Number of times per second to update CPU and Mem widgets [default: 1].
