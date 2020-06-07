package colorschemes

//revive:disable
const (
	Bold int = 1 << (iota + 9)
	Underline
	Reverse
)

//revive:enable

/*
Colorscheme defines colors and fonts used by TUI elements. The standard
256 terminal colors are supported.

For int values, -1 = clear

Colors may be combined with 'Bold', 'Underline', or 'Reverse' by using
bitwise OR ('|') and the name of the Color. For example, to get bold red
labels, you would use 'Labels: 2 | Bold'
*/
type Colorscheme struct {
	// Name is the key used to look up the colorscheme, e.g. as provided by the user
	Name string
	// Who created the color scheme
	Author string

	// Foreground color
	Fg int
	// Background color
	Bg int

	// BorderLabel is the color of the widget title label
	BorderLabel int
	// BorderLine is the color of the widget border
	BorderLine int

	// CPULines define the colors used for the CPU activity graph, in
	// order, for each core. Should add at least 8 here; they're
	// selected in order, with wrapping.
	CPULines []int

	// BattLines define the colors used for the battery history graph.
	// Should add at least 2; they're selected in order, with wrapping.
	BattLines []int

	// MemLines define the colors used for the memory histograph.
	// Should add at least 2 (physical & swap); they're selected in order,
	// with wrapping.
	MemLines []int

	// ProcCursor is used as the color for the color bar in the process widget
	ProcCursor int

	// SparkLine determines the color of the data line in spark graphs
	Sparkline int

	// DiskBar is the color of the disk gauge bars (currently unused,
	// as there's no disk gauge widget)
	DiskBar int

	// TempLow determines the color of the temperature number when it's under
	// a certain threshold
	TempLow int
	// TempHigh determines the color of the temperature number when it's over
	// a certain threshold
	TempHigh int
}
