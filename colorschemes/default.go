package colorschemes

var Default = Colorscheme{
	Default: ColorTeal,

	BorderLabel: ColorTeal,
	// BorderLine:  6,
	BorderLine: ColorAqua,

	// CPULines: []int{4, 3, 2, 1, 5, 6, 7, 8},
	CPULines: []int32{ColorBlue, ColorOlive, ColorGreen, ColorRed},

	// MainMem: 5,
	MainMem: ColorPurple,
	// SwapMem: 11,
	SwapMem: ColorYellow,

	// ProcCursor: 4,
	ProcCursor: ColorBlue,

	Sparkline: ColorBlue,
	// Sparkline: 4,

	// DiskBar: 7,
	DiskBar: ColorTeal,

	// TempLow:  2,
	TempLow: ColorGreen,
	// TempHigh: 1,
	TempHigh: ColorRed,
}
