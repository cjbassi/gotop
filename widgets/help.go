package widgets

import (
	"image"
	"strings"

	ui "github.com/gizak/termui/v3"
)

const KEYBINDS = `
Quit: q or <C-c>

Process navigation:
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
  - dd: kill selected process or group of processes with SIGTERM (15)
  - d3: kill selected process or group of processes with SIGQUIT (3)
  - d9: kill selected process or group of processes with SIGKILL (9)

Process sorting:
  - c: CPU
  - m: Mem
  - p: PID

Process filtering:
  - /: start editing filter
  - (while editing):
    - <Enter>: accept filter
    - <C-c> and <Escape>: clear filter

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
	for _, line := range strings.Split(KEYBINDS, "\n") {
		if textWidth < len(line) {
			textWidth = len(line) + 2
		}
	}
	textHeight := strings.Count(KEYBINDS, "\n") + 1
	x := (termWidth - textWidth) / 2
	y := (termHeight - textHeight) / 2

	self.Block.SetRect(x, y, textWidth+x, textHeight+y)
}

func (self *HelpMenu) Draw(buf *ui.Buffer) {
	self.Block.Draw(buf)

	for y, line := range strings.Split(KEYBINDS, "\n") {
		for x, rune := range line {
			buf.SetCell(
				ui.NewCell(rune, ui.Theme.Default),
				image.Pt(self.Inner.Min.X+x, self.Inner.Min.Y+y-1),
			)
		}
	}
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
