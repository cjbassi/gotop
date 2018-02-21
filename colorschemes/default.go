package colorschemes

var DefaultCS = Colorscheme{
	Name:   "Default",
	Author: "Caleb Bassi",

	Bg: -1,

	Border{
		Labels: 0,
		Line:   0,
	},

	CPU{
		Lines: []int{0, 0, 0, 0},
	},

	Mem{
		Main: 0,
		Swap: 0,
	},

	Proc{
		Cursor: 5,
	},

	Sparkline{
		Graph: 10,
	},
}
