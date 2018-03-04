# gotop

![image](https://github.com/cjbassi/gotop/blob/master/demo.gif)

Another TUI graphical activity monitor, inspired by [gtop](https://github.com/aksakalli/gtop) and [vtop](https://github.com/MrRio/vtop), this time written in [Go](https://golang.org/)!  
Built with [gopsutil](https://github.com/shirou/gopsutil), [drawille-go](https://github.com/exrook/drawille-go), and a modified version of [termui](https://github.com/gizak/termui).


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

* Quit: `q` or `Ctrl-c`
* Navigation:
    * `<up>`/`<down>` and `j`/`k`: up and down
    * `C-d` and `C-u`: up and down half a page
    * `C-f` and `C-b`: up and down a full page
    * `gg` and `G`: jump to top and bottom
* Process Sorting:
    * `c`: CPU
    * `m`: Mem
    * `p`: PID
* `<tab>`: toggle process grouping
* `dd`: kill the selected process or process group
* `?`: toggles keybind help menu


### Mouse

* click to select process
* mouse wheel to scroll through processes


## Colorschemes

A different Colorscheme can be set with the `-c` flag followed its name. You can find them in the `colorschemes` folder.
Feel free to add a new one. You can use 256 colors, bold, underline, and reverse. You can see the template and get more info [here](https://github.com/cjbassi/gotop/blob/master/colorschemes/template.go) and see the default colorscheme as an example [here](https://github.com/cjbassi/gotop/blob/master/colorschemes/default.go).


## TODO

* Network Usage
    - increase height of sparkline depending on widget size
* Process List
    - memory total goes above 100%
* Graphs
    - zooming in and out of graphs
    - command line option to set polling interval for CPU and mem
* general
    - command line option to only show processes, CPU, and mem
    - gopsutil cross-compiling issue on linux_amd64 for darwin_i386
* try to get the drawille fork merged upstream
* termui code cleanup
    - termui buffers should ignore characters set outside the widget area
        - ignore writes or give an error?
    - termui Blocks should be indexed at 0, and maybe change their X and Y variables too
    - draw borders and label after widget contents
    - only merge buffers that are within the original's area
