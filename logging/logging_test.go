package logging

import (
	"io/ioutil"
	//"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xxxserxxx/gotop/v3"
)

func TestLogging(t *testing.T) {
	tdn := "testdir"
	path, err := filepath.Abs(tdn)
	defer os.RemoveAll(path)
	c := gotop.Config{
		MaxLogSize: 300,
		LogDir:     path,
		LogFile:    "errors.log",
	}
	wc, err := New(c)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	defer wc.Close()
	ds := make([]byte, 100)
	for i, _ := range ds {
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
