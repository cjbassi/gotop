package termui

import (
	"sync"

	tb "github.com/nsf/termbox-go"
)

// Bufferer should be implemented by all renderable components.
type Bufferer interface {
	Buffer() *Buffer
	GetXOffset() int
	GetYOffset() int
}

// Render renders all Bufferers in the given order to termbox, then asks termbox to print the screen.
func Render(bs ...Bufferer) {
	var wg sync.WaitGroup
	for _, b := range bs {
		wg.Add(1)
		go func(b Bufferer) {
			defer wg.Done()
			buf := b.Buffer()
			// set cells in buf
			for p, c := range buf.CellMap {
				if p.In(buf.Area) {
					tb.SetCell(p.X+b.GetXOffset(), p.Y+b.GetYOffset(), c.Ch, tb.Attribute(c.Fg)+1, tb.Attribute(c.Bg)+1)
				}
			}
		}(b)
	}

	wg.Wait()
	tb.Flush()
}

// Clear clears the screen with the default Bg color.
func Clear() {
	tb.Clear(tb.ColorDefault+1, tb.Attribute(Theme.Bg)+1)
}
