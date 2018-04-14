package colorschemes

/*
	The standard 256 terminal colors are supported.

	-1 = clear

	You can combine a color with 'Bold', 'Underline', or 'Reverse' by using bitwise OR ('|') and the name of the attribute.
	For example, to get Bold red Labels, you would do 'Labels: 2 | Bold'.

	Once you've created a colorscheme, add an entry for it in the `handleColorscheme` function
	in `gotop.go`.
*/

type Colorscheme struct {
	Default int32

	BorderLabel int32
	BorderLine  int32

	// should add at least 8 here
	CPULines []int32

	MainMem int32
	SwapMem int32

	ProcCursor int32

	Sparkline int32

	DiskBar int32

	// colors the temperature number a different color if it's over a certain threshold
	TempLow  int32
	TempHigh int32
}

const (
	ColorBlack   = 0x000000
	ColorMaroon  = 0x800000
	ColorGreen   = 0x008000
	ColorOlive   = 0x808000
	ColorNavy    = 0x000080
	ColorPurple  = 0x800080
	ColorTeal    = 0x008080
	ColorSilver  = 0xC0C0C0
	ColorGray    = 0x808080
	ColorRed     = 0xFF0000
	ColorLime    = 0x00FF00
	ColorYellow  = 0xFFFF00
	ColorBlue    = 0x0000FF
	ColorFuchsia = 0xFF00FF
	ColorAqua    = 0x00FFFF
	ColorWhite   = 0xFFFFFF
)
