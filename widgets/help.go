package widgets

import (
	"image"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/xxxserxxx/lingo"
)

var tr lingo.Translations
var keyBinds string

type HelpMenu struct {
	ui.Block
}

func NewHelpMenu(tra lingo.Translations) *HelpMenu {
	tr = tra
	keyBinds = tr.Value("widgets.help")
	return &HelpMenu{
		Block: *ui.NewBlock(),
	}
}

func (help *HelpMenu) Resize(termWidth, termHeight int) {
	textWidth := 53
	for _, line := range strings.Split(keyBinds, "\n") {
		if textWidth < len(line) {
			textWidth = len(line) + 2
		}
	}
	textHeight := strings.Count(keyBinds, "\n") + 1
	x := (termWidth - textWidth) / 2
	y := (termHeight - textHeight) / 2

	help.Block.SetRect(x, y, textWidth+x, textHeight+y)
}

func (help *HelpMenu) Draw(buf *ui.Buffer) {
	help.Block.Draw(buf)

	for y, line := range strings.Split(keyBinds, "\n") {
		for x, rune := range line {
			buf.SetCell(
				ui.NewCell(rune, ui.Theme.Default),
				image.Pt(help.Inner.Min.X+x, help.Inner.Min.Y+y-1),
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
