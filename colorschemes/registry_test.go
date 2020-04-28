package colorschemes

import (
	"github.com/shibukawa/configdir"
	"path/filepath"
	"reflect"
	"testing"
)

func TestColorRegistry(t *testing.T) {
	colors := []string{"default", "default-dark", "solarized", "solarized16-dark", "solarized16-light", "monokai", "vice"}
	zeroCS := Colorscheme{}
	cd := configdir.New("", "gotop")
	cd.LocalPath, _ = filepath.Abs(".")
	for _, cn := range colors {
		c, e := FromName(cd, cn)
		if e != nil {
			t.Errorf("unexpected error fetching built-in color %s: %s", cn, e)
		}
		if reflect.DeepEqual(c, zeroCS) {
			t.Error("expected a colorscheme, but got back a zero value.")
		}
	}
}
