package gotop

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xxxserxxx/gotop/v4/widgets"
)

// FIXME This is totally broken since the updates
func TestParse(t *testing.T) {
	tests := []struct {
		i string
		f func(c Config, e error)
	}{
		{
			i: "graphhorizontalscale",
			f: func(c Config, e error) {
				assert.Error(t, e, "invalid graphhorizontalscale syntax")
			},
		},
		{
			i: "helpvisible=true=false",
			f: func(c Config, e error) {
				assert.NotNil(t, e)
			},
		},
		{
			i: "GRAPHHORIZONTALSCALE=1\nhelpVisible=true",
			f: func(c Config, e error) {
				assert.Nil(t, e, "unexpected error")
				assert.Equal(t, 1, c.GraphHorizontalScale)
			},
		},
		{
			i: "graphhorizontalscale=a",
			f: func(c Config, e error) {
				assert.Error(t, e, "expected invalid value for graphhorizontalscale")
			},
		},
		{
			i: "helpvisible=a",
			f: func(c Config, e error) {
				assert.Error(t, e, "expected invalid value for helpvisible")
			},
		},
		{
			i: "helpvisible=true\nupdateinterval=30\naveragecpu=true\nPerCPULoad=true\ntempscale=F\nstatusbar=true\nnetinterface=eth0\nlayout=minimal\nmaxlogsize=200",
			f: func(c Config, e error) {
				assert.Nil(t, e, "unexpected error")
				assert.Equal(t, true, c.HelpVisible)
				assert.Equal(t, time.Duration(30), c.UpdateInterval)
				assert.Equal(t, true, c.AverageLoad)
				assert.Equal(t, true, c.PercpuLoad)
				assert.Equal(t, widgets.TempScale(70), c.TempScale)
				assert.Equal(t, true, c.Statusbar)
				assert.Equal(t, "eth0", c.NetInterface)
				assert.Equal(t, "minimal", c.Layout)
				assert.Equal(t, int64(200), c.MaxLogSize)
			},
		},
	}
	for _, tc := range tests {
		in := strings.NewReader(tc.i)
		c := NewConfig()
		e := load(in, &c)
		tc.f(c, e)
	}
}
