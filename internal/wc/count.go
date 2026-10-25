// Package wc counts lines, words, characters and bytes in a stream.
package wc

import (
	"bufio"
	"io"
	"unicode"
)

// Counts holds the totals for one input.
type Counts struct {
	Lines int
	Words int
	Chars int
	Bytes int
}

// Count reads r to EOF and returns its totals. It streams the input through a
// fixed-size buffer, so memory use does not grow with the size of the input or
// the length of any line.
func Count(r io.Reader) (Counts, error) {
	var c Counts
	br := bufio.NewReader(r)
	inWord := false
	for {
		ch, size, err := br.ReadRune()
		if err == io.EOF {
			return c, nil
		}
		if err != nil {
			return Counts{}, err
		}
		c.Bytes += size
		c.Chars++
		if ch == '\n' {
			c.Lines++
		}
		if unicode.IsSpace(ch) {
			inWord = false
		} else if !inWord {
			inWord = true
			c.Words++
		}
	}
}
