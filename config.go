package gotop

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/xxxserxxx/gotop/colorschemes"
	"github.com/xxxserxxx/gotop/widgets"
)

// TODO: Cross-compiling for darwin, openbsd requiring native procs & temps
// TODO: Merge #184 or #177 degree symbol (BartWillems:master, fleaz:master)
// TODO: Merge #169 % option for network use (jrswab:networkPercentage)
// TODO: Merge #167 configuration file (jrswab:configFile111)
// TODO: Merge #157 FreeBSD fixes & Nvidia GPU support (kraust:master)
// TODO: Merge #156 Added temperatures for NVidia GPUs (azak-azkaran:master)
// TODO: Merge #135 linux console font (cmatsuoka:console-font)
// TODO: Export Prometheus metrics @feature
// TODO: Virtual devices from Prometeus metrics @feature
type Config struct {
	ConfigDir string
	LogDir    string
	LogFile   string

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
}

func Parse(in io.Reader, conf *Config) error {
	r := bufio.NewScanner(in)
	var lineNo int
	for r.Scan() {
		l := strings.TrimSpace(r.Text())
		kv := strings.Split(l, "=")
		if len(kv) != 2 {
			return fmt.Errorf("bad config file syntax; should be KEY=VALUE, was %s", l)
		}
		key := strings.ToLower(kv[0])
		switch key {
		default:
			return fmt.Errorf("invalid config option %s", key)
		case "configdir":
			conf.ConfigDir = kv[1]
		case "logdir":
			conf.LogDir = kv[1]
		case "logfile":
			conf.LogFile = kv[1]
		case "graphhorizontalscale":
			iv, err := strconv.Atoi(kv[1])
			if err != nil {
				return err
			}
			conf.GraphHorizontalScale = iv
		case "helpvisible":
			bv, err := strconv.ParseBool(kv[1])
			if err != nil {
				return fmt.Errorf("line %d: %v", lineNo, err)
			}
			conf.HelpVisible = bv
		case "colorscheme":
			cs, err := colorschemes.FromName(conf.ConfigDir, kv[1])
			if err != nil {
				return fmt.Errorf("line %d: %v", lineNo, err)
			}
			conf.Colorscheme = cs
		case "updateinterval":
			iv, err := strconv.Atoi(kv[1])
			if err != nil {
				return err
			}
			conf.UpdateInterval = time.Duration(iv)
		case "averagecpu":
			bv, err := strconv.ParseBool(kv[1])
			if err != nil {
				return fmt.Errorf("line %d: %v", lineNo, err)
			}
			conf.AverageLoad = bv
		case "percpuload":
			bv, err := strconv.ParseBool(kv[1])
			if err != nil {
				return fmt.Errorf("line %d: %v", lineNo, err)
			}
			conf.PercpuLoad = bv
		case "tempscale":
			iv, err := strconv.Atoi(kv[1])
			if err != nil {
				return err
			}
			conf.TempScale = widgets.TempScale(iv)
		case "statusbar":
			bv, err := strconv.ParseBool(kv[1])
			if err != nil {
				return fmt.Errorf("line %d: %v", lineNo, err)
			}
			conf.Statusbar = bv
		case "netinterface":
			conf.NetInterface = kv[1]
		case "layout":
			conf.Layout = kv[1]
		case "maxlogsize":
			iv, err := strconv.Atoi(kv[1])
			if err != nil {
				return err
			}
			conf.MaxLogSize = int64(iv)
		case "export":
			conf.ExportPort = kv[1]
		}
	}

	return nil
}
