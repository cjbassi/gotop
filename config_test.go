package gotop

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xxxserxxx/gotop/v3/widgets"
)

func TestParse(t *testing.T) {
	tests := []struct {
		i string
		f func(c Config, e error)
	}{
		{
			i: "logdir",
			f: func(c Config, e error) {
				assert.Error(t, e, "invalid config syntax")
			},
		},
		{
			i: "logdir=foo=bar",
			f: func(c Config, e error) {
				assert.NotNil(t, e)
			},
		},
		{
			i: "foo=bar",
			f: func(c Config, e error) {
				assert.NotNil(t, e)
			},
		},
		{
			i: "configdir=abc\nlogdir=bar\nlogfile=errors",
			f: func(c Config, e error) {
				assert.Nil(t, e, "unexpected error")
				assert.Equal(t, "abc", c.ConfigDir)
				assert.Equal(t, "bar", c.LogDir)
				assert.Equal(t, "errors", c.LogFile)
			},
		},
		{
			i: "CONFIGDIR=abc\nloGdir=bar\nlogFILe=errors",
			f: func(c Config, e error) {
				assert.Nil(t, e, "unexpected error")
				assert.Equal(t, "abc", c.ConfigDir)
				assert.Equal(t, "bar", c.LogDir)
				assert.Equal(t, "errors", c.LogFile)
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
			i: "helpvisible=true\nupdateinterval=30\naveragecpu=true\nPerCPULoad=true\ntempscale=100\nstatusbar=true\nnetinterface=eth0\nlayout=minimal\nmaxlogsize=200",
			f: func(c Config, e error) {
				assert.Nil(t, e, "unexpected error")
				assert.Equal(t, true, c.HelpVisible)
				assert.Equal(t, time.Duration(30), c.UpdateInterval)
				assert.Equal(t, true, c.AverageLoad)
				assert.Equal(t, true, c.PercpuLoad)
				assert.Equal(t, widgets.TempScale(100), c.TempScale)
				assert.Equal(t, true, c.Statusbar)
				assert.Equal(t, "eth0", c.NetInterface)
				assert.Equal(t, "minimal", c.Layout)
				assert.Equal(t, int64(200), c.MaxLogSize)
			},
		},
	}
	for _, tc := range tests {
		in := strings.NewReader(tc.i)
		c := Config{}
		e := Parse(in, &c)
		tc.f(c, e)
	}
}
