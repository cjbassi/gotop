package colorschemes

func init() {
	register("monokai", Colorscheme{
		Fg: 249,
		Bg: -1,

		BorderLabel: 249,
		BorderLine:  239,

		CPULines: []int{81, 70, 208, 197, 249, 141, 221, 186},

		BattLines: []int{81, 70, 208, 197, 249, 141, 221, 186},

		MemLines: []int{208, 186, 81, 70, 208, 197, 249, 141, 221, 186},

		ProcCursor: 197,

		Sparkline: 81,

		DiskBar: 102,

		TempLow:  70,
		TempHigh: 208,
	})
}
