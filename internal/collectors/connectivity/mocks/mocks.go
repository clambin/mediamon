// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package mocks

import (
	"github.com/clambin/mediamon/v2/iplocator"
	mock "github.com/stretchr/testify/mock"
)

// NewLocator creates a new instance of Locator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLocator(t interface {
	mock.TestingT
	Cleanup(func())
}) *Locator {
	mock := &Locator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// Locator is an autogenerated mock type for the Locator type
type Locator struct {
	mock.Mock
}

type Locator_Expecter struct {
	mock *mock.Mock
}

func (_m *Locator) EXPECT() *Locator_Expecter {
	return &Locator_Expecter{mock: &_m.Mock}
}

// Locate provides a mock function for the type Locator
func (_mock *Locator) Locate(s string) (iplocator.Location, error) {
	ret := _mock.Called(s)

	if len(ret) == 0 {
		panic("no return value specified for Locate")
	}

	var r0 iplocator.Location
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(string) (iplocator.Location, error)); ok {
		return returnFunc(s)
	}
	if returnFunc, ok := ret.Get(0).(func(string) iplocator.Location); ok {
		r0 = returnFunc(s)
	} else {
		r0 = ret.Get(0).(iplocator.Location)
	}
	if returnFunc, ok := ret.Get(1).(func(string) error); ok {
		r1 = returnFunc(s)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// Locator_Locate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Locate'
type Locator_Locate_Call struct {
	*mock.Call
}

// Locate is a helper method to define mock.On call
//   - s string
func (_e *Locator_Expecter) Locate(s interface{}) *Locator_Locate_Call {
	return &Locator_Locate_Call{Call: _e.mock.On("Locate", s)}
}

func (_c *Locator_Locate_Call) Run(run func(s string)) *Locator_Locate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 string
		if args[0] != nil {
			arg0 = args[0].(string)
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *Locator_Locate_Call) Return(location iplocator.Location, err error) *Locator_Locate_Call {
	_c.Call.Return(location, err)
	return _c
}

func (_c *Locator_Locate_Call) RunAndReturn(run func(s string) (iplocator.Location, error)) *Locator_Locate_Call {
	_c.Call.Return(run)
	return _c
}
