package termui

// Color is an integer in the range -1 to 255.
// -1 is clear, while 0-255 are xterm 256 colors.
type Color int

// ColorDefault = clear
const ColorDefault = -1

// Copied from termbox. Attributes that can be bitwise OR'ed with a color.
const (
	AttrBold Color = 1 << (iota + 9)
	AttrUnderline
	AttrReverse
)

// A Colorscheme represents the current look-and-feel of the dashboard.
type Colorscheme struct {
	Fg Color
	Bg Color

	LabelFg  Color
	LabelBg  Color
	BorderFg Color
	BorderBg Color

	Sparkline   Color
	LineGraph   Color
	TableCursor Color
	GaugeColor  Color
}

var Theme = Colorscheme{
	Fg: 7,
	Bg: -1,

	LabelFg:  7,
	LabelBg:  -1,
	BorderFg: 6,
	BorderBg: -1,

	Sparkline:   4,
	LineGraph:   0,
	TableCursor: 4,
	GaugeColor:  7,
}
