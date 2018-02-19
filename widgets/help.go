package widgets

import (
	"strings"

	ui "github.com/cjbassi/gotop/termui"
)

const KEYBINDS = `
Quit: 'q' or 'Ctrl-c'

Navigation:
  - '<up>'/'<down>' and 'j'/'k': up and down
  - 'C-d' and 'C-u': up and down half a page
  - 'C-f' and 'C-b': up and down a full page
  - 'gg' and 'G': jump to top and bottom

Process Sorting:
  - 'c': CPU
  - 'm': Mem
  - 'p': PID

'<tab>': toggle process grouping
'dd': kill the selected process or process group
'<left>'/'<right>' and 'h'/'l': ...
'u': update gotop
`

type HelpMenu struct {
	ui.Block
}

func NewHelpMenu() *HelpMenu {
	block := *ui.NewBlock()
	block.X = 48
	block.Y = 17
	return &HelpMenu{block}
}

func (hm *HelpMenu) Buffer() *ui.Buffer {
	buf := hm.Block.Buffer()

	for y, line := range strings.Split(KEYBINDS, "\n") {
		for x, char := range line {
			buf.SetCell(x+1, y, ui.NewCell(char, ui.ColorWhite, ui.ColorDefault))
		}
	}

	buf.SetAreaXY(100, 100)

	return buf
}
