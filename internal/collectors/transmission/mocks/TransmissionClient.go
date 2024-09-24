// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	transmissionrpc "github.com/hekmon/transmissionrpc/v3"
)

// TransmissionClient is an autogenerated mock type for the TransmissionClient type
type TransmissionClient struct {
	mock.Mock
}

type TransmissionClient_Expecter struct {
	mock *mock.Mock
}

func (_m *TransmissionClient) EXPECT() *TransmissionClient_Expecter {
	return &TransmissionClient_Expecter{mock: &_m.Mock}
}

// SessionArgumentsGetAll provides a mock function with given fields: ctx
func (_m *TransmissionClient) SessionArgumentsGetAll(ctx context.Context) (transmissionrpc.SessionArguments, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for SessionArgumentsGetAll")
	}

	var r0 transmissionrpc.SessionArguments
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (transmissionrpc.SessionArguments, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) transmissionrpc.SessionArguments); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(transmissionrpc.SessionArguments)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransmissionClient_SessionArgumentsGetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SessionArgumentsGetAll'
type TransmissionClient_SessionArgumentsGetAll_Call struct {
	*mock.Call
}

// SessionArgumentsGetAll is a helper method to define mock.On call
//   - ctx context.Context
func (_e *TransmissionClient_Expecter) SessionArgumentsGetAll(ctx interface{}) *TransmissionClient_SessionArgumentsGetAll_Call {
	return &TransmissionClient_SessionArgumentsGetAll_Call{Call: _e.mock.On("SessionArgumentsGetAll", ctx)}
}

func (_c *TransmissionClient_SessionArgumentsGetAll_Call) Run(run func(ctx context.Context)) *TransmissionClient_SessionArgumentsGetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *TransmissionClient_SessionArgumentsGetAll_Call) Return(sessionArgs transmissionrpc.SessionArguments, err error) *TransmissionClient_SessionArgumentsGetAll_Call {
	_c.Call.Return(sessionArgs, err)
	return _c
}

func (_c *TransmissionClient_SessionArgumentsGetAll_Call) RunAndReturn(run func(context.Context) (transmissionrpc.SessionArguments, error)) *TransmissionClient_SessionArgumentsGetAll_Call {
	_c.Call.Return(run)
	return _c
}

// SessionStats provides a mock function with given fields: ctx
func (_m *TransmissionClient) SessionStats(ctx context.Context) (transmissionrpc.SessionStats, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for SessionStats")
	}

	var r0 transmissionrpc.SessionStats
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (transmissionrpc.SessionStats, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) transmissionrpc.SessionStats); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(transmissionrpc.SessionStats)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransmissionClient_SessionStats_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SessionStats'
type TransmissionClient_SessionStats_Call struct {
	*mock.Call
}

// SessionStats is a helper method to define mock.On call
//   - ctx context.Context
func (_e *TransmissionClient_Expecter) SessionStats(ctx interface{}) *TransmissionClient_SessionStats_Call {
	return &TransmissionClient_SessionStats_Call{Call: _e.mock.On("SessionStats", ctx)}
}

func (_c *TransmissionClient_SessionStats_Call) Run(run func(ctx context.Context)) *TransmissionClient_SessionStats_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *TransmissionClient_SessionStats_Call) Return(stats transmissionrpc.SessionStats, err error) *TransmissionClient_SessionStats_Call {
	_c.Call.Return(stats, err)
	return _c
}

func (_c *TransmissionClient_SessionStats_Call) RunAndReturn(run func(context.Context) (transmissionrpc.SessionStats, error)) *TransmissionClient_SessionStats_Call {
	_c.Call.Return(run)
	return _c
}

// NewTransmissionClient creates a new instance of TransmissionClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransmissionClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *TransmissionClient {
	mock := &TransmissionClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}