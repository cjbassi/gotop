package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	//_ "net/http/pprof"

	ui "github.com/gizak/termui/v3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shibukawa/configdir"
	flag "github.com/xxxserxxx/opflag"

	"github.com/xxxserxxx/gotop/v4"
	"github.com/xxxserxxx/gotop/v4/colorschemes"
	"github.com/xxxserxxx/gotop/v4/devices"
	"github.com/xxxserxxx/gotop/v4/layout"
	"github.com/xxxserxxx/gotop/v4/logging"
	"github.com/xxxserxxx/gotop/v4/widgets"
	w "github.com/xxxserxxx/gotop/v4/widgets"
)

const (
	graphHorizontalScaleDelta = 3
	defaultUI                 = "2:cpu\ndisk/1 2:mem/2\ntemp\n2:net 2:procs"
	minimalUI                 = "cpu\nmem procs"
	batteryUI                 = "cpu/2 batt/1\ndisk/1 2:mem/2\ntemp\nnet procs"
	procsUI                   = "cpu 4:procs\ndisk\nmem\nnet"
	kitchensink               = "3:cpu/2 3:mem/1\n4:temp/1 3:disk/2\npower\n3:net 3:procs"
)

var (
	// Version of the program; set during build from git tags
	Version = "0.0.0"
	// BuildDate when the program was compiled; set during build
	BuildDate    = "Hadean"
	conf         gotop.Config
	help         *w.HelpMenu
	bar          *w.StatusBar
	statusbar    bool
	stderrLogger = log.New(os.Stderr, "", 0)
)

func parseArgs(conf *gotop.Config) error {
	cds := conf.ConfigDir.QueryFolders(configdir.All)
	cpaths := make([]string, len(cds))
	for i, p := range cds {
		cpaths[i] = p.Path
	}
	help := flag.BoolP("help", "h", false, "Show this screen.")
	color := flag.StringP("color", "c", conf.Colorscheme.Name, "Set a colorscheme.")
	flag.IntVarP(&conf.GraphHorizontalScale, "graphscale", "S", conf.GraphHorizontalScale, "Graph scale factor, >0")
	version := flag.BoolP("version", "v", false, "Print version and exit.")
	versioN := flag.BoolP("", "V", false, "Print version and exit.")
	flag.BoolVarP(&conf.PercpuLoad, "percpu", "p", conf.PercpuLoad, "Show each CPU in the CPU widget.")
	flag.BoolVarP(&conf.AverageLoad, "averagecpu", "a", conf.AverageLoad, "Show average CPU in the CPU widget.")
	fahrenheit := flag.BoolP("fahrenheit", "f", conf.TempScale == 'F', "Show temperatures in fahrenheit.Show temperatures in fahrenheit.")
	flag.BoolVarP(&conf.Statusbar, "statusbar", "s", conf.Statusbar, "Show a statusbar with the time.")
	flag.DurationVarP(&conf.UpdateInterval, "rate", "r", conf.UpdateInterval, "Number of times per second to update CPU and Mem widgets.")
	flag.StringVarP(&conf.Layout, "layout", "l", conf.Layout, `Name of layout spec file for the UI. Use "-" to pipe.`)
	flag.StringVarP(&conf.NetInterface, "interface", "i", "all", "Select network interface. Several interfaces can be defined using comma separated values. Interfaces can also be ignored using `!`")
	flag.StringVarP(&conf.ExportPort, "export", "x", conf.ExportPort, "Enable metrics for export on the specified port.")
	flag.BoolVarP(&conf.Mbps, "mbps", "", conf.Mbps, "Show network rate as mbps.")
	// FIXME Where did this go??
	//conf.Band = flag.IntP("bandwidth", "B", 100, "Specify the number of bits per seconds.")
	flag.BoolVar(&conf.Test, "test", conf.Test, "Runs tests and exits with success/failure code.")
	list := flag.String("list", "", `List <devices|layouts|colorschemes|paths|keys>
         devices: Prints out device names for filterable widgets
         layouts: Lists build-in layouts
         colorschemes: Lists built-in colorschemes
         paths: List out configuration file search paths
         keys: Show the keyboard bindings.`)
	wc := flag.Bool("write-config", false, "Write out a default config file.")
	flag.SortFlags = false
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if *version || *versioN {
		fmt.Printf("gotop %s (%s)\n", Version, BuildDate)
		os.Exit(0)
	}
	if *help {
		flag.Usage()
		os.Exit(0)
	}
	cs, err := colorschemes.FromName(conf.ConfigDir, *color)
	if err != nil {
		return err
	}
	conf.Colorscheme = cs
	if *fahrenheit {
		conf.TempScale = 'F'
	} else {
		conf.TempScale = 'C'
	}
	if *list != "" {
		switch *list {
		case "layouts":
			fmt.Println(_layouts)
		case "colorschemes":
			fmt.Println(_colorschemes)
		case "paths":
			fmt.Println("Loadable colorschemes & layouts, and the config file, are searched for, in order:")
			paths := make([]string, 0)
			for _, d := range conf.ConfigDir.QueryFolders(configdir.All) {
				paths = append(paths, d.Path)
			}
			fmt.Println(strings.Join(paths, "\n"))
			fmt.Printf("\nThe log file is in %s\n", filepath.Join(conf.ConfigDir.QueryCacheFolder().Path, logging.LOGFILE))
		case "devices":
			listDevices()
		case "keys":
			fmt.Println(widgets.KEYBINDS)
		default:
			fmt.Printf("Unknown option \"%s\"; try layouts, colorschemes, keys, paths, or devices\n", *list)
			os.Exit(1)
		}
		os.Exit(0)
	}
	if *wc {
		path, err := conf.Write()
		if err != nil {
			fmt.Printf("Failed to write configuration file: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("Config written to %s\n", path)
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
				case "b":
					if grid.Net != nil {
						grid.Net.Mbps = !grid.Net.Mbps
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

// TODO: Add fans
// TODO: mpd visualizer widget
// TODO: Add tab completion for Linux https://gist.github.com/icholy/5314423
// TODO: state:merge #135 linux console font (cmatsuoka/console-font)
// TODO: Abstract out the UI toolkit.  mum4k/termdash, VladimirMarkelov/clui, gcla/gowid, rivo/tview, marcusolsson/tui-go might work better for some OS/Archs. Performance/memory use comparison would be interesting.
// TODO: all of the go vet stuff, more unit tests, benchmarks, finish remote.
// TODO: color bars for memory, a-la bashtop
func main() {
	// For performance testing
	//go func() {
	//	log.Fatal(http.ListenAndServe(":7777", nil))
	//}()

	// This is just to make sure gotop returns a useful exit code, but also
	// executes all defer statements and so cleans up before exit.  Sort of
	// annoying work-around for a lack of a clean way to exit Go programs
	// with exit codes.
	ec := run()
	if ec > 0 {
		if ec < 2 {
			fmt.Printf("errors encountered; check the log file %s\n", filepath.Join(conf.ConfigDir.QueryCacheFolder().Path, logging.LOGFILE))
		}
	}
	os.Exit(ec)
}

func run() int {
	conf := gotop.NewConfig()
	// Find the config file; look in (1) local, (2) user, (3) global
	err := conf.Load()
	if err != nil {
		fmt.Printf("failed to parse config file: %s\n", err)
		return 2
	}
	// Override with command line arguments
	err = parseArgs(&conf)
	if err != nil {
		fmt.Printf("parsing CLI args: %s\n", err)
		return 2
	}

	logfile, err := logging.New(conf)
	if err != nil {
		fmt.Printf("failed to setup log file: %v\n", err)
		return 2
	}
	defer logfile.Close()

	lstream, err := getLayout(conf)
	if err != nil {
		stderrLogger.Print(err)
		return 1
	}
	ly := layout.ParseLayout(lstream)

	if conf.Test {
		return runTests(conf)
	}

	errs := devices.Startup(conf.ExtensionVars)
	if len(errs) > 0 {
		for _, err := range errs {
			stderrLogger.Print(err)
		}
		return 1
	}
	if err = ui.Init(); err != nil {
		stderrLogger.Print(err)
		return 1
	}
	defer ui.Close()

	setDefaultTermuiColors(conf) // done before initializing widgets to allow inheriting colors
	help = w.NewHelpMenu()
	if statusbar {
		bar = w.NewStatusBar()
	}

	grid, err := layout.Layout(ly, conf)
	if err != nil {
		stderrLogger.Print(err)
		return 1
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

func runTests(conf gotop.Config) int {
	fmt.Printf("PASS")
	return 0
}

func listDevices() {
	ms := devices.Domains
	sort.Strings(ms)
	for _, m := range ms {
		fmt.Printf("%s:\n", m)
		for _, d := range devices.Devices(m, true) {
			fmt.Printf("\t%s\n", d)
		}
	}
}

const _layouts = `Built-in layouts:
   default
   minimal
   battery
   kitchensink`
const _colorschemes = `Built-in colorschemes:
   default
   default-dark (for white background)
   solarized
   solarized16-dark
   solarized16-light
   monokai
   vice`
