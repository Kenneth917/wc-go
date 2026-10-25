package wc

import (
	"strings"
	"testing"
)

func TestCount(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Counts
	}{
		{"empty", "", Counts{}},
		{"one line", "hello\n", Counts{Lines: 1, Words: 1, Chars: 6, Bytes: 6}},
		{"several words", "the quick brown fox\n", Counts{Lines: 1, Words: 4, Chars: 20, Bytes: 20}},
		{"no trailing newline is not a line", "a b", Counts{Lines: 0, Words: 2, Chars: 3, Bytes: 3}},
		{"multi-byte chars differ from bytes", "héllo\n", Counts{Lines: 1, Words: 1, Chars: 6, Bytes: 7}},
		{"runs of whitespace are one separator", "a \t\n  b\n", Counts{Lines: 2, Words: 2, Chars: 8, Bytes: 8}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Count(strings.NewReader(tt.input))
			if err != nil {
				t.Fatalf("Count returned error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Count(%q) = %+v, want %+v", tt.input, got, tt.want)
			}
		})
	}
}
