package colorschemes

func init() {
	register("default", Colorscheme{
		Fg: 7,
		Bg: -1,

		BorderLabel: 7,
		BorderLine:  6,

		CPULines: []int{4, 3, 2, 1, 5, 6, 7, 8},

		BattLines: []int{4, 3, 2, 1, 5, 6, 7, 8},

		MemLines: []int{5, 11, 4, 3, 2, 1, 6, 7, 8},

		ProcCursor: 4,

		Sparklines: [2]int{4, 5},

		DiskBar: 7,

		TempLow:  2,
		TempHigh: 1,
	})
}
