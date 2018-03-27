package widgets

import (
	"strings"

	ui "github.com/cjbassi/gotop/termui"
)

const KEYBINDS = `
Quit: q or <C-c>

Process Navigation
  - <up>/<down> and j/k: up and down
  - <C-d> and <C-u>: up and down half a page
  - <C-f> and <C-b>: up and down a full page
  - gg and G: jump to top and bottom

Process Sorting
  - c: CPU
  - m: Mem
  - p: PID

<tab>: toggle process grouping
dd: kill the selected process or process group

h and l: zoom in and out of CPU and Mem graphs
`

type HelpMenu struct {
	*ui.Block
}

func NewHelpMenu() *HelpMenu {
	block := ui.NewBlock()
	block.X = 48                                   // width - 1
	block.Y = 17                                   // height - 1
	block.XOffset = (ui.Body.Width - block.X) / 2  // X coordinate
	block.YOffset = (ui.Body.Height - block.Y) / 2 // Y coordinate
	return &HelpMenu{block}
}

func (self *HelpMenu) Buffer() *ui.Buffer {
	buf := self.Block.Buffer()

	for y, line := range strings.Split(KEYBINDS, "\n") {
		for x, char := range line {
			buf.SetCell(x+1, y, ui.NewCell(char, ui.Color(7), self.Bg))
		}
	}

	return buf
}
