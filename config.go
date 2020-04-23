package gotop

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/shibukawa/configdir"
	"github.com/xxxserxxx/gotop/v3/colorschemes"
	"github.com/xxxserxxx/gotop/v3/widgets"
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
	r := bufio.NewScanner(bytes.NewReader(in))
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
			iv, err := strconv.Atoi(kv[1])
			if err != nil {
				return err
			}
			conf.TempScale = widgets.TempScale(iv)
		case statusbar:
			bv, err := strconv.ParseBool(kv[1])
			if err != nil {
				return fmt.Errorf("line %d: %v", lineNo, err)
			}
			conf.Statusbar = bv
		case netinterface:
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
	fmt.Fprintf(buff, "%s=%d\n", graphhorizontalscale, c.GraphHorizontalScale)
	fmt.Fprintf(buff, "%s=%t\n", helpvisible, c.HelpVisible)
	fmt.Fprintf(buff, "%s=%s\n", colorscheme, c.Colorscheme.Name)
	fmt.Fprintf(buff, "%s=%d\n", updateinterval, c.UpdateInterval)
	fmt.Fprintf(buff, "%s=%t\n", averagecpu, c.AverageLoad)
	fmt.Fprintf(buff, "%s=%t\n", percpuload, c.PercpuLoad)
	fmt.Fprintf(buff, "%s=%d\n", tempscale, c.TempScale)
	fmt.Fprintf(buff, "%s=%t\n", statusbar, c.Statusbar)
	fmt.Fprintf(buff, "%s=%s\n", netinterface, c.NetInterface)
	fmt.Fprintf(buff, "%s=%s\n", layout, c.Layout)
	fmt.Fprintf(buff, "%s=%d\n", maxlogsize, c.MaxLogSize)
	fmt.Fprintf(buff, "%s=%s\n", export, c.ExportPort)
	fmt.Fprintf(buff, "%s=%t\n", mbps, c.Mbps)
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
