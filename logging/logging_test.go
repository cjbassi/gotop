package logging

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/shibukawa/configdir"
	"github.com/stretchr/testify/assert"
	gotop "github.com/xxxserxxx/gotop/v4"
)

func TestLogging(t *testing.T) {
	c := gotop.NewConfig()
	c.ConfigDir = configdir.New("", "gotoptest")
	c.MaxLogSize = 300
	path := c.ConfigDir.QueryCacheFolder().Path
	var err error
	defer os.RemoveAll(path)
	wc, err := New(c)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	defer wc.Close()
	ds := make([]byte, 100)
	for i := range ds {
		ds[i] = 'x'
	}

	// Base case -- empty log file
	td, err := ioutil.ReadDir(path)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, 1, len(td))

	for i := 1; i < 6; i++ {
		wc.Write(ds)
		wc.Write(ds)
		wc.Write(ds)
		wc.Write([]byte{'\n'}) // max... + 1
		td, err = ioutil.ReadDir(path)
		assert.NoError(t, err)
		k := i
		if k > 4 {
			k = 4
		}
		assert.Equal(t, k, len(td))
	}
}
