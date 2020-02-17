package colorschemes

// This is a neutral version of the Solarized 256-color palette. The exception
// is that the one grey color uses the average of base0 and base00, which are
// already middle of the road.
func init() {
	register("solarized", Colorscheme{
		Fg: -1,
		Bg: -1,

		BorderLabel: -1,
		BorderLine:  37,

		CPULines: []int{61, 33, 37, 64, 125, 160, 166, 136},

		BattLines: []int{61, 33, 37, 64, 125, 160, 166, 136},

		MainMem: 125,
		SwapMem: 166,

		ProcCursor: 136,

		Sparkline: 33,

		DiskBar: 243,

		TempLow:  64,
		TempHigh: 160,
	})
}
