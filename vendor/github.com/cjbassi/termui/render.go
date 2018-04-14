package termui

import (
	"sync"
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
					screen.SetContent(p.X+b.GetXOffset(), p.Y+b.GetYOffset(), c.Ch, []rune{}, c.Style)
				}
			}
		}(b)
	}
	wg.Wait()
	screen.Show()
}

// Clear clears the screen with the default Bg color.
func Clear() {
	screen.Clear()
}
