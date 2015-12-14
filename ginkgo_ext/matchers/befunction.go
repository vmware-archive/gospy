package matchers

import (
	"reflect"
	"fmt"
)

func fptr(f interface{}) uintptr {
	return reflect.ValueOf(f).Pointer()
}

type _BeFunctionMatcher struct {
	expected interface{}
}

func (matcher *_BeFunctionMatcher) Match(actual interface{}) (success bool, err error) {
	success = fptr(actual) == fptr(matcher.expected)
	return
}

func (matcher *_BeFunctionMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nto be function\n\t%#v", actual, matcher.expected)
}

func (matcher *_BeFunctionMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nnot to be function\n\t%#v", actual, matcher.expected)
}



