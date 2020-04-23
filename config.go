package gotop

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/shibukawa/configdir"
	"github.com/xxxserxxx/gotop/v4/colorschemes"
	"github.com/xxxserxxx/gotop/v4/widgets"
)

type Config struct {
	ConfigDir configdir.ConfigDir

	GraphHorizontalScale int
	HelpVisible          bool
	Colorscheme          colorschemes.Colorscheme

	UpdateInterval time.Duration
	AverageLoad    bool
	PercpuLoad     bool
	Statusbar      bool
	TempScale      widgets.TempScale
	NetInterface   string
	Layout         string
	MaxLogSize     int64
	ExportPort     string
	Mbps           bool
	Temps          []string

	Test bool
}

func (conf *Config) Load() error {
	var in []byte
	var err error
	cfn := "gotop.conf"
	folder := conf.ConfigDir.QueryFolderContainsFile(cfn)
	if folder != nil {
		if in, err = folder.ReadFile(cfn); err != nil {
			return err
		}
	} else {
		return nil
	}
	return load(bytes.NewReader(in), conf)
}

func load(in io.Reader, conf *Config) error {
	r := bufio.NewScanner(in)
	var lineNo int
	for r.Scan() {
		l := strings.TrimSpace(r.Text())
		if l[0] == '#' {
			continue
		}
		kv := strings.Split(l, "=")
		if len(kv) != 2 {
			return fmt.Errorf("bad config file syntax; should be KEY=VALUE, was %s", l)
		}
		key := strings.ToLower(kv[0])
		switch key {
		default:
			return fmt.Errorf("invalid config option %s", key)
		case "configdir":
			log.Printf("configdir is deprecated.  Ignored configdir=%s", kv[1])
		case "logdir":
			log.Printf("logdir is deprecated.  Ignored logdir=%s", kv[1])
		case "logfile":
			log.Printf("logfile is deprecated.  Ignored logfile=%s", kv[1])
		case graphhorizontalscale:
			iv, err := strconv.Atoi(kv[1])
			if err != nil {
				return err
			}
			conf.GraphHorizontalScale = iv
		case helpvisible:
			bv, err := strconv.ParseBool(kv[1])
			if err != nil {
				return fmt.Errorf("line %d: %v", lineNo, err)
			}
			conf.HelpVisible = bv
		case colorscheme:
			cs, err := colorschemes.FromName(conf.ConfigDir, kv[1])
			if err != nil {
				return fmt.Errorf("line %d: %v", lineNo, err)
			}
			conf.Colorscheme = cs
		case updateinterval:
			iv, err := strconv.Atoi(kv[1])
			if err != nil {
				return err
			}
			conf.UpdateInterval = time.Duration(iv)
		case averagecpu:
			bv, err := strconv.ParseBool(kv[1])
			if err != nil {
				return fmt.Errorf("line %d: %v", lineNo, err)
			}
			conf.AverageLoad = bv
		case percpuload:
			bv, err := strconv.ParseBool(kv[1])
			if err != nil {
				return fmt.Errorf("line %d: %v", lineNo, err)
			}
			conf.PercpuLoad = bv
		case tempscale:
			switch kv[1] {
			case "C":
				conf.TempScale = 'C'
			case "F":
				conf.TempScale = 'F'
			default:
				conf.TempScale = 'C'
				return fmt.Errorf("invalid TempScale value %s", kv[1])
			}
		case statusbar:
			bv, err := strconv.ParseBool(kv[1])
			if err != nil {
				return fmt.Errorf("line %d: %v", lineNo, err)
			}
			conf.Statusbar = bv
		case netinterface:
			// FIXME this should be a comma-separated list
			conf.NetInterface = kv[1]
		case layout:
			conf.Layout = kv[1]
		case maxlogsize:
			iv, err := strconv.Atoi(kv[1])
			if err != nil {
				return err
			}
			conf.MaxLogSize = int64(iv)
		case export:
			conf.ExportPort = kv[1]
		case mbps:
			conf.Mbps = true
		case temperatures:
			conf.Temps = strings.Split(kv[1], ",")
		}
	}

	return nil
}

func (c *Config) Write() (string, error) {
	cfn := "gotop.conf"
	ds := c.ConfigDir.QueryFolders(configdir.Global)
	if len(ds) == 0 {
		ds = c.ConfigDir.QueryFolders(configdir.Local)
		if len(ds) == 0 {
			return "", fmt.Errorf("error locating config folders")
		}
	}
	marshalled := marshal(c)
	err := ds[0].WriteFile(cfn, marshalled)
	if err != nil {
		return "", err
	}
	return filepath.Join(ds[0].Path, cfn), nil
}

func marshal(c *Config) []byte {
	buff := bytes.NewBuffer(nil)
	fmt.Fprintln(buff, "# Scale graphs to this level; 7 is the default, 2 is zoomed out.")
	fmt.Fprintf(buff, "%s=%d\n", graphhorizontalscale, c.GraphHorizontalScale)
	fmt.Fprintln(buff, "# If true, start the UI with the help visible")
	fmt.Fprintf(buff, "%s=%t\n", helpvisible, c.HelpVisible)
	fmt.Fprintln(buff, "# The color scheme to use.  See `--list colorschemes`")
	fmt.Fprintf(buff, "%s=%s\n", colorscheme, c.Colorscheme.Name)
	fmt.Fprintln(buff, "# How frequently to update the UI, in nanoseconds")
	fmt.Fprintf(buff, "%s=%d\n", updateinterval, c.UpdateInterval)
	fmt.Fprintln(buff, "# If true, show the average CPU load")
	fmt.Fprintf(buff, "%s=%t\n", averagecpu, c.AverageLoad)
	fmt.Fprintln(buff, "# If true, show load per CPU")
	fmt.Fprintf(buff, "%s=%t\n", percpuload, c.PercpuLoad)
	fmt.Fprintln(buff, "# Temperature units. C for Celcius, F for Fahrenheit")
	fmt.Fprintf(buff, "%s=%c\n", tempscale, c.TempScale)
	fmt.Fprintln(buff, "# If true, display a status bar")
	fmt.Fprintf(buff, "%s=%t\n", statusbar, c.Statusbar)
	fmt.Fprintln(buff, "# The network interface to monitor")
	fmt.Fprintf(buff, "%s=%s\n", netinterface, c.NetInterface)
	fmt.Fprintln(buff, "# A layout name. See `--list layouts`")
	fmt.Fprintf(buff, "%s=%s\n", layout, c.Layout)
	fmt.Fprintln(buff, "# The maximum log file size, in bytes")
	fmt.Fprintf(buff, "%s=%d\n", maxlogsize, c.MaxLogSize)
	fmt.Fprintln(buff, "# If set, export data as Promethius metrics on the interface:port.\n# E.g., `:8080` (colon is required, interface is not)")
	if c.ExportPort == "" {
		fmt.Fprint(buff, "#")
	}
	fmt.Fprintf(buff, "%s=%s\n", export, c.ExportPort)
	fmt.Fprintln(buff, "# Display network IO in mpbs if true")
	fmt.Fprintf(buff, "%s=%t\n", mbps, c.Mbps)
	fmt.Fprintln(buff, "# A list of enabled temp sensors.  See `--list devices`")
	if len(c.Temps) == 0 {
		fmt.Fprint(buff, "#")
	}
	fmt.Fprintf(buff, "%s=%s\n", temperatures, strings.Join(c.Temps, ","))
	return buff.Bytes()
}

const (
	graphhorizontalscale = "graphhorizontalscale"
	helpvisible          = "helpvisible"
	colorscheme          = "colorscheme"
	updateinterval       = "updateinterval"
	averagecpu           = "averagecpu"
	percpuload           = "percpuload"
	tempscale            = "tempscale"
	statusbar            = "statusbar"
	netinterface         = "netinterface"
	layout               = "layout"
	maxlogsize           = "maxlogsize"
	export               = "metricsexportport"
	mbps                 = "mbps"
	temperatures         = "temperatures"
)
