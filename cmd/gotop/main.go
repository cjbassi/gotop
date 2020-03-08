package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"plugin"
	"strconv"
	"strings"
	"syscall"
	"time"

	docopt "github.com/docopt/docopt.go"
	ui "github.com/gizak/termui/v3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shibukawa/configdir"

	"github.com/xxxserxxx/gotop/v3"
	"github.com/xxxserxxx/gotop/v3/colorschemes"
	"github.com/xxxserxxx/gotop/v3/layout"
	"github.com/xxxserxxx/gotop/v3/logging"
	w "github.com/xxxserxxx/gotop/v3/widgets"
)

const (
	appName = "gotop"

	graphHorizontalScaleDelta = 3
	defaultUI                 = "2:cpu\ndisk/1 2:mem/2\ntemp\n2:net 2:procs"
	minimalUI                 = "cpu\nmem procs"
	batteryUI                 = "cpu/2 batt/1\ndisk/1 2:mem/2\ntemp\nnet procs"
	procsUI                   = "cpu 4:procs\ndisk\nmem\nnet"
	kitchensink               = "3:cpu/2 3:mem/1\n4:temp/1 3:disk/2\npower\n3:net 3:procs"
)

var (
	// TODO: Set this at compile time; having to check this in sucks.
	Version      = "3.5.0"
	conf         gotop.Config
	help         *w.HelpMenu
	bar          *w.StatusBar
	statusbar    bool
	stderrLogger = log.New(os.Stderr, "", 0)
)

// TODO: Add tab completion for Linux https://gist.github.com/icholy/5314423
// TODO: state:merge #135 linux console font (cmatsuoka/console-font)
// TODO: state:deferred 157 FreeBSD fixes & Nvidia GPU support (kraust/master). Significant CPU use impact for NVidia changes.
// TODO: Virtual devices from Prometeus metrics @feature
// TODO: Abstract out the UI toolkit.  mum4k/termdash, VladimirMarkelov/clui, gcla/gowid, rivo/tview, marcusolsson/tui-go might work better for some OS/Archs. Performance/memory use comparison would be interesting.
func parseArgs(conf *gotop.Config) error {
	usage := `
Usage: gotop [options]

Options:
  -c, --color=NAME        Set a colorscheme.
  -h, --help              Show this screen.
  -m, --minimal           Only show CPU, Mem and Process widgets. Overrides -l. (DEPRECATED, use -l minimal)
  -r, --rate=RATE         Number of times per second to update CPU and Mem widgets [default: 1].
  -V, --version           Print version and exit.
  -p, --percpu            Show each CPU in the CPU widget.
  -a, --averagecpu        Show average CPU in the CPU widget.
  -f, --fahrenheit        Show temperatures in fahrenheit.
  -s, --statusbar         Show a statusbar with the time.
  -b, --battery           Show battery level widget ('minimal' turns off). (DEPRECATED, use -l battery)
  -B, --bandwidth=bits	  Specify the number of bits per seconds.
  -l, --layout=NAME       Name of layout spec file for the UI.  Looks first in $XDG_CONFIG_HOME/gotop, then as a path.  Use "-" to pipe.
  -i, --interface=NAME    Select network interface [default: all]. Several interfaces can be defined using comma separated values. Interfaces can also be ignored using !  
  -x, --export=PORT       Enable metrics for export on the specified port.
  -X, --extensions=NAMES  Enables the listed extensions.  This is a comma-separated list without the .so suffix. The current and config directories will be searched.  
      --test              Runs tests and exits with success/failure code  
      --print-paths       List out the paths that gotop will look for gotop.conf, layouts, color schemes, and extensions  
	  --print-keys        Show the keyboard bindings  


Built-in layouts:
  default
  minimal
  battery
  kitchensink

Colorschemes:
  default
  default-dark (for white background)
  solarized
  solarized16-dark
  solarized16-light
  monokai
  vice
`

	var err error
	conf.Colorscheme, err = colorschemes.FromName(conf.ConfigDir, "default")
	if err != nil {
		return err
	}

	args, err := docopt.ParseArgs(usage, os.Args[1:], Version)
	if err != nil {
		return err
	}

	if val, _ := args["--layout"]; val != nil {
		conf.Layout = val.(string)
	}
	if val, _ := args["--color"]; val != nil {
		cs, err := colorschemes.FromName(conf.ConfigDir, val.(string))
		if err != nil {
			return err
		}
		conf.Colorscheme = cs
	}
	if args["--averagecpu"].(bool) {
		conf.AverageLoad, _ = args["--averagecpu"].(bool)
	}
	if args["--percpu"].(bool) {
		conf.PercpuLoad, _ = args["--percpu"].(bool)
	}
	if args["--statusbar"].(bool) {
		statusbar, _ = args["--statusbar"].(bool)
	}
	if args["--battery"].(bool) {
		conf.Layout = "battery"
	}
	if args["--minimal"].(bool) {
		conf.Layout = "minimal"
	}
	if val, _ := args["--export"]; val != nil {
		conf.ExportPort = val.(string)
	}
	if val, _ := args["--rate"]; val != nil {
		rateStr, _ := val.(string)
		rate, err := strconv.ParseFloat(rateStr, 64)
		if err != nil {
			return fmt.Errorf("invalid rate parameter")
		}
		if rate < 1 {
			conf.UpdateInterval = time.Second * time.Duration(1/rate)
		} else {
			conf.UpdateInterval = time.Second / time.Duration(rate)
		}
	}
	if val, _ := args["--fahrenheit"]; val != nil {
		fahrenheit, _ := val.(bool)
		if fahrenheit {
			conf.TempScale = w.Fahrenheit
		}
	}
	if val, _ := args["--interface"]; val != nil {
		conf.NetInterface, _ = args["--interface"].(string)
	}
	if val, _ := args["--extensions"]; val != nil {
		exs, _ := args["--extensions"].(string)
		conf.Extensions = strings.Split(exs, ",")
	}
	if val, _ := args["--test"]; val != nil {
		conf.Test = val.(bool)
	}
	if args["--print-paths"].(bool) {
		paths := make([]string, 0)
		for _, d := range conf.ConfigDir.QueryFolders(configdir.All) {
			paths = append(paths, d.Path)
		}
		fmt.Println(strings.Join(paths, "\n"))
		os.Exit(0)
	}
	if args["--print-keys"].(bool) {
		fmt.Println(`
Quit: q or <C-c>
Process navigation:
    k and <Up>: up
    j and <Down>: down
    <C-u>: half page up
    <C-d>: half page down
    <C-b>: full page up
    <C-f>: full page down
    gg and <Home>: jump to top
    G and <End>: jump to bottom
Process actions:
    <Tab>: toggle process grouping
    dd: kill selected process or group of processes with SIGTERM
    d3: kill selected process or group of processes with SIGQUIT
    d9: kill selected process or group of processes with SIGKILL
Process sorting
    c: CPU
    m: Mem
    p: PID
Process filtering:
    /: start editing filter
    (while editing):
        <Enter> accept filter
        <C-c> and <Escape>: clear filter
CPU and Mem graph scaling:
    h: scale in
    l: scale out
?: toggles keybind help menu`)
		os.Exit(0)
	}

	return nil
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
			if grid.Proc != nil && grid.Proc.HandleEvent(e) {
				ui.Render(grid.Proc)
				break
			}
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
				case "/":
					if grid.Proc != nil {
						grid.Proc.SetEditingFilter(true)
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

func makeConfig() gotop.Config {
	cd := configdir.New("", appName)
	cd.LocalPath, _ = filepath.Abs(".")
	conf = gotop.Config{
		ConfigDir:            cd,
		GraphHorizontalScale: 7,
		HelpVisible:          false,
		UpdateInterval:       time.Second,
		AverageLoad:          false,
		PercpuLoad:           true,
		TempScale:            w.Celsius,
		Statusbar:            false,
		NetInterface:         w.NET_INTERFACE_ALL,
		MaxLogSize:           5000000,
		Layout:               "default",
	}
	return conf
}

// TODO: mpd visualizer widget
func main() {
	// This is just to make sure gotop returns a useful exit code, but also
	// executes all defer statements and so cleans up before exit.  Sort of
	// annoying work-around for a lack of a clean way to exit Go programs
	// with exit codes.
	ec := run()
	os.Exit(ec)
}

func run() int {
	// Set up default config
	conf := makeConfig()
	// Find the config file; look in (1) local, (2) user, (3) global
	err := conf.Load()
	if err != nil {
		stderrLogger.Printf("failed to parse config file: %s", err)
	}
	// Override with command line arguments
	err = parseArgs(&conf)
	if err != nil {
		stderrLogger.Fatalf("failed to parse cli args: %v", err)
	}

	logfile, err := logging.New(conf)
	if err != nil {
		fmt.Printf("failed to setup log file: %v\n", err)
		return 1
	}
	defer logfile.Close()

	lstream, err := getLayout(conf)
	if err != nil {
		fmt.Printf("failed to find layou: %s\n", err)
		return 1
	}
	ly := layout.ParseLayout(lstream)

	err = loadExtensions(conf)
	if err != nil {
		return 1
	}

	if conf.Test {
		return runTests(conf)
	}

	if err := ui.Init(); err != nil {
		stderrLogger.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	setDefaultTermuiColors(conf) // done before initializing widgets to allow inheriting colors
	help = w.NewHelpMenu()
	if statusbar {
		bar = w.NewStatusBar()
	}

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

	if conf.ExportPort != "" {
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			http.ListenAndServe(conf.ExportPort, nil)
		}()
	}
	eventLoop(conf, grid)
	return 0
}

func getLayout(conf gotop.Config) (io.Reader, error) {
	switch conf.Layout {
	case "-":
		return os.Stdin, nil
	case "default":
		return strings.NewReader(defaultUI), nil
	case "minimal":
		return strings.NewReader(minimalUI), nil
	case "battery":
		return strings.NewReader(batteryUI), nil
	case "procs":
		return strings.NewReader(procsUI), nil
	case "kitchensink":
		return strings.NewReader(kitchensink), nil
	default:
		folder := conf.ConfigDir.QueryFolderContainsFile(conf.Layout)
		if folder == nil {
			paths := make([]string, 0)
			for _, d := range conf.ConfigDir.QueryFolders(configdir.Existing) {
				paths = append(paths, d.Path)
			}
			return nil, fmt.Errorf("unable find layout file %s in %s", conf.Layout, strings.Join(paths, ", "))
		}
		lo, err := folder.ReadFile(conf.Layout)
		if err != nil {
			return nil, err
		}
		return strings.NewReader(string(lo)), nil
	}
}

func loadExtensions(conf gotop.Config) error {
	var hasError bool
	for _, ex := range conf.Extensions {
		exf := ex + ".so"
		fn := exf
		folder := conf.ConfigDir.QueryFolderContainsFile(fn)
		if folder == nil {
			paths := make([]string, 0)
			for _, d := range conf.ConfigDir.QueryFolders(configdir.Existing) {
				paths = append(paths, d.Path)
			}
			log.Printf("unable find extension %s in %s", fn, strings.Join(paths, ", "))
			continue
		}
		fp := filepath.Join(folder.Path, fn)
		p, err := plugin.Open(fp)
		if err != nil {
			hasError = true
			log.Printf(err.Error())
			continue
		}
		init, err := p.Lookup("Init")
		if err != nil {
			hasError = true
			log.Printf(err.Error())
			continue
		}
		initFunc, ok := init.(func())
		if !ok {
			hasError = true
			log.Printf(err.Error())
			continue
		}
		initFunc()
	}
	if hasError {
		folder := conf.ConfigDir.QueryFolderContainsFile(logging.LOGFILE)
		var err error
		if folder == nil {
			err = fmt.Errorf("error initializing requested plugins\n")
		} else {
			err = fmt.Errorf("error initializing requested plugins; check the log file %s\n", filepath.Join(folder.Path, logging.LOGFILE))
		}
		return err
	}
	return nil
}

func runTests(conf gotop.Config) int {
	fmt.Printf("PASS")
	return 0
}
