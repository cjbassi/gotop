package gotop

import (
	"io"
	"time"

	"github.com/xxxserxxx/gotop/colorschemes"
	"github.com/xxxserxxx/gotop/widgets"
)

// TODO: Cross-compiling for darwin, openbsd requiring native procs & temps
type Config struct {
	ConfigDir string
	LogDir    string
	LogPath   string

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
}
