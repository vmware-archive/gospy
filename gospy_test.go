package gospy_test

import (
	. "github.com/cfmobile/gospy"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"fmt"
)

const (
	kOriginalStringReturn = "original string value"
	kOriginalFloatReturn = float64(123.45)
)

var _ = Describe("GoSpy", func() {
	var subject *GoSpy

	var functionToSpy func(string, int, bool) (string, float64)

	BeforeEach(func() {
	    subject = nil
		functionToSpy = func(string, int, bool) (string, float64) {
			return kOriginalStringReturn, kOriginalFloatReturn
		}
	})

	Describe("Constructors", func() {
		var panicked bool

		panicRecover := func() {
			if r := recover(); r != nil {
				panicked = true
				fmt.Println("Panic: ", r)
			} else {
				panicked = false
			}
		}

		BeforeEach(func() {
		    panicked = false
		})

	    Describe("Spy", func() {

	        Context("when calling Spy with a valid function pointer", func() {
				BeforeEach(func() {
					defer panicRecover()
				    subject = Spy(&functionToSpy)
				})

				It("should not have panicked", func() {
				    Expect(panicked).To(BeFalse())
				})

				It("should have returned a valid *GoSpy object", func() {
				    Expect(subject).NotTo(BeNil())
				})
	        })
	    })
	})
})
