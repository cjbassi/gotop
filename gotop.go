package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cjbassi/gotop/colorschemes"
	ui "github.com/cjbassi/gotop/termui"
	w "github.com/cjbassi/gotop/widgets"
	"github.com/docopt/docopt-go"
)

const VERSION = "1.0.1"

var (
	// for when terminal is resized
	resized = make(chan bool, 1)

	// for when help menu is toggled
	helpToggled = make(chan bool, 1)
	// whether help menu is toggled
	helpStatus = false

	// proc widget takes longer to load, wait to render until it loads data
	procLoaded = make(chan bool, 1)
	// used to render the proc widget whenever a key is pressed for it
	keyPressed = make(chan bool, 1)

	colorscheme = colorschemes.Default

	cpu  *w.CPU
	mem  *w.Mem
	proc *w.Proc
	net  *w.Net
	disk *w.Disk
	temp *w.Temp

	help *w.HelpMenu
)

// Sets up docopt which is a command line argument parser
func cliArguments() {
	usage := `
Usage: gotop [options]

Options:
  -c, --color <name>    Set a colorscheme.
  -h, --help		    Show this screen.
  -v, --version         Show version.

Colorschemes:
  default
  solarized
  monokai
`

	args, _ := docopt.ParseArgs(usage, os.Args[1:], VERSION)

	if val, _ := args["--color"]; val != nil {
		handleColorscheme(val.(string))
	}
}

func handleColorscheme(cs string) {
	switch cs {
	case "monokai":
		colorscheme = colorschemes.Monokai
	case "solarized":
		colorscheme = colorschemes.Solarized
	case "default":
		colorscheme = colorschemes.Default
	default:
		fmt.Fprintf(os.Stderr, "error: colorscheme not recognized\n")
		os.Exit(1)
	}
}

func setupGrid() {
	ui.Body.Cols = 12
	ui.Body.Rows = 12

	ui.Body.Set(0, 0, 12, 4, cpu)

	ui.Body.Set(0, 4, 4, 6, disk)
	ui.Body.Set(0, 6, 4, 8, temp)
	ui.Body.Set(4, 4, 12, 8, mem)

	ui.Body.Set(0, 8, 6, 12, net)
	ui.Body.Set(6, 8, 12, 12, proc)
}

func keyBinds() {
	// quits
	ui.On("q", "C-c", func(e ui.Event) {
		ui.StopLoop()
	})

	// toggles help menu
	ui.On("?", func(e ui.Event) {
		helpToggled <- true
		helpStatus = !helpStatus
	})
	// hides help menu
	ui.On("<escape>", func(e ui.Event) {
		if helpStatus {
			helpToggled <- true
			helpStatus = false
		}
	})
}

func termuiColors() {
	ui.Theme.Fg = ui.Color(colorscheme.Fg)
	ui.Theme.Bg = ui.Color(colorscheme.Bg)
	ui.Theme.LabelFg = ui.Color(colorscheme.BorderLabel)
	ui.Theme.LabelBg = ui.Color(colorscheme.Bg)
	ui.Theme.BorderFg = ui.Color(colorscheme.BorderLine)
	ui.Theme.BorderBg = ui.Color(colorscheme.Bg)

	ui.Theme.TableCursor = ui.Color(colorscheme.ProcCursor)
	ui.Theme.Sparkline = ui.Color(colorscheme.Sparkline)
	ui.Theme.BarColor = ui.Color(colorscheme.DiskBar)
	ui.Theme.TempLow = ui.Color(colorscheme.TempLow)
	ui.Theme.TempHigh = ui.Color(colorscheme.TempHigh)
}

func widgetColors() {
	// memory widget colors
	mem.LineColor["Main"] = ui.Color(colorscheme.MainMem)
	mem.LineColor["Swap"] = ui.Color(colorscheme.SwapMem)

	// cpu widget colors
	LineColor := make(map[string]ui.Color)
	for i := 0; i < len(cpu.Data); i++ {
		LineColor[fmt.Sprintf("CPU%d", i+1)] = ui.Color(colorscheme.CPULines[i])
	}
	cpu.LineColor = LineColor
}

func main() {
	cliArguments()

	keyBinds()

	// need to do this before initializing widgets so that they can inherit the colors
	termuiColors()

	cpu = w.NewCPU()
	mem = w.NewMem()
	proc = w.NewProc(procLoaded, keyPressed)
	net = w.NewNet()
	disk = w.NewDisk()
	temp = w.NewTemp()

	widgetColors()

	// blocks till loaded
	<-procLoaded

	// inits termui
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	setupGrid()

	// load help widget after init termui/termbox so that it has access to terminal size
	help = w.NewHelpMenu()

	ui.On("resize", func(e ui.Event) {
		ui.Body.Width, ui.Body.Height = e.Width, e.Height
		ui.Body.Resize()

		help.XOffset = (ui.Body.Width - help.X) / 2
		help.YOffset = (ui.Body.Height - help.Y) / 2

		resized <- true
	})

	// All rendering done here
	go func() {
		ui.Render(ui.Body)
		drawTick := time.NewTicker(time.Second)
		for {
			select {
			case <-helpToggled:
				if helpStatus {
					ui.Clear()
					ui.Render(help)
				} else {
					ui.Render(ui.Body)
				}
			case <-resized:
				if !helpStatus {
					ui.Clear()
					ui.Render(ui.Body)
				} else if helpStatus {
					ui.Clear()
					ui.Render(help)
				}
			case <-keyPressed:
				if !helpStatus {
					ui.Render(proc)
				}
			case <-drawTick.C:
				if !helpStatus {
					ui.Render(ui.Body)
				}
			}
		}
	}()

	// handles kill signal
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		ui.StopLoop()
	}()

	ui.Loop()
}
