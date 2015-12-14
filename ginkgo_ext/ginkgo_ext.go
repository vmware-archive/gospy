package ginkgo_ext

import (
	"github.com/cfmobile/gospy"
	. "github.com/onsi/ginkgo"
)

func GSpy(target interface{}) *gospy.GoSpy {
	var spy *gospy.GoSpy

	BeforeEach(func() {
	    spy = gospy.Spy(target)
	})

	AfterEach(func() {
	    spy.Restore()
	})

	return spy
}

func GSpyAndFake(target interface{}) *gospy.GoSpy {
	var spy *gospy.GoSpy

	BeforeEach(func() {
	    spy = gospy.SpyAndFake(target)
	})

	AfterEach(func() {
	    spy.Restore()
	})

	return spy
}

func GSpyAndFakeWithReturn(target interface{}, fakeReturnValues ...interface{}) *gospy.GoSpy {
	var spy *gospy.GoSpy

	BeforeEach(func() {
	    spy = gospy.SpyAndFakeWithReturn(target, fakeReturnValues...)
	})

	AfterEach(func() {
	    spy.Restore()
	})

	return spy
}

func GSpyAndFakeWithFunc(target interface{}, fakeFunc interface{}) *gospy.GoSpy {
	var spy *gospy.GoSpy

	BeforeEach(func() {
	    spy = gospy.SpyAndFakeWithFunc(target, fakeFunc)
	})

	AfterEach(func() {
	    spy.Restore()
	})

	return spy
}
