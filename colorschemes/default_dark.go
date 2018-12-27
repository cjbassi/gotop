package colorschemes

var DefaultDark = Colorscheme{
	Fg: 235,
	Bg: -1,

	BorderLabel: 235,
	BorderLine:  6,

	CPULines: []int{4, 3, 2, 1, 5, 6, 7, 8},

	BattLines: []int{4, 3, 2, 1, 5, 6, 7, 8},

	MainMem: 5,
	SwapMem: 3,

	ProcCursor: 33,

	Sparkline: 4,

	DiskBar: 252,

	TempLow:  2,
	TempHigh: 1,
}
