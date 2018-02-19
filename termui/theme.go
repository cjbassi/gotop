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

// 0 <= r,g,b <= 5
func ColorRGB(r, g, b int) Attribute {
	within := func(n int) int {
		if n < 0 {
			return 0
		}
		if n > 5 {
			return 5
		}
		return n
	}

	r, b, g = within(r), within(b), within(g)
	return Attribute(0x0f + 36*r + 6*g + b)
}
