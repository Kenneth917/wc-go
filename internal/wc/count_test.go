package wc_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Kenneth917/wc-go/internal/wc"
)

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
})
