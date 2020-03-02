package layout

// NOT MY FAULT.  Some dependency already pulled in testify -- 13kLOC
import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestParsing(t *testing.T) {
	tests := []struct {
		i string
		f func(l layout)
	}{
		{"cpu", func(l layout) {
			assert.Equal(t, 1, len(l.Rows))
			assert.Equal(t, 1, len(l.Rows[0]))
		}},
		{"   cpu   \ndisk/1     mem/3\ntemp   \nnet    procs", func(l layout) {
			assert.Equal(t, 4, len(l.Rows))
			assert.Equal(t, 1, len(l.Rows[0]))
			assert.Equal(t, 2, len(l.Rows[1]))
			assert.Equal(t, 1, len(l.Rows[2]))
			assert.Equal(t, 2, len(l.Rows[3]))
		}},
		{"cpu\ndisk/1 mem/3\ntemp\nnet procs", func(l layout) {
			assert.Equal(t, 4, len(l.Rows))
			// 1
			assert.Equal(t, 1, len(l.Rows[0]))
			assert.Equal(t, 1.0, l.Rows[0][0].Weight)
			assert.Equal(t, 1, l.Rows[0][0].Height)
			// 2
			assert.Equal(t, 2, len(l.Rows[1]))
			assert.Equal(t, 1.0/4, l.Rows[1][0].Weight)
			assert.Equal(t, 1, l.Rows[1][0].Height)
			assert.Equal(t, 3.0/4, l.Rows[1][1].Weight)
			assert.Equal(t, 1, l.Rows[1][1].Height)
			// 3
			assert.Equal(t, 1, len(l.Rows[2]))
			assert.Equal(t, 1.0, l.Rows[2][0].Weight)
			assert.Equal(t, 1, l.Rows[2][0].Height)
			// 4
			assert.Equal(t, 2, len(l.Rows[3]))
			assert.Equal(t, 0.5, l.Rows[3][0].Weight)
			assert.Equal(t, 1, l.Rows[3][0].Height)
			assert.Equal(t, 0.5, l.Rows[3][1].Weight)
			assert.Equal(t, 1, l.Rows[3][1].Height)
		}},
		{"2:cpu\ndisk\nmem", func(l layout) {
			assert.Equal(t, 3, len(l.Rows))
			assert.Equal(t, 1, len(l.Rows[0]))
			assert.Equal(t, 2, l.Rows[0][0].Height)
			assert.Equal(t, 1, len(l.Rows[1]))
			assert.Equal(t, 1, l.Rows[1][0].Height)
			assert.Equal(t, 1, len(l.Rows[2]))
			assert.Equal(t, 1, l.Rows[2][0].Height)
		}},
		{"2:cpu disk\nmem", func(l layout) {
			assert.Equal(t, 2, len(l.Rows))
			assert.Equal(t, 2, len(l.Rows[0]))
			assert.Equal(t, 2, l.Rows[0][0].Height)
			assert.Equal(t, 1, l.Rows[0][1].Height)
			assert.Equal(t, 1, len(l.Rows[1]))
			assert.Equal(t, 1, l.Rows[1][0].Height)
		}},
		{"cpu 2:disk\nmem", func(l layout) {
			assert.Equal(t, 2, len(l.Rows))
			assert.Equal(t, 2, len(l.Rows[0]))
			assert.Equal(t, 1, l.Rows[0][0].Height)
			assert.Equal(t, 2, l.Rows[0][1].Height)
			assert.Equal(t, 1, len(l.Rows[1]))
			assert.Equal(t, 1, l.Rows[1][0].Height)
		}},
		{"cpu disk\n2:mem", func(l layout) {
			assert.Equal(t, 2, len(l.Rows))
			assert.Equal(t, 2, len(l.Rows[0]))
			assert.Equal(t, 1, l.Rows[0][0].Height)
			assert.Equal(t, 1, l.Rows[0][1].Height)
			assert.Equal(t, 1, len(l.Rows[1]))
			assert.Equal(t, 2, l.Rows[1][0].Height)
		}},
		{"cpu 2:disk/3\nmem", func(l layout) {
			assert.Equal(t, 2, len(l.Rows))
			assert.Equal(t, 2, len(l.Rows[0]))
			assert.Equal(t, 1, l.Rows[0][0].Height)
			assert.Equal(t, 1.0/4, l.Rows[0][0].Weight)
			assert.Equal(t, 2, l.Rows[0][1].Height)
			assert.Equal(t, 3.0/4, l.Rows[0][1].Weight)
			assert.Equal(t, 1, len(l.Rows[1]))
			assert.Equal(t, 1, l.Rows[1][0].Height)
			assert.Equal(t, 1.0, l.Rows[1][0].Weight)
		}},
		{"2:cpu disk\nmem/3", func(l layout) {
			assert.Equal(t, 2, len(l.Rows))
			assert.Equal(t, 2, len(l.Rows[0]))
			assert.Equal(t, 2, l.Rows[0][0].Height)
			assert.Equal(t, 0.5, l.Rows[0][0].Weight)
			assert.Equal(t, 1, l.Rows[0][1].Height)
			assert.Equal(t, 0.5, l.Rows[0][1].Weight)
			assert.Equal(t, 1, len(l.Rows[1]))
			assert.Equal(t, 1, l.Rows[1][0].Height)
			assert.Equal(t, 1.0, l.Rows[1][0].Weight)
		}},
		{"cpu/2 mem/1 6:procs\n3:temp/1 2:disk/2\npower\nnet procs", func(l layout) {
			assert.Equal(t, 4, len(l.Rows))
			// First row
			assert.Equal(t, 3, len(l.Rows[0]))
			assert.Equal(t, 1, l.Rows[0][0].Height)
			assert.Equal(t, 0.5, l.Rows[0][0].Weight)
			assert.Equal(t, 1, l.Rows[0][1].Height)
			assert.Equal(t, 0.25, l.Rows[0][1].Weight)
			assert.Equal(t, 6, l.Rows[0][2].Height)
			assert.Equal(t, 0.25, l.Rows[0][2].Weight)
			// Second row
			assert.Equal(t, 2, len(l.Rows[1]))
			assert.Equal(t, 3, l.Rows[1][0].Height)
			assert.Equal(t, 1/3.0, l.Rows[1][0].Weight)
			assert.Equal(t, 2, l.Rows[1][1].Height)
			assert.Equal(t, 2/3.0, l.Rows[1][1].Weight)
			// Third row
			assert.Equal(t, 1, len(l.Rows[2]))
			assert.Equal(t, 1, l.Rows[2][0].Height)
			assert.Equal(t, 1.0, l.Rows[2][0].Weight)
			// Fourth row
			assert.Equal(t, 2, len(l.Rows[3]))
			assert.Equal(t, 1, l.Rows[3][0].Height)
			assert.Equal(t, 0.5, l.Rows[3][0].Weight)
			assert.Equal(t, 1, l.Rows[3][1].Height)
			assert.Equal(t, 0.5, l.Rows[3][1].Weight)
		}},
	}

	for _, tc := range tests {
		in := strings.NewReader(tc.i)
		l := ParseLayout(in)
		tc.f(l)
	}
}
