// Code generated by mockery v2.43.1. DO NOT EDIT.

package mocks

import (
	iplocator "github.com/clambin/mediamon/v2/pkg/iplocator"
	mock "github.com/stretchr/testify/mock"
)

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

// Locate provides a mock function with given fields: _a0
func (_m *Locator) Locate(_a0 string) (iplocator.Location, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Locate")
	}

	var r0 iplocator.Location
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (iplocator.Location, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) iplocator.Location); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(iplocator.Location)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
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
//   - _a0 string
func (_e *Locator_Expecter) Locate(_a0 interface{}) *Locator_Locate_Call {
	return &Locator_Locate_Call{Call: _e.mock.On("Locate", _a0)}
}

func (_c *Locator_Locate_Call) Run(run func(_a0 string)) *Locator_Locate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *Locator_Locate_Call) Return(_a0 iplocator.Location, _a1 error) *Locator_Locate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Locator_Locate_Call) RunAndReturn(run func(string) (iplocator.Location, error)) *Locator_Locate_Call {
	_c.Call.Return(run)
	return _c
}

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