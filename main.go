package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/cjbassi/gotop/colorschemes"
	w "github.com/cjbassi/gotop/widgets"
	ui "github.com/cjbassi/termui"
	"github.com/docopt/docopt-go"
)

const VERSION = "1.2.7"

var (
	termResized = make(chan bool, 1)

	helpToggled = make(chan bool, 1)
	helpVisible = false

	// proc widget takes longer to load, wait to render until it loads data
	procLoaded = make(chan bool, 1)
	// used to render the proc widget whenever a key is pressed for it
	keyPressed = make(chan bool, 1)
	// used to render cpu and mem when zoom has changed
	zoomed = make(chan bool, 1)

	colorscheme = colorschemes.Default

	minimal      = false
	interval     = time.Second
	zoom         = 7
	zoomInterval = 3

	cpu  *w.CPU
	mem  *w.Mem
	proc *w.Proc
	net  *w.Net
	disk *w.Disk
	temp *w.Temp

	help *w.HelpMenu
)

func cliArguments() {
	usage := `
Usage: gotop [options]

Options:
  -c, --color=NAME      Set a colorscheme.
  -h, --help            Show this screen.
  -m, --minimal         Only show CPU, Mem and Process widgets.
  -r, --rate=RATE       Number of times per second to update CPU and Mem widgets [default: 1].
  -v, --version         Show version.

Colorschemes:
  default
  default-dark (for white background)
  solarized
  monokai
`

	args, _ := docopt.ParseArgs(usage, os.Args[1:], VERSION)

	if val, _ := args["--color"]; val != nil {
		handleColorscheme(val.(string))
	}

	minimal, _ = args["--minimal"].(bool)

	rateStr, _ := args["--rate"].(string)
	rate, _ := strconv.ParseFloat(rateStr, 64)
	if rate < 1 {
		interval = time.Second * time.Duration(1/rate)
	} else {
		interval = time.Second / time.Duration(rate)
	}
}

func handleColorscheme(cs string) {
	switch cs {
	case "default":
		colorscheme = colorschemes.Default
	case "solarized":
		colorscheme = colorschemes.Solarized
	case "monokai":
		colorscheme = colorschemes.Monokai
	case "default-dark":
		colorscheme = colorschemes.DefaultDark
	default:
		fmt.Fprintf(os.Stderr, "error: colorscheme not recognized\n")
		os.Exit(1)
	}
}

func setupGrid() {
	ui.Body.Cols = 12
	ui.Body.Rows = 12

	if minimal {
		ui.Body.Set(0, 0, 12, 6, cpu)
		ui.Body.Set(0, 6, 6, 12, mem)
		ui.Body.Set(6, 6, 12, 12, proc)
	} else {
		ui.Body.Set(0, 0, 12, 4, cpu)

		ui.Body.Set(0, 4, 4, 6, disk)
		ui.Body.Set(0, 6, 4, 8, temp)
		ui.Body.Set(4, 4, 12, 8, mem)

		ui.Body.Set(0, 8, 6, 12, net)
		ui.Body.Set(6, 8, 12, 12, proc)
	}
}

func keyBinds() {
	// quits
	ui.On("q", "<C-c>", func(e ui.Event) {
		ui.StopLoop()
	})

	// toggles help menu
	ui.On("?", func(e ui.Event) {
		helpToggled <- true
		helpVisible = !helpVisible
	})
	// hides help menu
	ui.On("<escape>", func(e ui.Event) {
		if helpVisible {
			helpToggled <- true
			helpVisible = false
		}
	})

	ui.On("h", func(e ui.Event) {
		zoom += zoomInterval
		cpu.Zoom = zoom
		mem.Zoom = zoom
		zoomed <- true
	})
	ui.On("l", func(e ui.Event) {
		if zoom > zoomInterval {
			zoom -= zoomInterval
			cpu.Zoom = zoom
			mem.Zoom = zoom
			zoomed <- true
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
	ui.Theme.GaugeColor = ui.Color(colorscheme.DiskBar)
}

func widgetColors() {
	mem.LineColor["Main"] = ui.Color(colorscheme.MainMem)
	mem.LineColor["Swap"] = ui.Color(colorscheme.SwapMem)

	LineColor := make(map[string]ui.Color)
	if cpu.Count <= 8 {
		for i := 0; i < len(cpu.Data); i++ {
			LineColor[fmt.Sprintf("CPU%d", i)] = ui.Color(colorscheme.CPULines[i])
		}
	} else {
		LineColor["Average"] = ui.Color(colorscheme.CPULines[0])
	}
	cpu.LineColor = LineColor

	if !minimal {
		temp.TempLow = ui.Color(colorscheme.TempLow)
		temp.TempHigh = ui.Color(colorscheme.TempHigh)
	}
}

func main() {
	cliArguments()

	keyBinds()

	// need to do this before initializing widgets so that they can inherit the colors
	termuiColors()

	cpu = w.NewCPU(interval, zoom)
	mem = w.NewMem(interval, zoom)
	proc = w.NewProc(procLoaded, keyPressed)
	if !minimal {
		net = w.NewNet()
		disk = w.NewDisk()
		temp = w.NewTemp()
	}

	widgetColors()

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

	ui.On("<resize>", func(e ui.Event) {
		ui.Body.Width, ui.Body.Height = e.Width, e.Height
		ui.Body.Resize()

		help.XOffset = (ui.Body.Width - help.X) / 2
		help.YOffset = (ui.Body.Height - help.Y) / 2

		termResized <- true
	})

	// all rendering done here
	go func() {
		ui.Render(ui.Body)
		drawTick := time.NewTicker(interval)
		for {
			select {
			case <-helpToggled:
				if helpVisible {
					ui.Clear()
					ui.Render(help)
				} else {
					ui.Render(ui.Body)
				}
			case <-termResized:
				if !helpVisible {
					ui.Clear()
					ui.Render(ui.Body)
				} else if helpVisible {
					ui.Clear()
					ui.Render(help)
				}
			case <-keyPressed:
				if !helpVisible {
					ui.Render(proc)
				}
			case <-zoomed:
				if !helpVisible {
					ui.Render(ui.Body)
				}
			case <-drawTick.C:
				if !helpVisible {
					ui.Render(ui.Body)
				}
			}
		}
	}()

	// handles os kill signal
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		ui.StopLoop()
	}()

	ui.Loop()
}
