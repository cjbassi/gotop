/*
	The standard 256 terminal colors are supported.

	-1 = clear

	You can combine a color with 'Bold', 'Underline', or 'Reverse' by using bitwise OR ('|') and the name of the Color.
	For example, to get Bold red Labels, you would do 'Labels: 2 | Bold'.

	Once you've created a colorscheme, add an entry for it in the `handleColorscheme` function in 'main.go'.
*/

package colorschemes

func init() {
	register("nord", Colorscheme{
		Name:   "A Nord Approximation",
		Author: "@jrswab",
		Fg:     254, // lightest
		Bg:     -1,

		BorderLabel: 254,
		BorderLine:  96, // Purple

		CPULines: []int{4, 3, 2, 1, 5, 6, 7, 8},

		BattLines: []int{4, 3, 2, 1, 5, 6, 7, 8},

		MainMem: 172, // Orange
		SwapMem: 221, // yellow

		ProcCursor: 31, // blue (nord9)

		Sparkline: 31,

		DiskBar: 254,

		TempLow:  64,  // green
		TempHigh: 167, // red
	})
}
