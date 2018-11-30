package widgets

import (
	"strings"

	ui "github.com/cjbassi/termui"
)

const KEYBINDS = `
Quit: q or <C-c>

Process navigation
  - k and <Up>: up
  - j and <Down>: down
  - <C-u>: half page up
  - <C-d>: half page down
  - <C-b>: full page up
  - <C-f>: full page down
  - gg and <Home>: jump to top
  - G and <End>: jump to bottom

Process actions:
  - <Tab>: toggle process grouping
  - dd: kill selected process or group of processes

Process sorting
  - c: CPU
  - m: Mem
  - p: PID

CPU and Mem graph scaling:
  - h: scale in
  - l: scale out
`

type HelpMenu struct {
	*ui.Block
}

func NewHelpMenu() *HelpMenu {
	block := ui.NewBlock()
	block.X = 51 // width - 1
	block.Y = 24 // height - 1
	return &HelpMenu{block}
}

func (self *HelpMenu) Buffer() *ui.Buffer {
	buf := self.Block.Buffer()

	self.Block.XOffset = (ui.Body.Width - self.Block.X) / 2  // X coordinate
	self.Block.YOffset = (ui.Body.Height - self.Block.Y) / 2 // Y coordinate

	for y, line := range strings.Split(KEYBINDS, "\n") {
		for x, char := range line {
			buf.SetCell(x+1, y, ui.NewCell(char, ui.Color(7), self.Bg))
		}
	}

	return buf
}
