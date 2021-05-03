package widgets

import (
	"strings"

	"github.com/gizak/termui/v3/widgets"
	"github.com/xxxserxxx/lingo/v2"
)

// Used by all widgets
var tr lingo.Translations

type HelpMenu struct {
	widgets.Paragraph
}

func NewHelpMenu(tra lingo.Translations) *HelpMenu {
	tr = tra
	help := &HelpMenu{
		Paragraph: *widgets.NewParagraph(),
	}
	help.Paragraph.Text = tra.Value("help.help")
	return help
}

func (help *HelpMenu) Resize(termWidth, termHeight int) {
	textWidth := 53
	var nlines int
	var line string
	for nlines, line = range strings.Split(help.Text, "\n") {
		if textWidth < len(line) {
			textWidth = len(line) + 2
		}
	}
	textHeight := nlines + 2
	x := (termWidth - textWidth) / 2
	y := (termHeight - textHeight) / 2

	help.Paragraph.SetRect(x, y, textWidth+x, textHeight+y)
}
