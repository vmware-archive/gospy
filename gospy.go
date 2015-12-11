package gospy

import (
	"reflect"
	"github.com/cfmobile/gmock"
	"fmt"
	"errors"
)

type ArgList []interface{}

type CallList []ArgList

type GoSpy struct {
	calls CallList
	mock  *gmock.GMock
}

func Spy(targetFuncVar interface{}) *GoSpy {
	spy := createSpy(targetFuncVar)
	defaultFn := spy.getDefaultFn()
	spy.setTargetFn(defaultFn)
	return spy
}

func SpyAndFake(targetFuncVar interface{}) *GoSpy {
	return SpyAndFakeWithReturn(targetFuncVar) //nil fakeReturnValues will create default
}

func SpyAndFakeWithReturn(targetFuncVar interface{}, fakeReturnValues ...interface{}) *GoSpy {
	spy := createSpy(targetFuncVar)
	fakeReturnFn := spy.getFnWithReturnValues(fakeReturnValues)
	spy.setTargetFn(fakeReturnFn)
	return spy
}

func SpyAndFakeWithFunc(targetFuncVar interface{}, mockFunc interface{}) *GoSpy {
	spy := createSpy(targetFuncVar)

	if err := mockFuncIsValid(targetFuncVar, mockFunc); err != nil {
		panic(err.Error())
	}

	fakeFuncFn := spy.getFnWithMockFunc(mockFunc)
	spy.setTargetFn(fakeFuncFn)
	return spy
}

func (self *GoSpy) Called() bool {
	return self.CallCount() > 0
}

func (self *GoSpy) CallCount() int {
	return len(self.calls)
}

func (self *GoSpy) Calls() CallList {
	return self.calls
}

func (self *GoSpy) ArgsForCall(callIndex uint) ArgList {
	return self.calls[callIndex]
}

func (self *GoSpy) Reset() {
	self.calls = nil
}

func (self *GoSpy) Restore() {
	self.mock.Restore()
}

func (self *GoSpy) setTargetFn(fn func(args []reflect.Value) []reflect.Value) {
	targetType := self.mock.GetTarget().Type()
	wrapperFn := func(args []reflect.Value) []reflect.Value {
		self.storeCall(args)
		return reflect.MakeFunc(targetType, fn).Call(args)
	}

	targetFn := reflect.MakeFunc(targetType, wrapperFn)
	self.mock.Replace(targetFn.Interface())
}

func (self *GoSpy) storeCall(arguments []reflect.Value) {
	var call ArgList
	for _, arg := range arguments {
		call = append(call, arg.Interface())
	}

	self.calls = append(self.calls, call)
}

func (self *GoSpy) getDefaultFn() (func(args []reflect.Value) []reflect.Value) {
	return self.mock.GetOriginal().Call
}

func (self *GoSpy) getFnWithReturnValues(fakeReturnValues []interface{}) (func(args []reflect.Value) []reflect.Value) {
	targetType := self.mock.GetTarget().Type()

	// Gets the expected number of return values from the target
	var numReturnValues = targetType.NumOut()

	if fakeReturnValues != nil && numReturnValues != len(fakeReturnValues) {
		panic("Invalid number of return values. Either specify the exact number of return values or none for defaults")
	}

	res := make([]reflect.Value, 0)

	// Builds slice of return values, if required
	for i := 0; i < numReturnValues; i++ {
		returnItem := reflect.New(targetType.Out(i))

		var returnElem = returnItem.Elem()

		// Gets value for return from fakeReturnValues, or leaves default constructed value if not available
		if fakeReturnValues != nil && fakeReturnValues[i] != nil {
			returnElem.Set(reflect.ValueOf(fakeReturnValues[i]))
		}

		res = append(res, returnElem)
	}

	return func([]reflect.Value) []reflect.Value {
		return res
	}
}

func (self *GoSpy) getFnWithMockFunc(mockFunc interface{}) (func(args []reflect.Value) []reflect.Value) {
	mockFuncValue := reflect.ValueOf(mockFunc)
	if !mockFuncValue.IsValid() {
		targetType := self.mock.GetTarget().Type()
		mockFuncValue = reflect.Zero(targetType)
	}

	return mockFuncValue.Call
}

func createSpy(targetFuncPtr interface{}) *GoSpy {
	if !targetIsValid(targetFuncPtr) {
		panic("Spy target has to be the pointer to a Func variable")
	}

	spy := &GoSpy{calls: nil, mock: gmock.CreateMockWithTarget(targetFuncPtr)}

	return spy
}

func targetIsValid(target interface{}) bool {
	// Target has to be a ptr to a function
	targetValue := reflect.ValueOf(target)
	return targetValue.Kind() == reflect.Ptr &&
	targetValue.Elem().Kind() == reflect.Func
}

func mockFuncIsValid(target interface{}, mockFunc interface{}) error {
	targetType := reflect.ValueOf(target).Type().Elem() // target is a *func()
	mockFuncType := reflect.ValueOf(mockFunc).Type()

	if targetType != mockFuncType {
		return errors.New(fmt.Sprintf("Fake function has to have the same signature as the target [target: %+v, mock: %+v]", targetType, mockFuncType))
	}

	return nil
}
