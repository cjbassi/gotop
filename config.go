package gotop

import (
	"io"
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
	TempScale      widgets.TempScale
	Battery        bool
	Statusbar      bool
	NetInterface   string
	Layout         io.Reader
	MaxLogSize     int64
}
