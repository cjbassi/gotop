package config

import (
	"time"

	"github.com/cjbassi/gotop/colorschemes"
	"github.com/cjbassi/gotop/src/widgets"
)

type Config struct {
	ConfigDir string
	LogDir    string
	LogPath   string

	GraphHorizontalScale int
	HelpVisible          bool
	Colorscheme          colorschemes.Colorscheme

	UpdateInterval time.Duration
	MinimalMode    bool
	AverageLoad    bool
	PercpuLoad     bool
	TempScale      widgets.TempScale
	Battery        bool
	Statusbar      bool
	NetInterface   string
}
