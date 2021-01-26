package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
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

	"github.com/VictoriaMetrics/metrics"
	jj "github.com/cloudfoundry-attic/jibber_jabber"
	ui "github.com/gizak/termui/v3"
	"github.com/jdkeke142/lingo-toml"
	"github.com/shibukawa/configdir"
	"github.com/xxxserxxx/opflag"

	"github.com/xxxserxxx/gotop/v4"
	"github.com/xxxserxxx/gotop/v4/colorschemes"
	"github.com/xxxserxxx/gotop/v4/devices"
	"github.com/xxxserxxx/gotop/v4/layout"
	"github.com/xxxserxxx/gotop/v4/logging"
	"github.com/xxxserxxx/gotop/v4/translations"
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
	tr           lingo.Translations
)

func parseArgs() error {
	cds := conf.ConfigDir.QueryFolders(configdir.All)
	cpaths := make([]string, len(cds))
	for i, p := range cds {
		cpaths[i] = p.Path
	}
	help := opflag.BoolP("help", "h", false, tr.Value("args.help"))
	color := opflag.StringP("color", "c", conf.Colorscheme.Name, tr.Value("args.color"))
	opflag.IntVarP(&conf.GraphHorizontalScale, "graphscale", "S", conf.GraphHorizontalScale, tr.Value("args.scale"))
	version := opflag.BoolP("version", "v", false, tr.Value("args.version"))
	versioN := opflag.BoolP("", "V", false, tr.Value("args.version"))
	opflag.BoolVarP(&conf.PercpuLoad, "percpu", "p", conf.PercpuLoad, tr.Value("args.percpu"))
	opflag.BoolVarP(&conf.AverageLoad, "averagecpu", "a", conf.AverageLoad, tr.Value("args.cpuavg"))
	fahrenheit := opflag.BoolP("fahrenheit", "f", conf.TempScale == 'F', tr.Value("args.temp"))
	opflag.BoolVarP(&conf.Statusbar, "statusbar", "s", conf.Statusbar, tr.Value("args.statusbar"))
	opflag.DurationVarP(&conf.UpdateInterval, "rate", "r", conf.UpdateInterval, tr.Value("args.rate"))
	opflag.StringVarP(&conf.Layout, "layout", "l", conf.Layout, tr.Value("args.layout"))
	opflag.StringVarP(&conf.NetInterface, "interface", "i", "all", tr.Value("args.net"))
	opflag.StringVarP(&conf.ExportPort, "export", "x", conf.ExportPort, tr.Value("args.export"))
	opflag.BoolVarP(&conf.Mbps, "mbps", "", conf.Mbps, tr.Value("args.mbps"))
	opflag.BoolVar(&conf.Test, "test", conf.Test, tr.Value("args.test"))
	opflag.StringP("", "C", "", tr.Value("args.conffile"))
	list := opflag.String("list", "", tr.Value("args.list"))
	wc := opflag.Bool("write-config", false, tr.Value("args.write"))
	opflag.SortFlags = false
	opflag.Usage = func() {
		fmt.Fprintf(os.Stderr, tr.Value("usage", os.Args[0]))
		opflag.PrintDefaults()
	}
	opflag.Parse()
	if *version || *versioN {
		fmt.Printf("gotop %s (%s)\n", Version, BuildDate)
		os.Exit(0)
	}
	if *help {
		opflag.Usage()
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
			fmt.Println(tr.Value("help.layouts"))
		case "colorschemes":
			fmt.Println(tr.Value("help.colorschemes"))
		case "paths":
			fmt.Println(tr.Value("help.paths"))
			paths := make([]string, 0)
			for _, d := range conf.ConfigDir.QueryFolders(configdir.All) {
				paths = append(paths, d.Path)
			}
			fmt.Println(strings.Join(paths, "\n"))
			fmt.Println()
			fmt.Println(tr.Value("help.log", filepath.Join(conf.ConfigDir.QueryCacheFolder().Path, logging.LOGFILE)))
		case "devices":
			listDevices()
		case "keys":
			fmt.Println(tr.Value("help.help"))
		case "widgets":
			fmt.Println(tr.Value("help.widgets"))
		case "langs":
			vs, err := translations.AssetDir("")
			if err != nil {
				return err
			}
			for _, v := range vs {
				v = strings.TrimSuffix(v, ".toml")
				fmt.Println(v)
			}
		default:
			fmt.Printf(tr.Value("error.unknownopt", *list))
			os.Exit(1)
		}
		os.Exit(0)
	}
	if *wc {
		path, err := conf.Write()
		if err != nil {
			fmt.Println(tr.Value("error.writefail", err.Error()))
			os.Exit(1)
		}
		fmt.Println(tr.Value("help.written", path))
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

func main() {
	// TODO: Make this an option, for performance testing
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
			logpath := filepath.Join(conf.ConfigDir.QueryCacheFolder().Path, logging.LOGFILE)
			fmt.Println(tr.Value("error.checklog", logpath))
			bs, _ := ioutil.ReadFile(logpath)
			fmt.Println(string(bs))
		}
	}
	os.Exit(ec)
}

func run() int {
	ling, err := lingo.New("en_US", "", translations.AssetFile())
	if err != nil {
		fmt.Printf("failed to load language files: %s\n", err)
		return 2
	}
	lang, err := jj.DetectIETF()
	if err != nil {
		lang = "en_US"
	}
	lang = strings.Replace(lang, "-", "_", -1)
	// Get the locale from the os
	tr = ling.TranslationsForLocale(lang)
	colorschemes.SetTr(tr)
	conf = gotop.NewConfig()
	conf.Tr = tr
	// Find the config file; look in (1) local, (2) user, (3) global
	// Check the last argument first
	fs := flag.NewFlagSet("config", flag.ContinueOnError)
	cfg := fs.String("C", "", tr.Value("configfile"))
	fs.SetOutput(bufio.NewWriter(nil))
	fs.Parse(os.Args[1:])
	if *cfg != "" {
		conf.ConfigFile = *cfg
	}
	err = conf.Load()
	if err != nil {
		fmt.Println(tr.Value("error.configparse", err.Error()))
		return 2
	}
	// Override with command line arguments
	err = parseArgs()
	if err != nil {
		fmt.Println(tr.Value("error.cliparse", err.Error()))
		return 2
	}

	logfile, err := logging.New(conf)
	if err != nil {
		fmt.Println(tr.Value("logsetup", err.Error()))
		return 2
	}
	defer logfile.Close()

	errs := devices.Startup(conf.ExtensionVars)
	if len(errs) > 0 {
		for _, err := range errs {
			stderrLogger.Print(err)
		}
		return 1
	}

	lstream, err := getLayout(conf)
	if err != nil {
		stderrLogger.Print(err)
		return 1
	}
	ly := layout.ParseLayout(lstream)

	if conf.Test {
		return runTests(conf)
	}

	if err = ui.Init(); err != nil {
		stderrLogger.Print(err)
		return 1
	}
	defer ui.Close()

	setDefaultTermuiColors(conf) // done before initializing widgets to allow inheriting colors
	help = w.NewHelpMenu(tr)
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

	// TODO https://godoc.org/github.com/VictoriaMetrics/metrics#Set
	if conf.ExportPort != "" {
		go func() {
			http.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
				metrics.WritePrometheus(w, true)
			})
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
			return nil, fmt.Errorf(tr.Value("error.findlayout", conf.Layout, strings.Join(paths, ", ")))
		}
		lo, err := folder.ReadFile(conf.Layout)
		if err != nil {
			return nil, err
		}
		return strings.NewReader(string(lo)), nil
	}
}

func runTests(_ gotop.Config) int {
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
