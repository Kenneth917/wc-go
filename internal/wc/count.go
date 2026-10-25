// Package wc counts lines, words, characters and bytes in a stream.
package wc

import "io"

// Counts holds the totals for one input.
type Counts struct {
	Lines int
	Words int
	Chars int
	Bytes int
}

// Count reads r to EOF and returns its totals.
func Count(r io.Reader) (Counts, error) {
	return Counts{}, nil // TODO: make the tests pass
}
