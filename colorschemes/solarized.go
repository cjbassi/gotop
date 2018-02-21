package colorschemes

var SolarizedCS = Colorscheme{
	Name:   "Default",
	Author: "Caleb Bassi",

	Fg: 250,
	Bg: -1,

	BorderLabel: 250,
	BorderLine:  37,

	CPULines: []int{64, 37, 33, 61, 125, 160, 166, 136},

	MainMem: 125,
	SwapMem: 166,

	ProcCursor: 136,

	Sparkline: 33,

	DiskBar: 245,

	TempLow:  64,
	TempHigh: 160,
}
