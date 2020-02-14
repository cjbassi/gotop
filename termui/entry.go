package termui

import (
	"image"
	"strings"
	"unicode/utf8"

	. "github.com/gizak/termui/v3"
	rw "github.com/mattn/go-runewidth"
	"github.com/xxxserxxx/gotop/utils"
)

const (
	ELLIPSIS = "…"
	CURSOR   = " "
)

type Entry struct {
	Block

	Style Style

	Label          string
	Value          string
	ShowWhenEmpty  bool
	UpdateCallback func(string)

	editing bool
}

func (self *Entry) SetEditing(editing bool) {
	self.editing = editing
}

func (self *Entry) update() {
	if self.UpdateCallback != nil {
		self.UpdateCallback(self.Value)
	}
}

// HandleEvent handles input events if the entry is being edited.
// Returns true if the event was handled.
func (self *Entry) HandleEvent(e Event) bool {
	if !self.editing {
		return false
	}
	if utf8.RuneCountInString(e.ID) == 1 {
		self.Value += e.ID
		self.update()
		return true
	}
	switch e.ID {
	case "<C-c>", "<Escape>":
		self.Value = ""
		self.editing = false
		self.update()
	case "<Enter>":
		self.editing = false
	case "<Backspace>":
		if self.Value != "" {
			r := []rune(self.Value)
			self.Value = string(r[:len(r)-1])
			self.update()
		}
	case "<Space>":
		self.Value += " "
		self.update()
	default:
		return false
	}
	return true
}

func (self *Entry) Draw(buf *Buffer) {
	if self.Value == "" && !self.editing && !self.ShowWhenEmpty {
		return
	}

	style := self.Style
	label := self.Label
	if self.editing {
		label += "["
		style = NewStyle(style.Fg, style.Bg, ModifierBold)
	}
	cursorStyle := NewStyle(style.Bg, style.Fg, ModifierClear)

	p := image.Pt(self.Min.X, self.Min.Y)
	buf.SetString(label, style, p)
	p.X += rw.StringWidth(label)

	tail := " "
	if self.editing {
		tail = "] "
	}

	maxLen := self.Max.X - p.X - rw.StringWidth(tail)
	if self.editing {
		maxLen -= 1 // for cursor
	}
	value := utils.TruncateFront(self.Value, maxLen, ELLIPSIS)
	buf.SetString(value, self.Style, p)
	p.X += rw.StringWidth(value)

	if self.editing {
		buf.SetString(CURSOR, cursorStyle, p)
		p.X += rw.StringWidth(CURSOR)
		if remaining := maxLen - rw.StringWidth(value); remaining > 0 {
			buf.SetString(strings.Repeat(" ", remaining), self.TitleStyle, p)
			p.X += remaining
		}
	}
	buf.SetString(tail, style, p)
}
