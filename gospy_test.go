package gospy_test

import (
	. "github.com/cfmobile/gospy"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	kOriginalStringReturn = "original string value"
	kOriginalFloatReturn = float64(123.45)
)

var _ = Describe("GoSpy", func() {
	var subject *GoSpy

	var functionToSpy func(string, int, bool) (string, float64)
	var panicked bool

	BeforeEach(func() {
	    subject = nil
		panicked = false
		functionToSpy = func(string, int, bool) (string, float64) {
			return kOriginalStringReturn, kOriginalFloatReturn
		}
	})

	panicRecover := func() {
		panicked = recover() != nil
	}

	Describe("Constructors", func() {

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

				It("should not affect the function's behaviour", func() {
					stringResult, floatResult := functionToSpy("something", 10, false)
					Expect(stringResult).To(Equal(kOriginalStringReturn))
					Expect(floatResult).To(Equal(kOriginalFloatReturn))
				})
	        })
	    })
	})

	Context("when a valid GoSpy object is created", func() {

		BeforeEach(func() {
			subject = Spy(&functionToSpy)
		})

		It("Called() should indicate that the function hasn't been called yet", func() {
		    Expect(subject.Called()).To(BeFalse())
		})

		It("CallCount() should indicate a call count of zero", func() {
		    Expect(subject.CallCount()).To(BeZero())
		})

		It("Calls() should indicate a nil call list", func() {
		    Expect(subject.Calls()).To(BeNil())
		})

		Context("when ArgsForCall() is called with no calls in the Spy", func() {
			BeforeEach(func() {
			    defer panicRecover()
				subject.ArgsForCall(0)
			})

			It("should panic", func() {
				Expect(panicked).To(BeTrue())
			})
		})

		goSpyResetTests := func() {
			Context("when Reset() is called", func() {
				BeforeEach(func() {
					subject.Reset()
				})

				It("should zero the call count", func() {
					Expect(subject.CallCount()).To(BeZero())
				})

				It("should return a nil call list", func() {
					Expect(subject.Calls()).To(BeNil())
				})

				It("should have reset the call indicator", func() {
					Expect(subject.Called()).To(BeFalse())
				})
			})
		}

		Context("and the monitored function is called once", func() {
			kFirstArg, kSecondArg, kThirdArg := "test value", 101, true

			BeforeEach(func() {
			    functionToSpy(kFirstArg, kSecondArg, kThirdArg)
			})

			It("Called() should indicate that the function was called", func() {
				Expect(subject.Called()).To(BeTrue())
			})

			It("CallCount() should indicate that a call was made", func() {
				Expect(subject.CallCount()).To(Equal(1))
			})

			It("Calls() should return a valid call list", func() {
			    Expect(subject.Calls()).NotTo(BeNil())
			})

			It("ArgsForCall() should return the arguments that were used in the call", func() {
			    Expect(subject.ArgsForCall(0)).To(Equal(ArgList{kFirstArg, kSecondArg, kThirdArg}))
			})

			goSpyResetTests()
		})

		Context("and the monitored function is called several times", func() {
			BeforeEach(func() {
			    functionToSpy("call 1", 1, true)
				functionToSpy("call 2", 2, false)
				functionToSpy("call 3", 3, true)
			})

			It("CallCount() should reflect the right number of calls", func() {
			    Expect(subject.CallCount()).To(Equal(3))
			})

			It("Calls() should return the arguments for each call in the order they were made", func() {
				expectedCallList := CallList{
					{"call 1", 1, true},
					{"call 2", 2, false},
					{"call 3", 3, true},
				}

				Expect(subject.Calls()).To(Equal(expectedCallList))
			})

			It("ArgsForCall(n) should return the arguments for the n-th call (0-based index) ", func() {
			    Expect(subject.ArgsForCall(0)).To(Equal(ArgList{"call 1", 1, true}))
			    Expect(subject.ArgsForCall(1)).To(Equal(ArgList{"call 2", 2, false}))
			    Expect(subject.ArgsForCall(2)).To(Equal(ArgList{"call 3", 3, true}))
			})

			goSpyResetTests()
		})
	})
})
