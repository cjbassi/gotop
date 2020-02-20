package colorschemes

import (
	"reflect"
	"testing"
)

func TestColorRegistry(t *testing.T) {
	colors := []string{"default", "default-dark", "solarized", "solarized16-dark", "solarized16-light", "monokai", "vice"}
	zeroCS := Colorscheme{}
	for _, cn := range colors {
		c, e := FromName("", cn)
		if e != nil {
			t.Errorf("unexpected error fetching built-in color %s: %s", cn, e)
		}
		if reflect.DeepEqual(c, zeroCS) {
			t.Error("expected a colorscheme, but got back a zero value.")
		}
	}
}
