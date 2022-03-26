package colorschemes

func init() {
	register("vice", Colorscheme{
		Fg: 231,
		Bg: -1,

		BorderLabel: 123,
		BorderLine:  102,

		CPULines: []int{212, 218, 123, 159, 229, 158, 183, 146},

		BattLines: []int{212, 218, 123, 159, 229, 158, 183, 146},

		MemLines: []int{201, 97, 212, 218, 123, 159, 229, 158, 183, 146},

		ProcCursor: 159,

		Sparklines: [2]int{183, 146},

		DiskBar: 158,

		TempLow:  49,
		TempHigh: 197,
	})
}
