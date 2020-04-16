package gotop

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/shibukawa/configdir"
	"github.com/xxxserxxx/gotop/v3/colorschemes"
	"github.com/xxxserxxx/gotop/v3/widgets"
)

// TODO: test, build, release [#119] [#120] [#121]
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
	Extensions     []string
	Mbps           bool

	Test bool
}

func (conf *Config) Load() error {
	var in io.Reader
	cfn := "gotop.conf"
	folder := conf.ConfigDir.QueryFolderContainsFile(cfn)
	if folder != nil {
		// FIXME: Shouldn't this be looking in folder??
		if cf, err := os.Open(cfn); err == nil {
			defer cf.Close()
		} else {
			return err
		}
	} else {
		return nil
	}
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
			log.Printf("configdir is deprecated.  Ignored configdir=%s", kv[1])
		case "logdir":
			log.Printf("logdir is deprecated.  Ignored logdir=%s", kv[1])
		case "logfile":
			log.Printf("logfile is deprecated.  Ignored logfile=%s", kv[1])
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
		case "extensions":
			conf.Extensions = strings.Split(kv[1], ",")
		case "mbps":
			conf.Mbps = true
		}
	}

	return nil
}
