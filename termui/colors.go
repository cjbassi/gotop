package termui

/* ---------------Port from termbox-go --------------------- */

// Attribute is printable cell's color and style.
type Attribute uint16

const (
	ColorDefault Attribute = iota
	ColorBlack
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
)

const NumberofColors = 8

const (
	AttrBold Attribute = 1 << (iota + 9)
	AttrUnderline
	AttrReverse
)

/* ----------------------- End ----------------------------- */
