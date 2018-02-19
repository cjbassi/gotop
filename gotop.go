package main

import (
	// "fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	ui "github.com/cjbassi/gotop/termui"
	w "github.com/cjbassi/gotop/widgets"
	"github.com/docopt/docopt-go"
)

const VERSION = "1.0"

var (
	resized = make(chan bool, 1)

	helpToggled = make(chan bool, 1)
	helpStatus  = false

	procLoaded = make(chan bool, 1)
	keyPressed = make(chan bool, 1)

	cpu  = w.NewCPU()
	mem  = w.NewMem()
	proc = w.NewProc(procLoaded, keyPressed)
	net  = w.NewNet()
	disk = w.NewDisk()
	temp = w.NewTemp()

	help = w.NewHelpMenu()
)

// Sets up docopt which is a command line argument parser
func docoptInit() {
	usage := `
Usage: gotop [options]

Options:
  -c, --color       Set a colorscheme.
  -h, --help        Show this screen.
  -u, --upgrade     Updates gotop if needed.
  -v, --version     Show version.

Colorschemes:
  default
`

	args, _ := docopt.ParseArgs(usage, os.Args[1:], VERSION)
	if val, _ := args["--upgrade"]; val.(bool) {
		updateGotop()
		os.Exit(0)
	}
	if val, _ := args["--color"]; val.(bool) {
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

func updateGotop() {
	cmd := exec.Command("sleep", "1")
	cmd.Run()
	return
}

func main() {
	docoptInit()

	keyBinds()

	<-procLoaded

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	setupGrid()

	ui.On("resize", func(e ui.Event) {
		ui.Body.Width, ui.Body.Height = e.Width, e.Height
		ui.Body.Resize()
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
