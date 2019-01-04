package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"syscall"
	"time"

	appdir "github.com/ProtonMail/go-appdir"
	"github.com/cjbassi/gotop/colorschemes"
	"github.com/cjbassi/gotop/src/logging"
	w "github.com/cjbassi/gotop/src/widgets"
	docopt "github.com/docopt/docopt.go"
	ui "github.com/gizak/termui"
)

var version = "1.7.1"

var (
	colorscheme  = colorschemes.Default
	minimal      = false
	interval     = time.Second
	zoom         = 7
	zoomInterval = 3
	helpVisible  = false
	averageLoad  = false
	battery      = false
	percpuLoad   = false
	fahrenheit   = false
	configDir    = appdir.New("gotop").UserConfig()
	logPath      = filepath.Join(configDir, "errors.log")
	stderrLogger = log.New(os.Stderr, "", 0)
	statusbar    = false
	termWidth    int
	termHeight   int

	cpu  *w.CPU
	batt *w.Batt
	mem  *w.Mem
	proc *w.Proc
	net  *w.Net
	disk *w.Disk
	temp *w.Temp
	help *w.HelpMenu
	grid *ui.Grid
)

func cliArguments() error {
	usage := `
Usage: gotop [options]

Options:
  -c, --color=NAME      Set a colorscheme.
  -h, --help            Show this screen.
  -m, --minimal         Only show CPU, Mem and Process widgets.
  -r, --rate=RATE       Number of times per second to update CPU and Mem widgets [default: 1].
  -v, --version         Print version and exit.
  -p, --percpu          Show each CPU in the CPU widget.
  -a, --averagecpu      Show average CPU in the CPU widget.
  -f, --fahrenheit      Show temperatures in fahrenheit.
  -s, --statusbar       Show a statusbar with the time.
  -b, --battery         Show battery level widget ('minimal' turns off).

Colorschemes:
  default
  default-dark (for white background)
  solarized
  monokai
`

	args, err := docopt.ParseArgs(usage, os.Args[1:], version)
	if err != nil {
		return err
	}

	if val, _ := args["--color"]; val != nil {
		if err := handleColorscheme(val.(string)); err != nil {
			return err
		}
	}
	averageLoad, _ = args["--averagecpu"].(bool)
	percpuLoad, _ = args["--percpu"].(bool)
	battery, _ = args["--battery"].(bool)

	minimal, _ = args["--minimal"].(bool)

	statusbar, _ = args["--statusbar"].(bool)

	rateStr, _ := args["--rate"].(string)
	rate, err := strconv.ParseFloat(rateStr, 64)
	if err != nil {
		return fmt.Errorf("invalid rate parameter")
	}
	if rate < 1 {
		interval = time.Second * time.Duration(1/rate)
	} else {
		interval = time.Second / time.Duration(rate)
	}
	fahrenheit, _ = args["--fahrenheit"].(bool)

	return nil
}

func handleColorscheme(cs string) error {
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
		if colorscheme, err := getCustomColorscheme(cs); err != nil {
			colorscheme = colorscheme
			return err
		}
	}
	return nil
}

// getCustomColorscheme	tries to read a custom json colorscheme from {configDir}/{name}.json
func getCustomColorscheme(name string) (colorschemes.Colorscheme, error) {
	var colorscheme colorschemes.Colorscheme
	filePath := filepath.Join(configDir, name+".json")
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		return colorscheme, fmt.Errorf("colorscheme file not found")
	}
	err = json.Unmarshal(dat, &colorscheme)
	if err != nil {
		return colorscheme, fmt.Errorf("could not parse colorscheme file")
	}
	return colorscheme, nil
}

func setupGrid() {
	grid = ui.NewGrid()
	grid.SetRect(0, 0, termWidth, termHeight)

	var barRow interface{}
	if minimal {
		rowHeight := 1.0 / 2
		if statusbar {
			rowHeight = 50.0 / 101
			barRow = ui.NewRow(1.0/101, w.NewStatusBar())
		}
		grid.Set(
			ui.NewRow(rowHeight, cpu),
			ui.NewRow(rowHeight,
				ui.NewCol(1.0/2, mem),
				ui.NewCol(1.0/2, proc),
			),
			barRow,
		)
	} else {
		rowHeight := 1.0 / 3
		if statusbar {
			rowHeight = 50.0 / 151
			barRow = ui.NewRow(1.0/151, w.NewStatusBar())
		}
		var cpuRow ui.GridItem
		if battery {
			cpuRow = ui.NewRow(rowHeight,
				ui.NewCol(2.0/3, cpu),
				ui.NewCol(1.0/3, batt),
			)
		} else {
			cpuRow = ui.NewRow(rowHeight, cpu)
		}
		grid.Set(
			cpuRow,
			ui.NewRow(rowHeight,
				ui.NewCol(1.0/3,
					ui.NewRow(1.0/2, disk),
					ui.NewRow(1.0/2, temp),
				),
				ui.NewCol(2.0/3, mem),
			),
			ui.NewRow(rowHeight,
				ui.NewCol(1.0/2, net),
				ui.NewCol(1.0/2, proc),
			),
			barRow,
		)
	}
}

func termuiColors() {
	ui.Theme.Default = ui.AttrPair{ui.Attribute(colorscheme.Fg), ui.Attribute(colorscheme.Bg)}
	ui.Theme.Block.Title = ui.AttrPair{ui.Attribute(colorscheme.BorderLabel), ui.Attribute(colorscheme.Bg)}
	ui.Theme.Block.Border = ui.AttrPair{ui.Attribute(colorscheme.BorderLine), ui.Attribute(colorscheme.Bg)}
}

func widgetColors() {
	mem.LineColor["Main"] = ui.Attribute(colorscheme.MainMem)
	mem.LineColor["Swap"] = ui.Attribute(colorscheme.SwapMem)

	proc.CursorColor = ui.Attribute(colorscheme.ProcCursor)

	var keys []string
	for key := range cpu.Data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	i := 0
	for _, v := range keys {
		if i >= len(colorscheme.CPULines) {
			// assuming colorscheme for CPU lines is not empty
			i = 0
		}
		c := colorscheme.CPULines[i]
		cpu.LineColor[v] = ui.Attribute(c)
		i++
	}

	if !minimal {
		if battery {
			var battKeys []string
			for key := range batt.Data {
				battKeys = append(battKeys, key)
			}
			sort.Strings(battKeys)
			i = 0 // Re-using variable from CPU
			for _, v := range battKeys {
				if i >= len(colorscheme.BattLines) {
					// assuming colorscheme for battery lines is not empty
					i = 0
				}
				c := colorscheme.BattLines[i]
				batt.LineColor[v] = ui.Attribute(c)
				i++
			}
		}

		temp.TempLow = ui.Attribute(colorscheme.TempLow)
		temp.TempHigh = ui.Attribute(colorscheme.TempHigh)

		net.Lines[0].LineColor = ui.Attribute(colorscheme.Sparkline)
		net.Lines[0].TitleColor = ui.Attribute(colorscheme.BorderLabel)
		net.Lines[1].LineColor = ui.Attribute(colorscheme.Sparkline)
		net.Lines[1].TitleColor = ui.Attribute(colorscheme.BorderLabel)
	}
}

func initWidgets() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		cpu = w.NewCPU(interval, zoom, averageLoad, percpuLoad)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		mem = w.NewMem(interval, zoom)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		proc = w.NewProc()
		wg.Done()
	}()
	if !minimal {
		if battery {
			wg.Add(1)
			go func() {
				batt = w.NewBatt(time.Minute, zoom)
				wg.Done()
			}()
		}
		wg.Add(1)
		go func() {
			net = w.NewNet()
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			disk = w.NewDisk()
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			temp = w.NewTemp(fahrenheit)
			wg.Done()
		}()
	}

	wg.Wait()
}

func eventLoop() {
	drawTicker := time.NewTicker(interval).C

	// handles kill signal sent to gotop
	sigTerm := make(chan os.Signal, 2)
	signal.Notify(sigTerm, os.Interrupt, syscall.SIGTERM)

	uiEvents := ui.PollEvents()

	previousKey := ""

	for {
		select {
		case <-sigTerm:
			return
		case <-drawTicker:
			if !helpVisible {
				ui.Render(grid)
			}
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "?":
				helpVisible = !helpVisible
				if helpVisible {
					ui.Clear()
					ui.Render(help)
				} else {
					ui.Render(grid)
				}
			case "h":
				if !helpVisible {
					zoom += zoomInterval
					cpu.Zoom = zoom
					mem.Zoom = zoom
					ui.Render(cpu, mem)
				}
			case "l":
				if !helpVisible {
					if zoom > zoomInterval {
						zoom -= zoomInterval
						cpu.Zoom = zoom
						mem.Zoom = zoom
						ui.Render(cpu, mem)
					}
				}
			case "<Escape>":
				if helpVisible {
					helpVisible = false
					ui.Render(grid)
				}
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				help.Resize(payload.Width, payload.Height)
				ui.Clear()
				if helpVisible {
					ui.Render(help)
				} else {
					ui.Render(grid)
				}

			case "<MouseLeft>":
				payload := e.Payload.(ui.Mouse)
				proc.Click(payload.X, payload.Y)
				ui.Render(proc)
			case "k", "<Up>", "<MouseWheelUp>":
				proc.Up()
				ui.Render(proc)
			case "j", "<Down>", "<MouseWheelDown>":
				proc.Down()
				ui.Render(proc)
			case "g", "<Home>":
				if previousKey == "g" {
					proc.Top()
					ui.Render(proc)
				}
			case "G", "<End>":
				proc.Bottom()
				ui.Render(proc)
			case "<C-d>":
				proc.HalfPageDown()
				ui.Render(proc)
			case "<C-u>":
				proc.HalfPageUp()
				ui.Render(proc)
			case "<C-f>":
				proc.PageDown()
				ui.Render(proc)
			case "<C-b>":
				proc.PageUp()
				ui.Render(proc)
			case "d":
				if previousKey == "d" {
					proc.Kill()
				}
			case "<Tab>":
				proc.Tab()
				ui.Render(proc)
			case "m", "c", "p":
				proc.ChangeSort(e)
				ui.Render(proc)
			}

			if previousKey == e.ID {
				previousKey = ""
			} else {
				previousKey = e.ID
			}
		}
	}
}

func setupLogging() (*os.File, error) {
	// make the config directory
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to make the configuration directory: %v", err)
	}
	// open the log file
	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0660)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	// log time, filename, and line number
	log.SetFlags(log.Ltime | log.Lshortfile)
	// log to file
	log.SetOutput(lf)

	return lf, nil
}

func main() {
	lf, err := setupLogging()
	if err != nil {
		stderrLogger.Fatalf("failed to setup logging: %v", err)
	}
	defer lf.Close()

	if err := cliArguments(); err != nil {
		stderrLogger.Fatalf("failed to parse cli args: %v", err)
	}

	if err := ui.Init(); err != nil {
		stderrLogger.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	logging.StderrToLogfile(lf)

	termWidth, termHeight = ui.TerminalSize()

	termuiColors() // need to do this before initializing widgets so that they can inherit the colors
	initWidgets()
	widgetColors()
	help = w.NewHelpMenu()
	help.Resize(termWidth, termHeight)

	setupGrid()
	ui.Render(grid)

	eventLoop()
}
