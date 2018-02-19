package termui

import (
	"sync"

	tb "github.com/nsf/termbox-go"
)

var renderJobs chan []Bufferer

// So that only one render function can flush/write to the screen at a time
// var renderLock sync.Mutex

// Bufferer should be implemented by all renderable components. Bufferers can render a Buffer.
type Bufferer interface {
	Buffer() *Buffer
	GetXOffset() int
	GetYOffset() int
}

// Render renders all Bufferer in the given order from left to right, right could overlap on left ones.
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
					tb.SetCell(p.X+b.GetXOffset(), p.Y+b.GetYOffset(), c.Ch, tb.Attribute(c.Fg), tb.Attribute(c.Bg))
				}
			}
		}(b)
	}

	// renderLock.Lock()

	wg.Wait()
	tb.Flush()
	// renderLock.Unlock()
}

func Clear() {
	tb.Clear(tb.ColorDefault, tb.Attribute(Theme.Bg))
}
