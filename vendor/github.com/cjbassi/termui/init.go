package termui

import (
	"github.com/gdamore/tcell"
)

var screen tcell.Screen

// Init initializes termui library. This function should be called before any others.
// After initialization, the library must be finalized by 'Close' function.
func Init() error {
	var err error

	screen, err = tcell.NewScreen()
	if err != nil {
		return err
	}
	if err = screen.Init(); err != nil {
		return err
	}
	screen.EnableMouse()

	Body = NewGrid()
	Body.Width, Body.Height = screen.Size()

	On("<resize>", func(e Event) {
		screen.Sync()
		Body.Width, Body.Height = e.Width, e.Height
		Body.Resize()
	})

	return nil
}

// Close finalizes termui library.
// It should be called after successful initialization when termui's functionality isn't required anymore.
func Close() {
	screen.Fini()
}
