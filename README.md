#GoSpy [![Build Status](https://travis-ci.org/cfmobile/gospy.svg?branch=master)](https://travis-ci.org/cfmobile/gospy)

Go testing utility that lets you monitor calls to a function, verify which arguments were passed in and mock its behaviour.

This was created with unit testing in mind, to make it easier to verify interactions with dependencies and isolate components. Inspired by [Counterfeiter](https://github.com/maxbrunsfeld/counterfeiter) and [Cedar's Doubles](https://github.com/pivotal/cedar/wiki/Writing-specs)

Compatibility verified with Go 1.4 and up.

##Installation

Just use go get

```
  go get github.com/cfmobile/gospy
```

##API

####Constructors
#####`func Spy(targetFuncPtr interface{}) *GoSpy`

Basic constructor of a `GoSpy` object.
This constructor doesn't modify the behaviour of the target function.

**`targetFuncPtr`** has to be a pointer to a function. Any other type will cause the constructor to panic.

**Returns:** a pointer to a new `GoSpy` object. After the constructor returns, the target function will be monitored for calls until `spy.Restore()` is called.


#####`func SpyAndFake(targetFuncPtr interface{}) *GoSpy`

Constructor of a `GoSpy` object that modifies the target's behaviour to just return the Zero value for each of it's return values.

**`targetFuncPtr`** has to be a pointer to a function. Any other type will cause the constructor to panic.

**Returns:** a pointer to a new `GoSpy` object. After the constructor returns, the target function have it's behaviour modified and will be monitored for calls until `spy.Restore()` is called.

#####`func SpyAndFakeWithReturn(targetFuncPtr interface{}, fakeReturnValues ...interface{}) *GoSpy`

Constructor of a `GoSpy` object that modifies the target's behaviour to return a matching set of mock return values specified by the user.

**`targetFuncPtr`** has to be a pointer to a function. Any other type will cause the constructor to panic.

**`fakeReturnValues`** (variadic) has to be a set of values of the same type and appearing in the same order as the return values of the target.

**Note 1:** Passing a mock value of a non-matching type, passing an incomplete list of mock values (one for each return value) or passing the list in the wrong order will all cause this constructor to panic.

**Note 2:** Passing no mock values at all will cause this constructor to behave like `SpyAndFake()`

**Returns:** a pointer to a new `GoSpy` object. After the constructor returns, the target function have it's behaviour modified and will be monitored for calls until `spy.Restore()` is called.

#####`func SpyAndFakeWithFunc(targetFuncPtr interface{}, mockFunc interface{}) *GoSpy`

Constructor of a `GoSpy` object that modifies the target's behaviour to be replaced by another specified mock function (with identical signature).

**`targetFuncPtr`** has to be a pointer to a function. Any other type will cause the constructor to panic.

**`mockFunc`** has to be a function with the exact same signature as the target's, with alternate behaviour.

**Note 1:** Passing nil or a mock function with a different signature will cause this constructor to panic.

**Returns:** a pointer to a new `GoSpy` object. After the constructor returns, the target function have it's behaviour modified and will be monitored for calls until `spy.Restore()` is called.
