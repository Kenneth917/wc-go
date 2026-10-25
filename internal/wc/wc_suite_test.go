package wc_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Wc Suite")
}
