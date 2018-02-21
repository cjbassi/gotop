package termui

var Theme = DefaultTheme

var DefaultTheme = ColorScheme{
	Fg: ColorWhite,
	Bg: ColorDefault,

	LabelFg:  ColorWhite,
	LabelBg:  ColorDefault,
	BorderFg: ColorCyan,
	BorderBg: ColorDefault,

	SparkLine:   ColorBlue,
	LineGraph:   ColorDefault,
	TableCursor: ColorBlue,
}

// A ColorScheme represents the current look-and-feel of the dashboard.
type ColorScheme struct {
	Fg Attribute
	Bg Attribute

	LabelFg  Attribute
	LabelBg  Attribute
	BorderFg Attribute
	BorderBg Attribute

	SparkLine   Attribute
	LineGraph   Attribute
	TableCursor Attribute
}
