package colorschemes

// This scheme assumes the terminal already uses Solarized. Only DiskBar is
// different between dark/light.
func init() {
	register("solarized16-light", Colorscheme{
		Fg: -1,
		Bg: -1,

		BorderLabel: -1,
		BorderLine:  6,

		CPULines: []int{13, 4, 6, 2, 5, 1, 9, 3},

		BattLines: []int{13, 4, 6, 2, 5, 1, 9, 3},

		MemLines: []int{5, 9, 13, 4, 6, 2, 1, 3},

		ProcCursor: 4,

		Sparklines: [2]int{4, 5},

		DiskBar: 11, // base00

		TempLow:  2,
		TempHigh: 1,
	})
}
