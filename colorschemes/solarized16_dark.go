package colorschemes

// This scheme assumes the terminal already uses Solarized. Only DiskBar is
// different between dark/light.
func init() {
	register("solarized16-dark", Colorscheme{
		Fg: -1,
		Bg: -1,

		BorderLabel: -1,
		BorderLine:  6,

		CPULines: []int{13, 4, 6, 2, 5, 1, 9, 3},

		BattLines: []int{13, 4, 6, 2, 5, 1, 9, 3},

		MainMem: 5,
		SwapMem: 9,

		ProcCursor: 4,

		Sparkline: 4,

		DiskBar: 12, // base0

		TempLow:  2,
		TempHigh: 1,
	})
}
