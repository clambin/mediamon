// Code generated by mockery v2.36.1. DO NOT EDIT.

package mocks

import (
	context "context"

	xxxarr "github.com/clambin/mediaclients/xxxarr"
	mock "github.com/stretchr/testify/mock"
)

// SonarrClient is an autogenerated mock type for the SonarrClient type
type SonarrClient struct {
	mock.Mock
}

type SonarrClient_Expecter struct {
	mock *mock.Mock
}

func (_m *SonarrClient) EXPECT() *SonarrClient_Expecter {
	return &SonarrClient_Expecter{mock: &_m.Mock}
}

// GetCalendar provides a mock function with given fields: ctx
func (_m *SonarrClient) GetCalendar(ctx context.Context) ([]xxxarr.SonarrCalendar, error) {
	ret := _m.Called(ctx)

	var r0 []xxxarr.SonarrCalendar
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]xxxarr.SonarrCalendar, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []xxxarr.SonarrCalendar); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]xxxarr.SonarrCalendar)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SonarrClient_GetCalendar_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCalendar'
type SonarrClient_GetCalendar_Call struct {
	*mock.Call
}

// GetCalendar is a helper method to define mock.On call
//   - ctx context.Context
func (_e *SonarrClient_Expecter) GetCalendar(ctx interface{}) *SonarrClient_GetCalendar_Call {
	return &SonarrClient_GetCalendar_Call{Call: _e.mock.On("GetCalendar", ctx)}
}

func (_c *SonarrClient_GetCalendar_Call) Run(run func(ctx context.Context)) *SonarrClient_GetCalendar_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *SonarrClient_GetCalendar_Call) Return(_a0 []xxxarr.SonarrCalendar, _a1 error) *SonarrClient_GetCalendar_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SonarrClient_GetCalendar_Call) RunAndReturn(run func(context.Context) ([]xxxarr.SonarrCalendar, error)) *SonarrClient_GetCalendar_Call {
	_c.Call.Return(run)
	return _c
}

// GetEpisodeByID provides a mock function with given fields: ctx, id
func (_m *SonarrClient) GetEpisodeByID(ctx context.Context, id int) (xxxarr.SonarrEpisode, error) {
	ret := _m.Called(ctx, id)

	var r0 xxxarr.SonarrEpisode
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (xxxarr.SonarrEpisode, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) xxxarr.SonarrEpisode); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(xxxarr.SonarrEpisode)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SonarrClient_GetEpisodeByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetEpisodeByID'
type SonarrClient_GetEpisodeByID_Call struct {
	*mock.Call
}

// GetEpisodeByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id int
func (_e *SonarrClient_Expecter) GetEpisodeByID(ctx interface{}, id interface{}) *SonarrClient_GetEpisodeByID_Call {
	return &SonarrClient_GetEpisodeByID_Call{Call: _e.mock.On("GetEpisodeByID", ctx, id)}
}

func (_c *SonarrClient_GetEpisodeByID_Call) Run(run func(ctx context.Context, id int)) *SonarrClient_GetEpisodeByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int))
	})
	return _c
}

func (_c *SonarrClient_GetEpisodeByID_Call) Return(_a0 xxxarr.SonarrEpisode, _a1 error) *SonarrClient_GetEpisodeByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SonarrClient_GetEpisodeByID_Call) RunAndReturn(run func(context.Context, int) (xxxarr.SonarrEpisode, error)) *SonarrClient_GetEpisodeByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetHealth provides a mock function with given fields: ctx
func (_m *SonarrClient) GetHealth(ctx context.Context) ([]xxxarr.SonarrHealth, error) {
	ret := _m.Called(ctx)

	var r0 []xxxarr.SonarrHealth
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]xxxarr.SonarrHealth, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []xxxarr.SonarrHealth); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]xxxarr.SonarrHealth)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SonarrClient_GetHealth_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetHealth'
type SonarrClient_GetHealth_Call struct {
	*mock.Call
}

// GetHealth is a helper method to define mock.On call
//   - ctx context.Context
func (_e *SonarrClient_Expecter) GetHealth(ctx interface{}) *SonarrClient_GetHealth_Call {
	return &SonarrClient_GetHealth_Call{Call: _e.mock.On("GetHealth", ctx)}
}

func (_c *SonarrClient_GetHealth_Call) Run(run func(ctx context.Context)) *SonarrClient_GetHealth_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *SonarrClient_GetHealth_Call) Return(_a0 []xxxarr.SonarrHealth, _a1 error) *SonarrClient_GetHealth_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SonarrClient_GetHealth_Call) RunAndReturn(run func(context.Context) ([]xxxarr.SonarrHealth, error)) *SonarrClient_GetHealth_Call {
	_c.Call.Return(run)
	return _c
}

// GetQueue provides a mock function with given fields: ctx
func (_m *SonarrClient) GetQueue(ctx context.Context) ([]xxxarr.SonarrQueue, error) {
	ret := _m.Called(ctx)

	var r0 []xxxarr.SonarrQueue
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]xxxarr.SonarrQueue, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []xxxarr.SonarrQueue); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]xxxarr.SonarrQueue)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SonarrClient_GetQueue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetQueue'
type SonarrClient_GetQueue_Call struct {
	*mock.Call
}

// GetQueue is a helper method to define mock.On call
//   - ctx context.Context
func (_e *SonarrClient_Expecter) GetQueue(ctx interface{}) *SonarrClient_GetQueue_Call {
	return &SonarrClient_GetQueue_Call{Call: _e.mock.On("GetQueue", ctx)}
}

func (_c *SonarrClient_GetQueue_Call) Run(run func(ctx context.Context)) *SonarrClient_GetQueue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *SonarrClient_GetQueue_Call) Return(_a0 []xxxarr.SonarrQueue, _a1 error) *SonarrClient_GetQueue_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SonarrClient_GetQueue_Call) RunAndReturn(run func(context.Context) ([]xxxarr.SonarrQueue, error)) *SonarrClient_GetQueue_Call {
	_c.Call.Return(run)
	return _c
}

// GetSeries provides a mock function with given fields: ctx
func (_m *SonarrClient) GetSeries(ctx context.Context) ([]xxxarr.SonarrSeries, error) {
	ret := _m.Called(ctx)

	var r0 []xxxarr.SonarrSeries
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]xxxarr.SonarrSeries, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []xxxarr.SonarrSeries); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]xxxarr.SonarrSeries)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SonarrClient_GetSeries_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSeries'
type SonarrClient_GetSeries_Call struct {
	*mock.Call
}

// GetSeries is a helper method to define mock.On call
//   - ctx context.Context
func (_e *SonarrClient_Expecter) GetSeries(ctx interface{}) *SonarrClient_GetSeries_Call {
	return &SonarrClient_GetSeries_Call{Call: _e.mock.On("GetSeries", ctx)}
}

func (_c *SonarrClient_GetSeries_Call) Run(run func(ctx context.Context)) *SonarrClient_GetSeries_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *SonarrClient_GetSeries_Call) Return(_a0 []xxxarr.SonarrSeries, _a1 error) *SonarrClient_GetSeries_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SonarrClient_GetSeries_Call) RunAndReturn(run func(context.Context) ([]xxxarr.SonarrSeries, error)) *SonarrClient_GetSeries_Call {
	_c.Call.Return(run)
	return _c
}

// GetSystemStatus provides a mock function with given fields: ctx
func (_m *SonarrClient) GetSystemStatus(ctx context.Context) (xxxarr.SonarrSystemStatus, error) {
	ret := _m.Called(ctx)

	var r0 xxxarr.SonarrSystemStatus
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (xxxarr.SonarrSystemStatus, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) xxxarr.SonarrSystemStatus); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(xxxarr.SonarrSystemStatus)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SonarrClient_GetSystemStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSystemStatus'
type SonarrClient_GetSystemStatus_Call struct {
	*mock.Call
}

// GetSystemStatus is a helper method to define mock.On call
//   - ctx context.Context
func (_e *SonarrClient_Expecter) GetSystemStatus(ctx interface{}) *SonarrClient_GetSystemStatus_Call {
	return &SonarrClient_GetSystemStatus_Call{Call: _e.mock.On("GetSystemStatus", ctx)}
}

func (_c *SonarrClient_GetSystemStatus_Call) Run(run func(ctx context.Context)) *SonarrClient_GetSystemStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *SonarrClient_GetSystemStatus_Call) Return(_a0 xxxarr.SonarrSystemStatus, _a1 error) *SonarrClient_GetSystemStatus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SonarrClient_GetSystemStatus_Call) RunAndReturn(run func(context.Context) (xxxarr.SonarrSystemStatus, error)) *SonarrClient_GetSystemStatus_Call {
	_c.Call.Return(run)
	return _c
}

// NewSonarrClient creates a new instance of SonarrClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSonarrClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *SonarrClient {
	mock := &SonarrClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}