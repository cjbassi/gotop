package widgets

import (
	"image"
	"strings"

	ui "github.com/gizak/termui"
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
	ui.Block
}

func NewHelpMenu() *HelpMenu {
	return &HelpMenu{
		Block: *ui.NewBlock(),
	}
}

func (self *HelpMenu) Resize(termWidth, termHeight int) {
	textWidth := 53
	textHeight := 22
	x := (termWidth - textWidth) / 2
	y := (termHeight - textHeight) / 2

	self.Block.SetRect(x, y, textWidth+x, textHeight+y)
}

func (self *HelpMenu) Draw(buf *ui.Buffer) {
	self.Block.Draw(buf)

	for y, line := range strings.Split(KEYBINDS, "\n") {
		for x, rune := range line {
			buf.SetCell(
				ui.NewCell(rune, ui.NewStyle(7)),
				image.Pt(self.Inner.Min.X+x, self.Inner.Min.Y+y-1),
			)
		}
	}
}
