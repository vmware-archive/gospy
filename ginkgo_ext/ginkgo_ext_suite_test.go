package ginkgo_ext_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGinkgoExt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GinkgoExt Suite")
}
