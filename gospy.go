package gospy

import (
	"reflect"
	"github.com/cfmobile/gmock"
	"fmt"
)

type ParamList []interface{}

type CallList []ParamList

type GoSpy struct {
	Called  bool
	Calls   CallList
	mock   *gmock.GMock
}

func (self *GoSpy) Restore() {
	self.mock.Restore()
}

func (self *GoSpy) Reset() {
	self.Called = false
	self.Calls = nil
}

func (self *GoSpy) storeCall(arguments []reflect.Value) {
	self.Called = true;

	var call ParamList
	for _, arg := range arguments {
		call = append(call, arg.Interface())
	}

	self.Calls = append(self.Calls, call)
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
	targetType := reflect.ValueOf(targetFuncVar).Type()
	mockFuncType := reflect.ValueOf(mockFunc).Type()

	if targetType != mockFuncType {
		panic(fmt.Sprintf("Fake function has to have the same signature as the target [target: %+v, mock: %+v]", targetType, mockFuncType))
	}

	spy := createSpy(targetFuncVar)
	fakeFuncFn := spy.getFnWithMockFunc(mockFunc)
	spy.setTargetFn(fakeFuncFn)
	return spy
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

func (self *GoSpy) getDefaultFn() (func(args []reflect.Value) []reflect.Value) {
	return self.mock.GetOriginal().Call
}

func (self *GoSpy) getFnWithReturnValues(fakeReturnValues []interface{}) (func(args []reflect.Value) []reflect.Value) {
	// Gets the expected number of return values from the target
	var numReturnValues = self.mock.GetTarget().Type().NumOut()

	if fakeReturnValues != nil && numReturnValues != len(fakeReturnValues) {
		panic("Invalid number of return values. Either specify all or none")
	}

	return func(args []reflect.Value) (results []reflect.Value) {
		// Builds slice of return values, if required
		targetType := self.mock.GetTarget().Type()
		for i := 0; i < numReturnValues; i++ {
			returnItem := reflect.New(targetType.Out(i))

			var returnElem = returnItem.Elem()

			// Gets value for return from fakeReturnValues, or leaves default constructed value if not available
			if fakeReturnValues != nil && fakeReturnValues[i] != nil {
				returnElem.Set(reflect.ValueOf(fakeReturnValues[i]))
			}

			results = append(results, returnElem)
		}

		return results
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

	spy := &GoSpy{Called: false, Calls: nil, mock: gmock.CreateMockWithTarget(targetFuncPtr)}

	return spy
}

func targetIsValid(target interface{}) bool {
	// Target has to be a ptr to a function
	targetValue := reflect.ValueOf(target)
	return targetValue.Kind() == reflect.Ptr &&
		   targetValue.Elem().Kind() == reflect.Func
}