// Command wc prints line, word, character and byte counts for files or stdin.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/Kenneth917/wc-go/internal/wc"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("wc", flag.ContinueOnError)
	fs.SetOutput(stderr)
	lines := fs.Bool("l", false, "print the line count")
	words := fs.Bool("w", false, "print the word count")
	chars := fs.Bool("m", false, "print the character count")
	bytes := fs.Bool("c", false, "print the byte count")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if !*lines && !*words && !*chars && !*bytes {
		*lines, *words, *bytes = true, true, true // the default, as in wc
	}
	show := func(c wc.Counts, name string) {
		var out string
		for _, col := range []struct {
			on bool
			n  int
		}{{*lines, c.Lines}, {*words, c.Words}, {*chars, c.Chars}, {*bytes, c.Bytes}} {
			if col.on {
				out += fmt.Sprintf("%8d", col.n)
			}
		}
		if name != "" {
			out += " " + name
		}
		fmt.Fprintln(stdout, out)
	}

	if fs.NArg() == 0 {
		c, err := wc.Count(stdin)
		if err != nil {
			fmt.Fprintf(stderr, "wc: %v\n", err)
			return 1
		}
		show(c, "")
		return 0
	}

	status := 0
	var total wc.Counts
	for _, name := range fs.Args() {
		c, err := countFile(name)
		if err != nil {
			fmt.Fprintf(stderr, "wc: %v\n", err)
			status = 1
			continue
		}
		show(c, name)
		total.Lines += c.Lines
		total.Words += c.Words
		total.Chars += c.Chars
		total.Bytes += c.Bytes
	}
	if fs.NArg() > 1 {
		show(total, "total")
	}
	return status
}

func countFile(name string) (wc.Counts, error) {
	f, err := os.Open(name)
	if err != nil {
		return wc.Counts{}, err
	}
	defer f.Close()
	return wc.Count(f)
}
