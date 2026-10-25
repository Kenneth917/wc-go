package wc_test

import (
	"errors"
	"io"
	"runtime"
	"strings"
	"testing/iotest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Kenneth917/wc-go/internal/wc"
)

// repeating yields pattern n times without ever holding all of it in memory,
// so tests can simulate inputs far larger than they could afford to allocate.
type repeating struct {
	pattern   []byte
	remaining int // repetitions left
	offset    int // position inside the current repetition
}

func newRepeating(pattern string, times int) io.Reader {
	return &repeating{pattern: []byte(pattern), remaining: times}
}

func (r *repeating) Read(p []byte) (int, error) {
	if r.remaining == 0 {
		return 0, io.EOF
	}
	n := 0
	for n < len(p) && r.remaining > 0 {
		c := copy(p[n:], r.pattern[r.offset:])
		n += c
		r.offset += c
		if r.offset == len(r.pattern) {
			r.offset = 0
			r.remaining--
		}
	}
	return n, nil
}

var _ = Describe("Count", func() {
	count := func(input string) wc.Counts {
		got, err := wc.Count(strings.NewReader(input))
		Expect(err).NotTo(HaveOccurred())
		return got
	}

	Context("when the input is empty", func() {
		It("reports zero for everything", func() {
			Expect(count("")).To(Equal(wc.Counts{}))
		})
	})

	Context("when the input is a single line of one word", func() {
		It("reports one line, one word and every character", func() {
			Expect(count("hello\n")).To(Equal(wc.Counts{Lines: 1, Words: 1, Chars: 6, Bytes: 6}))
		})
	})

	Context("when a line has several words", func() {
		It("counts each word", func() {
			Expect(count("the quick brown fox\n").Words).To(Equal(4))
		})
	})

	Context("when the last line has no trailing newline", func() {
		It("does not count it as a line, because lines are newline characters", func() {
			Expect(count("a b")).To(Equal(wc.Counts{Lines: 0, Words: 2, Chars: 3, Bytes: 3}))
		})
	})

	Context("when the text has multi-byte characters", func() {
		It("counts characters and bytes separately", func() {
			got := count("héllo\n")
			Expect(got.Chars).To(Equal(6))
			Expect(got.Bytes).To(Equal(7))
		})
	})

	Context("when words are separated by runs of spaces, tabs and newlines", func() {
		It("treats each run as a single separator", func() {
			Expect(count("a \t\n  b\n")).To(Equal(wc.Counts{Lines: 2, Words: 2, Chars: 8, Bytes: 8}))
		})
	})

	Context("when the reader delivers one byte at a time", func() {
		It("gives the same answer as a normal read, even when a character is split across reads", func() {
			input := "héllo wörld\nfoo bar\n"
			slow, err := wc.Count(iotest.OneByteReader(strings.NewReader(input)))
			Expect(err).NotTo(HaveOccurred())
			Expect(slow).To(Equal(count(input)))
			Expect(slow).To(Equal(wc.Counts{Lines: 2, Words: 4, Chars: 20, Bytes: 22}))
		})

		It("counts a word that spans several reads only once", func() {
			slow, err := wc.Count(iotest.OneByteReader(strings.NewReader("extraordinary")))
			Expect(err).NotTo(HaveOccurred())
			Expect(slow.Words).To(Equal(1))
		})
	})

	Context("when the underlying reader fails", func() {
		It("returns the error instead of a partial answer", func() {
			boom := errors.New("disk on fire")
			_, err := wc.Count(iotest.ErrReader(boom))
			Expect(err).To(MatchError(boom))
		})
	})

	Context("when the input is a single enormous line with no newline", func() {
		It("is not limited by a maximum line length", func() {
			// 8 MiB on one line: well past bufio.Scanner's 64 KiB default token limit.
			const size = 8 << 20
			got, err := wc.Count(newRepeating("x", size))
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(wc.Counts{Lines: 0, Words: 1, Chars: size, Bytes: size}))
		})
	})

	Context("when the input is far larger than the memory we are willing to use", func() {
		It("counts it in constant memory", func() {
			// 10 million repetitions of a 6-byte line is 60 MB. Reading it all into
			// memory would allocate at least that much; streaming should need a tiny fraction.
			const reps = 10_000_000
			var before, after runtime.MemStats
			runtime.GC()
			runtime.ReadMemStats(&before)

			got, err := wc.Count(newRepeating("ab cd\n", reps))

			runtime.ReadMemStats(&after)
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(wc.Counts{Lines: reps, Words: 2 * reps, Chars: 6 * reps, Bytes: 6 * reps}))
			Expect(after.TotalAlloc-before.TotalAlloc).To(BeNumerically("<", 8<<20),
				"allocated more than 8 MiB, so the input is probably being buffered whole")
		})
	})
})
