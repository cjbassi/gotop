package termui

// Attribute is printable cell's color and style.
type Attribute int

// Define basic terminal colors
const (
	// ColorDefault clears the color
	ColorDefault Attribute = iota - 1
	ColorBlack
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
)

// These can be bitwise ored to modify cells
const (
	AttrBold Attribute = 1 << (iota + 9)
	AttrUnderline
	AttrReverse
)

// AttrPair holds a cell's Fg and Bg
type AttrPair struct {
	Fg Attribute
	Bg Attribute
}
