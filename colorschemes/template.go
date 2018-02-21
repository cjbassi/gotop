package colorschemes

/*
	the standard 256 terminal colors are supported

	-1 = clear

	You can combine a color with Bold, Underline, or Reverse by using bitwise OR ('|').
	For example, to get Bold red Labels, you would do 'Labels: 2 | Bold'
*/

// Ignore this
const (
	Bold int = 1 << (iota + 9)
	Underline
	Reverse
)

type Colorscheme struct {
	Name   string
	Author string

	Bg int

	Border struct {
		Labels int
		Line   int
	}

	CPU struct {
		Lines []int
	}

	Mem struct {
		Main int
		Swap int
	}

	Proc struct {
		Cursor int
	}

	Sparkline struct {
		Graph int
	}
}
