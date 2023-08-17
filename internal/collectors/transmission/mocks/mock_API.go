// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	context "context"

	transmission "github.com/clambin/mediaclients/transmission"
	mock "github.com/stretchr/testify/mock"
)

// API is an autogenerated mock type for the API type
type API struct {
	mock.Mock
}

type API_Expecter struct {
	mock *mock.Mock
}

func (_m *API) EXPECT() *API_Expecter {
	return &API_Expecter{mock: &_m.Mock}
}

// GetSessionParameters provides a mock function with given fields: ctx
func (_m *API) GetSessionParameters(ctx context.Context) (transmission.SessionParameters, error) {
	ret := _m.Called(ctx)

	var r0 transmission.SessionParameters
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (transmission.SessionParameters, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) transmission.SessionParameters); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(transmission.SessionParameters)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// API_GetSessionParameters_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSessionParameters'
type API_GetSessionParameters_Call struct {
	*mock.Call
}

// GetSessionParameters is a helper method to define mock.On call
//   - ctx context.Context
func (_e *API_Expecter) GetSessionParameters(ctx interface{}) *API_GetSessionParameters_Call {
	return &API_GetSessionParameters_Call{Call: _e.mock.On("GetSessionParameters", ctx)}
}

func (_c *API_GetSessionParameters_Call) Run(run func(ctx context.Context)) *API_GetSessionParameters_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *API_GetSessionParameters_Call) Return(_a0 transmission.SessionParameters, _a1 error) *API_GetSessionParameters_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *API_GetSessionParameters_Call) RunAndReturn(run func(context.Context) (transmission.SessionParameters, error)) *API_GetSessionParameters_Call {
	_c.Call.Return(run)
	return _c
}

// GetSessionStatistics provides a mock function with given fields: ctx
func (_m *API) GetSessionStatistics(ctx context.Context) (transmission.SessionStats, error) {
	ret := _m.Called(ctx)

	var r0 transmission.SessionStats
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (transmission.SessionStats, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) transmission.SessionStats); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(transmission.SessionStats)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// API_GetSessionStatistics_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSessionStatistics'
type API_GetSessionStatistics_Call struct {
	*mock.Call
}

// GetSessionStatistics is a helper method to define mock.On call
//   - ctx context.Context
func (_e *API_Expecter) GetSessionStatistics(ctx interface{}) *API_GetSessionStatistics_Call {
	return &API_GetSessionStatistics_Call{Call: _e.mock.On("GetSessionStatistics", ctx)}
}

func (_c *API_GetSessionStatistics_Call) Run(run func(ctx context.Context)) *API_GetSessionStatistics_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *API_GetSessionStatistics_Call) Return(stats transmission.SessionStats, err error) *API_GetSessionStatistics_Call {
	_c.Call.Return(stats, err)
	return _c
}

func (_c *API_GetSessionStatistics_Call) RunAndReturn(run func(context.Context) (transmission.SessionStats, error)) *API_GetSessionStatistics_Call {
	_c.Call.Return(run)
	return _c
}

// NewAPI creates a new instance of API. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAPI(t interface {
	mock.TestingT
	Cleanup(func())
}) *API {
	mock := &API{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
