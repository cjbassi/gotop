package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/cjbassi/gotop/src/colorschemes"
	w "github.com/cjbassi/gotop/src/widgets"
	ui "github.com/cjbassi/termui"
	"github.com/docopt/docopt-go"
)

var version = "1.4.0"

var (
	termResized = make(chan bool, 1)

	helpToggled = make(chan bool, 1)
	helpVisible = false

	wg sync.WaitGroup
	// used to render the proc widget whenever a key is pressed for it
	keyPressed = make(chan bool, 1)
	// used to render cpu and mem when zoom has changed
	zoomed = make(chan bool, 1)

	colorscheme = colorschemes.Default

	minimal      = false
	widgetCount  = 6
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

	args, _ := docopt.ParseArgs(usage, os.Args[1:], version)

	if val, _ := args["--color"]; val != nil {
		handleColorscheme(val.(string))
	}

	minimal, _ = args["--minimal"].(bool)
	if minimal {
		widgetCount = 3
	}

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

// load widgets asynchronously but wait till they are all finished
func initWidgets() {
	wg.Add(widgetCount)

	go func() {
		cpu = w.NewCPU(interval, zoom)
		wg.Done()
	}()
	go func() {
		mem = w.NewMem(interval, zoom)
		wg.Done()
	}()
	go func() {
		proc = w.NewProc(keyPressed)
		wg.Done()
	}()
	if !minimal {
		go func() {
			net = w.NewNet()
			wg.Done()
		}()
		go func() {
			disk = w.NewDisk()
			wg.Done()
		}()
		go func() {
			temp = w.NewTemp()
			wg.Done()
		}()
	}

	wg.Wait()
}

func main() {
	cliArguments()

	keyBinds()

	// need to do this before initializing widgets so that they can inherit the colors
	termuiColors()

	initWidgets()

	widgetColors()

	help = w.NewHelpMenu()

	// inits termui
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	setupGrid()

	ui.On("<resize>", func(e ui.Event) {
		ui.Body.Width, ui.Body.Height = e.Width, e.Height
		ui.Body.Resize()

		termResized <- true
	})

	// all rendering done here
	go func() {
		ui.Render(ui.Body)
		drawTick := time.NewTicker(interval)
		for {
			if helpVisible {
				select {
				case <-helpToggled:
					ui.Render(ui.Body)
				case <-termResized:
					ui.Clear()
					ui.Render(help)
				}
			} else {
				select {
				case <-helpToggled:
					ui.Clear()
					ui.Render(help)
				case <-termResized:
					ui.Clear()
					ui.Render(ui.Body)
				case <-keyPressed:
					ui.Render(proc)
				case <-zoomed:
					ui.Render(cpu, mem)
				case <-drawTick.C:
					ui.Render(ui.Body)
				}
			}
		}
	}()

	// handles kill signal sent to gotop
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		ui.StopLoop()
	}()

	ui.Loop()
}
