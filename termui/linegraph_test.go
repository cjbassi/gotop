package termui

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestLess(t *testing.T) {
	tests := []struct {
		a, b string
		e    bool
	}{
		{a: "abc", b: "def", e: true},
		{a: "abc", b: "abc", e: true},
		{a: "def", b: "abc", e: false},
		{a: "1", b: "10", e: true},
		{a: "1", b: "2", e: true},
		{a: "a2", b: "2", e: false},
		{a: "a2", b: "a10", e: true},
		{a: "a20", b: "a2", e: false},
		{a: "abc", b: "abc20", e: true},
		{a: "a1", b: "abc", e: true},
		{a: "a20", b: "abc", e: true},
		{a: "abc20", b: "abc", e: false},
		{a: "abc20", b: "def2", e: true},
		{a: "abc20", b: "abc2", e: false},
		{a: "abc20", b: "abc20", e: true},
		{a: "abc30", b: "abc20", e: false},
		{a: "abc20a", b: "abc20", e: false},
		{a: "abc20", b: "abc20a", e: true},
		{a: "abc20", b: "abc2a", e: false},
		{a: "abc20", b: "abc3a", e: false},
		{a: "abc20", b: "abc2abc", e: false},
	}
	for _, k := range tests {
		n := numbered([]string{k.a, k.b})
		g := n.Less(0, 1)
		assert.Equal(t, k.e, g, "%s < %s %t", k.a, k.b, g)
	}
}

func TestSort(t *testing.T) {
	tests := []struct {
		in, ex numbered
	}{
		{
			in: numbered{"abc", "def", "abc", "def", "abc", "1", "10", "1", "2", "a2", "2", "a2", "a10", "a20", "a2", "abc20", "def2", "abc20", "abc2", "abc20", "abc20", "abc30", "abc20", "abc20a", "abc20", "abc20", "abc20a", "abc20", "abc2a", "abc"},
			ex: numbered{"1", "1", "2", "2", "10", "a2", "a2", "a2", "a10", "a20", "abc", "abc", "abc", "abc", "abc2", "abc2a", "abc20", "abc20", "abc20", "abc20", "abc20", "abc20", "abc20", "abc20", "abc20a", "abc20a", "abc30", "def", "def", "def2"},
		},
		{
			in: numbered{"CPU12", "CPU11", "CPU9", "CPU3", "CPU4", "CPU0", "CPU6", "CPU7", "CPU8", "CPU5", "CPU10", "CPU1", "CPU2"},
			ex: numbered{"CPU0", "CPU1", "CPU2", "CPU3", "CPU4", "CPU5", "CPU6", "CPU7", "CPU8", "CPU9", "CPU10", "CPU11", "CPU12"},
		},
	}

	for _, k := range tests {
		sort.Sort(k.in)
		assert.Equal(t, k.ex, k.in)
	}
}

func BenchmarkLess(b *testing.B) {
	tests := []numbered{
		numbered([]string{"abc", "def"}),
		numbered([]string{"abc", "abc"}),
		numbered([]string{"def", "abc"}),
		numbered([]string{"1", "10"}),
		numbered([]string{"1", "2"}),
		numbered([]string{"a2", "2"}),
		numbered([]string{"a2", "a10"}),
		numbered([]string{"a20", "a2"}),
		numbered([]string{"abc", "abc20"}),
		numbered([]string{"a1", "abc"}),
		numbered([]string{"a20", "abc"}),
		numbered([]string{"abc20", "abc"}),
		numbered([]string{"abc20", "def2"}),
		numbered([]string{"abc20", "abc2"}),
		numbered([]string{"abc20", "abc20"}),
		numbered([]string{"abc30", "abc20"}),
		numbered([]string{"abc20a", "abc20"}),
		numbered([]string{"abc20", "abc20a"}),
		numbered([]string{"abc20", "abc2a"}),
		numbered([]string{"abc20", "abc3a"}),
		numbered([]string{"abc20", "abc2abc"}),
	}
	for i := 0; i < b.N; i++ {
		for _, t := range tests {
			t.Less(0, 1)
		}
	}
}
