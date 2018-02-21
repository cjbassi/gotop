package colorschemes

/*
	the standard 256 terminal colors are supported

	-1 = clear

	You can combine a color with Bold, Underline, or Reverse by using bitwise OR ('|').
	For example, to get Bold red Labels, you would do 'Labels: 2 | Bold'
*/

// Ignore this
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

	CPULines []int

	MainMem int
	SwapMem int

	ProcCursor int

	Sparkline int

	DiskBar int

	// Temperature colors depending on if it's over a certain threshold
	TempLow  int
	TempHigh int
}
