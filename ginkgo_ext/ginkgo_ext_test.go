package ginkgo_ext_test

import (
	. "github.com/cfmobile/gospy/ginkgo_ext"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"fmt"
	. "github.com/cfmobile/gospy/ginkgo_ext/matchers"
)

const (
	kOriginalString = "original string return value"
	kOriginalInt = 12345
)

var _ = Describe("Ginkgo Extensions", func() {
	myFunc := func(string, int) (string, int) {
		return kOriginalString, kOriginalInt
	}

	var originalFunc func(string, int) (string, int)

	BeforeEach(func() {
	    originalFunc = myFunc
	})

	Describe("GSpy", func() {
		It("before GSpy, the function should be unmodified", func() {
			Expect(myFunc).To(BeFunction(originalFunc))
		})

		Context("when GSpy is built", func() {
		    spy := GSpy(&myFunc)

			It("should have modified the function", func() {
			    Expect(myFunc).NotTo(BeFunction(originalFunc))
			})

			FContext("when the target function is called", func() {
				var stringResult string
				var intResult int

				BeforeEach(func() {
					fmt.Println(spy)
			        stringResult, intResult = myFunc("asd", 10)
			    })

				It("should have monitored the function call", func() {
				    Expect(spy.CallCount()).To(Equal(1))
					Expect(spy.ArgsForCall(0)).To(MatchArgs("asd", 10))
				})
			})

			AfterEach(func() {
			    // It should have restored the function in the AfterEach
				Expect(myFunc).To(BeFunction(originalFunc))
			})


		})
	})

})
