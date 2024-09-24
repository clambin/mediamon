// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	iplocator "github.com/clambin/mediamon/v2/iplocator"
	mock "github.com/stretchr/testify/mock"
)

// IPLocator is an autogenerated mock type for the IPLocator type
type IPLocator struct {
	mock.Mock
}

type IPLocator_Expecter struct {
	mock *mock.Mock
}

func (_m *IPLocator) EXPECT() *IPLocator_Expecter {
	return &IPLocator_Expecter{mock: &_m.Mock}
}

// Locate provides a mock function with given fields: _a0
func (_m *IPLocator) Locate(_a0 string) (iplocator.Location, error) {
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

// IPLocator_Locate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Locate'
type IPLocator_Locate_Call struct {
	*mock.Call
}

// Locate is a helper method to define mock.On call
//   - _a0 string
func (_e *IPLocator_Expecter) Locate(_a0 interface{}) *IPLocator_Locate_Call {
	return &IPLocator_Locate_Call{Call: _e.mock.On("Locate", _a0)}
}

func (_c *IPLocator_Locate_Call) Run(run func(_a0 string)) *IPLocator_Locate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *IPLocator_Locate_Call) Return(_a0 iplocator.Location, _a1 error) *IPLocator_Locate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IPLocator_Locate_Call) RunAndReturn(run func(string) (iplocator.Location, error)) *IPLocator_Locate_Call {
	_c.Call.Return(run)
	return _c
}

// NewIPLocator creates a new instance of IPLocator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIPLocator(t interface {
	mock.TestingT
	Cleanup(func())
}) *IPLocator {
	mock := &IPLocator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
