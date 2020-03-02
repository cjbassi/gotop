package utils

import "testing"

const (
	ELLIPSIS = "…"
)

func TestTruncateFront(t *testing.T) {
	tests := []struct {
		s      string
		w      int
		prefix string
		want   string
	}{
		{"", 0, ELLIPSIS, ""},
		{"", 1, ELLIPSIS, ""},
		{"", 10, ELLIPSIS, ""},

		{"abcdef", 0, ELLIPSIS, ELLIPSIS},
		{"abcdef", 1, ELLIPSIS, ELLIPSIS},
		{"abcdef", 2, ELLIPSIS, ELLIPSIS + "f"},
		{"abcdef", 5, ELLIPSIS, ELLIPSIS + "cdef"},
		{"abcdef", 6, ELLIPSIS, "abcdef"},
		{"abcdef", 10, ELLIPSIS, "abcdef"},

		{"abcdef", 0, "...", "..."},
		{"abcdef", 1, "...", "..."},
		{"abcdef", 3, "...", "..."},
		{"abcdef", 4, "...", "...f"},
		{"abcdef", 5, "...", "...ef"},
		{"abcdef", 6, "...", "abcdef"},
		{"abcdef", 10, "...", "abcdef"},

		{"｟full～width｠", 15, ".", "｟full～width｠"},
		{"｟full～width｠", 14, ".", ".full～width｠"},
		{"｟full～width｠", 13, ".", ".ull～width｠"},
		{"｟full～width｠", 10, ".", ".～width｠"},
		{"｟full～width｠", 9, ".", ".width｠"},
		{"｟full～width｠", 8, ".", ".width｠"},
		{"｟full～width｠", 3, ".", ".｠"},
		{"｟full～width｠", 2, ".", "."},
	}

	for _, test := range tests {
		if got := TruncateFront(test.s, test.w, test.prefix); got != test.want {
			t.Errorf("TruncateFront(%q, %d, %q) = %q; want %q", test.s, test.w, test.prefix, got, test.want)
		}
	}
}
