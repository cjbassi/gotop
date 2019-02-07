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
	docopt "github.com/docopt/docopt.go"
	ui "github.com/gizak/termui"

	"github.com/cjbassi/gotop/colorschemes"
	"github.com/cjbassi/gotop/src/logging"
	w "github.com/cjbassi/gotop/src/widgets"
)

const (
	version = "2.0.1"

	graphHorizontalScaleDelta = 3
)

var (
	configDir = appdir.New("gotop").UserConfig()
	logDir    = appdir.New("gotop").UserLogs()
	logPath   = filepath.Join(logDir, "errors.log")

	stderrLogger = log.New(os.Stderr, "", 0)

	graphHorizontalScale = 7
	helpVisible          = false

	colorscheme    = colorschemes.Default
	updateInterval = time.Second
	minimalMode    = false
	averageLoad    = false
	percpuLoad     = false
	fahrenheit     = false
	battery        = false
	statusbar      = false

	renderLock sync.RWMutex

	cpu  *w.CPU
	batt *w.Batt
	mem  *w.Mem
	proc *w.Proc
	net  *w.Net
	disk *w.Disk
	temp *w.Temp
	help *w.HelpMenu
	grid *ui.Grid
	bar  *w.StatusBar
)

func parseArgs() error {
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

	minimalMode, _ = args["--minimal"].(bool)

	statusbar, _ = args["--statusbar"].(bool)

	rateStr, _ := args["--rate"].(string)
	rate, err := strconv.ParseFloat(rateStr, 64)
	if err != nil {
		return fmt.Errorf("invalid rate parameter")
	}
	if rate < 1 {
		updateInterval = time.Second * time.Duration(1/rate)
	} else {
		updateInterval = time.Second / time.Duration(rate)
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
		custom, err := getCustomColorscheme(cs)
		if err != nil {
			return err
		}
		colorscheme = custom
	}
	return nil
}

// getCustomColorscheme	tries to read a custom json colorscheme from <configDir>/<name>.json
func getCustomColorscheme(name string) (colorschemes.Colorscheme, error) {
	var colorscheme colorschemes.Colorscheme
	filePath := filepath.Join(configDir, name+".json")
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		return colorscheme, fmt.Errorf("failed to read colorscheme file: %v", err)
	}
	err = json.Unmarshal(dat, &colorscheme)
	if err != nil {
		return colorscheme, fmt.Errorf("failed to parse colorscheme file: %v", err)
	}
	return colorscheme, nil
}

func setupGrid() {
	grid = ui.NewGrid()

	if minimalMode {
		grid.Set(
			ui.NewRow(1.0/2, cpu),
			ui.NewRow(1.0/2,
				ui.NewCol(1.0/2, mem),
				ui.NewCol(1.0/2, proc),
			),
		)
	} else {
		var cpuRow ui.GridItem
		if battery {
			cpuRow = ui.NewRow(1.0/3,
				ui.NewCol(2.0/3, cpu),
				ui.NewCol(1.0/3, batt),
			)
		} else {
			cpuRow = ui.NewRow(1.0/3, cpu)
		}
		grid.Set(
			cpuRow,
			ui.NewRow(1.0/3,
				ui.NewCol(1.0/3,
					ui.NewRow(1.0/2, disk),
					ui.NewRow(1.0/2, temp),
				),
				ui.NewCol(2.0/3, mem),
			),
			ui.NewRow(1.0/3,
				ui.NewCol(1.0/2, net),
				ui.NewCol(1.0/2, proc),
			),
		)
	}
}

func setDefaultTermuiColors() {
	ui.Theme.Default = ui.NewStyle(ui.Color(colorscheme.Fg), ui.Color(colorscheme.Bg))
	ui.Theme.Block.Title = ui.NewStyle(ui.Color(colorscheme.BorderLabel), ui.Color(colorscheme.Bg))
	ui.Theme.Block.Border = ui.NewStyle(ui.Color(colorscheme.BorderLine), ui.Color(colorscheme.Bg))
}

func setWidgetColors() {
	mem.LineColor["Main"] = ui.Color(colorscheme.MainMem)
	mem.LineColor["Swap"] = ui.Color(colorscheme.SwapMem)

	proc.CursorColor = ui.Color(colorscheme.ProcCursor)

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
		cpu.LineColor[v] = ui.Color(c)
		i++
	}

	if !minimalMode {
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
				batt.LineColor[v] = ui.Color(c)
				i++
			}
		}

		temp.TempLow = ui.Color(colorscheme.TempLow)
		temp.TempHigh = ui.Color(colorscheme.TempHigh)

		net.Lines[0].LineColor = ui.Color(colorscheme.Sparkline)
		net.Lines[0].TitleColor = ui.Color(colorscheme.BorderLabel)
		net.Lines[1].LineColor = ui.Color(colorscheme.Sparkline)
		net.Lines[1].TitleColor = ui.Color(colorscheme.BorderLabel)
	}
}

func initWidgets() {
	cpu = w.NewCPU(&renderLock, updateInterval, graphHorizontalScale, averageLoad, percpuLoad)
	mem = w.NewMem(&renderLock, updateInterval, graphHorizontalScale)
	proc = w.NewProc(&renderLock)
	help = w.NewHelpMenu()
	if !minimalMode {
		if battery {
			batt = w.NewBatt(&renderLock, graphHorizontalScale)
		}
		net = w.NewNet(&renderLock)
		disk = w.NewDisk(&renderLock)
		temp = w.NewTemp(&renderLock, fahrenheit)
	}
	if statusbar {
		bar = w.NewStatusBar()
	}
}

func render(drawable ...ui.Drawable) {
	renderLock.Lock()
	ui.Render(drawable...)
	renderLock.Unlock()
}

func eventLoop() {
	drawTicker := time.NewTicker(updateInterval).C

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
				render(grid)
			}
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "?":
				helpVisible = !helpVisible
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				termWidth, termHeight := payload.Width, payload.Height
				if statusbar {
					grid.SetRect(0, 0, termWidth, termHeight-1)
					bar.SetRect(0, termHeight-1, termWidth, termHeight)
				} else {
					grid.SetRect(0, 0, payload.Width, payload.Height)
				}
				help.Resize(payload.Width, payload.Height)
				ui.Clear()
			}

			if helpVisible {
				switch e.ID {
				case "?":
					ui.Clear()
					ui.Render(help)
				case "<Escape>":
					helpVisible = false
					render(grid)
				case "<Resize>":
					ui.Render(help)
				}
			} else {
				switch e.ID {
				case "?":
					render(grid)
				case "h":
					graphHorizontalScale += graphHorizontalScaleDelta
					cpu.HorizontalScale = graphHorizontalScale
					mem.HorizontalScale = graphHorizontalScale
					render(cpu, mem)
				case "l":
					if graphHorizontalScale > graphHorizontalScaleDelta {
						graphHorizontalScale -= graphHorizontalScaleDelta
						cpu.HorizontalScale = graphHorizontalScale
						mem.HorizontalScale = graphHorizontalScale
						render(cpu, mem)
					}
				case "<Resize>":
					render(grid)
					if statusbar {
						render(bar)
					}
				case "<MouseLeft>":
					payload := e.Payload.(ui.Mouse)
					proc.Click(payload.X, payload.Y)
					render(proc)
				case "k", "<Up>", "<MouseWheelUp>":
					proc.Up()
					render(proc)
				case "j", "<Down>", "<MouseWheelDown>":
					proc.Down()
					render(proc)
				case "<Home>":
					proc.Top()
					render(proc)
				case "g":
					if previousKey == "g" {
						proc.Top()
						render(proc)
					}
				case "G", "<End>":
					proc.Bottom()
					render(proc)
				case "<C-d>":
					proc.HalfPageDown()
					render(proc)
				case "<C-u>":
					proc.HalfPageUp()
					render(proc)
				case "<C-f>":
					proc.PageDown()
					render(proc)
				case "<C-b>":
					proc.PageUp()
					render(proc)
				case "d":
					if previousKey == "d" {
						proc.Kill()
					}
				case "<Tab>":
					proc.Tab()
					render(proc)
				case "m", "c", "p":
					proc.ChangeSort(e)
					render(proc)
				}

				if previousKey == e.ID {
					previousKey = ""
				} else {
					previousKey = e.ID
				}
			}

		}
	}
}

func setupLogfile() (*os.File, error) {
	// create the log directory
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to make the log directory: %v", err)
	}
	// open the log file
	logfile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0660)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	// log time, filename, and line number
	log.SetFlags(log.Ltime | log.Lshortfile)
	// log to file
	log.SetOutput(logfile)

	return logfile, nil
}

func main() {
	if err := parseArgs(); err != nil {
		stderrLogger.Fatalf("failed to parse cli args: %v", err)
	}

	logfile, err := setupLogfile()
	if err != nil {
		stderrLogger.Fatalf("failed to setup log file: %v", err)
	}
	defer logfile.Close()

	if err := ui.Init(); err != nil {
		stderrLogger.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	logging.StderrToLogfile(logfile)

	setDefaultTermuiColors() // done before initializing widgets to allow inheriting colors
	initWidgets()
	setWidgetColors()

	setupGrid()

	termWidth, termHeight := ui.TerminalDimensions()
	if statusbar {
		grid.SetRect(0, 0, termWidth, termHeight-1)
	} else {
		grid.SetRect(0, 0, termWidth, termHeight)
	}
	help.Resize(termWidth, termHeight)

	ui.Render(grid)
	if statusbar {
		bar.SetRect(0, termHeight-1, termWidth, termHeight)
		ui.Render(bar)
	}

	eventLoop()
}
