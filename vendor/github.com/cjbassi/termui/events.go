package termui

import (
	// "strconv"

	"github.com/gdamore/tcell"
)

/*
here's the list of events which you can assign handlers too using the `On` function:
	mouse events:
		<MouseLeft> <MouseRight> <MouseMiddle>
		<MouseWheelUp> <MouseWheelDown>
	keyboard events:
		any uppercase or lowercase letter or a set of two letters like j or jj or J or JJ
		<C-d> etc
		<M-d> etc
		<up> <down> <left> <right>
		<insert> <delete> <home> <end> <previous> <next>
		<backspace> <tab> <enter> <escape> <space>
		<C-<space>> etc
	terminal events:
		<resize>
*/

var eventStream = EventStream{
	make(map[string]func(Event)),
	"",
	make(chan bool, 1),
	make(chan tcell.Event),
}

type EventStream struct {
	eventHandlers map[string]func(Event)
	prevKey       string // previous keypress
	stopLoop      chan bool
	eventQueue    chan tcell.Event // list of events from termbox
}

// Event is a copy of termbox.Event that only contains the fields we need.
type Event struct {
	Key    string
	Width  int
	Height int
	MouseX int
	MouseY int
}

// handleEvent calls the approriate callback function if there is one.
func handleEvent(e tcell.Event) {
	// if e.Type == tcell.EventError {
	// 	panic(e.Err)
	// }

	ne := convertEvent(e)

	if val, ok := eventStream.eventHandlers[ne.Key]; ok {
		val(ne)
	}

	// if val, ok := eventStream.eventHandlers[ne.Key]; ok {
	// 	val(ne)
	// 	eventStream.prevKey = ""
	// } else { // check if the last 2 keys form a key combo with a handler
	// 	// if this is a keyboard event and the previous event was unhandled
	// 	if e.Type == tcell.EventKey && eventStream.prevKey != "" {
	// 		combo := eventStream.prevKey + ne.Key
	// 		if val, ok := eventStream.eventHandlers[combo]; ok {
	// 			ne.Key = combo
	// 			val(ne)
	// 			eventStream.prevKey = ""
	// 		} else {
	// 			eventStream.prevKey = ne.Key
	// 		}
	// 	} else {
	// 		eventStream.prevKey = ne.Key
	// 	}
	// }
}

// Loop gets events from termbox and passes them off to handleEvent.
// Stops when StopLoop is called.
func Loop() {
	go func() {
		for {
			eventStream.eventQueue <- screen.PollEvent()
		}
	}()

	for {
		select {
		case <-eventStream.stopLoop:
			return
		case e := <-eventStream.eventQueue:
			handleEvent(e)
		}
	}
}

// StopLoop stops the event loop.
func StopLoop() {
	eventStream.stopLoop <- true
}

// On assigns event names to their handlers. Takes a string, strings, or a slice of strings, and a function.
func On(things ...interface{}) {
	function := things[len(things)-1].(func(Event))
	for _, thing := range things {
		if value, ok := thing.(string); ok {
			eventStream.eventHandlers[value] = function
		}
		if value, ok := thing.([]string); ok {
			for _, name := range value {
				eventStream.eventHandlers[name] = function
			}
		}
	}
}

// // convertEventKey converts a tcell keyboard event to a more friendly string format.
// // Combines modifiers into the string instead of having them as additional fields in an event.
// func convertEventKey(e tcell.Event) string {
// 	k := string(e.Ch)
// 	pre := ""
// 	mod := ""

// 	if e.Mod == tcell.ModAlt {
// 		mod = "<M-"
// 	}
// 	if e.Ch == 0 {
// 		if e.Key > 0xFFFF-12 {
// 			k = "<f" + strconv.Itoa(0xFFFF-int(e.Key)+1) + ">"
// 		} else if e.Key > 0xFFFF-25 {
// 			ks := []string{"<insert>", "<delete>", "<home>", "<end>", "<previous>", "<next>", "<up>", "<down>", "<left>", "<right>"}
// 			k = ks[0xFFFF-int(e.Key)-12]
// 		}

// 		if e.Key <= 0x7F {
// 			pre = "<C-"
// 			k = string('a' - 1 + int(e.Key))
// 			kmap := map[tcell.Key][2]string{
// 				tcell.KeyCtrlSpace:     {"C-", "<space>"},
// 				tcell.KeyBackspace:     {"", "<backspace>"},
// 				tcell.KeyTab:           {"", "<tab>"},
// 				tcell.KeyEnter:         {"", "<enter>"},
// 				tcell.KeyEsc:           {"", "<escape>"},
// 				tcell.KeyCtrlBackslash: {"C-", "\\"},
// 				tcell.KeyCtrlSlash:     {"C-", "/"},
// 				tcell.KeySpace:         {"", "<space>"},
// 				tcell.KeyCtrl8:         {"C-", "8"},
// 			}
// 			if sk, ok := kmap[e.Key]; ok {
// 				pre = sk[0]
// 				k = sk[1]
// 			}
// 		}
// 	}

// 	if pre != "" {
// 		k += ">"
// 	}

// 	return pre + mod + k
// }

func convertMouseValue(e *tcell.EventMouse) string {
	switch e.Buttons() {
	case tcell.Button1:
		return "<MouseLeft>"
	case tcell.Button2:
		return "<MouseMiddle>"
	case tcell.Button3:
		return "<MouseRight>"
	case tcell.WheelUp:
		return "<MouseWheelUp>"
	case tcell.WheelDown:
		return "<MouseWheelDown>"
	}
	return ""
}

// convertEvent turns a termbox event into a termui event.
func convertEvent(e interface{}) Event {
	var ne Event

	switch e.(type) {
	case *tcell.EventKey:
		ne = Event{
			// Key: convertEventKey(e.Key()),
			Key: "hi",
		}
	case *tcell.EventMouse:
		me := e.(*tcell.EventMouse)
		x, y := me.Position()
		ne = Event{
			Key:    convertMouseValue(me),
			MouseX: x,
			MouseY: y,
		}
	case *tcell.EventResize:
		width, height := screen.Size()
		ne = Event{
			Key:    "<resize>",
			Width:  width,
			Height: height,
		}
	}

	return ne
}
