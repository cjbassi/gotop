package widgets

import termui "github.com/gizak/termui/v3"

type Scalable interface {
	termui.Drawable
	Scale(i int)
}
