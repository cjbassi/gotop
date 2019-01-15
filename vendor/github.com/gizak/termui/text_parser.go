// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import (
	"strings"

	wordwrap "github.com/mitchellh/go-wordwrap"
)

type textParser struct {
	baseAttrs AttrPair
	plainTx   []rune
	markers   []marker
}

type marker struct {
	st int
	ed int
	fg Attribute
	bg Attribute
}

var colorMap = map[string]Attribute{
	"red":     ColorRed,
	"blue":    ColorBlue,
	"black":   ColorBlack,
	"cyan":    ColorCyan,
	"yellow":  ColorYellow,
	"white":   ColorWhite,
	"default": ColorDefault,
	"green":   ColorGreen,
	"magenta": ColorMagenta,
}

var attributeMap = map[string]Attribute{
	"bold":      AttrBold,
	"underline": AttrUnderline,
	"reverse":   AttrReverse,
}

// AddColorMap allows users to add/override the string to attribute mapping
func AddColorMap(str string, attr Attribute) {
	colorMap[str] = attr
}

// readAttributes translates strings like `fg-red,fg-bold,bg-white` to fg and bg Attribute
func (tp *textParser) readAttributes(s string) (Attribute, Attribute) {
	fg := tp.baseAttrs.Fg
	bg := tp.baseAttrs.Bg

	updateAttr := func(a Attribute, attrs []string) Attribute {
		for _, s := range attrs {
			// replace the color
			if c, ok := colorMap[s]; ok {
				a &= 0xFF00 // erase clr 0 ~ 8 bits
				a |= c      // set clr
			}
			// add attrs
			if c, ok := attributeMap[s]; ok {
				a |= c
			}
		}
		return a
	}

	ss := strings.Split(s, ",")
	fgs := []string{}
	bgs := []string{}
	for _, v := range ss {
		subs := strings.Split(v, "-")
		if len(subs) > 1 {
			if subs[0] == "fg" {
				fgs = append(fgs, subs[1])
			} else if subs[0] == "bg" {
				bgs = append(bgs, subs[1])
			}
			// else maybe error somehow?
		}
	}

	fg = updateAttr(fg, fgs)
	bg = updateAttr(bg, bgs)
	return fg, bg
}

// parse streams and parses text into normalized text and render sequence.
func (tp *textParser) parse(str string) {
	rs := []rune(str)
	normTx := []rune{}
	square := []rune{}
	brackt := []rune{}
	accSquare := false
	accBrackt := false
	cntSquare := 0

	reset := func() {
		square = []rune{}
		brackt = []rune{}
		accSquare = false
		accBrackt = false
		cntSquare = 0
	}
	// pipe stacks into normTx and clear
	rollback := func() {
		normTx = append(normTx, square...)
		normTx = append(normTx, brackt...)
		reset()
	}
	// chop first and last
	chop := func(s []rune) []rune {
		return s[1 : len(s)-1]
	}

	for i, r := range rs {
		switch {
		// stacking brackt
		case accBrackt:
			brackt = append(brackt, r)
			if ')' == r {
				fg, bg := tp.readAttributes(string(chop(brackt)))
				st := len(normTx)
				ed := len(normTx) + len(square) - 2
				tp.markers = append(tp.markers, marker{st, ed, fg, bg})
				normTx = append(normTx, chop(square)...)
				reset()
			} else if i+1 == len(rs) {
				rollback()
			}
		// stacking square
		case accSquare:
			switch {
			// squares closed and followed by a '('
			case cntSquare == 0 && '(' == r:
				accBrackt = true
				brackt = append(brackt, '(')
			// squares closed but not followed by a '('
			case cntSquare == 0:
				rollback()
				if '[' == r {
					accSquare = true
					cntSquare = 1
					brackt = append(brackt, '[')
				} else {
					normTx = append(normTx, r)
				}
			// hit the end
			case i+1 == len(rs):
				square = append(square, r)
				rollback()
			case '[' == r:
				cntSquare++
				square = append(square, '[')
			case ']' == r:
				cntSquare--
				square = append(square, ']')
			// normal char
			default:
				square = append(square, r)
			}
		// stacking normTx
		default:
			if '[' == r {
				accSquare = true
				cntSquare = 1
				square = append(square, '[')
			} else {
				normTx = append(normTx, r)
			}
		}
	}

	tp.plainTx = normTx
}

func WrapText(cs []Cell, wl int) []Cell {
	tmpCell := make([]Cell, len(cs))
	copy(tmpCell, cs)

	// get the plaintext
	plain := CellsToString(cs)

	// wrap
	plainWrapped := wordwrap.WrapString(plain, uint(wl))

	// find differences and insert
	finalCell := tmpCell // finalcell will get the inserts and is what is returned

	plainRune := []rune(plain)
	plainWrappedRune := []rune(plainWrapped)
	trigger := "go"
	plainRuneNew := plainRune

	for trigger != "stop" {
		plainRune = plainRuneNew
		for i := range plainRune {
			if plainRune[i] == plainWrappedRune[i] {
				trigger = "stop"
			} else if plainRune[i] != plainWrappedRune[i] && plainWrappedRune[i] == 10 {
				trigger = "go"
				cell := Cell{10, AttrPair{0, 0}}
				j := i - 0

				// insert a cell into the []Cell in correct position
				tmpCell[i] = cell

				// insert the newline into plain so we avoid indexing errors
				plainRuneNew = append(plainRune, 10)
				copy(plainRuneNew[j+1:], plainRuneNew[j:])
				plainRuneNew[j] = plainWrappedRune[j]

				// restart the inner for loop until plain and plain wrapped are
				// the same; yeah, it's inefficient, but the text amounts
				// should be small
				break

			} else if plainRune[i] != plainWrappedRune[i] &&
				plainWrappedRune[i-1] == 10 && // if the prior rune is a newline
				plainRune[i] == 32 { // and this rune is a space
				trigger = "go"
				// need to delete plainRune[i] because it gets rid of an extra
				// space
				plainRuneNew = append(plainRune[:i], plainRune[i+1:]...)
				break

			} else {
				trigger = "stop" // stops the outer for loop
			}
		}
	}

	finalCell = tmpCell

	return finalCell
}

func ParseText(s string, baseAttrs AttrPair) []Cell {
	tp := textParser{
		baseAttrs: baseAttrs,
	}
	tp.parse(s)
	cs := make([]Cell, len(tp.plainTx))
	for i := range cs {
		cs[i] = Cell{tp.plainTx[i], baseAttrs}
	}
	for _, mrk := range tp.markers {
		for i := mrk.st; i < mrk.ed; i++ {
			cs[i].Attrs.Fg = mrk.fg
			cs[i].Attrs.Bg = mrk.bg
		}
	}

	return cs
}
