package utils

import (
	rw "github.com/mattn/go-runewidth"
)

func TruncateFront(s string, w int, prefix string) string {
	if rw.StringWidth(s) <= w {
		return s
	}
	r := []rune(s)
	pw := rw.StringWidth(prefix)
	w -= pw
	width := 0
	i := len(r) - 1
	for ; i >= 0; i-- {
		cw := rw.RuneWidth(r[i])
		width += cw
		if width > w {
			break
		}
	}
	return prefix + string(r[i+1:len(r)])
}
