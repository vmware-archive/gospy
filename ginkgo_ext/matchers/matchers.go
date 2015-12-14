package matchers

import (
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/matchers"
	"github.com/onsi/gomega/types"
	"github.com/cfmobile/gospy"
)

func BeFunction(expected interface{}) OmegaMatcher {
	return &_BeFunctionMatcher{expected}
}

var ContainFunction = func(expected interface{}) OmegaMatcher {
	return ContainElement(BeFunction(expected))
}

func MatchArgs(expected ...interface{}) types.GomegaMatcher {
	return &matchers.EqualMatcher{Expected: gospy.ArgList(expected)}
}
