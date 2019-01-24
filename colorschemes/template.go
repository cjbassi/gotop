package colorschemes

/*
	The standard 256 terminal colors are supported.

	-1 = clear

	You can combine a color with 'Bold', 'Underline', or 'Reverse' by using bitwise OR ('|') and the name of the Color.
	For example, to get Bold red Labels, you would do 'Labels: 2 | Bold'.

	Once you've created a colorscheme, add an entry for it in the `handleColorscheme` function in 'main.go'.
*/

const (
	Bold int = 1 << (iota + 9)
	Underline
	Reverse
)

type Colorscheme struct {
	Name   string
	Author string

	Fg int
	Bg int

	BorderLabel int
	BorderLine  int

	// should add at least 8 here
	CPULines []int

	BattLines []int

	MainMem int
	SwapMem int

	ProcCursor int

	Sparkline int

	DiskBar int

	// colors the temperature number a different color if it's over a certain threshold
	TempLow  int
	TempHigh int
}
