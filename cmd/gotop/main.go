package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	docopt "github.com/docopt/docopt.go"
	ui "github.com/gizak/termui/v3"

	"github.com/xxxserxxx/gotop"
	"github.com/xxxserxxx/gotop/colorschemes"
	"github.com/xxxserxxx/gotop/layout"
	"github.com/xxxserxxx/gotop/logging"
	"github.com/xxxserxxx/gotop/utils"
	w "github.com/xxxserxxx/gotop/widgets"
)

const (
	appName = "gotop"
	version = "3.1.1"

	graphHorizontalScaleDelta = 3
	defaultUI                 = "cpu\ndisk/1 2:mem/2\ntemp\nnet procs"
	minimalUI                 = "cpu\nmem procs"
)

var (
	conf         gotop.Config
	help         *w.HelpMenu
	bar          *w.StatusBar
	statusbar    bool
	stderrLogger = log.New(os.Stderr, "", 0)
)

func parseArgs() (gotop.Config, error) {
	usage := `
Usage: gotop [options]

Options:
  -c, --color=NAME        Set a colorscheme.
  -h, --help              Show this screen.
  -m, --minimal           Only show CPU, Mem and Process widgets. Overrides -l
  -r, --rate=RATE         Number of times per second to update CPU and Mem widgets [default: 1].
  -V, --version           Print version and exit.
  -p, --percpu            Show each CPU in the CPU widget.
  -a, --averagecpu        Show average CPU in the CPU widget.
  -f, --fahrenheit        Show temperatures in fahrenheit.
  -s, --statusbar         Show a statusbar with the time.
  -b, --battery           Show battery level widget ('minimal' turns off).
  -i, --interface=NAME    Select network interface [default: all].
  -B, --bandwidth=bits	  Specify the number of bits per seconds.
  -l, --layout=NAME       Name of layout spec file for the UI.  Looks first in $XDG_CONFIG_HOME/gotop, then as a path.  Use "-" to pipe.

Colorschemes:
  default
  default-dark (for white background)
  solarized
  solarized16-dark
  solarized16-light
  monokai
  vice
`

	ld := utils.GetLogDir(appName)
	cd := utils.GetConfigDir(appName)
	conf = gotop.Config{
		ConfigDir:            cd,
		LogDir:               ld,
		LogPath:              filepath.Join(ld, "errors.log"),
		GraphHorizontalScale: 7,
		HelpVisible:          false,
		Colorscheme:          colorschemes.Default,
		UpdateInterval:       time.Second,
		AverageLoad:          false,
		PercpuLoad:           false,
		TempScale:            w.Celsius,
		Battery:              false,
		Statusbar:            false,
		NetInterface:         w.NET_INTERFACE_ALL,
	}

	args, err := docopt.ParseArgs(usage, os.Args[1:], version)
	if err != nil {
		return conf, err
	}

	if val, _ := args["--layout"]; val != nil {
		s := val.(string)
		if s == "-" {
			conf.Layout = os.Stdin
		} else {
			fp := filepath.Join(cd, s)
			conf.Layout, err = os.Open(fp)
			if err != nil {
				conf.Layout, err = os.Open(s)
				if err != nil {
					stderrLogger.Fatalf("Unable to open layout file %s", fp)
				}
			}
		}
	} else {
		conf.Layout = strings.NewReader(defaultUI)
	}
	if val, _ := args["--color"]; val != nil {
		cs, err := handleColorscheme(val.(string))
		if err != nil {
			return conf, err
		}
		conf.Colorscheme = cs
	}
	conf.AverageLoad, _ = args["--averagecpu"].(bool)
	conf.PercpuLoad, _ = args["--percpu"].(bool)
	conf.Battery, _ = args["--battery"].(bool)
	statusbar, _ = args["--statusbar"].(bool)

	if args["--minimal"].(bool) {
		conf.Layout = strings.NewReader(minimalUI)
	}
	rateStr, _ := args["--rate"].(string)
	rate, err := strconv.ParseFloat(rateStr, 64)
	if err != nil {
		return conf, fmt.Errorf("invalid rate parameter")
	}
	if rate < 1 {
		conf.UpdateInterval = time.Second * time.Duration(1/rate)
	} else {
		conf.UpdateInterval = time.Second / time.Duration(rate)
	}
	fahrenheit, _ := args["--fahrenheit"].(bool)
	if fahrenheit {
		conf.TempScale = w.Fahrenheit
	}

	return conf, nil
}

func handleColorscheme(c string) (colorschemes.Colorscheme, error) {
	var cs colorschemes.Colorscheme
	switch c {
	case "default":
		cs = colorschemes.Default
	case "solarized":
		cs = colorschemes.Solarized
	case "solarized16-light":
		cs = colorschemes.Solarized16Light
	case "solarized16-dark":
		cs = colorschemes.Solarized16Dark
	case "monokai":
		cs = colorschemes.Monokai
	case "vice":
		cs = colorschemes.Vice
	case "default-dark":
		cs = colorschemes.DefaultDark
	case "nord":
		cs = colorschemes.Nord
	default:
		custom, err := getCustomColorscheme(conf, c)
		if err != nil {
			return cs, err
		}
		cs = custom
	}
	return cs, nil
}

// getCustomColorscheme	tries to read a custom json colorscheme from <configDir>/<name>.json
func getCustomColorscheme(c gotop.Config, name string) (colorschemes.Colorscheme, error) {
	var cs colorschemes.Colorscheme
	filePath := filepath.Join(c.ConfigDir, name+".json")
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		return cs, fmt.Errorf("failed to read colorscheme file: %v", err)
	}
	err = json.Unmarshal(dat, &cs)
	if err != nil {
		return cs, fmt.Errorf("failed to parse colorscheme file: %v", err)
	}
	return cs, nil
}

func setDefaultTermuiColors(c gotop.Config) {
	ui.Theme.Default = ui.NewStyle(ui.Color(c.Colorscheme.Fg), ui.Color(c.Colorscheme.Bg))
	ui.Theme.Block.Title = ui.NewStyle(ui.Color(c.Colorscheme.BorderLabel), ui.Color(c.Colorscheme.Bg))
	ui.Theme.Block.Border = ui.NewStyle(ui.Color(c.Colorscheme.BorderLine), ui.Color(c.Colorscheme.Bg))
}

func eventLoop(c gotop.Config, grid *layout.MyGrid) {
	drawTicker := time.NewTicker(c.UpdateInterval).C

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
			if !c.HelpVisible {
				ui.Render(grid)
				if statusbar {
					ui.Render(bar)
				}
			}
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "?":
				c.HelpVisible = !c.HelpVisible
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

			if c.HelpVisible {
				switch e.ID {
				case "?":
					ui.Clear()
					ui.Render(help)
				case "<Escape>":
					c.HelpVisible = false
					ui.Render(grid)
				case "<Resize>":
					ui.Render(help)
				}
			} else {
				switch e.ID {
				case "?":
					ui.Render(grid)
				case "h":
					c.GraphHorizontalScale += graphHorizontalScaleDelta
					for _, item := range grid.Lines {
						item.Scale(c.GraphHorizontalScale)
					}
					ui.Render(grid)
				case "l":
					if c.GraphHorizontalScale > graphHorizontalScaleDelta {
						c.GraphHorizontalScale -= graphHorizontalScaleDelta
						for _, item := range grid.Lines {
							item.Scale(c.GraphHorizontalScale)
							ui.Render(item)
						}
					}
				case "<Resize>":
					ui.Render(grid)
					if statusbar {
						ui.Render(bar)
					}
				case "<MouseLeft>":
					if grid.Proc != nil {
						payload := e.Payload.(ui.Mouse)
						grid.Proc.HandleClick(payload.X, payload.Y)
						ui.Render(grid.Proc)
					}
				case "k", "<Up>", "<MouseWheelUp>":
					if grid.Proc != nil {
						grid.Proc.ScrollUp()
						ui.Render(grid.Proc)
					}
				case "j", "<Down>", "<MouseWheelDown>":
					if grid.Proc != nil {
						grid.Proc.ScrollDown()
						ui.Render(grid.Proc)
					}
				case "<Home>":
					if grid.Proc != nil {
						grid.Proc.ScrollTop()
						ui.Render(grid.Proc)
					}
				case "g":
					if grid.Proc != nil {
						if previousKey == "g" {
							grid.Proc.ScrollTop()
							ui.Render(grid.Proc)
						}
					}
				case "G", "<End>":
					if grid.Proc != nil {
						grid.Proc.ScrollBottom()
						ui.Render(grid.Proc)
					}
				case "<C-d>":
					if grid.Proc != nil {
						grid.Proc.ScrollHalfPageDown()
						ui.Render(grid.Proc)
					}
				case "<C-u>":
					if grid.Proc != nil {
						grid.Proc.ScrollHalfPageUp()
						ui.Render(grid.Proc)
					}
				case "<C-f>":
					if grid.Proc != nil {
						grid.Proc.ScrollPageDown()
						ui.Render(grid.Proc)
					}
				case "<C-b>":
					if grid.Proc != nil {
						grid.Proc.ScrollPageUp()
						ui.Render(grid.Proc)
					}
				case "d":
					if grid.Proc != nil {
						if previousKey == "d" {
							grid.Proc.KillProc("SIGTERM")
						}
					}
				case "3":
					if grid.Proc != nil {
						if previousKey == "d" {
							grid.Proc.KillProc("SIGQUIT")
						}
					}
				case "9":
					if grid.Proc != nil {
						if previousKey == "d" {
							grid.Proc.KillProc("SIGKILL")
						}
					}
				case "<Tab>":
					if grid.Proc != nil {
						grid.Proc.ToggleShowingGroupedProcs()
						ui.Render(grid.Proc)
					}
				case "m", "c", "p":
					if grid.Proc != nil {
						grid.Proc.ChangeProcSortMethod(w.ProcSortMethod(e.ID))
						ui.Render(grid.Proc)
					}
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

func setupLogfile(c gotop.Config) (*os.File, error) {
	// create the log directory
	if err := os.MkdirAll(c.LogDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to make the log directory: %v", err)
	}
	// open the log file
	logfile, err := os.OpenFile(c.LogPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0660)
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
	conf, err := parseArgs()
	if err != nil {
		stderrLogger.Fatalf("failed to parse cli args: %v", err)
	}

	logfile, err := setupLogfile(conf)
	if err != nil {
		stderrLogger.Fatalf("failed to setup log file: %v", err)
	}
	defer logfile.Close()

	if err := ui.Init(); err != nil {
		stderrLogger.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	logging.StderrToLogfile(logfile)

	setDefaultTermuiColors(conf) // done before initializing widgets to allow inheriting colors
	help = w.NewHelpMenu()
	if statusbar {
		bar = w.NewStatusBar()
	}

	ly := layout.ParseLayout(conf.Layout)
	grid, err := layout.Layout(ly, conf)
	if err != nil {
		stderrLogger.Fatalf("failed to initialize termui: %v", err)
	}

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

	eventLoop(conf, grid)
}
