package termui

import (
	tb "github.com/nsf/termbox-go"
)

// Init initializes termui library. This function should be called before any others.
// After initialization, the library must be finalized by 'Close' function.
func Init() error {
	if err := tb.Init(); err != nil {
		return err
	}
	tb.SetInputMode(tb.InputEsc | tb.InputMouse)
	tb.SetOutputMode(tb.Output256)

	Body = NewGrid()
	Body.Width, Body.Height = tb.Size()

	return nil
}

// Close finalizes termui library.
// It should be called after successful initialization when termui's functionality isn't required anymore.
func Close() {
	tb.Close()
}
