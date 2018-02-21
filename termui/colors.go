package termui

type Color int

const ColorDefault = -1

const (
	AttrBold Color = 1 << (iota + 9)
	AttrUnderline
	AttrReverse
)

var Theme = DefaultTheme

var DefaultTheme = Colorscheme{
	Fg: 7,
	Bg: -1,

	LabelFg:  7,
	LabelBg:  -1,
	BorderFg: 6,
	BorderBg: -1,

	SparkLine:   4,
	LineGraph:   -1,
	TableCursor: 4,
}

// A ColorScheme represents the current look-and-feel of the dashboard.
type Colorscheme struct {
	Fg Color
	Bg Color

	LabelFg  Color
	LabelBg  Color
	BorderFg Color
	BorderBg Color

	SparkLine   Color
	LineGraph   Color
	TableCursor Color
}
